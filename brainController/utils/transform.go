package utils

import (
	apiv1 "k8s.io/api/core/v1"
)

func Transform(env map[string]string) []apiv1.EnvVar {
	var result []apiv1.EnvVar
	for k, v := range env {
		t := apiv1.EnvVar{
			Name:  k,
			Value: v,
		}
		result = append(result, t)
	}
	result = append(result, apiv1.EnvVar{
		Name: "DYNACONF_NODE",
		ValueFrom: &apiv1.EnvVarSource{
			FieldRef: &apiv1.ObjectFieldSelector{
				FieldPath: "spec.nodeName",
			},
		},
	})
	return result
}
