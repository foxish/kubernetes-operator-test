package main

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/foxish/kubernetes-operator-test/testutil"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/home/ramanathana/.kube/config")
	if err != nil {
		fmt.Println(err, "build config from flags failed")
	}

	cli, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err, "creating new kube-client failed")
	}

	namespace, err := testutil.CreateNamespace(cli, "fox")
	if err != nil {
		fmt.Println(nil, err, namespace)
	}

	err = testutil.DeleteNamespace(cli, "fox")
	if err != nil {
		fmt.Println(err)
	}
}
