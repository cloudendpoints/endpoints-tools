// Copyright 2016 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
///////////////////////////////////////////////////////////////////////////
package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"k8s.io/client-go/1.5/pkg/api"
	versioned "k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/labels"
)

var endpointsCmd = &cobra.Command{
	Use:   "endpoints [Kubernetes service]",
	Short: "Describe ESP endpoints for a Kubernetes service",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Println("Please specify Kubernetes service name")
			os.Exit(-1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		out, err := GetESPEndpoints(name)
		if err != nil {
			log.Println(err)
			os.Exit(-1)
		}
		bytes, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(bytes))
	},
}

func init() {
	RootCmd.AddCommand(endpointsCmd)
}

// GetESPEndpoints collects endpoints information
func GetESPEndpoints(name string) (map[string]map[string]string, error) {
	// collect all services running ESP
	list, err := GetESPServices(name)
	if err != nil {
		return nil, err
	}

	out := map[string]map[string]string{}

	for _, svc := range list {
		ends, err := GetEndpoints(svc)
		if err != nil {
			log.Println(err)
			os.Exit(-1)
		}
		out[svc.Name] = ends
	}

	return out, nil
}

// GetESPServices for a Kubernetes service
func GetESPServices(name string) ([]*versioned.Service, error) {
	label := labels.SelectorFromSet(labels.Set(map[string]string{
		AnnotationManagedService: name,
	}))
	options := api.ListOptions{LabelSelector: label}
	list, err := clientset.Core().Services(namespace).List(options)
	if err != nil {
		return nil, err
	}

	var out = make([]*versioned.Service, len(list.Items))
	for i := 0; i < len(list.Items); i++ {
		out[i] = &list.Items[i]
	}
	return out, nil
}

// GetEndpoints retrieves endpoints information for an ESP service
func GetEndpoints(svc *versioned.Service) (map[string]string, error) {
	out := map[string]string{}

	out[AnnotationConfigId] = svc.Annotations[AnnotationConfigId]
	out[AnnotationConfigName] = svc.Annotations[AnnotationConfigName]
	out[AnnotationDeploymentType] = svc.Annotations[AnnotationDeploymentType]

	if svc.Spec.Type == versioned.ServiceTypeNodePort {
		for _, port := range svc.Spec.Ports {
			out[port.Name] = "NODE_IP:" + strconv.Itoa(int(port.NodePort))
		}
	} else if svc.Spec.Type == versioned.ServiceTypeClusterIP {
		for _, port := range svc.Spec.Ports {
			out[port.Name] = svc.Name + ":" + strconv.Itoa(int(port.Port))
		}
	} else if svc.Spec.Type == versioned.ServiceTypeLoadBalancer {
		var address string
		var err error
		ok := repeat(func() bool {
			log.Println("Retrieving address of the service " + svc.Name)
			svc, err = clientset.Core().Services(namespace).Get(svc.Name)
			if err != nil {
				return false
			}

			ingress := svc.Status.LoadBalancer.Ingress
			if len(ingress) == 0 {
				return false
			}

			if ingress[0].IP != "" {
				address = ingress[0].IP
				return true
			} else if ingress[0].Hostname != "" {
				address = ingress[0].Hostname
				return true
			}

			return false
		})

		if !ok {
			return nil, errors.New("Failed to retrieve IP of the service")
		}

		for _, port := range svc.Spec.Ports {
			out[port.Name] = address + ":" + strconv.Itoa(int(port.Port))
		}
	} else {
		return nil, errors.New("Cannot handle service type")
	}

	return out, nil
}

const MaxTries = 10

// Repeat until success (function returns true) up to MaxTries
func repeat(f func() bool) bool {
	try := 0
	delay := 2 * time.Second
	result := false
	for !result && try < MaxTries {
		if try > 0 {
			log.Println("Waiting for next attempt: ", delay)
			time.Sleep(delay)
			delay = 2 * delay
			log.Println("Repeat attempt #", try+1)
		}
		result = f()
		try = try + 1
	}

	if !result {
		log.Println("Failed all attempts")
	}

	return result
}
