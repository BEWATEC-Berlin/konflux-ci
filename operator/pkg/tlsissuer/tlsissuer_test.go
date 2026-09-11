/*
Copyright 2026 BEWATEC-Berlin.

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

package tlsissuer

import (
	"testing"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	certmanagermeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"

	konfluxv1alpha1 "github.com/konflux-ci/konflux-ci/operator/api/v1alpha1"
)

func TestConfigureCertificate(t *testing.T) {
	testCases := []struct {
		name          string
		configuration *konfluxv1alpha1.TLSIssuerConfiguration
		want          certmanagermeta.IssuerReference
	}{
		{
			name: "preserves the manifest issuer without a configuration",
			want: certmanagermeta.IssuerReference{
				Group: "cert-manager.io", Kind: "ClusterIssuer", Name: "upstream-issuer",
			},
		},
		{
			name: "preserves the manifest issuer for managed mode",
			configuration: &konfluxv1alpha1.TLSIssuerConfiguration{
				Mode: konfluxv1alpha1.TLSIssuerModeManagedCluster,
			},
			want: certmanagermeta.IssuerReference{
				Group: "cert-manager.io", Kind: "ClusterIssuer", Name: "upstream-issuer",
			},
		},
		{
			name: "uses a namespace local issuer",
			configuration: &konfluxv1alpha1.TLSIssuerConfiguration{
				Mode: konfluxv1alpha1.TLSIssuerModeNamespaceLocal,
			},
			want: certmanagermeta.IssuerReference{
				Group: "cert-manager.io", Kind: "Issuer", Name: "component-selfsigned-issuer",
			},
		},
		{
			name: "uses the named external cluster issuer",
			configuration: &konfluxv1alpha1.TLSIssuerConfiguration{
				Mode:                  konfluxv1alpha1.TLSIssuerModeExistingCluster,
				ExistingClusterIssuer: "platform-ca",
			},
			want: certmanagermeta.IssuerReference{
				Group: "cert-manager.io", Kind: "ClusterIssuer", Name: "platform-ca",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			certificate := &certmanagerv1.Certificate{
				Spec: certmanagerv1.CertificateSpec{
					IssuerRef: certmanagermeta.IssuerReference{
						Group: "cert-manager.io", Kind: "ClusterIssuer", Name: "upstream-issuer",
					},
				},
			}

			ConfigureCertificate(certificate, testCase.configuration, "component-selfsigned-issuer")
			if certificate.Spec.IssuerRef != testCase.want {
				t.Fatalf("issuerRef = %#v, want %#v", certificate.Spec.IssuerRef, testCase.want)
			}
		})
	}
}
