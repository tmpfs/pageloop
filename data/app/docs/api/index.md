---
title: API Documentation
lang: en
description: Describes the REST API endpoints.
keywords: web, editor, documentation, api, rest
leader: |
  Here is all the information you need to interact with the REST API.
---

<div class="api">

To make requests use your favourite command line program or the [probe](/tools/api/probe/ "API Probe") tool. You can also [browse](/tools/api/browser/ "API Browser") the GET requests.

All requests should have an `application/json` content type header. URLs are shown relative to the base `/api/` path.

Containers and applications are referenced by name which may only contain alphanumeric digits and the hyphen, they may not begin with a hyphen.

## GET /

Get the list of containers.

A container is a collection of web applications.

## GET /{container}/

Get the list of applications in a container.

An application contains a set of source files that are published to a public location.

## PUT /{container}/

Create a new application.

To create a new application you must give it a valid name, 
mountpoint URL and description so a request body would look like:

```json
{
  "name": "new-app",
  "url": "/new-app/",
  "description": "New application"
}
```

You may optionally specify a `template` object to initialize the new 
application with all the files in the referenced template application.

You need to give the `container` and `application` names for the template:

```json
{
  "name": "new-app",
  "url": "/new-app/",
  "description": "New application",
  "template": {
    "container": "template",
    "application": "pure"
  }
}
```

## GET /{container}/{application}/

Get an application.

## DELETE /{container}/{application}/

Remove an application.

Removing an application indicates a complete deletion, it is 
irreversible.

The application is unmounted and the application mountpoint is 
deleted before the updated configuration is written to disc. Finally 
the published and source files are deleted.

## GET /{container}/{application}/files/

Get the list of files for an application.

## PUT /{container}/{application}/files/{url}

Create a file for an application, if the file already exists an
error is returned. Note that because file extensions can be changed
when files are published conflicts are also detected on the published
file name. For example, if you have an existing file
named `document.html` and try to create a file named `document.md`
it is an error as the published URLs would conflict.

If the file is considered to be a page it is also added to the list
of pages for the application.

Syncs the source file to disc and publishes an updated
version of the file to the public URL.

You can create a file using the content from an existing file template,
to do so you should send a `Content-Type` header using one of the following
MIME types:

* `template/markdown+partial`
* `template/markdown+standalone`
* `template/html+standalone`
* `template/html+layout`
* `template/html+partial`

When creating new files from templates there is no need to send a
request body.

## POST /{container}/{application}/files/{url}

Update a file, the file must already exist.

If a `Location` header is sent in the request the operation is a rename 
(the request body is ignored) and the target file is renamed to the URL 
given in the `Location` header.

Otherwise the operation is to update file content from the request body, 
in which case it is an error if the request MIME type does not match the
existing MIME type for the file.

Syncs the source file to disc and publishes an updated version of the 
file to the public URL.

## GET /{container}/{application}/files/{url}

Get file information for the file URL.

## GET /{container}/{application}/pages/

Get the list of pages for an application.

## GET /{container}/{application}/pages/{url}

Get page information for the page URL.

</div>