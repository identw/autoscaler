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
	"strconv"
	"encoding/json"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"io"
	"io/ioutil"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"k8s.io/klog"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
)

const (
	version = "v0.0.1"
)

var (
	hetznerProviderIDPrefix string
)

// newManage create manager for manage nodeGroups
func newManager(configReader io.Reader, cp cloudprovider.NodeGroupDiscoveryOptions) (*Manager, error) {
	cfg := &Config{}
	if configReader != nil {
		body, err := ioutil.ReadAll(configReader)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, cfg)
		if err != nil {
			return nil, err
		}
	}

	if cfg.Token == "" {
		return nil, errors.New("Config hetzner error: access token is not provided")
	}
	if cfg.Endpoint == "" {
		return nil, errors.New("Config hetzner error: endpoint is not provided")
	}
	if cfg.ProviderPrefix == "" {
		hetznerProviderIDPrefix = "hcloud://"
	} else {
		hetznerProviderIDPrefix = cfg.ProviderPrefix
	}
	if hetznerProviderIDPrefix != "hcloud://" && hetznerProviderIDPrefix != "hetzner://" {
		return nil, fmt.Errorf("Config hetzner error: provider_prefix should be equal to either hcloud:// or hetzner://, but it is equal: %v", hetznerProviderIDPrefix)
	}

	// parse  pools from --nodes param
	for _, pool := range cp.NodeGroupSpecs {
		nodes := strings.Split(pool, ":")
		poolName := nodes[2]
		if _, ok := cfg.Pools[poolName]; !ok {
			return nil, fmt.Errorf("Config hetzner pool error: '%s' not found in pools in cloud-config file", poolName)
		}
		cfg.Pools[poolName].MinNodes, _ = strconv.Atoi(nodes[0])
		cfg.Pools[poolName].MaxNodes, _ = strconv.Atoi(nodes[1])
		
		// SSHKeys
		if len(cfg.Pools[poolName].SSHKeys) == 0 && len(cfg.SSHKeys) == 0 {
			return nil, fmt.Errorf("Config hetzner pool error: SSHKeys not defined")
		}
		if len(cfg.Pools[poolName].SSHKeys) == 0 {
			cfg.Pools[poolName].SSHKeys = cfg.SSHKeys
		}

		// InstanceType
		if cfg.Pools[poolName].InstanceType == "" && cfg.InstanceType == "" {
			cfg.Pools[poolName].InstanceType = "cx41"
		}
		if cfg.Pools[poolName].InstanceType == "" {
			cfg.Pools[poolName].InstanceType = cfg.InstanceType
		}
		if _, ok := InstanceTypes[cfg.Pools[poolName].InstanceType]; !ok {
			return nil, fmt.Errorf("Config hetzner pool error: Instance type '%s' not found", cfg.Pools[poolName].InstanceType)
		}

		// Location
		if cfg.Pools[poolName].Location == "" && cfg.Location == "" {
			cfg.Pools[poolName].Location = "nbg1"
		}
		if cfg.Pools[poolName].Location == "" {
			cfg.Pools[poolName].Location = cfg.Location
		}
		if _, ok := Locations[cfg.Pools[poolName].Location]; !ok {
			return nil, fmt.Errorf("Config hetzner pool error: Location '%s' not found", cfg.Pools[poolName].Location)
		}

		// Image
		if cfg.Pools[poolName].Image == (Image{}) && cfg.Image == (Image{}) {
			cfg.Pools[poolName].Image = Image{
				ID: 168855,
				Name: "ubuntu-18.04",
				Type: "system",
			}
		}
		if cfg.Pools[poolName].Image == (Image{}) {
			cfg.Pools[poolName].Image = cfg.Image
		}

		// CloudInit
		if cfg.Pools[poolName].CloudInit == "" && cfg.CloudInit == "" {
			return nil, fmt.Errorf("Config hetzner pool error: CloudInit not defined")
		}
		if cfg.Pools[poolName].CloudInit == "" {
			cfg.Pools[poolName].CloudInit = cfg.CloudInit
		}
		cloudinit, err := base64.StdEncoding.DecodeString(cfg.Pools[poolName].CloudInit)
		if err != nil {
			return nil, fmt.Errorf("Config hetzner pool error: CloudInit not base64. Please encode CloudInit in base64. Base64 decode error: %v", err)
		}
		cfg.Pools[poolName].CloudInit = string(cloudinit)
	}

	opts := []hcloud.ClientOption{
		hcloud.WithToken(cfg.Token),
		hcloud.WithApplication("cluster-autoscaler", version),
		hcloud.WithEndpoint(cfg.Endpoint),
	}

	var cloudConfig CloudConfig
	cloudConfig.Pools = cfg.Pools

	hcloudClient := hcloud.NewClient(opts...)

	m := &Manager{
		client:     hcloudClient,
		nodeGroups: make([]*NodeGroup, 0),
		cloudConfig: &cloudConfig,
	}

	return m, nil
}

// Refresh refreshes the cache holding the nodegroups. This is called by the CA
// based on the `--scan-interval`. By default it's 10 seconds.
func (m *Manager) Refresh() error {
	nodePools, err:= listNodePools(m)
	if err != nil {
		return err
	}

	var group []*NodeGroup
	for _, nodePool := range nodePools {
		if !nodePool.AutoScale {
			continue
		}

		klog.V(4).Infof("adding node pool: %q name: %s min: %d max: %d",
			nodePool.ID, nodePool.Name, nodePool.MinNodes, nodePool.MaxNodes)

		group = append(group, &NodeGroup{
			id:          nodePool.ID,
			client:      m.client,
			cloudConfig: m.cloudConfig,
			nodePool:    nodePool,
			minSize:     nodePool.MinNodes,
			maxSize:     nodePool.MaxNodes,
		})
	}

	if len(group) == 0 {
		klog.V(4).Info("cluster-autoscaler is disabled. no node pools are configured")
	}

	m.nodeGroups = group
	return nil
}