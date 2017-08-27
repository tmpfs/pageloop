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
+ `-p, --publish=[dir] {=public}` Directory for application builds
+ `-h, --help` Display help and exit
+ `--version` Print the version and exit

# Configuration

Use a YAML configuration file to control the service behaviour.

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
