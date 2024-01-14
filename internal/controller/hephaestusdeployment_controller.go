/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	//"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	operatorv1 "kubiki.amocna/operator/api/v1"
)

// HephaestusDeploymentReconciler reconciles a HephaestusDeployment object
type HephaestusDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func createOrUpdateDeployment(r *HephaestusDeploymentReconciler, ctx context.Context, deploymentCreatingFunc func(operatorv1.HephaestusDeployment, *appsv1.Deployment), hephaestusDeployment operatorv1.HephaestusDeployment, suffix string) error {
	_ = log.FromContext(ctx)
	log.Log.Info("Deploying component", "suffix", suffix)
	// here we create namespace obj cuz it's needed for get
	var name = types.NamespacedName{
		Name:      hephaestusDeployment.Name + suffix,
		Namespace: hephaestusDeployment.Namespace,
	}
	// this var will store empty deployment if no deployment exists or existing one
	var deployment = &appsv1.Deployment{}
	if err := r.Get(ctx, name, deployment); err != nil {
		// need to create deployment
		if apierrors.IsNotFound(err) {
			deploymentCreatingFunc(hephaestusDeployment, deployment)
			log.Log.Info("Component not found, creating")
			if err := r.Create(ctx, deployment); err != nil {
				log.Log.Error(err, "unable to create Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
				return err
			}
		} else {
			return err
		}
	} else {
		// modifying old deployment which we have to do by mutating old object - I hate go
		log.Log.Info("Component found, updating")
		patch := client.MergeFrom(deployment.DeepCopy())
		deploymentCreatingFunc(hephaestusDeployment, deployment)
		if err := r.Patch(ctx, deployment, patch); err != nil {
			log.Log.Error(err, "unable to update Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return err
		}
	}

	return nil
}

//+kubebuilder:rbac:groups=operator.kubiki.amocna,resources=hephaestusdeployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operator.kubiki.amocna,resources=hephaestusdeployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operator.kubiki.amocna,resources=hephaestusdeployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HephaestusDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *HephaestusDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var hephaestusDeployment operatorv1.HephaestusDeployment
	if err := r.Get(ctx, req.NamespacedName, &hephaestusDeployment); err != nil {
		log.Log.Error(err, "unable to fetch Hesphaestus Deployment")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Log.Info("Reconciling Test Hesphaestus Deployment", "Hesphaestus Deployment", hephaestusDeployment)

	// persistent volume
	persistentVolume := getPersistentVolumeDeployment(hephaestusDeployment)
	if err := r.Create(ctx, &persistentVolume); err != nil {
		log.Log.Error(err, "unable to create persistent volume Deployment, will continue the operation", "persistent volume", persistentVolume)
	} else {
		log.Log.Info("Created PV", "PV", persistentVolume.Name)
	}

	// persistent volume claim
	volumeDeployment := getVolumeDeployment(hephaestusDeployment)
	if err := r.Create(ctx, &volumeDeployment); err != nil {
		log.Log.Error(err, "unable to create volume Deployment, will continue the operation", "volume", volumeDeployment)
	}
	log.Log.Info("Created PVC", "PVC", volumeDeployment.Name)

	//config-map
	if hephaestusDeployment.Spec.HephaestusGuiConfigMapRaw != nil {
		configMap := getConfigMap(hephaestusDeployment)
		if err := r.Create(ctx, &configMap); err != nil {
			log.Log.Error(err, "Unable to create config map, will continue the operation", "config map", configMap)
		}
		log.Log.Info("Created Config map", "config map", configMap.Name)
	}
	//gui
	log.FromContext(ctx).Info("GUI Version is ", "HephaestusGuiVersion", hephaestusDeployment.Spec.HephaestusGuiVersion)

	if hephaestusDeployment.Spec.HephaestusGuiVersion == "" {
		log.Log.Info("GUI Version is not set")
	} else {
		log.Log.Info("GUI Version is set", "HephaestusGuiVersion", hephaestusDeployment.Spec.HephaestusGuiVersion)
	}

	if err := createOrUpdateDeployment(r, ctx, getGuiDeployment, hephaestusDeployment, "-gui-deployment"); err != nil {
		return ctrl.Result{}, err
	}
	log.Log.Info("Created Gui")

	//gui-service
	guiService := getGuiService(hephaestusDeployment)
	if err := r.Create(ctx, &guiService); err != nil {
		log.Log.Error(err, "unable to create Gui Service, will continue the operation", "GuiService.Namespace", guiService.Namespace, "GuiService.Name", guiService.Name)
	}
	log.Log.Info("Created Gui Service", "GuiService.Namespace", guiService.Namespace, "GuiService.Name", guiService.Name)

	//execution-controller
	log.FromContext(ctx).Info("Execution Controller Image is ", "ExecutionControllerImage", hephaestusDeployment.Spec.ExecutionControllerImage)

	if hephaestusDeployment.Spec.ExecutionControllerImage == "" {
		log.Log.Info("Execution Controller Image is not set")
	} else {
		log.Log.Info("Execution Controller Image is set", "ExecutionControllerImage", hephaestusDeployment.Spec.ExecutionControllerImage)
	}

	if err := createOrUpdateDeployment(r, ctx, getExecutionControllerDeployment, hephaestusDeployment, "-execution-controller-deployment"); err != nil {
		return ctrl.Result{}, err
	}
	log.Log.Info("Created Execution Controller")

	//execution-controller-service
	executionControllerService := getExecutionControllerService(hephaestusDeployment)
	if err := r.Create(ctx, &executionControllerService); err != nil {
		log.Log.Error(err, "unable to create Execution Controller Service, will continue the operation", "ExecutionControllerService.Namespace", executionControllerService.Namespace, "ExecutionControllerService.Name", executionControllerService.Name)
	}
	log.Log.Info("Created Execution Controller Service", "ExecutionControllerService.Namespace", executionControllerService.Namespace, "ExecutionControllerService.Name", executionControllerService.Name)

	//metrics-adapter
	log.FromContext(ctx).Info("Metrics Adapter Image is ", "MetricsAdapterImage", hephaestusDeployment.Spec.MetricsAdapterImage)
	if hephaestusDeployment.Spec.MetricsAdapterImage == "" {
		log.Log.Info("Metrics Adapter Image is not set")
	} else {
		log.Log.Info("Metrics Adapter Image is set", "MetricsAdapterImage", hephaestusDeployment.Spec.MetricsAdapterImage)
	}

	if err := createOrUpdateDeployment(r, ctx, getMetricsAdapterDeployment, hephaestusDeployment, "-metrics-adapter-deployment"); err != nil {
		return ctrl.Result{}, err
	}
	log.Log.Info("Created Metrics Adapter")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HephaestusDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1.HephaestusDeployment{}).
		Complete(r)
}
