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

package controllers

import (
	"context"
	"fmt"

	managedcontrolplanev1 "github.com/alexander-demicev/k3s-CPaaS/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// SetupWithManager sets up the controller with the Manager.
func (r *ManagedControlPlaneReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&managedcontrolplanev1.ManagedControlPlane{}).
		Complete(r)
}

// ManagedControlPlaneReconciler reconciles a ManagedControlPlane object
type ManagedControlPlaneReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=managedcontrolplane.k3s.control-plane.io,resources=managedcontrolplanes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=managedcontrolplane.k3s.control-plane.io,resources=managedcontrolplanes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=managedcontrolplane.k3s.control-plane.io,resources=managedcontrolplanes/finalizers,verbs=update

func (r *ManagedControlPlaneReconciler) Reconcile(ctx context.Context, req reconcile.Request) (ctrl.Result, error) {
	logger := ctrl.LoggerFrom(ctx)

	logger.Info("Reconciling ManagedControlPlane")

	mcp := &managedcontrolplanev1.ManagedControlPlane{}
	if err := r.Get(ctx, req.NamespacedName, mcp); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("Object was not found, registration client has to create it")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, fmt.Errorf("failed to get machine registration object: %w", err)
	}

	currentReplicas := 0 // TODO: get current replicas
	desiredReplicas := mcp.Spec.Replicas

	switch {
	// We are creating the first replica
	case currentReplicas < desiredReplicas && currentReplicas == 0:
		// Create new control plane w/ init
		return r.initControlPlane(ctx, mcp)
	// We are scaling up
	case currentReplicas < desiredReplicas && currentReplicas > 0:
	// Create a new control plane w/ join
	// ...
	// We are scaling down
	case currentReplicas > desiredReplicas:
	}

	return ctrl.Result{}, nil
}

func (r *ManagedControlPlaneReconciler) initControlPlane(ctx context.Context, mcp *managedcontrolplanev1.ManagedControlPlane) (ctrl.Result, error) {
	if err := createService(ctx, r.Client, mcp); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create service: %w", err)
	}

	if err := createDeployment(ctx, r.Client, mcp); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create deployment: %w", err)
	}

	return ctrl.Result{}, nil
}

func createService(ctx context.Context, cl client.Client, mcp *managedcontrolplanev1.ManagedControlPlane) error {
	internalTrafficPolicy := corev1.ServiceInternalTrafficPolicyCluster
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: mcp.Name,
		},
		Spec: corev1.ServiceSpec{
			ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeCluster,
			InternalTrafficPolicy: &internalTrafficPolicy,
			IPFamilies:            []corev1.IPFamily{corev1.IPv4Protocol},
			Ports: []corev1.ServicePort{
				{
					Name:       "k3s-server",
					Port:       mcp.Spec.Service.Port,
					NodePort:   mcp.Spec.Service.Port,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: int32(mcp.Spec.Service.Port)},
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"managedControlPlane": mcp.Name,
			},
			Type:            corev1.ServiceTypeNodePort,
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
	return cl.Create(ctx, service)
}

func createDeployment(ctx context.Context, cl client.Client, mcp *managedcontrolplanev1.ManagedControlPlane) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      names.SimpleNameGenerator.GenerateName(mcp.Name + "-"),
			Namespace: mcp.Namespace,
			Labels: map[string]string{
				"managedControlPlane": mcp.Name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"managedControlPlane": mcp.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Args:            []string{"server"},
							ImagePullPolicy: corev1.PullAlways,
							Image:           "rancher/k3s:v1.26.1-k3s1", // TODO: customize image
							Name:            "k3s-control-plane",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "node-token",
									MountPath: "/etc/node-token",
								},
								{
									Name:      "etcd-certs",
									MountPath: "/etc/etcd-certs",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "node-token",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName:  "node-token",
									DefaultMode: pointer.Int32(420),
									Items: []corev1.KeyToPath{
										{
											Key:  "node-token",
											Path: "node-token",
										},
									},
								},
							},
						},
						{
							Name: "etcd-certs",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName:  "etcd-certs",
									DefaultMode: pointer.Int32(420),
									Items: []corev1.KeyToPath{
										{
											Key:  "ca.crt",
											Path: "ca.crt",
										},
										{
											Key:  "server.crt",
											Path: "server.crt",
										},
										{
											Key:  "server.key",
											Path: "server.key",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return cl.Create(ctx, deployment)
}
