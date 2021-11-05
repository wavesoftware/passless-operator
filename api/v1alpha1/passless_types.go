/*
Copyright 2021.

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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PassLessSpec defines the desired state of PassLess
type PassLessSpec map[string]PassLessEntry

// PassLessEntry contains a configuration for each secret value to be generated
type PassLessEntry struct {
	// Version defines an incremental number of a passless secret. To change secret
	// increment this number.
	Version uint `json:"version,omitempty"`

	// Scope defines a type of the passless secret.
	Scope string `json:"scope,omitempty"`

	// Length defines a length of the passless secret.
	Length uint8 `json:"length,omitempty"`
}

// PassLessStatus defines the observed state of PassLess
type PassLessStatus string

const (
	// Dirty is when secret isn't in par with passless yet
	Dirty PassLessStatus = "Dirty"

	// Ready is when secret has been reconciled
	Ready PassLessStatus = "Ready"

	// Blocked is when whe can't create secret as there is another user
	// created secret in the way
	Blocked PassLessStatus = "Blocked"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PassLess is the Schema for the passlesses API
type PassLess struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PassLessSpec   `json:"spec,omitempty"`
	Status PassLessStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PassLessList contains a list of PassLess
type PassLessList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PassLess `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PassLess{}, &PassLessList{})
}
