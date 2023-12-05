package common

import (
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

type Replicator interface {
	Run()
	Synced() bool
	NamespaceAdded(ns *v1.Namespace)
}

func PreviouslyPresentKeys(object *metav1.ObjectMeta) (map[string]struct{}, bool) {
	keyList, ok := object.Annotations[ReplicatedKeysAnnotation]
	if !ok {
		return nil, false
	}

	keys := strings.Split(keyList, ",")
	out := make(map[string]struct{})

	for _, k := range keys {
		out[k] = struct{}{}
	}

	return out, true
}

func BuildStrictRegex(regex string) string {
	reg := strings.TrimSpace(regex)
	if !strings.HasPrefix(reg, "^") {
		reg = "^" + reg
	}
	if !strings.HasSuffix(reg, "$") {
		reg = reg + "$"
	}
	return reg
}

func JSONPatchPathEscape(annotation string) string {
	return strings.ReplaceAll(annotation, "/", "~1")
}

func BuildOwnerReferences(objectMeta *metav1.ObjectMeta, apiVersion, kind string) []metav1.OwnerReference {
	blockOwnerDeletion := false
	isController := false
	return []metav1.OwnerReference{
		{
			APIVersion:         apiVersion,
			Kind:               kind,
			BlockOwnerDeletion: &blockOwnerDeletion,
			Controller:         &isController,
			Name:               objectMeta.Name,
			UID:                objectMeta.UID,
		},
	}
}

func GetGVK(obj runtime.Object) (string, string, error) {
	gvk, _, err := scheme.Scheme.ObjectKinds(obj)
	if err != nil {
		return "", "", err
	}
	return gvk[0].GroupVersion().String(), gvk[0].Kind, nil
}
