package util

import(
  "strconv"
)

var(
  // TODO: KiB, MiB etc
  Extensions map[string] int64
  Kb int64 = 1000
)

func PrettyBytes(num int64) string {
  for k, v := range Extensions {
    if num >= v {
      // Don't do decimal precision for small values
      if v <= Kb {
        return strconv.FormatFloat(float64(num) / float64(v), 'f', 0, 64) + k
      }
      return strconv.FormatFloat(float64(num) / float64(v), 'f', 2, 64) + k
    }
  }
  return string(num)
}

func init() {
  Extensions = make(map[string]int64)
  Extensions["PB"] = Kb * Kb * Kb * Kb * Kb
  Extensions["TB"] = Kb * Kb * Kb * Kb
  Extensions["GB"] = Kb * Kb * Kb
  Extensions["MB"] = Kb * Kb
  Extensions["KB"] = Kb
}
