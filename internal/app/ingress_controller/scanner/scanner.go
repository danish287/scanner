// Scanner pulls information from the kubernetes cluster using the API running locally on the machine.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "log"
	"time"

	"net/http"

	"os/exec"
)

// MyServices contains
type MyServices struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink        string `json:"selfLink"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Items []Service `json:"items"`
}

// Service contains
type Service struct {
	Metadata struct {
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		SelfLink          string    `json:"selfLink"`
		UID               string    `json:"uid"`
		ResourceVersion   string    `json:"resourceVersion"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Labels            struct {
			Component string `json:"component"`
			Provider  string `json:"provider"`
		} `json:"labels"`
	} `json:"metadata"`
	Spec struct {
		Ports []struct {
			Name       string `json:"name"`
			Protocol   string `json:"protocol"`
			Port       int    `json:"port"`
			TargetPort int    `json:"targetPort"`
		} `json:"ports"`
		ClusterIP       string `json:"clusterIP"`
		Type            string `json:"type"`
		SessionAffinity string `json:"sessionAffinity"`
	} `json:"spec"`
	Status struct {
		LoadBalancer struct {
		} `json:"loadBalancer"`
	} `json:"status"`
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

// Ruleset stores
var Ruleset []Rules

func main() {

	// run the kubectl proxy without TLS credentials
	exec.Command("kubectl", "proxy", "--insecure-skip-tls-verify").Start()
	GetServices()
}

// GetServices gets all of the services in our cluster from the API
func GetServices() {

	// api/v1/namespaces/revature.com/services/
	// request information of services from k8s API
	var reqServices MyServices
	var currServices 
	serviceURL := "http://localhost:8001/api/v1/services"
	body := GetResponse(serviceURL)
	
	// unmarshall body of the request and populate structure currServices with information of current services from K8S API
	err := json.Unmarshal(body, &reqServices)
	if err != nil { fmt.Println("FIRST ERROR", err)}
	fmt.Println("\n\nUNMARSHALL", reqServices.Items)
	fmt.Println("L:ENGHT", len(reqServices.Items))

}

// GetDeployments
func GetDeployments() {
	// GET /apis/apps/v1/namespaces/{namespace}/deployments/{name}
	// request information of services from k8s API
	var reqServices MyServices
	var currServices 
	serviceURL := "http://localhost:8001/apis/apps/v1/deployments"
	body := GetResponse(serviceURL)
	
	// unmarshall body of the request and populate structure currServices with information of current services from K8S API
	err := json.Unmarshal(body, &reqServices)
	if err != nil { fmt.Println("FIRST ERROR", err)}
	fmt.Println("\n\nUNMARSHALL", reqServices.Items)
	fmt.Println("L:ENGHT", len(reqServices.Items))

}

// GetIngress 
func GetIngress() {
	// GET /apis/extensions/v1beta1/namespaces/{namespace}/ingresses/{name}


	
}

// GetResponse will request response from Kubernates API
func GetResponse(requestURL string) (respBody []byte) {

	// create a new instance of client & create new request to retrieve info from k8s API
	client := http.Client{}
	apiReq, err := http.NewRequest("GET", requestURL, nil)
	if err != nil { fmt.Println("FIRST ERROR", err)}

	// client do request: send HTTP request & recieve HTTP response
	response, err := client.Do(apiReq)
	if err != nil { fmt.Println("FIRST ERROR", err)}

	// read body of the reponse recieved from k8s API and defer closing body until end
	respBody, err = ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil { fmt.Println("FIRST ERROR", err)}

	return 
}
