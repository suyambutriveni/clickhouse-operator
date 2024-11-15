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
	"sort"
	"sync"

	"github.com/altinity/clickhouse-operator/pkg/apis/common/types"
	"github.com/altinity/clickhouse-operator/pkg/util"
	"github.com/altinity/clickhouse-operator/pkg/version"
)

const (
	maxActions = 10
	maxErrors  = 10
	maxTaskIDs = 10
)

// Possible CR statuses
const (
	StatusInProgress  = "InProgress"
	StatusCompleted   = "Completed"
	StatusAborted     = "Aborted"
	StatusTerminating = "Terminating"
)

// Status defines status section of the custom resource.
//
// Note: application level reads and writes to Status fields should be done through synchronized getter/setter functions.
// While all of these fields need to be exported for JSON and YAML serialization/deserialization, we can at least audit
// that application logic sticks to the synchronized getter/setters by auditing whether all explicit Go field-level
// accesses are strictly within _this_ source file OR the generated deep copy source file.
type Status struct {
	CHOpVersion            string                  `json:"chop-version,omitempty"           yaml:"chop-version,omitempty"`
	CHOpCommit             string                  `json:"chop-commit,omitempty"            yaml:"chop-commit,omitempty"`
	CHOpDate               string                  `json:"chop-date,omitempty"              yaml:"chop-date,omitempty"`
	CHOpIP                 string                  `json:"chop-ip,omitempty"                yaml:"chop-ip,omitempty"`
	ClustersCount          int                     `json:"clusters,omitempty"               yaml:"clusters,omitempty"`
	ShardsCount            int                     `json:"shards,omitempty"                 yaml:"shards,omitempty"`
	ReplicasCount          int                     `json:"replicas,omitempty"               yaml:"replicas,omitempty"`
	HostsCount             int                     `json:"hosts,omitempty"                  yaml:"hosts,omitempty"`
	Status                 string                  `json:"status,omitempty"                 yaml:"status,omitempty"`
	TaskID                 string                  `json:"taskID,omitempty"                 yaml:"taskID,omitempty"`
	TaskIDsStarted         []string                `json:"taskIDsStarted,omitempty"         yaml:"taskIDsStarted,omitempty"`
	TaskIDsCompleted       []string                `json:"taskIDsCompleted,omitempty"       yaml:"taskIDsCompleted,omitempty"`
	Action                 string                  `json:"action,omitempty"                 yaml:"action,omitempty"`
	Actions                []string                `json:"actions,omitempty"                yaml:"actions,omitempty"`
	Error                  string                  `json:"error,omitempty"                  yaml:"error,omitempty"`
	Errors                 []string                `json:"errors,omitempty"                 yaml:"errors,omitempty"`
	HostsUpdatedCount      int                     `json:"hostsUpdated,omitempty"           yaml:"hostsUpdated,omitempty"`
	HostsAddedCount        int                     `json:"hostsAdded,omitempty"             yaml:"hostsAdded,omitempty"`
	HostsUnchangedCount    int                     `json:"hostsUnchanged,omitempty"         yaml:"hostsUnchanged,omitempty"`
	HostsFailedCount       int                     `json:"hostsFailed,omitempty"            yaml:"hostsFailed,omitempty"`
	HostsCompletedCount    int                     `json:"hostsCompleted,omitempty"         yaml:"hostsCompleted,omitempty"`
	HostsDeletedCount      int                     `json:"hostsDeleted,omitempty"           yaml:"hostsDeleted,omitempty"`
	HostsDeleteCount       int                     `json:"hostsDelete,omitempty"            yaml:"hostsDelete,omitempty"`
	Pods                   []string                `json:"pods,omitempty"                   yaml:"pods,omitempty"`
	PodIPs                 []string                `json:"pod-ips,omitempty"                yaml:"pod-ips,omitempty"`
	FQDNs                  []string                `json:"fqdns,omitempty"                  yaml:"fqdns,omitempty"`
	Endpoint               string                  `json:"endpoint,omitempty"               yaml:"endpoint,omitempty"`
	NormalizedCR           *ClickHouseInstallation `json:"normalized,omitempty"             yaml:"normalized,omitempty"`
	NormalizedCRCompleted  *ClickHouseInstallation `json:"normalizedCompleted,omitempty"    yaml:"normalizedCompleted,omitempty"`
	HostsWithTablesCreated []string                `json:"hostsWithTablesCreated,omitempty" yaml:"hostsWithTablesCreated,omitempty"`
	UsedTemplates          []*TemplateRef          `json:"usedTemplates,omitempty"          yaml:"usedTemplates,omitempty"`

	mu sync.RWMutex `json:"-" yaml:"-"`
}

