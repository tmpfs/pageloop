package blocks

import (
  . "os"
)

type Application struct {
  Title string `json:"title"`
  Pages []Page `json:"pages"`
  Files []File `json:"files"`
}
