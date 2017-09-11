package adapter

import (
  "net/http"
  "reflect"
)

var(
  Routes []*ActionMap
)

// Initialize the action list with route actions and action definitions.
func init() {
  push := func(action *Action, def *CommandDefinition) {
    Routes = append(Routes, &ActionMap{Action: action, CommandDefinition: def})
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
  push(NewAction(OperationRead, ""),
    &CommandDefinition{MethodName: "ReadMeta", Status: http.StatusOK})
  // GET /templates
  push(NewAction(OperationRead, "/templates"),
    &CommandDefinition{MethodName: "ReadApplicationTemplates", Status: http.StatusOK})
  // GET /jobs
  push(NewAction(OperationRead, "/jobs"),
    &CommandDefinition{MethodName: "ReadActiveJobs", Status: http.StatusOK})
  // GET /jobs/{id}
  push(NewAction(OperationRead, "/jobs/*"),
    &CommandDefinition{MethodName: "ReadJob", Status: http.StatusOK, Arguments: contextArg})
  // DELETE /jobs/{id}
  push(NewAction(OperationDelete, "/jobs/*"),
    &CommandDefinition{MethodName: "DeleteJob", Status: http.StatusOK, Arguments: contextArg})
  // GET /apps
  push(NewAction(OperationRead, "/apps"),
    &CommandDefinition{MethodName: "ReadHost", Status: http.StatusOK})
  // GET /apps/{container}
  push(NewAction(OperationRead, "/apps/*"),
    &CommandDefinition{MethodName: "ReadContainer", Status: http.StatusOK, Arguments: containerArg})
  // PUT /apps/{container}
  push(NewAction(OperationCreate, "/apps/*"),
    &CommandDefinition{MethodName: "CreateApp", Status: http.StatusCreated, Arguments: containerArg})
  // GET /apps/{container}/{application}
  push(NewAction(OperationRead, "/apps/*/*"),
    &CommandDefinition{MethodName: "ReadApplication", Status: http.StatusOK, Arguments: applicationArg, Index: 1})
  // DELETE /apps/{container}/{application}
  push(NewAction(OperationDelete, "/apps/*/*"),
    &CommandDefinition{MethodName: "DeleteApp", Status: http.StatusOK, Arguments: applicationArg})
  // GET /apps/{container}/{application}/files
  push(NewAction(OperationRead, "/apps/*/*/files"),
    &CommandDefinition{MethodName: "ReadApplicationFiles", Status: http.StatusOK, Arguments: applicationArg})
  // GET /apps/{container}/{application}/pages
  push(NewAction(OperationRead, "/apps/*/*/pages"),
    &CommandDefinition{MethodName: "ReadApplicationPages", Status: http.StatusOK, Arguments: applicationArg})
  // GET /apps/{container}/{application}/files/{url}
  push(NewAction(OperationRead, "/apps/*/*/files/*"),
    &CommandDefinition{MethodName: "ReadFile", Status: http.StatusOK, Arguments: fileArg})
  // GET /apps/{container}/{application}/pages/{url}
  push(NewAction(OperationRead, "/apps/*/*/pages/*"),
    &CommandDefinition{MethodName: "ReadPage", Status: http.StatusOK, Arguments: fileArg})
  // DELETE /apps/{container}/{application}/files/{url}
  push(NewAction(OperationDelete, "/apps/*/*/files/*"),
    &CommandDefinition{MethodName: "DeleteFile", Status: http.StatusOK, Arguments: fileArg})
  // PUT /apps/{container}/{application}/tasks/{name}
  push(NewAction(OperationCreate, "/apps/*/*/tasks/*"),
    &CommandDefinition{MethodName: "RunTask", Status: http.StatusAccepted, Arguments: fileArg})
}