// FillStatusParams is a struct used to fill status params
type FillStatusParams struct {
	CHOpIP              string
	ClustersCount       int
	ShardsCount         int
	HostsCount          int
	TaskID              string
	HostsUpdatedCount   int
	HostsAddedCount     int
	HostsUnchangedCount int
	HostsCompletedCount int
	HostsDeleteCount    int
	HostsDeletedCount   int
	Pods                []string
	FQDNs               []string
	Endpoint            string
	NormalizedCR        *ClickHouseInstallation
}

// Fill is a synchronized setter for a fairly large number of fields. We take a struct type "params" argument to avoid
// confusion of similarly typed positional arguments, and to avoid defining a lot of separate synchronized setters
// for these fields that are typically all set together at once (during "fills").
func (s *Status) Fill(params *FillStatusParams) {
	doWithWriteLock(s, func(s *Status) {
		// We always set these (build-hardcoded) version fields.
		s.CHOpVersion = version.Version
		s.CHOpCommit = version.GitSHA
		s.CHOpDate = version.BuiltAt

		// Now, set fields from the provided input.
		s.CHOpIP = params.CHOpIP
		s.ClustersCount = params.ClustersCount
		s.ShardsCount = params.ShardsCount
		s.HostsCount = params.HostsCount
		s.TaskID = params.TaskID
		s.HostsUpdatedCount = params.HostsUpdatedCount
		s.HostsAddedCount = params.HostsAddedCount
		s.HostsUnchangedCount = params.HostsUnchangedCount
		s.HostsCompletedCount = params.HostsCompletedCount
		s.HostsDeleteCount = params.HostsDeleteCount
		s.HostsDeletedCount = params.HostsDeletedCount
		s.Pods = params.Pods
		s.FQDNs = params.FQDNs
		s.Endpoint = params.Endpoint
		s.NormalizedCR = params.NormalizedCR
	})
}

// SetError sets status error
func (s *Status) SetError(err string) {
	doWithWriteLock(s, func(s *Status) {
		s.Error = err
	})
}

// SetAndPushError sets and pushes error into status
func (s *Status) SetAndPushError(err string) {
	doWithWriteLock(s, func(s *Status) {
		s.Error = err
		s.Errors = append([]string{err}, s.Errors...)
		if len(s.Errors) > maxErrors {
			s.Errors = s.Errors[:maxErrors]
		}
	})
}

// PushHostTablesCreated pushes host to the list of hosts with created tables
func (s *Status) PushHostTablesCreated(host string) {
	doWithWriteLock(s, func(s *Status) {
		if util.InArray(host, s.HostsWithTablesCreated) {
			return
		}
		s.HostsWithTablesCreated = append(s.HostsWithTablesCreated, host)
	})
}

// SyncHostTablesCreated syncs list of hosts with tables created with actual list of hosts
func (s *Status) SyncHostTablesCreated() {
	doWithWriteLock(s, func(s *Status) {
		if s.FQDNs == nil {
			return
		}
		s.HostsWithTablesCreated = util.IntersectStringArrays(s.HostsWithTablesCreated, s.FQDNs)
	})
}

// PushUsedTemplate pushes used template to the list of used templates
func (s *Status) PushUsedTemplate(templateRef *TemplateRef) {
	doWithWriteLock(s, func(s *Status) {
		s.UsedTemplates = append(s.UsedTemplates, templateRef)
	})
}

// GetUsedTemplatesCount gets used templates count
func (s *Status) GetUsedTemplatesCount() int {
	return getIntWithReadLock(s, func(s *Status) int {
		return len(s.UsedTemplates)
	})
}

// SetAction action setter
func (s *Status) SetAction(action string) {
	doWithWriteLock(s, func(s *Status) {
		s.Action = action
	})
}

