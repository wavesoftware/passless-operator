package test

import (
	"github.com/mitchellh/go-homedir"
	"github.com/wavesoftware/go-ensure"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// BuildDefaultClientConfig builds default client config for testing.
func BuildDefaultClientConfig() (*rest.Config, error) {
	kubeConfigPath, err := homedir.Expand("~/.kube/config")
	ensure.NoError(err)
	return BuildClientConfig(kubeConfigPath, "")
}

// BuildClientConfig builds client config for testing.
func BuildClientConfig(kubeConfigPath string, clusterName string) (*rest.Config, error) {
	overrides := clientcmd.ConfigOverrides{}
	// Override the cluster name if provided.
	if clusterName != "" {
		overrides.Context.Cluster = clusterName
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath},
		&overrides).ClientConfig()
}
