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

// MyDeployments contains
type MyDeployments struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink        string `json:"selfLink"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Items []Deployment `json:"items"`
}

// Deployment contains
type Deployment struct {
	Metadata struct {
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		SelfLink          string    `json:"selfLink"`
		UID               string    `json:"uid"`
		ResourceVersion   string    `json:"resourceVersion"`
		Generation        int       `json:"generation"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Annotations       struct {
			DeploymentKubernetesIoRevision              string `json:"deployment.kubernetes.io/revision"`
			KubectlKubernetesIoLastAppliedConfiguration string `json:"kubectl.kubernetes.io/last-applied-configuration"`
		} `json:"annotations"`
	} `json:"metadata"`
	Spec struct {
		Replicas int `json:"replicas"`
		Selector struct {
			MatchLabels struct {
				App string `json:"app"`
			} `json:"matchLabels"`
		} `json:"selector"`
		Template struct {
			Metadata struct {
				CreationTimestamp interface{} `json:"creationTimestamp"`
				Labels            struct {
					App string `json:"app"`
				} `json:"labels"`
			} `json:"metadata"`
			Spec struct {
				Containers []struct {
					Name  string `json:"name"`
					Image string `json:"image"`
					Ports []struct {
						ContainerPort int    `json:"containerPort"`
						Protocol      string `json:"protocol"`
					} `json:"ports"`
					Resources struct {
					} `json:"resources"`
					TerminationMessagePath   string `json:"terminationMessagePath"`
					TerminationMessagePolicy string `json:"terminationMessagePolicy"`
					ImagePullPolicy          string `json:"imagePullPolicy"`
				} `json:"containers"`
				RestartPolicy                 string `json:"restartPolicy"`
				TerminationGracePeriodSeconds int    `json:"terminationGracePeriodSeconds"`
				DNSPolicy                     string `json:"dnsPolicy"`
				SecurityContext               struct {
				} `json:"securityContext"`
				SchedulerName string `json:"schedulerName"`
			} `json:"spec"`
		} `json:"template"`
		Strategy struct {
			Type          string `json:"type"`
			RollingUpdate struct {
				MaxUnavailable string `json:"maxUnavailable"`
				MaxSurge       string `json:"maxSurge"`
			} `json:"rollingUpdate"`
		} `json:"strategy"`
		RevisionHistoryLimit    int `json:"revisionHistoryLimit"`
		ProgressDeadlineSeconds int `json:"progressDeadlineSeconds"`
	} `json:"spec"`
	Status struct {
		ObservedGeneration  int `json:"observedGeneration"`
		Replicas            int `json:"replicas"`
		UpdatedReplicas     int `json:"updatedReplicas"`
		UnavailableReplicas int `json:"unavailableReplicas"`
		Conditions          []struct {
			Type               string    `json:"type"`
			Status             string    `json:"status"`
			LastUpdateTime     time.Time `json:"lastUpdateTime"`
			LastTransitionTime time.Time `json:"lastTransitionTime"`
			Reason             string    `json:"reason"`
			Message            string    `json:"message"`
		} `json:"conditions"`
	} `json:"status"`
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

// ReqDeployments stores
var ReqDeployments MyDeployments

// ReqServices contians
var ReqServices MyServices

func main() {
	// run the kubectl proxy without TLS credentials
	exec.Command("kubectl", "proxy", "--insecure-skip-tls-verify").Start()
	GetServices()
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

//FindService will
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

// GetDeployments contains
func GetDeployments() {
	// GET /apis/apps/v1/namespaces/{namespace}/deployments/{name}
	// request information of services from k8s API
	serviceURL := "http://localhost:8001/apis/apps/v1/deployments"
	body := GetResponse(serviceURL)
	
	// unmarshall body of the request and populate structure currServices with information of current services from K8S API
	err := json.Unmarshal(body, &ReqDeployments)
	if err != nil { fmt.Println("FIRST ERROR", err)}

	fmt.Println("\n\nUNMARSHALL", ReqDeployments.Items)
	fmt.Println("LENGHT", len(ReqDeployments.Items))

}

//FindDeployment will
func FindDeployment() {
	deploymentLst := ReqDeployments.Items
	for i:=0; i < len(serviceLst); i++{
		currService := serviceLst[i]
		// if currService.Name == name we are looking for 
		// Route.ServiceName = currService..Metadata.Name
		// Route.ServicePort = currService.Spec.Ports.TargetPort
		// Route.ServiceIP = currService.Spec.ClusterIP
		fmt.Println("CURRENTSERVICE", currService.Metadata.Name)
	}
}

// GetIngress contains
func GetIngress() {
	// GET /apis/extensions/v1beta1/namespaces/{namespace}/ingresses/{name}
	
}


func Portal() {
	///apis/revature.com/v1/namespaces/default/portals/
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
