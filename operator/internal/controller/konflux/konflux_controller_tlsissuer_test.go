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

package konflux

import (
	"testing"

	konfluxv1alpha1 "github.com/konflux-ci/konflux-ci/operator/api/v1alpha1"
)

func TestTLSIssuerConfiguration(t *testing.T) {
	falseValue := false

	testCases := []struct {
		name  string
		owner *konfluxv1alpha1.Konflux
		want  *konfluxv1alpha1.TLSIssuerConfiguration
	}{
		{
			name:  "leaves legacy child specifications unchanged",
			owner: &konfluxv1alpha1.Konflux{},
		},
		{
			name: "propagates namespace local mode",
			owner: &konfluxv1alpha1.Konflux{Spec: konfluxv1alpha1.KonfluxSpec{
				CertManager: &konfluxv1alpha1.CertManagerConfig{CreateClusterIssuer: &falseValue},
			}},
			want: &konfluxv1alpha1.TLSIssuerConfiguration{
				Mode: konfluxv1alpha1.TLSIssuerModeNamespaceLocal,
			},
		},
		{
			name: "propagates an existing cluster issuer",
			owner: &konfluxv1alpha1.Konflux{Spec: konfluxv1alpha1.KonfluxSpec{
				CertManager: &konfluxv1alpha1.CertManagerConfig{
					CreateClusterIssuer:   &falseValue,
					ExistingClusterIssuer: "platform-ca",
				},
			}},
			want: &konfluxv1alpha1.TLSIssuerConfiguration{
				Mode:                  konfluxv1alpha1.TLSIssuerModeExistingCluster,
				ExistingClusterIssuer: "platform-ca",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := tlsIssuerConfiguration(testCase.owner)
			if (got == nil) != (testCase.want == nil) {
				t.Fatalf("tlsIssuerConfiguration() = %#v, want %#v", got, testCase.want)
			}
			if got != nil && *got != *testCase.want {
				t.Fatalf("tlsIssuerConfiguration() = %#v, want %#v", got, testCase.want)
			}
		})
	}
}
