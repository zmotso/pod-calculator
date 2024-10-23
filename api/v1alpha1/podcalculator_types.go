package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodCalculatorSpec defines the desired state of PodCalculator
type PodCalculatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of PodCalculator. Edit podcalculator_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// PodCalculatorStatus defines the observed state of PodCalculator
type PodCalculatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PodCalculator is the Schema for the podcalculators API
type PodCalculator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodCalculatorSpec   `json:"spec,omitempty"`
	Status PodCalculatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PodCalculatorList contains a list of PodCalculator
type PodCalculatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodCalculator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodCalculator{}, &PodCalculatorList{})
}
