package util

import(
  "fmt"
  "strconv"
)

var(
  // TODO: KiB, MiB etc
  Extensions map[string] int64
  Kb int64 = 1000
  Bytes string = "bytes"
)

func PrettyBytes(num int64) string {
  for k, v := range Extensions {
    if num >= v {
      if v == Kb {
        // Don't do decimal precision for KB values
        return strconv.FormatFloat(float64(num) / float64(v), 'f', 0, 64) + k
      }
      return strconv.FormatFloat(float64(num) / float64(v), 'f', 2, 64) + k
    }
  }

  return fmt.Sprintf("%d bytes", num)
}

func init() {
  Extensions = make(map[string]int64)
  Extensions["PB"] = Kb * Kb * Kb * Kb * Kb
  Extensions["TB"] = Kb * Kb * Kb * Kb
  Extensions["GB"] = Kb * Kb * Kb
  Extensions["MB"] = Kb * Kb
  Extensions["KB"] = Kb
}
