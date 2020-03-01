// Scanner pulls information from the kubernetes cluster using the API running locally on the machine.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"net/http"

	"os/exec"
)

// MyServices contains data for all services in the network
type MyServices struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink        string `json:"selfLink"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Items []Service `json:"items"`
}


// Service struct contains data pertaining to a service 
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

// Portal retrieves data of portal information  
type Portal struct {
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Metadata   struct {
			Annotations struct {
				KubectlKubernetesIoLastAppliedConfiguration string `json:"kubectl.kubernetes.io/last-applied-configuration"`
			} `json:"annotations"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
			Generation        int       `json:"generation"`
			Name              string    `json:"name"`
			Namespace         string    `json:"namespace"`
			ResourceVersion   string    `json:"resourceVersion"`
			SelfLink          string    `json:"selfLink"`
			UID               string    `json:"uid"`
		} `json:"metadata"`
		Spec struct {
			Portal   string `json:"portal"`
			Targetip string `json:"targetip"`
		} `json:"spec"`
	} `json:"items"`
	Kind     string `json:"kind"`
	Metadata struct {
		Continue        string `json:"continue"`
		ResourceVersion string `json:"resourceVersion"`
		SelfLink        string `json:"selfLink"`
	} `json:"metadata"`
}

// IngressData stores
type IngressData struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink        string `json:"selfLink"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Items []struct {
		Metadata struct {
			Name              string    `json:"name"`
			Namespace         string    `json:"namespace"`
			SelfLink          string    `json:"selfLink"`
			UID               string    `json:"uid"`
			ResourceVersion   string    `json:"resourceVersion"`
			Generation        int       `json:"generation"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
			Annotations       struct {
				NginxIngressKubernetesIoRewriteTarget string `json:"nginx.ingress.kubernetes.io/rewrite-target"`
			} `json:"annotations"`
		} `json:"metadata"`
		Spec struct {
			Rules []struct {
				HTTP struct {
					Paths []struct {
						Path    string `json:"path"`
						Backend struct {
							ServiceName string `json:"serviceName"`
							ServicePort int    `json:"servicePort"`
						} `json:"backend"`
					} `json:"paths"`
				} `json:"http"`
			} `json:"rules"`
		} `json:"spec"`
		Status struct {
			LoadBalancer struct {
			} `json:"loadBalancer"`
		} `json:"status"`
	} `json:"items"`
}

// Route stores desired route of _________
type Route struct {
	ServiceName string `json:"ServiceName"`
	ServicePort string `json:"ServicePort"`
	ServiceIP   string `json:"ServiceIP"`
}

// Rules stores ___________
type Rules struct {
	Protocol string `json:"Protocol"`
	Path     string `json:"Path"`
	Route    Route  `json:"Route"`
}

// Ruleset stores
var Ruleset []Rules

// ReqServices contians
var ReqServices MyServices

// TargetIP will store alternative IP address to dial if first one is not found
var TargetIP []string

func main() {
	// run the kubectl proxy without TLS credentials
	exec.Command("kubectl", "proxy", "--insecure-skip-tls-verify").Start()
	// GetServices()
	// GetTargetIP()
	GetIngress()
}

// GetServices gets all of the services in our cluster from the API
func GetServices() {

	// request information of services from k8s API
	serviceURL := "http://localhost:8001/api/v1/services"
	body := GetResponse(serviceURL)
	
	// unmarshall body of the request and populate structure currServices with information of current services from K8S API
	err := json.Unmarshal(body, &ReqServices)
	if err != nil { fmt.Println("FIRST ERROR", err)}
	FindService()

	fmt.Println("\n\nUNMARSHALL", ReqServices.Items)
	fmt.Println("LENGHT", len(ReqServices.Items))
}

//FindService searches list of services by 'name' to match 
func FindService() {

	serviceLst := ReqServices.Items
	for i:=0; i < len(serviceLst); i++{
		currService := serviceLst[i]
		// if currService.Name == name we are looking for 
		// Route.ServiceName = currService..Metadata.Name
		// Route.ServicePort = currService.Spec.Ports.Port
		// Route.ServiceIP = currService.Spec.ClusterIP
		fmt.Println("CURRENTSERVICE", currService.Metadata.Name)
	}
}


// GetIngress contains
func GetIngress() {
	
	// GET /apis/extensions/v1beta1/namespaces/{namespace}/ingresses/{name}
	// /apis/extensions/v1beta1/ingresses
	//DreamTeamIngress
	// items.spec, items.rules, items.http, items.path, items.sepc.ruleshttp.paths.backend.serviceport == serviceport, items.sepc.ruleshttp.paths.backend.servicename = servicename
	//serviceip == cluster ip
	var TargetData IngressData
	serviceURL := "http://localhost:8001/apis/extensions/v1beta1/ingresses"
	body := GetResponse(serviceURL)

	err := json.Unmarshal(body, &TargetData)
	if err != nil {fmt.Println(err)}
	fmt.Println("NAME", TargetData.Items[0].Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName)



}

// GetTargetIP will retrieve targetIP from the portal to provide an alternative IP address for proxy
func GetTargetIP() {
	// request information of services from k8s API
	var PortalData Portal
	serviceURL := "http://localhost:8001/apis/revature.com/v1/namespaces/default/portals/"
	body := GetResponse(serviceURL)
	
	// unmarshall body of the request and populate structure currServices with information of current services from K8S API
	err := json.Unmarshal(body, &PortalData)
	if err != nil { fmt.Println(err)}
	tIP := []string{PortalData.Items[0].Spec.Targetip}
	TargetIP = append(tIP)
	fmt.Printf("%T", TargetIP)
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
