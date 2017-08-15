---
title: Data Model
lang: en
description: Describes the internal data model.
keywords: web, editor, documentation
leader: This section describes the data model used internally.
---

Conceptually the hierarchy is quite simple:

```
host > container > application > (file|page)
```

A host contains a list of containers which has a collection of applications.
Each application contains a list of files and a list of pages; pages are the
source files that represent published web pages (HTML and Markdown documents).

#### Host

A host is the root of the hierarchy. Currently there is only a single host but
in the future this may change.

#### Container

A container is a group of related applications and is referenced by name.

Containers may be `protected` in which case the container and the applications
within the container may not be deleted.

The program configures some default containers:

* `system`: System applications (protected).
* `template`: Application templates (protected).
* `user`: Userspace applications.
* `sandbox`: Test playground.

#### Application

An application is the collection of source files for the web application,
typically it is loaded from the filesystem.

It has a group of files which are assigned URLs matching their path relative
to the application root. Markdown and HTML files need special treatment so
files ending with the `.html`, `.md` and `.markdown` extensions are added
to a collection of pages and have additional fields including a pointer to
the underlying file.

#### File

A file belongs to an application and is referenced by it's URL relative
to the application root. Files are typically created by loading from a
local filesystem but there is an abstraction that allows loading files
from other sources such as zip archives or remote servers.

#### Page

A page represents a source file that is editable by a user. Like the file
type it also belongs to an application and is referenced by relative URL.
