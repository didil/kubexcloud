/*


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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AppSpec defines the desired state of App
type AppSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	Replicas int32 `json:"replicas"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Containers []Container `json:"containers"`
}

// Container object
type Container struct {
	// +kubebuilder:validation:Required
	Image string `json:"image"`
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	Command []string `json:"command,omitempty"`

	Ports []Port `json:"ports,omitempty"`
}

// Port object
type Port struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	Number int32 `json:"number"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=TCP;UDP
	Protocol corev1.Protocol `json:"protocol"`

	// only valid for http (through TCP protocol)
	ExposeExternally bool `json:"exposeExternally"`
}

// AppStatus defines the observed state of App
type AppStatus struct {
	ExternalURL         string `json:"externalUrl,omitempty"`
	AvailableReplicas   int32  `json:"availableReplicas,omitempty"`
	UnavailableReplicas int32  `json:"unavailableReplicas,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="ExternalURL",type=string,JSONPath=`.status.externalUrl`
// +kubebuilder:printcolumn:name="AvailableReplicas",type=integer,JSONPath=`.status.availableReplicas`
// +kubebuilder:printcolumn:name="UnavailableReplicas",type=integer,JSONPath=`.status.unavailableReplicas`

// App is the Schema for the apps API
type App struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppSpec   `json:"spec,omitempty"`
	Status AppStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AppList contains a list of App
type AppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []App `json:"items"`
}

func init() {
	SchemeBuilder.Register(&App{}, &AppList{})
}
