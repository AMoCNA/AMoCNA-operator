package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	operatorv1 "kubiki.amocna/operator/api/v1"
)

func getGuiDeployment(hephaestusDeployment operatorv1.HephaestusDeployment, deployment *appsv1.Deployment) {
	var shouldMountConfigMap = hephaestusDeployment.Spec.HephaestusGuiConfigMapRaw != nil
	one := int32(1)
	port := getPortOrDefault(hephaestusDeployment.Spec.HephaestusGuiInternalPort, 8080)
	if shouldMountConfigMap {
		deployment.ObjectMeta = metav1.ObjectMeta{
			Name:      hephaestusDeployment.Name + "-gui-deployment",
			Namespace: hephaestusDeployment.Namespace,
		}
		deployment.Spec = appsv1.DeploymentSpec{
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
							Image: "hephaestusmetrics/gui:" + hephaestusDeployment.Spec.HephaestusGuiVersion,
							Env: []corev1.EnvVar{
								{
									Name:  "prometheus.address",
									Value: hephaestusDeployment.Spec.PrometheusAddress,
								},
								{
									Name:  "saved.path",
									Value: "/../storage/metrics/savedMetrics/metrics.json",
								},
								{
									Name:  "config.path",
									Value: "/../storage/metrics/configMetrics/metrics.json",
								},
								{
									Name:  "logs.path",
									Value: "/../storage/logs",
								},
							},
							ImagePullPolicy: corev1.PullPolicy("Always"),
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: port,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "storage",
									MountPath: "storage",
								},
								{
									Name:      "config-volume",
									MountPath: "storage/metrics/configMetrics",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: hephaestusDeployment.Name + "-volume-claim",
								},
							},
						},
						{
							Name: "config-volume",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: hephaestusDeployment.Name + "-config-map",
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "metrics.json",
											Path: "metrics.json",
										},
									},
								},
							},
						},
					},
				},
			},
		}

	} else {
		deployment.ObjectMeta = metav1.ObjectMeta{
			Name:      hephaestusDeployment.Name + "-gui-deployment",
			Namespace: hephaestusDeployment.Namespace,
		}
		deployment.Spec = appsv1.DeploymentSpec{
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
							Image: "hephaestusmetrics/gui:" + hephaestusDeployment.Spec.HephaestusGuiVersion,
							Env: []corev1.EnvVar{
								{
									Name:  "prometheus.address",
									Value: hephaestusDeployment.Spec.PrometheusAddress,
								},
								{
									Name:  "saved.path",
									Value: "/../storage/metrics/savedMetrics/metrics.json",
								},
								{
									Name:  "config.path",
									Value: "/../storage/metrics/configMetrics/metrics.json",
								},
								{
									Name:  "logs.path",
									Value: "/../storage/logs",
								},
							},
							ImagePullPolicy: corev1.PullPolicy("Always"),
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: port,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "storage",
									MountPath: "storage",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: hephaestusDeployment.Name + "-volume-claim",
								},
							},
						},
					},
				},
			},
		}

	}
}