// HasNormalizedCRCompleted is a checker
func (s *Status) HasNormalizedCRCompleted() bool {
	return s.GetNormalizedCRCompleted() != nil
}

// HasNormalizedCR is a checker
func (s *Status) HasNormalizedCR() bool {
	return s.GetNormalizedCR() != nil
}

// PushAction pushes action into status
func (s *Status) PushAction(action string) {
	doWithWriteLock(s, func(s *Status) {
		s.Actions = append([]string{action}, s.Actions...)
		trimActionsNoSync(s)
	})
}

// PushError sets and pushes error into status
func (s *Status) PushError(error string) {
	doWithWriteLock(s, func(s *Status) {
		s.Errors = append([]string{error}, s.Errors...)
		if len(s.Errors) > maxErrors {
			s.Errors = s.Errors[:maxErrors]
		}
	})
}

// SetPodIPs sets pod IPs
func (s *Status) SetPodIPs(podIPs []string) {
	doWithWriteLock(s, func(s *Status) {
		s.PodIPs = podIPs
	})
}

// HostDeleted increments deleted hosts counter
func (s *Status) HostDeleted() {
	doWithWriteLock(s, func(s *Status) {
		s.HostsDeletedCount++
	})
}

// HostUpdated increments updated hosts counter
func (s *Status) HostUpdated() {
	doWithWriteLock(s, func(s *Status) {
		s.HostsUpdatedCount++
	})
}

// HostAdded increments added hosts counter
func (s *Status) HostAdded() {
	doWithWriteLock(s, func(s *Status) {
		s.HostsAddedCount++
	})
}

// HostUnchanged increments unchanged hosts counter
func (s *Status) HostUnchanged() {
	doWithWriteLock(s, func(s *Status) {
		s.HostsUnchangedCount++
	})
}

// HostFailed increments failed hosts counter
func (s *Status) HostFailed() {
	doWithWriteLock(s, func(s *Status) {
		s.HostsFailedCount++
	})
}

// HostCompleted increments completed hosts counter
func (s *Status) HostCompleted() {
	doWithWriteLock(s, func(s *Status) {
		s.HostsCompletedCount++
	})
}

// ReconcileStart marks reconcile start
func (s *Status) ReconcileStart(deleteHostsCount int) {
	doWithWriteLock(s, func(s *Status) {
		if s == nil {
			return
		}
		s.Status = StatusInProgress
		s.HostsUpdatedCount = 0
		s.HostsAddedCount = 0
		s.HostsUnchangedCount = 0
		s.HostsCompletedCount = 0
		s.HostsDeletedCount = 0
		s.HostsDeleteCount = deleteHostsCount
		pushTaskIDStartedNoSync(s)
	})
}

// ReconcileComplete marks reconcile completion
func (s *Status) ReconcileComplete() {
	doWithWriteLock(s, func(s *Status) {
		if s == nil {
			return
		}
		s.Status = StatusCompleted
		s.Action = ""
		pushTaskIDCompletedNoSync(s)
	})
}

// ReconcileAbort marks reconcile abortion
func (s *Status) ReconcileAbort() {
	doWithWriteLock(s, func(s *Status) {
		if s == nil {
			return
		}
		s.Status = StatusAborted
		s.Action = ""
		pushTaskIDCompletedNoSync(s)
	})
}

// DeleteStart marks deletion start
func (s *Status) DeleteStart() {
	doWithWriteLock(s, func(s *Status) {
		if s == nil {
			return
		}
		s.Status = StatusTerminating
		s.HostsUpdatedCount = 0
		s.HostsAddedCount = 0
		s.HostsUnchangedCount = 0
		s.HostsCompletedCount = 0
		s.HostsDeletedCount = 0
		s.HostsDeleteCount = 0
		pushTaskIDStartedNoSync(s)
	})
}

