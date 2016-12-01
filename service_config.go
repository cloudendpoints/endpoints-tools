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
//
// Service config CLI
//
package main

import (
	"deploy"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var s deploy.Service

func check(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}

func main() {
	log.SetPrefix("[service_config] ")
	log.SetFlags(0)
	var deploy = &cobra.Command{
		Use: "deploy [config+]",
		Long: `Deploy service config files, which are one or combination of
OpenAPI, Google Service Config, or proto descriptor.
(supported extensions: .yaml, .yml, .json, .pb, .descriptor).`,
		Run: func(cmd *cobra.Command, configs []string) {
			check(s.Connect())
			out, err := s.Deploy(configs, "")
			check(err)
			bytes, err := out.MarshalJSON()
			check(err)
			fmt.Print(string(bytes))
		},
	}
	var fetch = &cobra.Command{
		Use: "fetch",
		Run: func(cmd *cobra.Command, args []string) {
			check(s.Connect())
			out, err := s.Fetch()
			check(err)
			bytes, err := out.MarshalJSON()
			check(err)
			fmt.Print(string(bytes))
		},
	}
	var delete = &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			check(s.Connect())
			check(s.Delete())
		},
	}
	var undelete = &cobra.Command{
		Use: "undelete",
		Run: func(cmd *cobra.Command, args []string) {
			check(s.Connect())
			check(s.Undelete())
		},
	}

	fetch.PersistentFlags().StringVarP(&s.Version,
		"version", "v", "",
		"API service config version, empty to use the latest")

	var root = &cobra.Command{}
	root.PersistentFlags().StringVarP(&s.Name,
		"service", "s", "",
		"API service name")
	root.PersistentFlags().StringVarP(&s.CredentialsFile,
		"creds", "k", "",
		"Service account private key JSON file")
	root.PersistentFlags().StringVarP(&s.ProducerProject,
		"project", "p", "",
		"Service producer project")

	root.AddCommand(deploy, fetch, delete, undelete)
	root.Execute()
}
