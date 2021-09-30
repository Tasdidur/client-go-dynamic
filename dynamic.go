package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"log"
	"os"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/dynamic"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main(){
	kubeconfig := flag.String("kubeconfig", "/home/office/.kube/config", "path")

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil{
		log.Fatal(err)
	}


	//CREATE SERVICE ================================================== BEGIN

	myservice := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind": "Service",
			"metadata" : map[string]interface{}{
				"name": "hel-svc-x",
			},
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"finder": "hel-x",
				},
				"ports": []map[string]interface{}{
					{
						"protocol": "TCP",
						"port" : 8003,
						"targetPort": 8081,
					},
				},
			},
		},
	}

	serviceRes := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	fmt.Println("creating ingress...\n")
	//
	resultS, err := client.Resource(serviceRes).Namespace("default").Create(context.TODO(),myservice,metav1.CreateOptions{})
	if err != nil{
		log.Println(err)
	}else{
		fmt.Printf("ingress created %q. \n", resultS.GetName())
	}


	//CREATE SERVICE ================================================== END

	//CREATE DEPLOYMENT ================================================== BEGIN

	prompt()
	mydeployment := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind": "Deployment",
			"metadata" : map[string]interface{}{
				"name": "hel-dep-x",
			},
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"finder": "hel-x",
					},
				},
				"replicas" : 2,
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"name": "hel-pod-x",
						"labels": map[string]interface{}{
							"finder": "hel-x",
						},
					},
					"spec": map[string]interface{}{
						"containers": []map[string]interface{}{
							{
								"name": "hel-container-x",
								"image": "tasdidur/test-ingress",
							},
						},
					},
				},
			},
		},
	}

	deploymentRes := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	fmt.Println("creating deployment...\n")
	resultD , err := client.Resource(deploymentRes).Namespace("default").Create(context.Background(),mydeployment,metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
	}else{
		fmt.Printf("deployment created %q. \n", resultD.GetName())
	}

	//CREATE DEPLOYMENT ================================================== END

	//CREATE INGRESS ================================================== BEGIN

	myingress := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "networking.k8s.io/v1beta1",
			"kind": "Ingress",
			"metadata":map[string]interface{}{
				"name": "test-ingress",
				"namespace": "default",
				"annotations": map[string]interface{}{
					"nginx.ingress.kubernetes.io/use-regex": "true",
				},
			},
			"spec": map[string]interface{}{
				"rules": []map[string]interface{}{
					{
						"host": "tasdid2.com",
						"http": map[string]interface{}{
							"paths": []map[string]interface{}{
								{
									"path": "/hi",
									"backend": map[string]interface{}{
										"serviceName": "hel-svc",
										"servicePort": 8001,
										},
								},
								{
									"path": "/hello",
									"backend": map[string]interface{}{
										"serviceName": "hel-svc",
										"servicePort": 8001,
										},
								},
								{
									"path": "/bye",
									"backend": map[string]interface{}{
										"serviceName": "hel-svc",
										"servicePort": 8001,
									},
								},
							},
						},
					},
				},
			},
		},
	}


	ingressRes := schema.GroupVersionResource{
		Group:    "networking.k8s.io",
		Version:  "v1beta1",
		Resource: "ingresses",
	}

	fmt.Println("creating ingress...\n")

	resultI , err := client.Resource(ingressRes).Namespace("default").Create(context.Background(),myingress,metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
	}else {

		fmt.Printf("ingress created %q. \n", resultI.GetName())
	}

	//CREATE INGRESS ================================================== END
}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}

func int32Ptr(i int32) *int32 { return &i }