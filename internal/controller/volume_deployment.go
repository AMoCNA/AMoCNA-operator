package controller

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	operatorv1 "kubiki.amocna/operator/api/v1"
)

func getVolumeDeployment(hephaestusDeployment operatorv1.HephaestusDeployment) corev1.PersistentVolumeClaim {
	var storageClassName = "standard"
	return corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      hephaestusDeployment.Name + "-volume-claim",
			Namespace: hephaestusDeployment.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassName,
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceName(corev1.ResourceStorage): resource.MustParse("25Mi"),
				},
			},
		},
	}
}
