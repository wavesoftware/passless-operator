package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// +genclient

// PassLessSpec defines the desired state of PassLess
type PassLessSpec map[string]PassLessEntry

// PassLessEntry contains a configuration for each secret value to be generated
type PassLessEntry struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Version defines a incremental number of a passless secret. To change secret
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
	Dirty   PassLessStatus = "Dirty"

	// Ready is when secret has been reconciled
	Ready   PassLessStatus = "Ready"

	// Blocked is when whe can't create secret as there is another user created secret in the way
	Blocked PassLessStatus = "Blocked"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PassLess is the Schema for the passlesses API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=passlesses,scope=Namespaced
type PassLess struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PassLessSpec   `json:"spec,omitempty"`
	Status PassLessStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PassLessList contains a list of PassLess
type PassLessList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PassLess `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PassLess{}, &PassLessList{})
}
