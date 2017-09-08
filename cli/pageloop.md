# Name

pageloop - collaborative realtime server

# Synopsis

```
[flags] [options]
```

# Description

Web server that serves system applications and
those specified in the supplied configuration file.

# Options

+ `-a, --addr=[val] {=:3577}` Set the bind address
+ `-c, --config=[file]` Load server configuration from a YAML file
+ `-h, --help` Display help and exit
+ `--version` Print the version and exit

# Configuration

Use a YAML configuration file to control the service behaviour.
The configuration file should be writable so that new applications
created using the user interface can be persisted.

The `addr` configuration field sets the server bind address, it may
be a port in the form :8080 or fully qualified host or IP address such as
0.0.0.0:8080.

When applications are mounted they are published to a directory which
is configured using the `publish` field, the directory must be writable.
If no publish directory is configured the default location is a `public` directory
relative to the current working directory.

Define a `mountpoints` list in the configuration file to specify applications
to load when the server starts. Each entry can contain the fields:

+ `url` Public URL mountpoint
+ `path` Path to the source files
+ `description` A short description of the application

Note that applications mounted from a user configuration file are appended
to the list of system mountpoints, you cannot control system applications.

# Publish

When applications are published they are written to the publish directory
with a namespace which is the container name and application name. Such that
system/editor is the editor application in the system container.

Container and application names must be unique. For applications the name is
derived from the basename of the path and it is an error if two applications
in the same container have the same name.
