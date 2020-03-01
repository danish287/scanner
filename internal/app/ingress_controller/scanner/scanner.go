// Scanner pulls information from the kubernetes cluster using the API running locally on the machine.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net/http"

	"os/exec"
)

// MyServices will
type MyServices struct {
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		Metadata struct {
			CreationTimestamp string `json:"creationTimestamp"`
			Labels            struct {
				Component string `json:"component"`
				Provider  string `json:"provider"`
			} `json:"labels"`
			Name            string `json:"name"`
			Namespace       string `json:"namespace"`
			ResourceVersion string `json:"resourceVersion"`
			SelfLink        string `json:"selfLink"`
			UID             string `json:"uid"`
		} `json:"metadata"`
		Spec struct {
			ClusterIP string `json:"clusterIP"`
			Ports     []struct {
				Name       string `json:"name"`
				Port       int    `json:"port"`
				Protocol   string `json:"protocol"`
				TargetPort int    `json:"targetPort"`
			} `json:"ports"`
			SessionAffinity string `json:"sessionAffinity"`
			Type            string `json:"type"`
		} `json:"spec"`
		Status struct {
			LoadBalancer struct{} `json:"loadBalancer"`
		} `json:"status"`
	} `json:"items"`
	Kind     string `json:"kind"`
	Metadata struct {
		ResourceVersion string `json:"resourceVersion"`
		SelfLink        string `json:"selfLink"`
	} `json:"metadata"`
}

// Route stores
type Route struct {
	ServiceName string `json:"ServiceName"`
	ServicePort string `json:"ServicePort"`
	ServiceIP   string `json:"ServiceIP"`
}

// Rules stores
type Rules struct {
	Protocol string `json:"Protocol"`
	Path     string `json:"Path"`
	Route    Route  `json:"Route"`
}

// RuleSet stores
var RuleSet []Rules

func main() {

	// run the kubectl proxy with TLS credentials
	exec.Command("kubectl", "proxy", "--insecure-skip-tls-verify").Start()
	GetServices()
}

// GetServices gets all of the services in our cluster from the API
func GetServices() {
	var currServices MyServices
	// sends request to Kubernates API to retreive services
	// apiReq, err := http.NewRequest("GET", "localhost:8001/api/v1/services", nil)

	response, err := http.Get("http://localhost:8001/api/v1/services")

	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	err = json.Unmarshal(body, &currServices)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(body)
	fmt.Println(currServices)
}
