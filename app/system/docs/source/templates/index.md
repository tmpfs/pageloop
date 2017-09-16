---
title: Templates
lang: en
description: Template processing.
keywords: web, editor, template, documentation
leader: |
  Templates allow you to write page data to the output files.
template:
  delims:
    left: <?
    right: ?>
---

HTML files are parsed as templates and support all of the features of
the `html/template` package. When a template is executed the context
(`.`) is set to the the page data map, see [page data](/docs/page-data/).

## Layouts

A simple layout mechanism allows you to share common elements in a
layout file which must be named `layout.html`. If a layout file is
found in the current directory it is used, otherwise all parent
directories until the root directory for the application are searched.

You can use symbolic links if you want to share layouts between
applications.

A layout file defines the page structure and loads the `content`
template by convention. So a minimal layout would look like:

```html
<?`<!doctype html>
<html lang="{{.lang}}">
  <head>
    <title>{{.title}}</title>
  </head>
  <body>
    <main>
      {{block "content" .}}{{end}}
    </main>
  </body>
</html>` | html ?>
```

Your page file should not declare the `content` template it is
declared automatically. An example page for this layout:

```html
<?`---
title: Page Title
lang: en
---
<p>Page content.</p>` | html ?>
```

If you need to disable layouts for a page you can set `layout`
in the page data:

```html
<?`---
title: Standalone Page
lang: en
layout: false
---
<html><!-- ... --></html>` | html ?>
```

And the page will be treated as standalone so you should supply
a full HTML document.

## Helper Functions

Some useful functions are exposed to the templates so that you can
create applications independent of the application mountpoint (URL)
and perform other actions such as embed markdown in your HTML templates.

### markdown

Use the markdown helper function to embed markdown in your HTML documents:

```html
<? `
<div>
{{markdown "This is some *markdown* embedded in HTML."}}
</div>
` | html ?>
```

For long sections of markdown use the raw string type:

```html
{{markdown `
  This is some *markdown* embedded in HTML with lots more
  content.
`}}
```

### root

Returns a  URL relative to the root of the application, for example:

```html
<? `<link rel="stylesheet" href='{{root "app.css"}}' />` | html ?>
```

For an application mounted at `/docs/` would return `/docs/app.css`.

### bytesize

Converts an int64 file size to a pretty string such as 12.2MB.

```html
{{file.size | bytesize}}
```

## Template Configuration

If you need to disable template parsing for a page set the `template`
page data option:

```html
<?`---
template: false
---` | html ?>
```

### Delimiters

Sometimes you may want to switch the template left and right delimiters
usually when you want to show examples using the default delimiters.

To change the delimiters for a page (note it only applies to the page
template not a corresponding layout) you can specify them in the page data.

An example that uses different delimiters for template parsing:

```yaml
---
template:
  delims:
    left: ${
    right: }
---
```
