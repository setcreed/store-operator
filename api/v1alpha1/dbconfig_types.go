/*
Copyright 2024.

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

// DbConfigSpec defines the desired state of DbConfig
type DbConfigSpec struct {
	// Number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1.
	// +optional
	//+kubebuilder:default:=1
	//+kubebuilder:validation:Minimum:=1
	Replicas int `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`

	// connect dsn
	//+kubebuilder:validation:Required
	Dsn string `json:"dsn,omitempty"`
	//+kubebuilder:default:=15
	//+kubebuilder:validation:Minimum:=1
	//+kubebuilder:validation:Maximum:=2000
	MaxOpenConn int `json:"maxOpenConn,omitempty"`
	//+kubebuilder:default:=600
	//+kubebuilder:validation:Minimum:=60
	MaxLifeTime int `json:"maxLifeTime,omitempty"`
	//+kubebuilder:default:=5
	//+kubebuilder:validation:Minimum:=1
	//+kubebuilder:validation:Maximum:=2000
	MaxIdleConn int `json:"maxIdleConn,omitempty"`
}

// DbConfigStatus defines the observed state of DbConfig
type DbConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Replicas int32  `json:"replicas,omitempty"`
	Ready    string `json:"ready,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.labelSelector
//+kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="The readiness of the CR"
//+kubebuilder:printcolumn:name="最大连接数",type="integer",JSONPath=".spec.maxOpenConn",description="最大连接数"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The time when the resource was created"

// +kubebuilder:resource:path=dbconfigs,scope=Namespaced,shortName=dc
// DbConfig is the Schema for the dbconfigs API
type DbConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DbConfigSpec   `json:"spec,omitempty"`
	Status DbConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DbConfigList contains a list of DbConfig
type DbConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DbConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DbConfig{}, &DbConfigList{})
}
