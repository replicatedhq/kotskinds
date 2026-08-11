package v1beta1tests

import (
	"testing"

	kotsv1beta1 "github.com/replicatedhq/kotskinds/apis/kots/v1beta1"
	kotsscheme "github.com/replicatedhq/kotskinds/client/kotsclientset/scheme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/scheme"
)

func Test_Terraform(t *testing.T) {
	data := `apiVersion: kots.io/v1beta1
kind: Terraform
metadata:
  name: aws-infra
spec:
  # name is the stable, vendor-defined module name
  name: "aws-infrastructure"

  # provider is the Terraform provider
  provider: "aws"

  # sourceRepository is the git repository URL containing the Terraform module source
  sourceRepository: "https://github.com/acme-corp/my-application"

  # sourcePath is the path within the git repository to the Terraform module
  sourcePath: "terraform/aws"

  # sourceRef is the git reference (tag, branch, or commit) to use
  sourceRef: "v1.0.0"

  # minTerraformVersion specifies the minimum Terraform version required
  minTerraformVersion: "1.0.0"

  # docs provides documentation for various user-defined steps
  docs:
    prerequisites: |
      Before installing this Terraform configuration, ensure you have:

      1. AWS CLI configured with appropriate credentials
      2. Terraform >= 1.0.0 installed
      3. Required IAM permissions for the resources being created
      4. A valid AWS region selected
      5. Network access to AWS APIs

      Verify your setup:
        aws sts get-caller-identity
        terraform version
    install: |
      This step initializes and applies the Terraform module to create AWS infrastructure.

      Steps:
      1. Initialize Terraform:
         terraform init

      2. Review the execution plan:
         terraform plan -out=tfplan

      3. Apply the changes:
         terraform apply tfplan

      The module will create the following AWS resources:
      - VPC and networking components
      - Security groups
      - EC2 instances or other compute resources
      - IAM roles and policies as needed
`

	kotsscheme.AddToScheme(scheme.Scheme)

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, gvk, err := decode([]byte(data), nil, nil)
	require.NoError(t, err)

	assert.Equal(t, "kots.io", gvk.Group)
	assert.Equal(t, "v1beta1", gvk.Version)
	assert.Equal(t, "Terraform", gvk.Kind)

	terraform := obj.(*kotsv1beta1.Terraform)

	assert.Equal(t, "aws-infra", terraform.Name)
	assert.Equal(t, "aws-infrastructure", terraform.Spec.Name)
	assert.Equal(t, "aws", terraform.Spec.Provider)
	assert.Equal(t, "https://github.com/acme-corp/my-application", terraform.Spec.SourceRepository)
	assert.Equal(t, "terraform/aws", terraform.Spec.SourcePath)
	assert.Equal(t, "v1.0.0", terraform.Spec.SourceRef)
	assert.Equal(t, "1.0.0", terraform.Spec.MinTerraformVersion)

	require.NotNil(t, terraform.Spec.Docs)
	assert.Contains(t, terraform.Spec.Docs, "prerequisites")
	assert.Contains(t, terraform.Spec.Docs, "install")
	assert.Contains(t, terraform.Spec.Docs["prerequisites"], "AWS CLI configured")
	assert.Contains(t, terraform.Spec.Docs["install"], "terraform init")
}

func Test_Terraform_Minimal(t *testing.T) {
	data := `apiVersion: kots.io/v1beta1
kind: Terraform
metadata:
  name: simple-terraform
spec:
  name: "simple-module"
  provider: "aws"
  sourceRepository: "https://github.com/vendor/simple-app"
`

	kotsscheme.AddToScheme(scheme.Scheme)

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, gvk, err := decode([]byte(data), nil, nil)
	require.NoError(t, err)

	assert.Equal(t, "kots.io", gvk.Group)
	assert.Equal(t, "v1beta1", gvk.Version)
	assert.Equal(t, "Terraform", gvk.Kind)

	terraform := obj.(*kotsv1beta1.Terraform)

	assert.Equal(t, "simple-terraform", terraform.Name)
	assert.Equal(t, "simple-module", terraform.Spec.Name)
	assert.Equal(t, "aws", terraform.Spec.Provider)
	assert.Equal(t, "https://github.com/vendor/simple-app", terraform.Spec.SourceRepository)
	assert.Empty(t, terraform.Spec.SourcePath)
	assert.Empty(t, terraform.Spec.SourceRef)
	assert.Empty(t, terraform.Spec.MinTerraformVersion)
	assert.Nil(t, terraform.Spec.Docs)
}

func Test_Terraform_WithDocsOnly(t *testing.T) {
	data := `apiVersion: kots.io/v1beta1
kind: Terraform
metadata:
  name: terraform-with-docs
spec:
  name: "documented-module"
  provider: "gcp"
  sourceRepository: "https://github.com/vendor/documented-app"
  sourcePath: "terraform/gcp"
  sourceRef: "v2.1.0"
  docs:
    setup: "Run terraform init and terraform apply to set up the infrastructure."
    teardown: "Run terraform destroy to remove all resources."
`

	kotsscheme.AddToScheme(scheme.Scheme)

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, gvk, err := decode([]byte(data), nil, nil)
	require.NoError(t, err)

	assert.Equal(t, "kots.io", gvk.Group)
	assert.Equal(t, "v1beta1", gvk.Version)
	assert.Equal(t, "Terraform", gvk.Kind)

	terraform := obj.(*kotsv1beta1.Terraform)

	assert.Equal(t, "terraform-with-docs", terraform.Name)
	assert.Equal(t, "documented-module", terraform.Spec.Name)
	assert.Equal(t, "gcp", terraform.Spec.Provider)
	assert.Equal(t, "https://github.com/vendor/documented-app", terraform.Spec.SourceRepository)
	assert.Equal(t, "terraform/gcp", terraform.Spec.SourcePath)
	assert.Equal(t, "v2.1.0", terraform.Spec.SourceRef)
	assert.Empty(t, terraform.Spec.MinTerraformVersion)

	require.NotNil(t, terraform.Spec.Docs)
	assert.Equal(t, "Run terraform init and terraform apply to set up the infrastructure.", terraform.Spec.Docs["setup"])
	assert.Equal(t, "Run terraform destroy to remove all resources.", terraform.Spec.Docs["teardown"])
}
