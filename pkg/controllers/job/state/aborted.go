/*
Copyright 2017 The Volcano Authors.

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

package state

import (
	vkv1 "hpw.cloud/volcano/pkg/apis/batch/v1alpha1"
)

type abortedState struct {
	job *vkv1.Job
}

func (as *abortedState) Execute(action vkv1.Action, reason string, msg string) (error) {
	switch action {
	case vkv1.ResumeJobAction:
		return SyncJob(as.job, func(status vkv1.JobStatus) vkv1.JobState {
			return vkv1.JobState{
				Phase:   vkv1.Restarting,
				Reason:  reason,
				Message: msg,
			}
		})
	default:
		return KillJob(as.job, nil)
	}
}
