package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	operatorv1 "kubiki.amocna/operator/api/v1"
)

func getExecutionControllerDeployment(hephaestusDeployment operatorv1.HephaestusDeployment) appsv1.Deployment {
	one := int32(1)
	port := getPortOrDefault(hephaestusDeployment.Spec.ExecutionControllerInternalPort, 8097)
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      hephaestusDeployment.Name + "-execution-controller-deployment",
			Namespace: hephaestusDeployment.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &one,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": hephaestusDeployment.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": hephaestusDeployment.Name,
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: hephaestusDeployment.Spec.ExecutionControllerServiceAccountName,
					Containers: []corev1.Container{
						{
							Name:            hephaestusDeployment.Name,
							Image:           hephaestusDeployment.Spec.ExecutionControllerImage,
							ImagePullPolicy: corev1.PullPolicy("Always"),
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: port,
								},
							},
						},
					},
				},
			},
		},
	}
}
