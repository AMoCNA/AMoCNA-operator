package controller

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	operatorv1 "kubiki.amocna/operator/api/v1"
)

func getExecutionControllerService(hephaestusDeployment operatorv1.HephaestusDeployment) corev1.Service {
	internalPort := getPortOrDefault(hephaestusDeployment.Spec.ExecutionControllerInternalPort, 8097)
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      hephaestusDeployment.Name + "-hephaestus-exec-ctrl-service",
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
			}},
		},
	}
}
