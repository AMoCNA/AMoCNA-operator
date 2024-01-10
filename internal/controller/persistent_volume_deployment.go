package controller

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	operatorv1 "kubiki.amocna/operator/api/v1"
)

func getPersistentVolumeDeployment(hephaestusDeployment operatorv1.HephaestusDeployment) corev1.PersistentVolume {
	return corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      hephaestusDeployment.Name + "-gui-pv",
			Namespace: hephaestusDeployment.Namespace,
		},
		Spec: corev1.PersistentVolumeSpec{
			StorageClassName: "standard",
			Capacity: corev1.ResourceList{
				corev1.ResourceName(corev1.ResourceStorage): resource.MustParse("25Mi"),
			},
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},

			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/mnt/hephaestus-gui-pv",
					Type: nil,
				},
			},
		},
	}
}
