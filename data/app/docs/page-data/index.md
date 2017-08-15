---
title: Page Data
lang: en
description: Describes how page data is parsed.
keywords: web, editor, page, data
leader: |
  Page data is arbitrary data associated with a page.
---

<p>{{relative "foo.css"}}</p>

Use page data to pass page titles or other information to
your templates, you can use the YAML and JSON formats.

When looking for page data files the source file extension is
replaced so given a file named `index.html` the corresponding
page data file will be `index.yml` or `index.json`.

#### Frontmatter

The easiest method for creating page data is to embed YAML
frontmatter in your HTML or Markdown source files.

When using the frontmatter technique the YAML document must
be delimited by `---`, for example:

```yaml
---
title: Page Title
lang: en
---
```

#### YAML File

If no frontmatter data is detected a standalone `.yml` file with the
same name as the page is loaded if it exists and the page data is
extracted from the parsed file. In this case you should not use the
`---` YAML document delimiters.

#### JSON File

If no YAML data is available and a `.json` file with the same name as
the page exists it is parsed and assigned to the page data.
