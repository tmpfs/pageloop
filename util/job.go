package util

import(
  "fmt"
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

// Type for job runners that are cancelable.
type JobAbort interface {
  Abort() error
}

// Job is a potentially long running background task
// such as executing an external command in a goroutine.
type Job struct {
  Id string `json:"id"`
  Runner JobRunner `json:"run"`
  Number uint64 `json:"num"`
  Timestamp int64 `json:"timestamp"`
  start time.Time `json:"-"`
  running bool `json:"active"`
  duration time.Duration
}

// Determine if the job is active.
func (j *Job) Running() bool {
  return j.running
}

// Determine if the job can be aborted.
func (j *Job) CanAbort() bool {
  _, ok := j.Runner.(JobAbort)
  return ok
}

// Job manager creates, starts and stops jobs and maintains
// a list of active jobs.
type JobManager struct {
  // List of active jobs
  Active []*Job
}

// Create a new job.
func (j *JobManager) NewJob(id string, runner JobRunner) *Job {
  job := &Job{Id: id, Runner: runner}
  JobCount++
  job.Number = JobCount
  return job
}

// Find a job by id that is currently active.
func (j *JobManager) ActiveJob(id string) *Job {
  for _, job := range j.Active {
    if job.Id == id && job.Running() {
      return job
    }
  }
  return nil
}

// Start a job.
func (j *JobManager) Start(job *Job) {
  job.running = true
  job.start = time.Now()
  job.Timestamp = job.start.Unix()
  j.Active = append(j.Active, job)
}

// Stop a job. The job is removed from the list
// of active jobs.
func (j *JobManager) Stop(job *Job) {
  job.running = false
  job.duration = time.Since(job.start)
  for i, cj := range j.Active {
    if job == cj {
      before := j.Active[0:i]
      after := j.Active[i+1:]
      j.Active = append(before, after...)
    }
  }
}

// Abort an active job.
//
// It is an error if the job is not running of if the job
// cannot be aborted.
func (j *JobManager) Abort(job *Job) error {
  if !job.CanAbort() {
    return fmt.Errorf(
      "Cannot abort job %s (%d), job is not cancelable", job.Id, job.Number)
  }
  if j.ActiveJob(job.Id) == nil {
    return fmt.Errorf(
      "Cannot abort job %s (%d), job not running", job.Id, job.Number)
  }

  defer j.Stop(job)
  if abortable, ok := job.Runner.(JobAbort); ok {
    return abortable.Abort()
  }
  return nil
}

// Create singleton job manager.
func init() {
  Jobs = &JobManager{}
  // Need to initialize the active jobs so that endpoints
  // return the empty array rather than null
  Jobs.Active = make([]*Job, 0)
}
