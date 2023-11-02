/*
Copyright 2023.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HephaestusDeploymentSpec defines the desired state of HephaestusDeployment
type HephaestusDeploymentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	HephaestusGuiVersion           string            `json:"hephaestusGuiVersion"`
	HephaestusGuiConfigMapRaw      map[string]string `json:"hephaestusGuiConfigMapRaw,omitempty"`
	HephaestusGuiConfigMapFilePath string            `json:"hephaestusGuiConfigMapFilePath,omitempty"`
	// +kubebuilder:validation:Required
	PrometheusAddress string `json:"prometheusAddress"`
	// +kubebuilder:validation:Required
	ExecutionControllerImage string `json:"executionControllerImage"`
	// +kubebuilder:validation:Required
	MetricsAdapterImage                   string `json:"metricsAdapterImage"`
	MetricsAdapterInternalPort            int32  `json:"metricsAdapterInternalPort,omitempty"`
	HephaestusGuiInternalPort             int32  `json:"hephaestusGuiInternalPort,omitempty"`
	HephaestusGuiExternalPort             int32  `json:"hephaestusGuiExternalPort,omitempty"`
	ExecutionControllerInternalPort       int32  `json:"executionControllerInternalPort,omitempty"`
	ExecutionControllerServiceAccountName string `json:"executionControllerServiceAccountName,omitempty"`
}

// HephaestusDeploymentStatus defines the observed state of HephaestusDeployment
type HephaestusDeploymentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HephaestusDeployment is the Schema for the hephaestusdeployments API
type HephaestusDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HephaestusDeploymentSpec   `json:"spec,omitempty"`
	Status HephaestusDeploymentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HephaestusDeploymentList contains a list of HephaestusDeployment
type HephaestusDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HephaestusDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HephaestusDeployment{}, &HephaestusDeploymentList{})
}
