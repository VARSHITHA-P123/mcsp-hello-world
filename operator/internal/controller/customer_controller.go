/*
Copyright 2026.
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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	mcspv1 "github.com/VARSHITHA-P123/mcsp-hello-world/operator/api/v1"
)

// CustomerReconciler reconciles a Customer object
type CustomerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=mcsp.mcsp.io,resources=customers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mcsp.mcsp.io,resources=customers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mcsp.mcsp.io,resources=customers/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=route.openshift.io,resources=routes,verbs=get;list;watch;create;update;patch;delete

func (r *CustomerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Step 1 — Get the Customer CR
	customer := &mcspv1.Customer{}
	err := r.Get(ctx, req.NamespacedName, customer)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Customer not found, might have been deleted")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	customerName := customer.Spec.CustomerName
	log.Info("Reconciling Customer", "customerName", customerName)

	// Step 2 — Create Namespace
	namespace := &corev1.Namespace{}
	err = r.Get(ctx, types.NamespacedName{Name: customerName}, namespace)
	if err != nil && errors.IsNotFound(err) {
		namespace = &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: customerName,
				Labels: map[string]string{
					"tenant":   customerName,
					"customer": "true",
				},
			},
		}
		err = r.Create(ctx, namespace)
		if err != nil {
			log.Error(err, "Failed to create namespace")
			return ctrl.Result{}, err
		}
		log.Info("Namespace created", "namespace", customerName)
	}

	// Step 3 — Add Image Puller Permissions
	roleBinding := &rbacv1.RoleBinding{}
	err = r.Get(ctx, types.NamespacedName{Name: customerName + "-image-puller", Namespace: "learning-workspace"}, roleBinding)
	if err != nil && errors.IsNotFound(err) {
		roleBinding = &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      customerName + "-image-puller",
				Namespace: "learning-workspace",
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     "system:image-puller",
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      "default",
					Namespace: customerName,
				},
			},
		}
		err = r.Create(ctx, roleBinding)
		if err != nil {
			log.Error(err, "Failed to create image puller rolebinding")
			return ctrl.Result{}, err
		}
		log.Info("Image puller permission added", "namespace", customerName)
	}

	// Step 4 — Create Deployment
	deployment := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: "mcsp-app", Namespace: customerName}, deployment)
	if err != nil && errors.IsNotFound(err) {
		replicas := customer.Spec.Replicas
		deployment = &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "mcsp-app",
				Namespace: customerName,
				Labels: map[string]string{
					"app":    "mcsp-app",
					"tenant": customerName,
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "mcsp-app",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app":    "mcsp-app",
							"tenant": customerName,
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "mcsp-app",
								Image: "image-registry.openshift-image-registry.svc:5000/learning-workspace/mcsp-hello-world:latest",
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: 8080,
									},
								},
								Env: []corev1.EnvVar{
									{
										Name:  "PORT",
										Value: "8080",
									},
									{
										Name:  "SCENARIO",
										Value: "2 of 3",
									},
									{
										Name:  "NAMESPACE",
										Value: customerName,
									},
								},
							},
						},
					},
				},
			},
		}
		err = r.Create(ctx, deployment)
		if err != nil {
			log.Error(err, "Failed to create deployment")
			return ctrl.Result{}, err
		}
		log.Info("Deployment created", "deployment", "mcsp-app")
	}

	// Step 5 — Create Service
	service := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: "mcsp-app", Namespace: customerName}, service)
	if err != nil && errors.IsNotFound(err) {
		service = &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "mcsp-app",
				Namespace: customerName,
				Labels: map[string]string{
					"app":    "mcsp-app",
					"tenant": customerName,
				},
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"app": "mcsp-app",
				},
				Ports: []corev1.ServicePort{
					{
						Name:       "http",
						Port:       80,
						TargetPort: intstr.FromInt(8080),
						Protocol:   corev1.ProtocolTCP,
					},
				},
				Type: corev1.ServiceTypeClusterIP,
			},
		}
		err = r.Create(ctx, service)
		if err != nil {
			log.Error(err, "Failed to create service")
			return ctrl.Result{}, err
		}
		log.Info("Service created", "service", "mcsp-app")
	}

	// Step 6 — Create Route using unstructured
	route := &unstructured.Unstructured{}
	route.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "route.openshift.io",
		Version: "v1",
		Kind:    "Route",
	})
	err = r.Get(ctx, types.NamespacedName{Name: "mcsp-app", Namespace: customerName}, route)
	if err != nil && errors.IsNotFound(err) {
		route = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "route.openshift.io/v1",
				"kind":       "Route",
				"metadata": map[string]interface{}{
					"name":      "mcsp-app",
					"namespace": customerName,
					"labels": map[string]interface{}{
						"app":    "mcsp-app",
						"tenant": customerName,
					},
				},
				"spec": map[string]interface{}{
					"to": map[string]interface{}{
						"kind": "Service",
						"name": "mcsp-app",
					},
					"port": map[string]interface{}{
						"targetPort": "http",
					},
					"tls": map[string]interface{}{
						"termination":                   "edge",
						"insecureEdgeTerminationPolicy": "Redirect",
					},
				},
			},
		}
		err = r.Create(ctx, route)
		if err != nil {
			log.Error(err, "Failed to create route")
			return ctrl.Result{}, err
		}
		log.Info("Route created", "route", "mcsp-app")
	}

	// Step 7 — Update Status
	customer.Status.Deployed = true
	customer.Status.Message = fmt.Sprintf("Customer %s successfully deployed", customerName)
	err = r.Status().Update(ctx, customer)
	if err != nil {
		log.Error(err, "Failed to update customer status")
		return ctrl.Result{}, err
	}

	log.Info("Customer reconciled successfully", "customerName", customerName)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CustomerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mcspv1.Customer{}).
		Named("customer").
		Complete(r)
}
