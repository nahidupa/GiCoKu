package main

import (
	"context"
	"flag"
	"fmt"

	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	//Print all nodes that are not ready state
	printNotReadyNodes(clientset)

	//Print last five events
	printLastFewEvents(clientset)

	//Print pods are not ready state
	printPodsAreNotReady(clientset)

	//Print pods that are restarted more than once
	printRestatedMorethanOncePods(clientset)

}

// A function that print all kubernetes nodes that are not ready state
func printNotReadyNodes(clientset *kubernetes.Clientset) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, node := range nodes.Items {
		notReadyCount := 0
		for _, condition := range node.Status.Conditions {
			if condition.Type == "Ready" {
				if condition.Status != "True" {
					notReadyCount++
				}
			}
		}
		if notReadyCount > 0 {
			fmt.Printf("Node %s is not ready: %d conditions are not ready\n", node.Name, notReadyCount)
		}
	}
}

// A function that print last five kubernetes events
func printLastFewEvents(clientset *kubernetes.Clientset) {
	events, err := clientset.CoreV1().Events("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Last five events:")
	for _, e := range events.Items[len(events.Items)-5:] {
		fmt.Printf("%s %s\n", e.FirstTimestamp, e.Message)
	}
}

// A function that print pods that are restarted more than once
func printRestatedMorethanOncePods(clientset *kubernetes.Clientset) {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, pod := range pods.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.RestartCount > 0 {
				fmt.Printf("Pod %s was restarted %d times\n", pod.Name, containerStatus.RestartCount)
			}
		}

	}
}

//A function name "printPodsAreNotReady" that print pods are not ready state
func printPodsAreNotReady(clientset *kubernetes.Clientset) {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, pod := range pods.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Ready != true {
				fmt.Printf("Pod %s is not ready and containerStatus.RestartCount is %d\n", pod.Name, containerStatus.RestartCount)

			}
		}
	}
}

func printEvents(clientset *kubernetes.Clientset) {
	events, err := clientset.CoreV1().Events("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, event := range events.Items {
		fmt.Printf("Event: %s\n", event.Message)
	}
}


//A function name "printPodsAreNotReady" that print pods are not ready state
func printPodsAreNotReady(clientset *kubernetes.Clientset) {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		
		panic(err)
	}
	for _, pod := range pods.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Ready != true {
				fmt.Printf("Pod %s is not ready and containerStatus.RestartCount is %d\n", pod.Name, containerStatus.RestartCount)

			}
		}
	}
}