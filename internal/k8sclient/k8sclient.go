/*
Copyright 2016 The Kubernetes Authors.
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

// Note: the example only works with the code within the same release/branch.
package k8sclient

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"webhook-proxy/internal/helper"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//
	// Uncomment to load all auth plugins
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func lastString(ss []string) string {
	return ss[len(ss)-1]
}

func Inconfig() (*kubernetes.Clientset, context.Context) {
	var kubeconfig *string
	var config *rest.Config
	var err error
	_, incluster := os.LookupEnv("INCLUSTER")
	if home := homedir.HomeDir(); home != "" && !incluster {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		flag.Parse()
		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	} else {
		// kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ctx := context.Background()

	return clientset, ctx
}

func GetTag(namespace string, deploymentName string, clientset *kubernetes.Clientset, ctx context.Context) (tag string, err error) {
	delimiterCharacter := helper.GetEnv("TAG_DELIMITER", ":")
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, deploymentName, v1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Deployment %s in namespace %s not found\n", deploymentName, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting deployment %s in namespace %s: %v\n",
			deploymentName, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		// panic(err.Error())
		fmt.Println(err)
	} else {
		fmt.Printf("Found deployment %s in namespace %s\n", deploymentName, namespace)
		imageName := deployment.Spec.Template.Spec.Containers[0].Image
		tag = lastString(strings.Split(imageName, delimiterCharacter))
		fmt.Printf("Tag: %s\n", tag)
	}

	return tag, err

}
