package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	// Parse the user's kubeconfig file
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	fmt.Printf("Kubeconfig path: %s\n", *kubeconfig)

	// Use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// Create the Kubernetes client to take action on pods and replicasets
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Create the Dynamic client to take action on our spaceship CRD
	// dyclient := dynamic.ConfigFor(config)
	// dynamic.ResourceInterface.Get(context.TODO())

	// Controller Loop
	for {
		fmt.Println()

		// LIST ALL PODS
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		for _, value := range pods.Items {
			fmt.Printf("%s, ", value.Name)
		}
		fmt.Println()

		// LIST ALL REPLICASETS
		rs, err := clientset.AppsV1().ReplicaSets("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d ReplicaSets in the cluster\n", len(rs.Items))
		for _, value := range rs.Items {
			fmt.Printf("Name=%s, ", value.Name)

			scale, err := clientset.AppsV1().ReplicaSets(value.Namespace).GetScale(context.TODO(), value.Name, metav1.GetOptions{})
			if err != nil {
				panic(err.Error())
			}
			// var spec
			// err := json.NewDecoder(scale.Spec).Decode(&spec)
			// fmt.Printf("Scale=%s\n", spec)
		}

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		// namespace := "default"
		// pod := "example-xxxxx"
		// _, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
		// if errors.IsNotFound(err) {
		// 	fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
		// } else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		// 	fmt.Printf("Error getting pod %s in namespace %s: %v\n",
		// 		pod, namespace, statusError.ErrStatus.Message)
		// } else if err != nil {
		// 	panic(err.Error())
		// } else {
		// 	fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		// }

		time.Sleep(10 * time.Second)
	}
}
