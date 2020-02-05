/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hetzner

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
	apiv1 "k8s.io/api/core/v1"
)

// Manager handles Hetzner communication and data caching of
// node groups (node pools in Hetzner)
type Manager struct {
	client *hcloud.Client
	nodeGroups []*NodeGroup
	cloudConfig *CloudConfig
}

// Config from --cloud-config file
type Config struct {
	Token string                         `json:"token"`
	Endpoint string                      `json:"endpoint"`
	ProviderPrefix string                `json:"provider_prefix"`
	SSHKeys []int                        `json:"ssh_keys"`
	InstanceType string                  `json:"instance_type"`
	Location string                      `json:"location"`
	Image Image	                         `json:"image"`
	CloudInit string                     `json:"cloud_init"`
	Pools map[string]*ConfigPool         `json:"pools"`
	// KubeBootstrap CloudInitTemplateData  `json:kube_bootstrap`
}

// ConfigPool config for node pool
type ConfigPool struct {
	NodeNamePrefix string           `json:"node_name_prefix"`
	SSHKeys []int                   `json:"ssh_keys"`
	InstanceType string             `json:"instance_type"`
	Location string                 `json:"location"`
	Image Image	                    `json:"image"`
	CloudInit string                `json:"cloud_init"`
	NodeLabels map[string]string    `json:"node_labels"`
	NodeTaints []apiv1.Taint        `json:"node_taints"`
	MinNodes int  
	MaxNodes int
}

// // CloudInitTemplateData represents the variables that can be used in cloudinit templates
// type CloudInitTemplateData struct {
// 	BootstrapTokenID     string
// 	BootstrapTokenSecret string
// 	APIServerEndpoint    string
// }

// CloudConfig for init nodes
type CloudConfig struct {
	Pools map[string]*ConfigPool
}

// Image - hetzner image https://godoc.org/github.com/hetznercloud/hcloud-go/hcloud#Image
type Image struct {
	ID int      `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`

}

// Node server https://godoc.org/github.com/hetznercloud/hcloud-go/hcloud#Server
type Node struct {
    ID string
    Name string
    Status hcloud.ServerStatus
}

// NodePool - abstraction for NodeGroups in Hetzner
type NodePool struct {
	ID string
	Name string
	Count int
    MinNodes int
	MaxNodes int
	AutoScale bool
	Nodes []*Node
	InstanceType string
	Location string
	NodeNamePrefix string
	Image Image
	SSHKeys []int
	CloudInit string
}

// NodeGroup implements cloudprovider.NodeGroup interface. NodeGroup contains
// configuration info and functions to control a set of nodes that have the
// same capacity and set of labels.
type NodeGroup struct {
	id string
	client *hcloud.Client
	nodePool *NodePool
	cloudConfig *CloudConfig
	minSize int
	maxSize int
}

var Locations = map[string]bool{
	"hel1": true,
	"nbg1": true,
	"fsn1": true,
}