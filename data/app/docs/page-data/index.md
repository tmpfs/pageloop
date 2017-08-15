---
title: Page Data
lang: en
description: Describes how page data is parsed.
keywords: web, editor, page, data
leader: Page data is arbitrary data associated with a page.
---

Use page data to pass page titles or other information to
your templates, you can use the YAML and JSON formats.

#### Frontmatter

The easiest method for creating page data is to embed YAML
frontmatter in your HTML or Markdown source files.

When using the frontmatter technique the YAML document must
be delimited by `---`:

```yaml
---
title: Page Title
lang: en
---
```
