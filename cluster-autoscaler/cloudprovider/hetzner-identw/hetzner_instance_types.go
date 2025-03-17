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

// InstanceType type server
type InstanceType struct {
	InstanceType string
	VCPU         int64
	Memory       int64
	Storage      int64
}

// InstanceTypes which types of servers
var InstanceTypes = map[string]*InstanceType{
	"cx11": {
		InstanceType: "cx11",
		VCPU:         1,
		Memory:       2000000000,
		Storage:      20000000000,
	},
	"cpx11": {
		InstanceType: "cpx11",
		VCPU:         2,
		Memory:       2000000000,
		Storage:      40000000000,
	},
	"ccx11": {
		InstanceType: "ccx11",
		VCPU:         2,
		Memory:       8000000000,
		Storage:      80000000000,
	},
	"ccx12": {
		InstanceType: "ccx12",
		VCPU:         2,
		Memory:       8000000000,
		Storage:      80000000000,
	},
	"ccx13": {
		InstanceType: "ccx13",
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
	"cx22": {
		InstanceType: "cx22",
		VCPU:         2,
		Memory:       4000000000,
		Storage:      40000000000,
	},
	"cpx21": {
		InstanceType: "cpx21",
		VCPU:         3,
		Memory:       4000000000,
		Storage:      80000000000,
	},
	"ccx21": {
		InstanceType: "ccx21",
		VCPU:         4,
		Memory:       16000000000,
		Storage:      160000000000,
	},
	"ccx22": {
		InstanceType: "ccx22",
		VCPU:         4,
		Memory:       16000000000,
		Storage:      160000000000,
	},
	"ccx23": {
		InstanceType: "ccx23",
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
	"cx32": {
		InstanceType: "cx32",
		VCPU:         4,
		Memory:       8000000000,
		Storage:      80000000000,
	},
	"cpx31": {
		InstanceType: "cpx31",
		VCPU:         4,
		Memory:       8000000000,
		Storage:      160000000000,
	},
	"ccx31": {
		InstanceType: "ccx31",
		VCPU:         8,
		Memory:       32000000000,
		Storage:      240000000000,
	},
	"ccx32": {
		InstanceType: "ccx32",
		VCPU:         8,
		Memory:       32000000000,
		Storage:      240000000000,
	},
	"ccx33": {
		InstanceType: "ccx33",
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
	"cx42": {
		InstanceType: "cx42",
		VCPU:         8,
		Memory:       16000000000,
		Storage:      160000000000,
	},
	"cpx41": {
		InstanceType: "cpx41",
		VCPU:         8,
		Memory:       16000000000,
		Storage:      240000000000,
	},
	"ccx41": {
		InstanceType: "ccx41",
		VCPU:         16,
		Memory:       64000000000,
		Storage:      360000000000,
	},
	"ccx42": {
		InstanceType: "ccx42",
		VCPU:         16,
		Memory:       64000000000,
		Storage:      360000000000,
	},
	"ccx43": {
		InstanceType: "ccx43",
		VCPU:         16,
		Memory:       64000000000,
		Storage:      360000000000,
	},
	"cx51": {
		InstanceType: "cx51",
		VCPU:         8,
		Memory:       32000000000,
		Storage:      240000000000,      
	},
	"cx52": {
		InstanceType: "cx52",
		VCPU:         16,
		Memory:       32000000000,
		Storage:      360000000000,
	},
	"cpx51": {
		InstanceType: "cpx51",
		VCPU:         16,
		Memory:       32000000000,
		Storage:      360000000000,      
		
	},
	"ccx51": {
		InstanceType: "ccx51",
		VCPU:         32,
		Memory:       128000000000,
		Storage:      600000000000,      
		
	},
	"ccx52": {
		InstanceType: "ccx52",
		VCPU:         32,
		Memory:       128000000000,
		Storage:      600000000000,
		
	},
	"ccx53": {
		InstanceType: "ccx53",
		VCPU:         32,
		Memory:       128000000000,
		Storage:      600000000000,
	},
	"ccx62": {
		InstanceType: "ccx62",
		VCPU:         48,
		Memory:       192000000000,
		Storage:      960000000000,      
	},
	"ccx63": {
		InstanceType: "ccx63",
		VCPU:         48,
		Memory:       192000000000,
		Storage:      960000000000,
	},
}