// Scanner pulls information from the kubernetes cluster using the API running locally on the machine.
package main

import (
	"fmt"

	"net/http"

	"os/exec"
)

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
	// sends request to Kubernates API to retreive services
	apiReq, err := http.NewRequest("GET", "localhost:8001api/v1/services", nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(apiReq)
}
