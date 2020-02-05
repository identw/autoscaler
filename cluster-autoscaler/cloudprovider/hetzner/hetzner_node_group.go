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
	"fmt"
	"math/rand"
	"time"
	"strconv"

	"k8s.io/klog"
	"github.com/hetznercloud/hcloud-go/hcloud"
	apiv1 "k8s.io/api/core/v1"
	kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// MaxSize returns maximum size of the node group.
func (n *NodeGroup) MaxSize() int {
	return n.maxSize
}

// MinSize returns minimum size of the node group.
func (n *NodeGroup) MinSize() int {
	return n.minSize
}

// TargetSize returns the current target size of the node group. It is possible
// that the number of nodes in Kubernetes is different at the moment but should
// be equal to Size() once everything stabilizes (new nodes finish startup and
// registration or removed nodes are deleted completely). Implementation
// required.
func (n *NodeGroup) TargetSize() (int, error) {
	return n.nodePool.Count, nil
}

// IncreaseSize increases the size of the node group. To delete a node you need
// to explicitly name it and use DeleteNode. This function should wait until
// node group size is updated. Implementation required.
func (n *NodeGroup) IncreaseSize(delta int) error {
	if delta <= 0 {
		return fmt.Errorf("delta must be positive, have: %d", delta)
	}

	targetSize := n.nodePool.Count + delta

	if targetSize > n.MaxSize() {
		return fmt.Errorf("size increase is too large. current: %d desired: %d max: %d",
			n.nodePool.Count, targetSize, n.MaxSize())
	}

	err := createNodes(n, delta)
	if err != nil {
		return err
	}
	n.nodePool.Count = targetSize
	return nil
}

// DeleteNodes deletes nodes from this node group (and also increasing the size
// of the node group with that). Error is returned either on failure or if the
// given node doesn't belong to this node group. This function should wait
// until node group size is updated. Implementation required.
func (n *NodeGroup) DeleteNodes(nodes []*apiv1.Node) error {
	for _, node := range nodes {
		nodeID := toNodeID(node.Spec.ProviderID)

		err := deleteNode(n.client, nodeID)
		if err != nil {
			return fmt.Errorf("deleting node failed for cluster, node pool: %q node: %q: %s", n.id, nodeID, err)
		}

		// decrement the count by one  after a successful delete
		n.nodePool.Count--
	}

	return nil
}

// DecreaseTargetSize decreases the target size of the node group. This function
// doesn't permit to delete any existing node and can be used only to reduce the
// request for new nodes that have not been yet fulfilled. Delta should be negative.
// It is assumed that cloud provider will not delete the existing nodes when there
// is an option to just decrease the target. Implementation required.
func (n *NodeGroup) DecreaseTargetSize(delta int) error {
	if delta >= 0 {
		return fmt.Errorf("delta must be negative, have: %d", delta)
	}

	targetSize := n.nodePool.Count + delta
	if targetSize <= n.MinSize() {
		return fmt.Errorf("size decrease is too small. current: %d desired: %d min: %d",
			n.nodePool.Count, targetSize, n.MinSize())
	}
	if delta >= 0 {
		return fmt.Errorf("size decrease must be negative")
	}
	klog.V(0).Infof("Decreasing target size by %d, %d->%d", delta, n.nodePool.Count, targetSize)
	n.nodePool.Count = targetSize
	return fmt.Errorf("could not decrease target size")
}

// Id returns an unique identifier of the node group.
func (n *NodeGroup) Id() string {
	return n.id
}

// Debug returns a string containing all information regarding this node group.
func (n *NodeGroup) Debug() string {
	return fmt.Sprintf("cluster ID: %s (min:%d max:%d)", n.Id(), n.MinSize(), n.MaxSize())
}

// Nodes returns a list of all nodes that belong to this node group.  It is
// required that Instance objects returned by this method have Id field set.
// Other fields are optional.
func (n *NodeGroup) Nodes() ([]cloudprovider.Instance, error) {
	if n.nodePool == nil {
		return nil, fmt.Errorf("nodePool is nil")
	}

	// TODO: get server from api for pool
	// n.nodePool, err = listNodePools(n.client)
	return toInstances(n.nodePool.Nodes), nil
}

