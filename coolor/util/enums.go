package util

import "golang.org/x/exp/constraints"

func BitAnd[T constraints.Integer](a, b T) bool {
  return a & b != 0
}
