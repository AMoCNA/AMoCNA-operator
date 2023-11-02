package controller

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	operatorv1 "kubiki.amocna/operator/api/v1"
)

func getConfigMap(hephaestusDeployment operatorv1.HephaestusDeployment) corev1.ConfigMap {
	return corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      hephaestusDeployment.Name + "-config-map",
			Namespace: hephaestusDeployment.Namespace,
		},
		Data: hephaestusDeployment.Spec.HephaestusGuiConfigMapRaw,
	}
}
