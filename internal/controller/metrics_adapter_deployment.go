package controller

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	operatorv1 "kubiki.amocna/operator/api/v1"
)

func getMetricsAdapterDeployment(hephaestusDeployment operatorv1.HephaestusDeployment) appsv1.Deployment {
	one := int32(1)
	guiPort := getPortOrDefault(hephaestusDeployment.Spec.HephaestusGuiInternalPort, 8080)
	executionControllerPort := getPortOrDefault(hephaestusDeployment.Spec.ExecutionControllerInternalPort, 8097)
	metricsAdapterPort := getPortOrDefault(hephaestusDeployment.Spec.MetricsAdapterInternalPort, 8085)
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      hephaestusDeployment.Name + "-metrics-adapter-deployment",
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
					Containers: []corev1.Container{
						{
							Name:  hephaestusDeployment.Name,
							Image: hephaestusDeployment.Spec.MetricsAdapterImage,
							Env: []corev1.EnvVar{
								{
									Name:  "backend",
									Value: "http://" + hephaestusDeployment.Name + "-hephaestus-gui-service." + hephaestusDeployment.Namespace + ":" + fmt.Sprint(guiPort),
								},
								{
									Name:  "kubernetes-management",
									Value: "http://" + hephaestusDeployment.Name + "-hephaestus-exec-ctrl-service." + hephaestusDeployment.Namespace + ":" + fmt.Sprint(executionControllerPort),
								},
							},
							ImagePullPolicy: corev1.PullPolicy("Always"),
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: metricsAdapterPort,
								},
							},
						},
					},
				},
			},
		},
	}
}
