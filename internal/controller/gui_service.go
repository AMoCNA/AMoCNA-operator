package controller

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	operatorv1 "kubiki.amocna/operator/api/v1"
)

func getGuiService(hephaestusDeployment operatorv1.HephaestusDeployment) corev1.Service {
	internalPort := getPortOrDefault(hephaestusDeployment.Spec.HephaestusGuiInternalPort, 8080)
	externalPort := getPortOrDefault(hephaestusDeployment.Spec.HephaestusGuiExternalPort, 31122)
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      hephaestusDeployment.Name + "-hephaestus-gui-service",
			Namespace: hephaestusDeployment.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": hephaestusDeployment.Name,
			},
			Type: "NodePort",
			Ports: []corev1.ServicePort{{
				Protocol: "TCP",
				Port:     internalPort,
				TargetPort: intstr.IntOrString{
					IntVal: internalPort,
				},
				NodePort: externalPort,
			}},
		},
	}
}
