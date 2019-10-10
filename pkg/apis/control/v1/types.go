package v1

import(
	 metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Scaling struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec ScalingSpec `json:"spec"`
}

type ScalingSpec struct {
	TargetDeployment string `json:"targetDeployment"`
	MinReplicas int `json:"minReplicas"`
	MaxReplicas int `json:"maxReplicas"`
	MetricType string `json:"metricType"`
	Step int `json:"step"`
	ScaleUp int`json:"scaleUp"`
	ScaleDown int `json:"scaleDown"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ScalingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items []Scaling `json:"items"`
}
