package model

import(
  "fmt"
  "os"
  "os/exec"
  "strings"
  "io/ioutil"
  "path/filepath"
  "gopkg.in/yaml.v2"
)

// TODO: inject application and remove from function signatures

const(
  BuildFileName = "build.yml"
)

var(
  defaultTask = "publish"
)

type BuildFile struct {
  App *Application `json:"-" yaml:"-"`

  // Allows disabling builds when the server boots
  Boot bool `json:"-" yaml:"boot"`

  // List of build tasks
  Tasks Tasks `json:"tasks" yaml:"tasks"`

  // Main build command to run
  Command string `json:"-" yaml:"publish"`
}

// Formal task declaration.
type Task struct {
  // A namespace for this command, eg: {container}/{application}
  Namespace string
  Key string
  Raw string
  Command string
  Arguments []string
  Cwd string
  Cmd *exec.Cmd
}

func (t *Task) Parse(raw string) {
  // TODO: split on regexp
  var parts []string = strings.Split(raw, " ")
  t.Raw = raw
  t.Command = parts[0]
  if len(parts) > 1 {
    t.Arguments = parts[1:]
  }
}

func (t *Task) Id() string {
  return t.Namespace + ":" + t.Key
}

// Execute an arbitrary command in a goroutine and invoke the
// done callback on completion.
func (t *Task) Run(done TaskComplete) {
  var cmd *exec.Cmd = exec.Command(t.Command, t.Arguments...)
  cmd.Dir = t.Cwd

  //cmd.Stdout = os.Stdout
  //cmd.Stderr = os.Stderr

  t.Cmd = cmd

  run := func(c chan error) {
    if err := cmd.Run(); err != nil {
      c <- err
    }
    c <- nil
  }

  c := make(chan error)

  listen := func(c chan error) {
    e := <- c
    done.Done(e, cmd, t)
  }

  go listen(c)
  go run(c)
}


type Tasks map[string]string

type TaskComplete interface {
  Done(err error, cmd *exec.Cmd, t *Task)
}

type DefaultTaskComplete struct {}

func (d *DefaultTaskComplete) Done(err error, cmd *exec.Cmd, t *Task) {
  if err != nil {
    os.Stderr.WriteString(err.Error() + "\n")
  } else {
    fmt.Printf("[task] %s (%s)\n", t.Raw, cmd.Dir)
  }
}

func ReadBuildFile (app *Application) (*BuildFile, error) {
  var err error
  var file *BuildFile = &BuildFile{App: app}
  var input string = filepath.Join(app.SourceDirectory(), BuildFileName)
  var content []byte
  if content, err = ioutil.ReadFile(input); err != nil {
    if !os.IsNotExist(err) {
      return nil, err
    }
  } else {
    if err = yaml.Unmarshal(content, file); err != nil {
      return nil, err
    }

    if file.Tasks == nil {
      file.Tasks = make(Tasks)
    }

    // Top-level declaration will override a `publish` task
    if file.Command != "" {
      file.Tasks[defaultTask] = file.Command
    }

    if file.Command == "" {
      return nil, fmt.Errorf("Build file %s does not contain a publish command", input)
    }
    return file, nil
  }

  return nil, nil
}

// Get a task command by string key.
func (b *BuildFile) TaskCommand(key string, ns string) (*Task, error) {
  var t *Task = &Task{Key: key, Namespace: ns}
  if raw, ok := b.Tasks[key]; ok {
    t.Parse(raw)
  } else {
    return nil, fmt.Errorf("Task not found %s (%s)", t.Key, t.Id())
  }
  // Set working directory for command execution
  t.Cwd = b.App.SourceDirectory()
  return t, nil
}

// Run a build task.
func (b *BuildFile) Run(key string, done TaskComplete) (*Task, error) {
  var err error
  var t *Task
  if t, err = b.TaskCommand(key, ""); err != nil {
    return nil, err 
  }
  t.Run(done)
  return t, nil
}

// Run the main build task.
func (b *BuildFile) Build(done TaskComplete) (*Task, error) {
  return b.Run(defaultTask, done)
}
