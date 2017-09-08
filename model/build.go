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

const BuildFileName = "build.yml"

var(
  defaultTask = "publish"
)

type BuildFile struct {
  // Allows disabling builds when the server boots
  Boot bool `json:"-" yaml:"boot"`

  // List of build tasks
  Tasks Tasks `json:"tasks" yaml:"tasks"`

  // Main build command to run
  Command string `json:"-" yaml:"publish"`
}

type Tasks map[string]string

type TaskComplete interface {
  Done(err error, cmd *exec.Cmd, raw string)
}

type DefaultTaskComplete struct {}

func (d *DefaultTaskComplete) Done(err error, cmd *exec.Cmd, raw string) {
  if err != nil {
    os.Stderr.WriteString(err.Error() + "\n")
  } else {
    fmt.Printf("[task] %s (%s)\n", raw, cmd.Dir)
  }
}

func ReadBuildFile (app *Application) (*BuildFile, error) {
  var err error
  var file *BuildFile = &BuildFile{}
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

    if file.Command == "" && file.Tasks[defaultTask] != "" {
      file.Command = file.Tasks[defaultTask]
    } else if file.Command != "" {
      file.Tasks[defaultTask] = file.Command
    }

    if file.Command == "" {
      return nil, fmt.Errorf("Build file %s does not contain a publish command", input)
    }
    return file, nil
  }

  return nil, nil
}

// Run an arbitrary command in a goroutine and invoke the
// done callback on completion.
func (b *BuildFile) Run(app *Application, raw string, done TaskComplete) {
  var cwd = app.SourceDirectory()
  var parts []string = strings.Split(raw, " ")
  var name string = parts[0]
  var args []string = parts[1:]
  var cmd *exec.Cmd = exec.Command(name, args...)
  cmd.Dir = cwd

  run := func(c chan error) {
    if err := cmd.Run(); err != nil {
      c <- err
    }
    c <- nil
  }

  c := make(chan error)

  listen := func(c chan error) {
    e := <- c
    done.Done(e, cmd, raw)
  }

  go listen(c)
  go run(c)
}

// Run the main build task asynchonously.
func (b *BuildFile) Build(app *Application, done TaskComplete) {
  var command string = b.Command
  b.Run(app, command, done)
}