// CopyFrom copies the state of a given Status f into the receiver Status of the call.
func (s *Status) CopyFrom(f *Status, opts types.CopyStatusOptions) {
	doWithWriteLock(s, func(s *Status) {
		doWithReadLock(f, func(from *Status) {
			if s == nil || from == nil {
				return
			}

			if opts.InheritableFields {
				s.TaskIDsStarted = from.TaskIDsStarted
				s.TaskIDsCompleted = from.TaskIDsCompleted
				s.Actions = from.Actions
				s.Errors = from.Errors
				s.HostsWithTablesCreated = from.HostsWithTablesCreated
			}

			if opts.Actions {
				s.Action = from.Action
				mergeActionsNoSync(s, from)
				s.HostsWithTablesCreated = nil
				if len(from.HostsWithTablesCreated) > 0 {
					s.HostsWithTablesCreated = append(s.HostsWithTablesCreated, from.HostsWithTablesCreated...)
				}
				s.UsedTemplates = nil
				if len(from.UsedTemplates) > 0 {
					s.UsedTemplates = append(s.UsedTemplates, from.UsedTemplates...)
				}
			}

			if opts.Errors {
				s.Error = from.Error
				s.Errors = util.MergeStringArrays(s.Errors, from.Errors)
				sort.Sort(sort.Reverse(sort.StringSlice(s.Errors)))
			}

			if opts.MainFields {
				s.CHOpVersion = from.CHOpVersion
				s.CHOpCommit = from.CHOpCommit
				s.CHOpDate = from.CHOpDate
				s.CHOpIP = from.CHOpIP
				s.ClustersCount = from.ClustersCount
				s.ShardsCount = from.ShardsCount
				s.ReplicasCount = from.ReplicasCount
				s.HostsCount = from.HostsCount
				s.Status = from.Status
				s.TaskID = from.TaskID
				s.TaskIDsStarted = from.TaskIDsStarted
				s.TaskIDsCompleted = from.TaskIDsCompleted
				s.Action = from.Action
				mergeActionsNoSync(s, from)
				s.Error = from.Error
				s.Errors = from.Errors
				s.HostsUpdatedCount = from.HostsUpdatedCount
				s.HostsAddedCount = from.HostsAddedCount
				s.HostsUnchangedCount = from.HostsUnchangedCount
				s.HostsCompletedCount = from.HostsCompletedCount
				s.HostsDeletedCount = from.HostsDeletedCount
				s.HostsDeleteCount = from.HostsDeleteCount
				s.Pods = from.Pods
				s.PodIPs = from.PodIPs
				s.FQDNs = from.FQDNs
				s.Endpoint = from.Endpoint
				s.NormalizedCR = from.NormalizedCR
			}

			if opts.Normalized {
				s.NormalizedCR = from.NormalizedCR
			}

			if opts.WholeStatus {
				s.CHOpVersion = from.CHOpVersion
				s.CHOpCommit = from.CHOpCommit
				s.CHOpDate = from.CHOpDate
				s.CHOpIP = from.CHOpIP
				s.ClustersCount = from.ClustersCount
				s.ShardsCount = from.ShardsCount
				s.ReplicasCount = from.ReplicasCount
				s.HostsCount = from.HostsCount
				s.Status = from.Status
				s.TaskID = from.TaskID
				s.TaskIDsStarted = from.TaskIDsStarted
				s.TaskIDsCompleted = from.TaskIDsCompleted
				s.Action = from.Action
				mergeActionsNoSync(s, from)
				s.Error = from.Error
				s.Errors = from.Errors
				s.HostsUpdatedCount = from.HostsUpdatedCount
				s.HostsAddedCount = from.HostsAddedCount
				s.HostsUnchangedCount = from.HostsUnchangedCount
				s.HostsCompletedCount = from.HostsCompletedCount
				s.HostsDeletedCount = from.HostsDeletedCount
				s.HostsDeleteCount = from.HostsDeleteCount
				s.Pods = from.Pods
				s.PodIPs = from.PodIPs
				s.FQDNs = from.FQDNs
				s.Endpoint = from.Endpoint
				s.NormalizedCR = from.NormalizedCR
				s.NormalizedCRCompleted = from.NormalizedCRCompleted
			}
		})
	})
}

// ClearNormalizedCR clears normalized CR in status
func (s *Status) ClearNormalizedCR() {
	doWithWriteLock(s, func(s *Status) {
		s.NormalizedCR = nil
	})
}

// SetNormalizedCompletedFromCurrentNormalized sets completed CR from current CR
func (s *Status) SetNormalizedCompletedFromCurrentNormalized() {
	doWithWriteLock(s, func(s *Status) {
		s.NormalizedCRCompleted = s.NormalizedCR
	})
}

