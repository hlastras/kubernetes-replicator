package common

import (
	"encoding/json"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func BuildOwnerReferences(objectMeta *metav1.ObjectMeta, typeMeta *metav1.TypeMeta) []metav1.OwnerReference {
	fmt.Println(">>>>>222BuildOwnerReferences")
	fmt.Println("APIVersion: ", typeMeta.APIVersion)
	fmt.Println("Kind: ", typeMeta.Kind)
	fmt.Println("Name: ", objectMeta.Name)
	fmt.Println("UID: ", objectMeta.UID)

	b, err := json.MarshalIndent(typeMeta, "", "  ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Println(string(b))

	blockOwnerDeletion := false
	isController := false
	return []metav1.OwnerReference{
		{
			APIVersion:         typeMeta.APIVersion,
			Kind:               typeMeta.Kind,
			BlockOwnerDeletion: &blockOwnerDeletion,
			Controller:         &isController,
			Name:               objectMeta.Name,
			UID:                objectMeta.UID,
		},
	}
}
