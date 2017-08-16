---
title: API Documentation
lang: en
description: Describes the REST API endpoints.
keywords: web, editor, documentation, api, rest
leader: |
  Here is all the information you need to interact with the REST API.
---

To make requests use your favourite command line program or the <a href="/tools/api/probe/" title="API Probe">probe</a> tool. You can also <a href="/tools/api/browser/" title="API Browser">browse</a> the GET requests.

All requests should have an <code>application/json</code> content type header. URLs are shown relative to the base <code>/api/</code> path.

Containers and applications are referenced by name which may only contain alphanumeric digits and the hyphen, they may not begin with a hyphen.

## GET /

<p>Get the list of containers.</p>
<p>A container is a collection of web applications.</p>

## GET /{container}/

<p>Get the list of applications in a container.</p>
<p>An application contains a set of source files that are published to a public location.</p>

## PUT /{container}/

<p>Create a new application.</p>
<p>To create a new application you must give it at least a valid name so the request body should be valid JSON such as:</p>
<pre>{"name": "application-name"}</pre>

## GET /{container}/{application}/

<p>Get an application.</p>

## DELETE /{container}/{application}/

<p>Remove an application.</p>

## GET /{container}/{application}/files/

<p>Get the list of files for an application.</p>

## GET /{container}/{application}/files/{url}

<p>Get file information for the file URL.</p>

## GET /{container}/{application}/pages/

<p>Get the list of pages for an application.</p>

## GET /{container}/{application}/pages/{url}

<p>Get page information for the page URL.</p>
