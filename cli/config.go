package cli

import (
	"deploy"
	"fmt"

	"github.com/spf13/cobra"
)

var s deploy.Service

func init() {
	var submit = &cobra.Command{
		Use:   "submit [config+]",
		Short: "Submit and rollout service config files",
		Long: `Submit and rollout service config files, which are one or combination of
OpenAPI, Google Service Config, or proto descriptor.
(supported extensions: .yaml, .yml, .json, .pb, .descriptor).`,
		Run: func(cmd *cobra.Command, configs []string) {
			check("Cannot connect to the API:", s.Connect())
			out, err := s.Deploy(configs, "")
			check("Failed to deploy service config:", err)
			bytes, err := out.MarshalJSON()
			check("Failed to marshal response:", err)
			fmt.Print(string(bytes))
		},
	}

	var fetch = &cobra.Command{
		Use:   "fetch",
		Short: "Fetch service config",
		Run: func(cmd *cobra.Command, args []string) {
			check("Cannot connect to the API:", s.Connect())
			out, err := s.Fetch()
			check("Failed to fetch service config:", err)
			bytes, err := out.MarshalJSON()
			check("Failed to marshal response:", err)
			fmt.Print(string(bytes))
		},
	}

	var delete = &cobra.Command{
		Use:   "delete",
		Short: "Delete service config",
		Run: func(cmd *cobra.Command, args []string) {
			check("Cannot connect to the API:", s.Connect())
			check("Failed to delete service config:", s.Delete())
		},
	}

	var undelete = &cobra.Command{
		Use:   "undelete",
		Short: "Undelete service config",
		Run: func(cmd *cobra.Command, args []string) {
			check("Cannot connect to the API:", s.Connect())
			check("Failed to undelete service config:", s.Undelete())
		},
	}

	fetch.PersistentFlags().StringVarP(&s.Version,
		"version", "v", "",
		"API service config version, empty to use the latest")

	var configCmd = &cobra.Command{
		Use:   "config [command]",
		Short: "Manage ESP service configuration",
	}
	configCmd.PersistentFlags().StringVarP(&s.Name,
		"service", "s", "",
		"API service name")
	configCmd.PersistentFlags().StringVarP(&s.CredentialsFile,
		"creds", "k", "",
		"Service account private key JSON file")
	configCmd.PersistentFlags().StringVarP(&s.ProducerProject,
		"project", "p", "",
		"Service producer project")
	configCmd.AddCommand(submit, fetch, delete, undelete)
	RootCmd.AddCommand(configCmd)
}
