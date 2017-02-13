# Extensible Service Proxy tools #

This repository hosts helper tools for the [Extensible Service Proxy](https://cloud.google.com/endpoints/).

These tools are not official Google products.

# Code organization #

This repository hosts two tools used by the Extensible Service Proxy:
* [`start_esp`](/start_esp) is a Python start-up script for the proxy. The script includes a generic nginx configuration template and fetching logic to retrieve service configuration from Google Service Management service.
* [helper script for Kubernetes deployment](#esp-cli) is a Golang command line utility that automates ESP injection as a sidecar container in Kubernetes deployments and configuring the proxy. This script is used by a single-line ESP tutorial, and shows how to use the start-up script.

## ESP CLI ##

ESP command line tool (`espcli`) lets you try out the Extensible Service Proxy on a Kubernetes cluster using a simple two-step deployment process. ESP CLI depends on `kubectl` and its active kube configuration. You also need a [Google Cloud Platform](http://cloud.google.com) project.

### Quick Example ###

First, bring up a [Kubernetes deployment](https://raw.githubusercontent.com/GoogleCloudPlatform/endpoints-samples/master/k8s/echo_http.yaml) consisting of `echo` service and `echo` pods on your GCP cluster:

    kubectl create -f echo_http.yaml

Next, run the following command to inject ESP container into `echo` pods and create a new service called `endpoints`:

    espcli deploy echo endpoints --project MY_PROJECT -e LoadBalancer

Replace `MY_PROJECT` above with the name of your GCP project.
Then make a request to the external IP address of `endpoints` service as follows:

    curl -d '{"message":"hello world"}' -H "content-type:application/json" http://ENDPOINTS_IP/echo

Go to the [Endpoints UI](https://console.cloud.google.com/endpoints) and see the detailed logging and monitoring information for the newly created API.

## Build instructions ##

We use [Bazel](https://bazel.io) to build the ESP tools.
To build `espcli`, run the following command:

    bazel build :espcli

This command fetches all required dependencies and produces a single binary
`bazel-bin/espcli`. To read more about its usage and flags,
please consult:

    bazel run :espcli -- help


## Contributing ##

Your contributions are welcome. Please follow the [contributor
guidelines](/CONTRIBUTING.md).
