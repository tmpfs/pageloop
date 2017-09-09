package pageloop

var(
  Jobs *JobManager
  JobCount uint64 = 0
)

// Job is a potentially long running background task
// such as executing an external command in a goroutine.
type Job struct {
  Name string
  id uint64
  running bool
}

type JobManager struct {
  Jobs []Job
}

func (j *JobManager) NewJob(name string) *Job {
  job := &Job{Name: name}
  JobCount++
  job.id = JobCount
  return job
}

func (j *JobManager) StartJob(job *Job) {
  job.running = true
}

func (j *JobManager) StopJob(job *Job) {
  job.running = false
}

func init() {
  Jobs = &JobManager{}
}
