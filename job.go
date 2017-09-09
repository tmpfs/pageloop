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

func (j *Job) Id() uint64 {
  return j.id
}

func (j *Job) Running() bool {
  return j.running
}

type JobManager struct {
  Jobs []*Job
}

func (j *JobManager) NewJob(name string) *Job {
  job := &Job{Name: name}
  JobCount++
  job.id = JobCount
  j.Jobs = append(j.Jobs, job)
  return job
}

func (j *JobManager) GetRunningJob(name string) *Job {
  for _, job := range j.Jobs {
    if job.Name == name && job.Running() {
      return job
    }
  }
  return nil
}

func (j *JobManager) Start(job *Job) {
  job.running = true
}

func (j *JobManager) Stop(job *Job) {
  job.running = false
}

func init() {
  Jobs = &JobManager{}
}
