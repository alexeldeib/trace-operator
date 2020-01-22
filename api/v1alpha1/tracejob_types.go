package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// TraceJobSpec defines the desired state of TraceJob
type TraceJobSpec struct {
	// Program is a string literal to evaluate as a bpftrace program.
	Program             string  `json:"program"`
	Hostname            string  `json:"hostname"`
	ServiceAccount      *string `json:"serviceAccount,omitempty"`      // +optional
	ImageNameTag        *string `json:"imageNameTag,omitempty"`        // +optional
	InitImageNameTag    *string `json:"initImageNameTag,omitempty"`    // +optional
	FetchHeaders        bool    `json:"fetchHeaders,omitempty"`        // +optional
	Deadline            *int64  `json:"deadline,omitempty"`            // +optional
	DeadlineGracePeriod *int64  `json:"deadlineGracePeriod,omitempty"` // +optional
}

// TraceJobStatus defines the observed state of TraceJob
type TraceJobStatus struct {
	// ID is a generated UUID for this object.
	ID    types.UID `json:"id,omitempty"`
	State *string   `json:"state,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// TraceJob is the Schema for the tracejobs API
type TraceJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TraceJobSpec   `json:"spec,omitempty"`
	Status TraceJobStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TraceJobList contains a list of TraceJob
type TraceJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TraceJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TraceJob{}, &TraceJobList{})
}
