package certs

import "fmt"

type Purpose string

const (
	// KubeconfigDataName is the key used to store a Kubeconfig in the secret's data field.
	KubeconfigDataName = "value"

	// TLSKeyDataName is the key used to store a TLS private key in the secret's data field.
	TLSKeyDataName = "tls.key"

	// TLSCrtDataName is the key used to store a TLS certificate in the secret's data field.
	TLSCrtDataName = "tls.crt"

	// Kubeconfig is the secret name suffix storing the Cluster Kubeconfig.
	Kubeconfig = Purpose("kubeconfig")

	// ClusterCA is the secret name suffix for APIServer CA.
	ClusterCA = Purpose("ca")

	// ClientClusterCA is the secret name suffix for APIServer CA.
	ClientClusterCA = Purpose("cca")

	// EtcdCA is the secret name suffix for the Etcd CA
	EtcdCA Purpose = "etcd"

	// ServiceAccount is the secret name suffix for the Service Account keys
	ServiceAccount Purpose = "sa"

	// FrontProxyCA is the secret name suffix for Front Proxy CA
	FrontProxyCA Purpose = "proxy"

	// APIServerEtcdClient is the secret name of user-supplied secret containing the apiserver-etcd-client key/cert
	APIServerEtcdClient Purpose = "apiserver-etcd-client"
)

var (
	// allSecretPurposes defines a lists with all the secret suffix used by Cluster API
	allSecretPurposes = []Purpose{Kubeconfig, ClusterCA, EtcdCA, ServiceAccount, FrontProxyCA, APIServerEtcdClient}
)

// Name returns the name of the secret for a cluster.
func Name(cluster string, suffix Purpose) string {
	return fmt.Sprintf("%s-%s", cluster, suffix)
}