// GetCHOpVersion gets operator version
func (s *Status) GetCHOpVersion() string {
	return getStringWithReadLock(s, func(s *Status) string {
		return s.CHOpVersion
	})
}

// GetCHOpCommit gets operator build commit
func (s *Status) GetCHOpCommit() string {
	return getStringWithReadLock(s, func(s *Status) string {
		return s.CHOpCommit
	})
}

// GetCHOpDate gets operator build date
func (s *Status) GetCHOpDate() string {
	return getStringWithReadLock(s, func(s *Status) string {
		return s.CHOpDate
	})
}

// GetCHOpIP gets operator pod's IP
func (s *Status) GetCHOpIP() string {
	return getStringWithReadLock(s, func(s *Status) string {
		return s.CHOpIP
	})
}

// GetClustersCount gets clusters count
func (s *Status) GetClustersCount() int {
	return getIntWithReadLock(s, func(s *Status) int {
		return s.ClustersCount
	})
}

// GetShardsCount gets shards count
func (s *Status) GetShardsCount() int {
	return getIntWithReadLock(s, func(s *Status) int {
		return s.ShardsCount
	})
}

// GetReplicasCount gets replicas count
func (s *Status) GetReplicasCount() int {
	return getIntWithReadLock(s, func(s *Status) int {
		return s.ReplicasCount
	})
}

// GetHostsCount gets hosts count
func (s *Status) GetHostsCount() int {
	return getIntWithReadLock(s, func(s *Status) int {
		return s.HostsCount
	})
}

// GetStatus gets status
func (s *Status) GetStatus() string {
	return getStringWithReadLock(s, func(s *Status) string {
		return s.Status
	})
}

// GetTaskID gets task ipd
func (s *Status) GetTaskID() string {
	return getStringWithReadLock(s, func(s *Status) string {
		return s.TaskID
	})
}

// GetTaskIDsStarted gets started task id
func (s *Status) GetTaskIDsStarted() []string {
	return getStringArrWithReadLock(s, func(s *Status) []string {
		return s.TaskIDsStarted
	})
}

// GetTaskIDsCompleted gets completed task id
func (s *Status) GetTaskIDsCompleted() []string {
	return getStringArrWithReadLock(s, func(s *Status) []string {
		return s.TaskIDsCompleted
	})
}

// GetAction gets last action
func (s *Status) GetAction() string {
	return getStringWithReadLock(s, func(s *Status) string {
		return s.Action
	})
}

// GetActions gets all actions
func (s *Status) GetActions() []string {
	return getStringArrWithReadLock(s, func(s *Status) []string {
		return s.Actions
	})
}

// GetError gets last error
func (s *Status) GetError() string {
	return getStringWithReadLock(s, func(s *Status) string {
		return s.Error
	})
}

// GetErrors gets all errors
func (s *Status) GetErrors() []string {
	return getStringArrWithReadLock(s, func(s *Status) []string {
		return s.Errors
	})
}

// GetHostsUpdatedCount gets updated hosts counter
func (s *Status) GetHostsUpdatedCount() int {
	return getIntWithReadLock(s, func(s *Status) int {
		return s.HostsUpdatedCount
	})
}

// GetHostsAddedCount gets added hosts counter
func (s *Status) GetHostsAddedCount() int {
	return getIntWithReadLock(s, func(s *Status) int {
		return s.HostsAddedCount
	})
}

// GetHostsUnchangedCount gets unchanged hosts counter
func (s *Status) GetHostsUnchangedCount() int {
	return getIntWithReadLock(s, func(s *Status) int {
		return s.HostsUnchangedCount
	})
}

// GetHostsFailedCount gets failed hosts counter
func (s *Status) GetHostsFailedCount() int {
	return getIntWithReadLock(s, func(s *Status) int {
		return s.HostsFailedCount
	})
}

// GetHostsCompletedCount gets completed hosts counter
func (s *Status) GetHostsCompletedCount() int {
	return getIntWithReadLock(s, func(s *Status) int {
		return s.HostsCompletedCount
	})
}

// GetHostsDeletedCount gets deleted hosts counter
func (s *Status) GetHostsDeletedCount() int {
	return getIntWithReadLock(s, func(s *Status) int {
		return s.HostsDeletedCount
	})
}

