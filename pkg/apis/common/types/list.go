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

package types

import (
	apps "k8s.io/api/apps/v1"
	"time"
)

type List map[string]string

func (l List) Get(name string) string {
	return l[name]
}

func markPodRestartedNow(sts *apps.StatefulSet) {
	v, _ := time.Now().UTC().MarshalText()
	// Instantiate new annotations if they don't exist yet, otherwise append the new annotation
	if sts.Spec.Template.Annotations == nil {
		sts.Spec.Template.Annotations = make(map[string]string)
	}
	sts.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = string(v)
}

