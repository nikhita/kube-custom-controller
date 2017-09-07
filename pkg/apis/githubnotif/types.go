package githubnotif

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GithubNotif struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   GithubNotifSpec
	Status GithubNotifStatus
}

type GithubNotifSpec struct {
	Message string
}

type GithubNotifStatus struct {
	Delivered bool
}

type GithubNotifList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []GithubNotif
}
