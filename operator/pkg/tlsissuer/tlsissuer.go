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

// Package tlsissuer applies the resolved Konflux TLS issuer configuration.
package tlsissuer

import (
	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	certmanagermeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"

	konfluxv1alpha1 "github.com/konflux-ci/konflux-ci/operator/api/v1alpha1"
)

// ConfigureCertificate rewrites a component Certificate only for the two
// locked-down modes. The managed-cluster mode retains the upstream manifest.
func ConfigureCertificate(certificate *certmanagerv1.Certificate, configuration *konfluxv1alpha1.TLSIssuerConfiguration, namespaceLocalIssuer string) {
	if configuration == nil {
		return
	}

	switch configuration.Mode {
	case konfluxv1alpha1.TLSIssuerModeNamespaceLocal:
		certificate.Spec.IssuerRef = certmanagermeta.IssuerReference{
			Group: "cert-manager.io",
			Kind:  "Issuer",
			Name:  namespaceLocalIssuer,
		}
	case konfluxv1alpha1.TLSIssuerModeExistingCluster:
		if configuration.ExistingClusterIssuer == "" {
			return
		}
		certificate.Spec.IssuerRef = certmanagermeta.IssuerReference{
			Group: "cert-manager.io",
			Kind:  "ClusterIssuer",
			Name:  configuration.ExistingClusterIssuer,
		}
	}
}
