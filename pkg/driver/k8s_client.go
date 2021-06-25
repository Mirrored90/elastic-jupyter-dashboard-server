package driver

import (
	kubeflowtkestackiov1alpha1 "github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var MyK8S client.Client
var err error

func init() {
	MyK8S, err = client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		panic("server run error: " + err.Error())
	}

	err = kubeflowtkestackiov1alpha1.AddToScheme(scheme.Scheme)
	if err != nil {
		panic("server run error: " + err.Error())
	}
}
