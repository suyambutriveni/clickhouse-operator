// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
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

package v1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ICustomResource interface {
	GetNamespace() string
	GetName() string
	GetLabels() map[string]string
	GetAnnotations() map[string]string
	GetRuntime() ICustomResourceRuntime
	GetRootServiceTemplate() (*ServiceTemplate, bool)
	GetObjectMeta() meta.Object

	WalkClusters(f func(cluster ICluster) error) []error
	WalkHosts(func(host *Host) error) []error
}

type ICluster interface {
	GetName() string
	GetRuntime() IClusterRuntime
	GetServiceTemplate() (*ServiceTemplate, bool)
	GetSecret() *ClusterSecret
	GetPDBMaxUnavailable() *Int32

	WalkShards(f func(index int, shard IShard) error) []error
	WalkHosts(func(host *Host) error) []error
}

type IShard interface {
	GetName() string
	GetRuntime() IShardRuntime
	GetServiceTemplate() (*ServiceTemplate, bool)
	GetInternalReplication() *StringBool
	HasWeight() bool
	GetWeight() int

	WalkHosts(func(host *Host) error) []error
}

type IReplica interface {
	GetName() string
}

type IHost interface {
	GetRuntime() IHostRuntime
}
