package pageloop

var(
  // Singleton job manager.
  Jobs *JobManager

  // Incremental job identifier.
  JobCount uint64 = 0
)

// Job is a potentially long running background task
// such as executing an external command in a goroutine.
type Job struct {
  Name string `json:"name"`
  id uint64 `json:"id"`
  running bool `json:"active"`
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
func (j *JobManager) NewJob(name string) *Job {
  job := &Job{Name: name}
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
  j.Jobs = append(j.Jobs, job)
}

// Stop a job. The job is removed from the list
// of active jobs.
func (j *JobManager) Stop(job *Job) {
  job.running = false
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
