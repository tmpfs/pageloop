---
title: API Documentation
lang: en
description: Describes the REST API endpoints.
keywords: web, editor, documentation, api, rest
leader: |
  Here is all the information you need to interact with the REST API.
---

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

To create a new application you must give it at least a valid name so the request body should be valid JSON such as:

```
{"name": "application-name"}
```

## GET /{container}/{application}/

Get an application.

## DELETE /{container}/{application}/

Remove an application.

## GET /{container}/{application}/files/

Get the list of files for an application.

## PUT /{container}/{application}/files/{url}

Create a file for an application, if the file already exists an 
error is returned.

Syncs the source file to disc and publishes an updated 
version of the file to the public URL.

## POST /{container}/{application}/files/{url}

Update file content for a file, the file must already exist.

Syncs the source file to disc and publishes an updated 
version of the file to the public URL.

## GET /{container}/{application}/files/{url}

Get file information for the file URL.

## GET /{container}/{application}/pages/

Get the list of pages for an application.

## GET /{container}/{application}/pages/{url}

Get page information for the page URL.
