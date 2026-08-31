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

package v1alpha1

import "testing"

func TestResolveTLSIssuerConfiguration(t *testing.T) {
	trueValue := true
	falseValue := false

	testCases := []struct {
		name     string
		create   *bool
		existing string
		want     TLSIssuerConfiguration
	}{
		{
			name: "defaults to managed cluster issuers",
			want: TLSIssuerConfiguration{Mode: TLSIssuerModeManagedCluster},
		},
		{
			name:   "uses managed cluster issuers when enabled",
			create: &trueValue,
			want:   TLSIssuerConfiguration{Mode: TLSIssuerModeManagedCluster},
		},
		{
			name:   "uses namespace local issuers when creation is disabled",
			create: &falseValue,
			want:   TLSIssuerConfiguration{Mode: TLSIssuerModeNamespaceLocal},
		},
		{
			name:     "uses an existing cluster issuer when configured",
			create:   &falseValue,
			existing: "platform-ca",
			want: TLSIssuerConfiguration{
				Mode:                  TLSIssuerModeExistingCluster,
				ExistingClusterIssuer: "platform-ca",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := ResolveTLSIssuerConfiguration(testCase.create, testCase.existing)
			if got != testCase.want {
				t.Fatalf("ResolveTLSIssuerConfiguration() = %#v, want %#v", got, testCase.want)
			}
		})
	}
}
