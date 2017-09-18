package adapter

import (
  "net/http"
  "reflect"
)

var(
  // Maps service names to action and command definitions
  Services map[string]*ActionMap
  // Routes for REST API requests
  Routes []*ActionMap
)

// Initialize the action list with route actions and command definitions.
func init() {

  Services = make(map[string]*ActionMap)

  push := func(method string, action *Action, def *CommandDefinition) {
    m := &ActionMap{Action: action, CommandDefinition: def}
    Services[method] = m
    Routes = append(Routes, m)
  }

  contextArg := func(b *CommandAdapter, action *Action) []reflect.Value {
    var args []reflect.Value
    args = append(args, reflect.ValueOf(action.Context))
    return args
  }

  containerArg := func(b *CommandAdapter, action *Action) []reflect.Value {
    var args []reflect.Value
    args = append(args, reflect.ValueOf(action.Context))
    return args
  }

  applicationArg := func(b *CommandAdapter, action *Action) []reflect.Value {
    var args []reflect.Value
    args = append(args, reflect.ValueOf(action.Context), reflect.ValueOf(action.Target))
    return args
  }

  fileArg := func(b *CommandAdapter, action *Action) []reflect.Value {
    var args []reflect.Value
    args = append(
      args,
      reflect.ValueOf(action.Context),
      reflect.ValueOf(action.Target),
      reflect.ValueOf(action.Item))
    return args
  }

  // GET /
  push("Core.Meta", NewAction(OperationRead, ""),
    &CommandDefinition{MethodName: "ReadMeta", Status: http.StatusOK})
  // GET /stats
  push("Core.Stats", NewAction(OperationRead, "/stats"),
    &CommandDefinition{MethodName: "ReadStats", Status: http.StatusOK})
  // GET /templates
  push("Template.ReadApplications", NewAction(OperationRead, "/templates"),
    &CommandDefinition{MethodName: "ReadApplicationTemplates", Status: http.StatusOK})
  // GET /jobs
  push("Jobs.ReadActiveJobs", NewAction(OperationRead, "/jobs"),
    &CommandDefinition{MethodName: "ReadActiveJobs", Status: http.StatusOK})
  // GET /jobs/{id}
  push("Jobs.ReadJob", NewAction(OperationRead, "/jobs/*"),
    &CommandDefinition{MethodName: "ReadJob", Status: http.StatusOK, Arguments: contextArg})
  // DELETE /jobs/{id}
  push("Jobs.DeleteJob", NewAction(OperationDelete, "/jobs/*"),
    &CommandDefinition{MethodName: "DeleteJob", Status: http.StatusOK, Arguments: contextArg})
  // GET /apps
  push("Container.List", NewAction(OperationRead, "/apps"),
    &CommandDefinition{MethodName: "ReadHost", Status: http.StatusOK})
  // GET /apps/{container}
  push("Container.Read", NewAction(OperationRead, "/apps/*"),
    &CommandDefinition{MethodName: "ReadContainer", Status: http.StatusOK, Arguments: containerArg})
  // PUT /apps/{container}
  push("Container.CreateApp", NewAction(OperationCreate, "/apps/*"),
    &CommandDefinition{MethodName: "CreateApp", Status: http.StatusCreated, Arguments: containerArg})
  // GET /apps/{container}/{application}
  push("Application.Read", NewAction(OperationRead, "/apps/*/*"),
    &CommandDefinition{MethodName: "ReadApplication", Status: http.StatusOK, Arguments: applicationArg, Index: 1})
  // DELETE /apps/{container}/{application}
  push("Application.Delete", NewAction(OperationDelete, "/apps/*/*"),
    &CommandDefinition{MethodName: "DeleteApp", Status: http.StatusOK, Arguments: applicationArg})
  // GET /apps/{container}/{application}/files
  push("Application.ReadFiles", NewAction(OperationRead, "/apps/*/*/files"),
    &CommandDefinition{MethodName: "ReadApplicationFiles", Status: http.StatusOK, Arguments: applicationArg})
  // GET /apps/{container}/{application}/pages
  push("Application.ReadPages", NewAction(OperationRead, "/apps/*/*/pages"),
    &CommandDefinition{MethodName: "ReadApplicationPages", Status: http.StatusOK, Arguments: applicationArg})
  // GET /apps/{container}/{application}/files/{url}
  push("File.Read", NewAction(OperationRead, "/apps/*/*/files/*"),
    &CommandDefinition{MethodName: "ReadFile", Status: http.StatusOK, Arguments: fileArg})
  // PUT /apps/{container}/{application}/files/{url}
  push("File.Create", NewAction(OperationCreate, "/apps/*/*/files/*"),
    &CommandDefinition{MethodName: "CreateFile", Status: http.StatusCreated, Arguments: fileArg})
  // POST /apps/{container}/{application}/files/{url}
  push("File.Update", NewAction(OperationUpdate, "/apps/*/*/files/*"),
    &CommandDefinition{MethodName: "UpdateFile", Status: http.StatusOK, Arguments: fileArg})
  // GET /apps/{container}/{application}/pages/{url}
  push("File.ReadPage", NewAction(OperationRead, "/apps/*/*/pages/*"),
    &CommandDefinition{MethodName: "ReadPage", Status: http.StatusOK, Arguments: fileArg})
  // DELETE /apps/{container}/{application}/files/
  push("Application.DeleteFiles", NewAction(OperationDelete, "/apps/*/*/files"),
    &CommandDefinition{MethodName: "DeleteFiles", Status: http.StatusOK, Arguments: applicationArg})
  // DELETE /apps/{container}/{application}/files/{url}
  push("File.Delete", NewAction(OperationDelete, "/apps/*/*/files/*"),
    &CommandDefinition{MethodName: "DeleteFile", Status: http.StatusOK, Arguments: fileArg})
  // PUT /apps/{container}/{application}/tasks/{name}
  push("Application.RunTask", NewAction(OperationCreate, "/apps/*/*/tasks/*"),
    &CommandDefinition{MethodName: "RunTask", Status: http.StatusAccepted, Arguments: fileArg})

  // Mutation handlers are overloaded requests that are
  // actions that are tested by the caller and mutated
  // they should never match directly.

  // POST /apps/{container}/{application}/files/{url} - UpdateFile mutation
  push("File.Move", NewAction(OperationUpdate, "/apps/*/*/files/*"),
    &CommandDefinition{MethodName: "MoveFile", Status: http.StatusOK, Arguments: fileArg})

  // PUT /apps/{container}/{application}/files/{url} - CreateFile mutation
  push("File.CreateTemplate", NewAction(OperationCreate, "/apps/*/*/files/*"),
    &CommandDefinition{MethodName: "CreateFileTemplate", Status: http.StatusCreated, Arguments: fileArg})
}
