package util


func GetI[T any](i any) (b T) {
  b, ok := i.(T)
  if !ok {
    // panic(ok)
  }
  return
}
