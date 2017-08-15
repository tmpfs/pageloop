package model

import (
  "regexp"
)

var TEMPLATE_FILE = regexp.MustCompile(`\.html?$`)
var MARKDOWN_FILE = regexp.MustCompile(`\.(md|markdown)?$`)
var FRONTMATTER = regexp.MustCompile(`^---\n`)
var FRONTMATTER_END = regexp.MustCompile(`---$`)