// TemplateNodeInfo returns a schedulernodeinfo.NodeInfo structure of an empty
// (as if just started) node. This will be used in scale-up simulations to
// predict what would a new node look like if a node group was expanded. The
// returned NodeInfo is expected to have a fully populated Node object, with
// all of the labels, capacity and allocatable information as well as all pods
// that are started on the node by default, using manifest (most likely only
// kube-proxy). Implementation optional.
func (n *NodeGroup) TemplateNodeInfo() (*schedulernodeinfo.NodeInfo, error) {
	node := apiv1.Node{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	postfixString := strconv.Itoa(r.Intn(10000) + 10000)
	nodeName := n.nodePool.NodeNamePrefix + "-" + postfixString

	node.ObjectMeta = metav1.ObjectMeta{
		Name:     nodeName,
		SelfLink: fmt.Sprintf("/api/v1/nodes/%s", nodeName),
		Labels:   map[string]string{},
	}

	node.Status = apiv1.NodeStatus{
		Capacity: apiv1.ResourceList{},
	}

	instanceType := InstanceTypes[n.nodePool.InstanceType]
	
	node.Status.Capacity[apiv1.ResourcePods] = *resource.NewQuantity(110, resource.DecimalSI)
	node.Status.Capacity[apiv1.ResourceCPU] = *resource.NewQuantity(instanceType.VCPU, resource.DecimalSI)
	node.Status.Capacity[apiv1.ResourceMemory] = *resource.NewQuantity(instanceType.Memory, resource.DecimalSI)
	node.Status.Capacity[apiv1.ResourceEphemeralStorage] = *resource.NewQuantity(instanceType.Storage, resource.DecimalSI)
	node.Status.Allocatable = node.Status.Capacity

	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	dcNumber := strconv.Itoa(r.Intn(10))

	node.Labels[apiv1.LabelHostname] = nodeName
	node.Labels[kubeletapis.LabelArch] = cloudprovider.DefaultArch
	node.Labels[kubeletapis.LabelOS] = cloudprovider.DefaultOS

	node.Labels[apiv1.LabelInstanceType] = n.nodePool.InstanceType
	node.Labels[apiv1.LabelZoneRegion] = n.nodePool.Location
	node.Labels[apiv1.LabelZoneFailureDomain] = n.nodePool.Location + "-dc" + dcNumber
	if len(n.cloudConfig.Pools[n.id].NodeLabels) != 0 {
		for labelKey, labelValue := range n.cloudConfig.Pools[n.id].NodeLabels {
			node.Labels[labelKey] = labelValue
		}
	}

	if len(n.cloudConfig.Pools[n.id].NodeTaints) != 0 {
		node.Spec = apiv1.NodeSpec{
			Taints: n.cloudConfig.Pools[n.id].NodeTaints,
		}
	}
	node.Status.Conditions = cloudprovider.BuildReadyConditions()

	nodeInfo := schedulernodeinfo.NewNodeInfo(cloudprovider.BuildKubeProxy(nodeName))
	nodeInfo.SetNode(&node)
	return nodeInfo, nil
}

// Exist checks if the node group really exists on the cloud provider side.
// Allows to tell the theoretical node group from the real one. Implementation
// required.
func (n *NodeGroup) Exist() bool {
	return n.nodePool != nil
}

// Create creates the node group on the cloud provider side. Implementation
// optional.
func (n *NodeGroup) Create() (cloudprovider.NodeGroup, error) {
	return nil, cloudprovider.ErrNotImplemented
}

// Delete deletes the node group on the cloud provider side.  This will be
// executed only for autoprovisioned node groups, once their size drops to 0.
// Implementation optional.
func (n *NodeGroup) Delete() error {
	return cloudprovider.ErrNotImplemented
}

// Autoprovisioned returns true if the node group is autoprovisioned. An
// autoprovisioned group was created by CA and can be deleted when scaled to 0.
func (n *NodeGroup) Autoprovisioned() bool {
	return false
}

// toInstances converts a slice of *godo.KubernetesNode to
// cloudprovider.Instance
func toInstances(nodes []*Node) []cloudprovider.Instance {
	instances := make([]cloudprovider.Instance, 0, len(nodes))
	for _, nd := range nodes {
		instances = append(instances, toInstance(nd))
	}
	return instances
}

// toInstance converts the given *godo.KubernetesNode to a
// cloudprovider.Instance
func toInstance(node *Node) cloudprovider.Instance {
	return cloudprovider.Instance{
		Id:     toProviderID(node.ID),
		Status: toInstanceStatus(node.Status),
	}
}

// toInstanceStatus converts the given *godo.KubernetesNodeStatus to a
// cloudprovider.InstanceStatus
func toInstanceStatus(nodeState hcloud.ServerStatus) *cloudprovider.InstanceStatus {
	st := &cloudprovider.InstanceStatus{}
	switch nodeState {
	case hcloud.ServerStatusInitializing, hcloud.ServerStatusStarting:
		st.State = cloudprovider.InstanceCreating
	case hcloud.ServerStatusRunning:
		st.State = cloudprovider.InstanceRunning
	case hcloud.ServerStatusDeleting:
		st.State = cloudprovider.InstanceDeleting
	default:
		st.ErrorInfo = &cloudprovider.InstanceErrorInfo{
			ErrorClass:   cloudprovider.OtherErrorClass,
			ErrorCode:    "no-code-digitalocean",
			ErrorMessage: "Error status node",
		}
	}

	return st
}
