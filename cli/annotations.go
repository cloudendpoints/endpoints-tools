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

// Prefix for ESP managed resources
const EndpointsPrefix = "endpoints-"

// Label to tag pods running ESP for a given k8s service
// Label to tag ESP services for a given k8s service
const AnnotationManagedService = "endpoints.googleapis.com/managed-service"

// Label to tag ESP pods for loosely coupled ESP services
// (to support multiple loose ESP deployments for a service)
const AnnotationEndpointsService = "endpoints.googleapis.com/endpoints"

// Label to tag ESP services with service config name and id
const AnnotationConfigName = "endpoints.googleapis.com/config-name"
const AnnotationConfigId = "endpoints.googleapis.com/config-id"

// Label to tag ESP services with the type of the deployment
const AnnotationDeploymentType = "endpoints.googleapis.com/deployment-type"
