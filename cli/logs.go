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
	"fmt"
	"log"
	"os"

	logging "google.golang.org/api/logging/v2beta1"

	"github.com/spf13/cobra"

	"k8s.io/client-go/1.5/pkg/api"
	versioned "k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/labels"
)

var logCmd = &cobra.Command{
	Use:   "logs [kubernetes service]",
	Short: "Collect ESP logs for a service",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Println("Please specify kubernetes service name")
			os.Exit(-1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if active {
			ExtractFromPods(name)
		}
		if cfg.CredentialsFile != "" || cfg.ProducerProject != "" {
			ExtractFromGCP(name)
		}
	},
}

var (
	active bool
)

func init() {
	RootCmd.AddCommand(logCmd)
	logCmd.PersistentFlags().BoolVarP(&active,
		"active", "a", true,
		"Query kubectl to fetch logs (both stdout and stderr)")
	logCmd.PersistentFlags().StringVarP(&cfg.CredentialsFile,
		"creds", "k", "",
		"Service account credentials JSON file")
	logCmd.PersistentFlags().StringVarP(&cfg.ProducerProject,
		"project", "p", "",
		"Service producer project (optional if you use service account credentials)")
}

// ExtractFromGCP Logging service
func ExtractFromGCP(name string) {
	log.Println("Extracting logs from Google Cloud Logging:")

	hc, err := cfg.GetClient(logging.CloudPlatformScope)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	client, err := logging.New(hc)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	// TODO: qualify by cluster_name using resource.labels.cluster_name
	filter := fmt.Sprintf(`resource.type="container"
		AND resource.labels.container_name="%s"
		AND resource.labels.namespace_id="%s"
		AND severity=ERROR`, EndpointsPrefix+name, namespace)
	log.Println("Filter: ", filter)
	req := &logging.ListLogEntriesRequest{
		OrderBy:    "timestamp asc",
		Filter:     filter,
		ProjectIds: []string{cfg.ProducerProject},
	}

	resp, err := client.Entries.List(req).Fields(
		"entries(resource/labels)",
		"entries(severity,textPayload,timestamp)",
	).Do()
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	log.Println(len(resp.Entries), "entries available")
	for _, entry := range resp.Entries {
		fmt.Printf("[%s,%s,%s]%s",
			entry.Severity, entry.Timestamp, entry.Resource.Labels["pod_id"],
			entry.TextPayload)
	}
}

// ExtractFromPods fetches logs from pods running ESP container
func ExtractFromPods(name string) {
	log.Println("Extracting logs from existing pods:")
	label := labels.SelectorFromSet(labels.Set(map[string]string{
		AnnotationManagedService: name,
	}))
	options := api.ListOptions{LabelSelector: label}
	list, err := clientset.Core().Pods(namespace).List(options)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	var pods = make([]string, len(list.Items))
	for i := 0; i < len(list.Items); i++ {
		pods[i] = list.Items[i].Name
	}

	log.Println(pods)
	for _, pod := range pods {
		PrintLogs(name, pod)
	}
}

// PrintLogs from all pods and containers
func PrintLogs(name, pod string) {
	log.Printf("[pod_id=%s]", pod)
	logOptions := &versioned.PodLogOptions{
		Container: EndpointsPrefix + name,
	}
	raw, err := clientset.Core().Pods(namespace).GetLogs(pod, logOptions).Do().Raw()
	if err != nil {
		log.Println("Request error", err)
	} else {
		fmt.Println("\n" + string(raw))
	}
}
