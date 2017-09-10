package util

import(
  "time"
)

var(
  // Singleton job manager.
  Jobs *JobManager

  // Incremental job identifier.
  JobCount uint64 = 0
)

// Contract for types that require notification when a job
// is completed.
type JobComplete interface {
  Done(err error, j *Job)
}

// Job runners start a job running and must return a job.
//
// They should assign themselves as the runner for the job so 
// that complete listeners can inspect the runner that started the job.
type JobRunner interface {
  Run(done JobComplete) (*Job, error)
}

// Job is a potentially long running background task
// such as executing an external command in a goroutine.
type Job struct {
  Name string `json:"name"`
  Runner JobRunner `json:"run"`
  id uint64 `json:"id"`
  running bool `json:"active"`
  start time.Time
  duration time.Duration
}

// Access the job number.
func (j *Job) Id() uint64 {
  return j.id
}

// Determine if the job is active.
func (j *Job) Running() bool {
  return j.running
}

// Job manager creates, starts and stops jobs and maintains
// a list of active jobs.
type JobManager struct {
  // List of active jobs
  Jobs []*Job
}

// Create a new job.
func (j *JobManager) NewJob(name string, runner JobRunner) *Job {
  job := &Job{Name: name, Runner: runner}
  JobCount++
  job.id = JobCount
  return job
}

// Find a job by name that is currently running.
func (j *JobManager) GetRunningJob(name string) *Job {
  for _, job := range j.Jobs {
    if job.Name == name && job.Running() {
      return job
    }
  }
  return nil
}

// Start a job.
func (j *JobManager) Start(job *Job) {
  job.running = true
  job.start = time.Now()
  j.Jobs = append(j.Jobs, job)
}

// Stop a job. The job is removed from the list
// of active jobs.
func (j *JobManager) Stop(job *Job) {
  job.running = false
  job.duration = time.Since(job.start)
  for i, cj := range j.Jobs {
    if job == cj {
      before := j.Jobs[0:i]
      after := j.Jobs[i+1:]
      j.Jobs = append(before, after...)
    }
  }
}

// Create singleton job manager.
func init() {
  Jobs = &JobManager{}
}
