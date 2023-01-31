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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ManagedControlPlaneSpec defines the desired state of ManagedControlPlane
type ManagedControlPlaneSpec struct {
	// Replicas is the number of control plane nodes
	Replicas *int `json:"replicas"`

	// Service is the service configuration for the control plane pods
	// +optional
	Service Service `json:"service,omitempty"`

	// K3SConfig is the configuration for the K3S server
	// +optional
	K3SConfig K3SConfig `json:"k3sConfig,omitempty"`
}

type Service struct {
	Port int `json:"port"`
}

type K3SConfig struct {
	// AgentConfig specifies configuration for the agent nodes
	// +optional
	AgentConfig K3SAgentConfig `json:"agentConfig,omitempty"`

	// ServerConfig specifies configuration for the agent nodes
	// +optional
	ServerConfig K3SServerConfig `json:"serverConfig,omitempty"`

	// Version specifies the k3s version
	// +optional
	Version string `json:"version,omitempty"`
}

type K3SServerConfig struct {
	// KubeAPIServerArgs is a customized flag for kube-apiserver process
	// +optional
	KubeAPIServerArgs []string `json:"kubeAPIServerArg,omitempty"`

	// KubeControllerManagerArgs is a customized flag for kube-controller-manager process
	// +optional
	KubeControllerManagerArgs []string `json:"kubeControllerManagerArgs,omitempty"`

	// TLSSan Add additional hostname or IP as a Subject Alternative Name in the TLS cert
	// +optional
	TLSSan []string `json:"tlsSan,omitempty"`

	// BindAddress k3s bind address (default: 0.0.0.0)
	// +optional
	BindAddress string `json:"bindAddress,omitempty"`

	// HttpsListenPort HTTPS listen port (default: 6443)
	// +optional
	HttpsListenPort string `json:"httpsListenPort,omitempty"`

	// AdvertiseAddress IP address that apiserver uses to advertise to members of the cluster (default: node-external-ip/node-ip)
	// +optional
	AdvertiseAddress string `json:"advertiseAddress,omitempty"`

	// AdvertisePort Port that apiserver uses to advertise to members of the cluster (default: listen-port) (default: 0)
	// +optional
	AdvertisePort string `json:"advertisePort,omitempty"`

	// ClusterCidr  Network CIDR to use for pod IPs (default: "10.42.0.0/16")
	// +optional
	ClusterCidr string `json:"clusterCidr,omitempty"`

	// ServiceCidr Network CIDR to use for services IPs (default: "10.43.0.0/16")
	// +optional
	ServiceCidr string `json:"serviceCidr,omitempty"`

	// ClusterDNS  Cluster IP for coredns service. Should be in your service-cidr range (default: 10.43.0.10)
	// +optional
	ClusterDNS string `json:"clusterDNS,omitempty"`

	// ClusterDomain Cluster Domain (default: "cluster.local")
	// +optional
	ClusterDomain string `json:"clusterDomain,omitempty"`

	// DisableComponents  specifies extra commands to run before k3s setup runs
	// +optional
	DisableComponents []string `json:"disableComponents,omitempty"`
}

type K3SAgentConfig struct {
	// NodeLabels  Registering and starting kubelet with set of labels
	// +optional
	NodeLabels []string `json:"nodeLabels,omitempty"`

	// NodeTaints Registering kubelet with set of taints
	// +optional
	NodeTaints []string `json:"nodeTaints,omitempty"`

	// TODO: take in a object or secret and write to file. this is not useful
	// PrivateRegistry  registry configuration file (default: "/etc/rancher/k3s/registries.yaml")
	// +optional
	PrivateRegistry string `json:"privateRegistry,omitempty"`

	// KubeletArgs Customized flag for kubelet process
	// +optional
	KubeletArgs []string `json:"kubeletArgs,omitempty"`

	// KubeProxyArgs Customized flag for kube-proxy process
	// +optional
	KubeProxyArgs []string `json:"kubeProxyArgs,omitempty"`

	// NodeName Name of the Node
	// +optional
	NodeName string `json:"nodeName,omitempty"`
}

// ManagedControlPlaneStatus defines the observed state of ManagedControlPlane
type ManagedControlPlaneStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ManagedControlPlane is the Schema for the managedcontrolplanes API
type ManagedControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagedControlPlaneSpec   `json:"spec,omitempty"`
	Status ManagedControlPlaneStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ManagedControlPlaneList contains a list of ManagedControlPlane
type ManagedControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagedControlPlane `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ManagedControlPlane{}, &ManagedControlPlaneList{})
}
