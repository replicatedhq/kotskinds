/*
Copyright 2019 Replicated, Inc..

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1beta1 "github.com/replicatedhq/kotskinds/client/kotsclientset/typed/kots/v1beta1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeKotsV1beta1 struct {
	*testing.Fake
}

func (c *FakeKotsV1beta1) Airgaps(namespace string) v1beta1.AirgapInterface {
	return &FakeAirgaps{c, namespace}
}

func (c *FakeKotsV1beta1) Applications(namespace string) v1beta1.ApplicationInterface {
	return &FakeApplications{c, namespace}
}

func (c *FakeKotsV1beta1) Configs(namespace string) v1beta1.ConfigInterface {
	return &FakeConfigs{c, namespace}
}

func (c *FakeKotsV1beta1) ConfigValueses(namespace string) v1beta1.ConfigValuesInterface {
	return &FakeConfigValueses{c, namespace}
}

func (c *FakeKotsV1beta1) HelmCharts(namespace string) v1beta1.HelmChartInterface {
	return &FakeHelmCharts{c, namespace}
}

func (c *FakeKotsV1beta1) Identities(namespace string) v1beta1.IdentityInterface {
	return &FakeIdentities{c, namespace}
}

func (c *FakeKotsV1beta1) IdentityConfigs(namespace string) v1beta1.IdentityConfigInterface {
	return &FakeIdentityConfigs{c, namespace}
}

func (c *FakeKotsV1beta1) IngressConfigs(namespace string) v1beta1.IngressConfigInterface {
	return &FakeIngressConfigs{c, namespace}
}

func (c *FakeKotsV1beta1) Installations(namespace string) v1beta1.InstallationInterface {
	return &FakeInstallations{c, namespace}
}

func (c *FakeKotsV1beta1) Licenses(namespace string) v1beta1.LicenseInterface {
	return &FakeLicenses{c, namespace}
}

func (c *FakeKotsV1beta1) LintConfigs(namespace string) v1beta1.LintConfigInterface {
	return &FakeLintConfigs{c, namespace}
}

func (c *FakeKotsV1beta1) ReplicatedHelmCharts(namespace string) v1beta1.ReplicatedHelmChartInterface {
	return &FakeReplicatedHelmCharts{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeKotsV1beta1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
