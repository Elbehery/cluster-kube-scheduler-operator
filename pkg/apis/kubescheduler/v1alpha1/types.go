package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	operatorsv1alpha1api "github.com/openshift/api/operator/v1alpha1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeSchedulerConfig provides information to configure kube-scheduler
type KubeSchedulerConfig struct {
	metav1.TypeMeta `json:",inline"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeSchedulerOperatorConfig provides information to configure an operator to manage kube-scheduler.
type KubeSchedulerOperatorConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   KubeSchedulerOperatorConfigSpec   `json:"spec"`
	Status KubeSchedulerOperatorConfigStatus `json:"status"`
}

type KubeSchedulerOperatorConfigSpec struct {
	operatorsv1alpha1api.OperatorSpec `json:",inline"`

	// userConfig holds a sparse config that the user wants for this component.  It only needs to be the overrides from the defaults
	// it will end up overlaying in the following order:
	// 1. hardcoded default
	// 2. this config
	UserConfig runtime.RawExtension `json:"userConfig"`

	// observedConfig holds a sparse config that controller has observed from the cluster state.  It exists in spec because
	// it causes action for the operator
	ObservedConfig runtime.RawExtension `json:"observedConfig"`
}

type KubeSchedulerOperatorConfigStatus struct {
	operatorsv1alpha1api.OperatorStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeSchedulerOperatorConfigList is a collection of items
type KubeSchedulerOperatorConfigList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items contains the items
	Items []KubeSchedulerOperatorConfig `json:"items"`
}
