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

package hetznerIdentw

import (
	"context"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/hetzner/hcloud-go/hcloud"
	"k8s.io/klog/v2"
	"math/rand"
	"strconv"
	"time"
)

// listNodePools get servers from API by labelSelector
// One label = one pool (nodeGroup)
func listNodePools(m *Manager) ([]*NodePool, error) {
	var hetznerPools []*NodePool
	pools := make(map[string]*NodePool)

	servers, err := m.cache.getServers()
	if err != nil {
		klog.Errorf("listNodePools() error get servers from cache, error: %v\n", err)
		return nil, err
	}

	// init Pools
	for poolName, pool := range m.cloudConfig.Pools {
		pools[poolName] = &NodePool{
			ID:             poolName,
			Name:           "Hetzner k8s autoscaler: " + poolName,
			Count:          0,
			MinNodes:       pool.MinNodes,
			MaxNodes:       pool.MaxNodes,
			AutoScale:      true,
			Nodes:          []*Node{},
			InstanceType:   pool.InstanceType,
			Location:       pool.Location,
			NodeNamePrefix: pool.NodeNamePrefix,
			Image:          pool.Image,
			SSHKeys:        pool.SSHKeys,
			CloudInit:      pool.CloudInit,
		}
	}

	for _, s := range servers {
		poolName := ""
		for name, _ := range m.cloudConfig.Pools {
			if _, ok := s.Labels[name]; ok {
				poolName = name
				break
			}
		}
		if poolName == "" {
			continue
		}
		var node Node
		node.ID = strconv.Itoa(int(s.ID))
		node.Name = s.Name
		node.Status = s.Status

		pools[poolName].Nodes = append(pools[poolName].Nodes, &node)
		pools[poolName].Count++
	}

	for _, pool := range pools {
		hetznerPools = append(hetznerPools, pool)
	}

	return hetznerPools, nil
}

// createNode crete server in Hetzner for nodePool(nodeGroup)
// TODO: check exists nodes to randomName != exists node
func createNode(n *NodeGroup) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	postfixString := strconv.Itoa(r.Intn(10000) + 10000)

	serverType := &hcloud.ServerType{
		Name: n.nodePool.InstanceType,
	}
	image := &hcloud.Image{
		ID:   n.nodePool.Image.ID,
		Name: n.nodePool.Image.Name,
		Type: hcloud.ImageType(n.nodePool.Image.Type),
	}
	location := &hcloud.Location{
		Name: n.nodePool.Location,
	}

	sshKeys := make([]*hcloud.SSHKey, 0, len(n.nodePool.SSHKeys))
	for _, sshKey := range n.nodePool.SSHKeys {
		sshKeys = append(sshKeys, &hcloud.SSHKey{ID: sshKey})
	}

	labels := make(map[string]string)
	labels[n.nodePool.ID] = ""
	serverCreateOpts := hcloud.ServerCreateOpts{
		Name:       n.nodePool.NodeNamePrefix + "-" + postfixString,
		ServerType: serverType,
		UserData:   n.nodePool.CloudInit,
		SSHKeys:    sshKeys,
		Image:      image,
		Location:   location,
		Labels:     labels,
	}

	_, _, err := n.client.Server.Create(context.Background(), serverCreateOpts)
	if err != nil {
		return err
	}
	return nil
}

// createNodes create multiple servers for nodeGroups
func createNodes(n *NodeGroup, amountNodes int) error {
	for i := 0; i < amountNodes; i++ {
		err := createNode(n)
		if err != nil {
			return err
		}
	}
	return nil
}

// deleteNode delete server by nodeID
func deleteNode(client *hcloud.Client, nodeID string) error {
	id, _ := strconv.Atoi(nodeID)
	server := &hcloud.Server{
		ID: int64(id),
	}

	_, err := client.Server.Delete(context.Background(), server)
	if err != nil {
		return err
	}
	return nil
}
