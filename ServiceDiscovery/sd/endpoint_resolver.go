package sd

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	K8S_PSM_LABEL      = "bytemesh.io.psm"
)

var (
	clientset *kubernetes.Clientset
)

type Endpoint struct {
	IP string
	Port string
}

func init() {
	if os.Getenv("bearToken") == "" {
		panic(errors.New("bearToken not provided"))
	}
	bearToken := os.Getenv("bearToken")
	tokenFilePath := "./eptoken"
	epTokenFile, err := os.Create(tokenFilePath)
	if err != nil {
		return
	}
	_, err = epTokenFile.WriteString(bearToken)
	if err != nil {
		return
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		return
	}
	config.BearerToken = bearToken
	config.BearerTokenFile = tokenFilePath
	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return
	}
}

func Lookup(name string) []Endpoint {
	timeout := int64(1000)
	selector := labels.NewSelector()
	nameReq, _ := labels.NewRequirement(K8S_PSM_LABEL, selection.Equals, []string{name})
	selector = selector.Add(*nameReq)
	epList, err := clientset.CoreV1().Endpoints("default").List(context.TODO(), metav1.ListOptions{
		TypeMeta:       metav1.TypeMeta{},
		LabelSelector:  selector.String(),
		TimeoutSeconds: &timeout,
	})
	if err != nil {
		log.Printf("Discovery error : %v\n", err)
		return nil
	}
	// empty endpoints
	if epList == nil || epList.Items == nil || len(epList.Items) == 0 {
		log.Printf("Discovery empty target instances : %v\n", err)
		return nil
	}
	eps := []Endpoint{}
	for _, ep := range epList.Items[0].Subsets {
		for _, addr := range ep.Addresses {
			for _, port := range ep.Ports {
				// use port0 as primary port, ignore others
				eps = append(eps, Endpoint{
					IP: addr.IP,
					Port: strconv.Itoa(int(port.Port)),
				})
			}
		}
	}
	return eps
}