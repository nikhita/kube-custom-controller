package github

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Comment struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   CommentSpec
	Status CommentStatus
}

type CommentSpec struct {
	Message string
}

type CommentStatus struct {
	Created bool
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CommentList struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Items []Comment
}
