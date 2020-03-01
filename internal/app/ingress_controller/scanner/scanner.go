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
	Items []ServiceList `json:"items"`
}

// ServiceList contains
type ServiceList struct {
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
	// currServices stores the unmarshalled JSON information of the k8s services 
	var currServices MyServices
	serviceURL := "http://localhost:8001/api/v1/services"
	// create a new instance of client
	client := http.Client{}
	// create new request to retrieve information from k8s API
	apiReq, err := http.NewRequest("GET", "http://localhost:8001/api/v1/services", nil)
	if err != nil {
		fmt.Println("FIRST ERROR", err)
	}

	// client does request, send HTTP request to recieve HTTP response
	response, errors := client.Do(apiReq)
	fmt.Println("RESPONSE STATUS", response.Status)
	if errors != nil {
		fmt.Println("SECOND ERROR", err)
	}

	// read body of the reponse recieved from oure request
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	
	// unmarshall body of the request and populate structure currServices with information of out current services
	err = json.Unmarshal(body, &currServices)
	fmt.Println("\n\nUNMARSHALL", currServices.Items)
	fmt.Println("L:ENGHT", len(currServices.Items))

	if err != nil {
		fmt.Println(err)
	}

}

func GetResponse(url) []byte {

}
