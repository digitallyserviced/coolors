package util

import "golang.org/x/exp/constraints"




func Clamped(val, min, max float64) (float64, bool) {
	clampd := val > max
	
	if val < min {
		clampd = true
	}
	return Clamp(val, min, max), clampd
}
func Max[T constraints.Ordered](a,b T) T {
  if a > b {
    return a
  }
  return b
}
func Min[T constraints.Ordered](a,b T) T {
  if a < b {
    return a
  }
  return b
}
func Clamp[T constraints.Ordered](val, min, max T) T {
	return Max(min, Min(val, max))
}

