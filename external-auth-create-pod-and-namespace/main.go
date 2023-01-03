package main

import (
	"context"
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"),
			"Absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "Absolute path to kubeconfig file")
	}
	flag.Parse()
	cs, err := newClientSet(kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	namespaceName := "client-go-experiments"
	podName := "client-go-experiment-nginx-pod"
	// Check for namespace, create if it doesn't exist
	var namespace *v1.Namespace
	if ns, err := getNamespace(cs, namespaceName); err != nil {
		namespace, err = createNamespace(cs, namespaceName)
		fmt.Printf("The \"%s\" namespace has been created\n", namespace.Name)
		if err != nil {
			panic(err.Error())
		}
	} else {
		namespace = ns
		fmt.Printf("The \"%s\" namespace already exists\n", namespace.Name)
	}

	// Check for pod in "client-go-experiments" namespace, create it if it doesn't exist
	var pod *v1.Pod
	if p, err := getPod(cs, podName, namespace.Name); err != nil {
		pod, err = createPod(cs, podName, namespace.Name)
		fmt.Printf("The \"%s\" pod has been created in the \"%s\" namespace\n", pod.Name, pod.Namespace)
		if err != nil {
			panic(err.Error())
		}
	} else {
		pod = p
		fmt.Printf("The \"%s\" pod already exists in the \"%s\" namespace\n", pod.Name, pod.Namespace)
	}
}

func newClientSet(kubeconfig *string) (clientset *kubernetes.Clientset, err error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(cfg)
}

func getNamespace(clientset *kubernetes.Clientset, name string) (namespace *v1.Namespace, err error) {
	return clientset.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
}

func createNamespace(clientset *kubernetes.Clientset, name string) (ns *v1.Namespace, err error) {
	newNS := v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return clientset.CoreV1().Namespaces().Create(context.TODO(), &newNS, metav1.CreateOptions{})
}

func getPod(clientset *kubernetes.Clientset, name string, namespace string) (pod *v1.Pod, err error) {
	return clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func createPod(clientset *kubernetes.Clientset, name string, namespace string) (pod *v1.Pod, err error) {
	containers := []v1.Container{
		{
			Name:  "nginx",
			Image: "nginx:latest",
		},
	}

	newPod := v1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1.PodSpec{
			Containers: containers,
		},
		Status: v1.PodStatus{},
	}
	return clientset.CoreV1().Pods(namespace).Create(context.TODO(), &newPod, metav1.CreateOptions{})
}
