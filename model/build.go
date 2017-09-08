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

type BuildFile struct {
  // Allows disabling builds when the server boots
  Boot bool
  // Main build command to run
  Command string `yaml:"build"`
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

    if file.Command == "" {
      return nil, fmt.Errorf("Build file %s does not contain a build command", input)
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
