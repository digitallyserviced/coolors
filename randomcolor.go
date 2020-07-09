package main

import (
	"github.com/gdamore/tcell"
	"math/rand"
)

func randRange(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func randomColor() tcell.Color {
	r := int32(randRange(0, 255))
	g := int32(randRange(0, 255))
	b := int32(randRange(0, 255))
	return tcell.NewRGBColor(r, g, b)
}
