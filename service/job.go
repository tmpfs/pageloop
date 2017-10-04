package service

import(
  "fmt"
  "net/http"
  . "github.com/tmpfs/pageloop/util"
)

type JobRequest struct {
  Id string `json:"id"`
}

type JobService struct {}

// List active jobs.
func (s *JobService) List(argv *VoidArgs, reply *ServiceReply) *StatusError {
  reply.Reply = Jobs.Active
  return nil
}

// Read a job.
func (s *JobService) Read(req *JobRequest, reply *ServiceReply) *StatusError {
  if job, err := LookupJob(req.Id); err != nil {
    return err
  } else {
    reply.Reply = job
  }
  return nil
}

// Abort an active job.
func(s *JobService) Delete(req *JobRequest, reply *ServiceReply) *StatusError {
  if job, err := LookupJob(req.Id); err != nil {
    return err
  } else {
    if err := Jobs.Abort(job); err != nil {
      return CommandError(http.StatusConflict, err.Error())
    }

    reply.Reply = job
    reply.Status = http.StatusAccepted

    // Accepted for processing
    fmt.Printf("[job:%d] aborted %s\n", job.Number, job.Id)
  }
  return nil
}

// Private

func LookupJob(id string) (*Job, *StatusError) {
  var job *Job = Jobs.ActiveJob(id)
  if job == nil {
    return nil, CommandError(http.StatusNotFound, "Job %s not found", id)
  }
  return job, nil
}
