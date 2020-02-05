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

// InstanceType type server
type InstanceType struct {
	InstanceType string
	VCPU         int64
	Memory       int64
	Storage      int64
}

// InstanceTypes which types of servers
var InstanceTypes = map[string]*InstanceType{
	"cx51": {
		InstanceType: "cx51",
		VCPU:         8,
		Memory:       32000000000,
		Storage:      240000000000,      
		
	},
	"cx41": {
		InstanceType: "cx41",
		VCPU:         4,
		Memory:       16000000000,
		Storage:      160000000000,
	},
	"cx31": {
		InstanceType: "cx31",
		VCPU:         2,
		Memory:       8000000000,
		Storage:      80000000000,
	},
	"cx21": {
		InstanceType: "cx21",
		VCPU:         2,
		Memory:       4000000000,
		Storage:      40000000000,
	},
	"cx11": {
		InstanceType: "cx11",
		VCPU:         1,
		Memory:       2000000000,
		Storage:      20000000000,
	},
	"cx51-ceph": {
		InstanceType: "cx51-ceph",
		VCPU:         8,
		Memory:       32000000000,
		Storage:      240000000000,      
		
	},
	"cx41-ceph": {
		InstanceType: "cx41-ceph",
		VCPU:         4,
		Memory:       16000000000,
		Storage:      160000000000,
	},
	"cx31-ceph": {
		InstanceType: "cx31-ceph",
		VCPU:         2,
		Memory:       8000000000,
		Storage:      80000000000,
	},
	"cx21-ceph": {
		InstanceType: "cx21-ceph",
		VCPU:         2,
		Memory:       4000000000,
		Storage:      40000000000,
	},
	"cx11-ceph": {
		InstanceType: "cx11-ceph",
		VCPU:         1,
		Memory:       2000000000,
		Storage:      20000000000,
	},
	"ccx51": {
		InstanceType: "ccx51",
		VCPU:         32,
		Memory:       128000000000,
		Storage:      600000000000,      
		
	},
	"ccx41": {
		InstanceType: "ccx41",
		VCPU:         16,
		Memory:       64000000000,
		Storage:      360000000000,
	},
	"ccx31": {
		InstanceType: "ccx31",
		VCPU:         8,
		Memory:       32000000000,
		Storage:      240000000000,
	},
	"ccx21": {
		InstanceType: "ccx21",
		VCPU:         4,
		Memory:       16000000000,
		Storage:      160000000000,
	},
	"ccx11": {
		InstanceType: "ccx11",
		VCPU:         2,
		Memory:       8000000000,
		Storage:      80000000000,
	},
}