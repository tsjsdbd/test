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
	vkv1 "volcano.sh/volcano/pkg/apis/batch/v1alpha1"
	"volcano.sh/volcano/pkg/controllers/job/apis"
)

type restartingState struct {
	job *apis.JobInfo
}

func (ps *restartingState) Execute(action vkv1.Action) error {
	return SyncJob(ps.job, func(status *vkv1.JobStatus) {
		phase := vkv1.Restarting

		// Get the maximum number of retries.
		maxRetry := DefaultMaxRetry
		if ps.job.Job.Spec.MaxRetry != 0 {
			maxRetry = ps.job.Job.Spec.MaxRetry
		}

		if status.RetryCount >= maxRetry {
			// Failed is the phase that the job is restarted failed reached the maximum number of retries.
			phase = vkv1.Failed
		} else {
			status.RetryCount++
			if status.Terminating == 0 {
				if status.Running >= ps.job.Job.Spec.MinAvailable {
					phase = vkv1.Running
				} else {
					phase = vkv1.Pending
				}
			}
		}

		status.State.Phase = phase
	})

}