// GetHostsDeleteCount gets hosts to be deleted counter
func (s *Status) GetHostsDeleteCount() int {
	return getIntWithReadLock(s, func(s *Status) int {
		return s.HostsDeleteCount
	})
}

// GetPods gets list of pods
func (s *Status) GetPods() []string {
	return getStringArrWithReadLock(s, func(s *Status) []string {
		return s.Pods
	})
}

// GetPodIPs gets list of pod ips
func (s *Status) GetPodIPs() []string {
	return getStringArrWithReadLock(s, func(s *Status) []string {
		return s.PodIPs
	})
}

// GetFQDNs gets list of all FQDNs of hosts
func (s *Status) GetFQDNs() []string {
	return getStringArrWithReadLock(s, func(s *Status) []string {
		return s.FQDNs
	})
}

// GetEndpoint gets API endpoint
func (s *Status) GetEndpoint() string {
	return getStringWithReadLock(s, func(s *Status) string {
		return s.Endpoint
	})
}

// GetNormalizedCR gets target CR
func (s *Status) GetNormalizedCR() *ClickHouseInstallation {
	return getInstallationWithReadLock(s, func(s *Status) *ClickHouseInstallation {
		return s.NormalizedCR
	})
}

// GetNormalizedCRCompleted gets completed CR
func (s *Status) GetNormalizedCRCompleted() *ClickHouseInstallation {
	return getInstallationWithReadLock(s, func(s *Status) *ClickHouseInstallation {
		return s.NormalizedCRCompleted
	})
}

// GetHostsWithTablesCreated gets hosts with created tables
func (s *Status) GetHostsWithTablesCreated() []string {
	return getStringArrWithReadLock(s, func(s *Status) []string {
		return s.HostsWithTablesCreated
	})
}

// Begin helpers

func doWithWriteLock(s *Status, f func(s *Status)) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	f(s)
}

func doWithReadLock(s *Status, f func(s *Status)) {
	if s == nil {
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	f(s)
}

func getIntWithReadLock(s *Status, f func(s *Status) int) int {
	var zeroVal int
	if s == nil {
		return zeroVal
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	return f(s)
}

func getStringWithReadLock(s *Status, f func(s *Status) string) string {
	var zeroVal string
	if s == nil {
		return zeroVal
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	return f(s)
}

func getInstallationWithReadLock(s *Status, f func(s *Status) *ClickHouseInstallation) *ClickHouseInstallation {
	var zeroVal *ClickHouseInstallation
	if s == nil {
		return zeroVal
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	return f(s)
}

func getStringArrWithReadLock(s *Status, f func(s *Status) []string) []string {
	emptyArr := make([]string, 0, 0)
	if s == nil {
		return emptyArr
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	return f(s)
}

// mergeActionsNoSync merges the actions of from into those of s (without synchronization, because synchronized
// functions call into this).
func mergeActionsNoSync(s *Status, from *Status) {
	s.Actions = util.MergeStringArrays(s.Actions, from.Actions)
	sort.Sort(sort.Reverse(sort.StringSlice(s.Actions)))
	trimActionsNoSync(s)
}

// trimActionsNoSync trims actions (without synchronization, because synchronized functions call into this).
func trimActionsNoSync(s *Status) {
	if len(s.Actions) > maxActions {
		s.Actions = s.Actions[:maxActions]
	}
}

// pushTaskIDStartedNoSync pushes task id into status
func pushTaskIDStartedNoSync(s *Status) {
	s.TaskIDsStarted = append([]string{s.TaskID}, s.TaskIDsStarted...)
	if len(s.TaskIDsStarted) > maxTaskIDs {
		s.TaskIDsStarted = s.TaskIDsStarted[:maxTaskIDs]
	}
}

// pushTaskIDCompletedNoSync pushes task id into status w/o sync
func pushTaskIDCompletedNoSync(s *Status) {
	s.TaskIDsCompleted = append([]string{s.TaskID}, s.TaskIDsCompleted...)
	if len(s.TaskIDsCompleted) > maxTaskIDs {
		s.TaskIDsCompleted = s.TaskIDsCompleted[:maxTaskIDs]
	}
}
