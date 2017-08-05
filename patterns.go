package blocks

import (
  "regexp"
)

var TEMPLATE_FILE = regexp.MustCompile(`\.html?$`)
var INDEX_FILE = regexp.MustCompile(`index\.html?$`)
var FRONTMATTER = regexp.MustCompile(`^---\n`)
var FRONTMATTER_END = regexp.MustCompile(`---$`)

