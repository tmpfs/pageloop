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

type Tasks map[string]string

type BuildFile struct {
  // Allows disabling builds when the server boots
  Boot bool `json:"-" yaml:"-"`

  // List of build tasks
  Tasks Tasks `json:"tasks" yaml:"tasks"`

  // Main build command to run
  Command string `json:"-" yaml:"publish"`
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

func (b *BuildFile) Build(app *Application) error {
  // TODO: run in a goroutine
  var cwd = app.SourceDirectory()
  var command string = b.Command
  var parts []string = strings.Split(command, " ")
  var name string = parts[0]
  var args []string = parts[1:]
  var cmd = exec.Command(name, args...)
  cmd.Dir = cwd
  return cmd.Run()
}
