package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodCalculatorSpec defines the desired state of PodCalculator
type PodCalculatorSpec struct {
	// Namespace where to count pods.
	// +required
	Namespace string `json:"namespace"`

	//LabelsSelector to search pods.
	// +optional
	// +nullable
	LabelsSelector map[string]string `json:"labelsSelector"`

	// Reference of the configmap where the pods count should be stored.
	// +optional
	ConfigmapRef *ConfigMapRef `json:"configmapRef"`
}

// PodCalculatorStatus defines the observed state of PodCalculator
type PodCalculatorStatus struct {
	// Pods count.
	// +required
	Count int `json:"count"`
}

type ConfigMapRef struct {
	// Namespace of configmap.
	// +required
	Namespace string `json:"namespace"`

	// Name of configmap.
	// +required
	Name string `json:"name"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Count",type="integer",JSONPath=".status.count",description="Pod count"

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
