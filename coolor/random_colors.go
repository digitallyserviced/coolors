package coolor

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/samber/lo"
)


func RandomNamedAnsiClusterShade(tcol Color, distance float64) Color {
	clusterColor := lo.Sample[*CoolorColorCluster](getNamedAnsiColors())
	return RandomShadeFromCluster(tcol, clusterColor, distance)
}

func RandomAnsiClusterShade(tcol Color, distance float64) Color {
	clusterColor := lo.Sample[*CoolorColorCluster](getBaseAnsiClusterColors())
	return RandomShadeFromCluster(tcol, clusterColor, distance)
}

func RandomShadeFromCluster(tcol Color, cluster *CoolorColorCluster, distance float64) Color {
	col2 := cluster.leadColor
	return RandomShadeFromColors(tcol, col2, distance)
}

func RandomShadeFromColors(tcol, tcol2 Color, distance float64) Color {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	max := (distance * r.Float64()) + 0.01
	return tcol.BlendLuv(tcol2, max)
}

func MakeRandomCoolor() *Coolor {
	col := tcell.NewRGBColor(int32(randRange(0, 255)), int32(randRange(0, 255)), int32(randRange(0, 255)))
  // var a *Coolor = col
  a := &Coolor{
  	Color: col,
  }
	return a
}

func MakeRandomColor() *tcell.Color {
	col := tcell.NewRGBColor(int32(randRange(0, 255)), int32(randRange(0, 255)), int32(randRange(0, 255)))
	return &col
}

func randRanger(min int, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63()))
  
  v := MapVal(r.Float64(), 0.0, 1.0, float64(min), float64(max))
	// return r.Intn(max-min+1) + min
  return int(v)
}

func randRange(min int, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63()))
	return r.Intn(max-min+1) + min
}

func RandomColor() Color {
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63()))
  r,g,b := rng.Float64(), rng.Float64(), rng.Float64()
	return Color{r, g, b}
}

func randomColor() tcell.Color {
	r := int32(randRange(0, 255))
	g := int32(randRange(0, 255))
	b := int32(randRange(0, 255))
	return tcell.NewRGBColor(r, g, b)
}

func GenerateRandomCoolors(count int) []*Coolor {
	tcols := make(Coolors, count)
	for i := range tcols {
		tcols[i] = MakeRandomCoolor()
	}
	return tcols
}
func GenerateRandomColors(count int) []tcell.Color {
	tcols := make([]tcell.Color, count)
	for i := range tcols {
		tcols[i] = *MakeRandomColor()
	}
	return tcols
}
func genRandomSeededColor() interface{} {
	rand.Seed(time.Now().UnixNano() + rand.Int63())
	tcol2 := MakeColorFromTcell(randomColor())
	return tcol2
}

