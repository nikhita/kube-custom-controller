package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GithubNotif struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              GithubNotifSpec   `json:"spec"`
	Status            GithubNotifStatus `json:"status,omitempty"`
}

type GithubNotifSpec struct {
	Message string `json:"message"`
}

type GithubNotifStatus struct {
	Delivered bool `json:"delivered"`
}

type GithubNotifList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []GithubNotif `json:"items"`
}
