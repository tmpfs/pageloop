---
title: Template Includes
lang: en
description: Including files in templates.
keywords: web, editor, template, include, documentation
leader: |
  Include files allow you to share templates between files.
template:
  delims:
    left: <?
    right: ?>
---

## Includes

For complex templating requirements you may need to load shared templates, to 
do so you can specify an `includes` array of files in the page meta data.

Imagine you have a table of contents that you need to include in different 
pages, you can define a `toc` block template like this:

```html
---
layout: false
---
{{block "toc" .}}
&lt;ul&gt;
  &lt;!-- list of items --&gt;
&lt;/ul&gt;
{{end}}
```

Then include the file and render the referenced template:

```html
---
includes:
  - toc.html
---
{{block "toc" .}}{{end}}
&lt;!-- rest of the document --&gt;
```

Note that you can use the same technique to include templates in 
markdown source files.