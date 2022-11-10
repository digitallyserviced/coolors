package coolor

import (
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/samber/lo"
)

type CoolorColorCluster struct {
	name      string
	colors    []*Color
	leadColor Color
}

type CoolorColorDistanceInfo struct {
	cluster  *CoolorColorCluster
	distance float64
}

func NewCoolorColorDistanceInfo(cluster *CoolorColorCluster) *CoolorColorDistanceInfo {
	ccdi := &CoolorColorDistanceInfo{
		cluster:  cluster,
		distance: 0,
	}
	return ccdi
}

type CoolorColorClusterInfo struct {
	color    *CoolorColor
	clusters []*CoolorColorDistanceInfo
}

func NewCoolorColorClusterInfo(cc *CoolorColor) *CoolorColorClusterInfo {
	ccci := &CoolorColorClusterInfo{
		color:    cc,
		clusters: make([]*CoolorColorDistanceInfo, 0),
	}
	return ccci
}

type ClusterPalettes []*CoolorColorCluster

var (
	maxClusters         int = 8
	ClusterPaletteTypes map[string]ClusterPalettes
	CoolorClusterColors ClusterPalettes
	// friendly base ansi colors
	baseAnsi       = []string{"#000000", "#c51e14", "#1dc121", "#c7c329", "#0a2fc4", "#c839c5", "#20c5c6", "#7c7c7c"}
	baseBrightAnsi = []string{"#686868", "#fd6f6b", "#67f86f", "#fffa72", "#6a76fb", "#fd7cfc", "#68fdfe", "#ffffff"}

	// strict/hard base ansi colors
	baseXterm = []string{"#000000", "#800000", "#008000", "#808000", "#000080", "#800080", "#008080", "#707070"}

	baseBrightXterm = []string{"#A9A9A9", "#FF0000", "#00FF00", "#FFFF00", "#0000FF", "#FF00FF", "#00FFFF", "#FFFFFF"}

	// base16 xterm color names
	baseXtermAnsiColorNames = []string{"black", "maroon", "green", "olive", "navy", "purple", "teal", "silver", "gray", "red", "lime", "yellow", "blue", "fuchsia", "aqua", "white"}

	// friendly base 16 colors with either dim or bright attributes added later
	baseAnsiNames = []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "grey"}
	// if 0-7 are considered the "base" or "normal" bright is 8-15
	brightAnsiPrefix = "bright"
	// some terms see 8-15 as considered the "base" or "normal" with dim being 0-7
	dimAnsiPrefix = "dim"
	// ColorLightGray:            0xD3D3D3,
	// ColorLightSlateGray:       0x778899,
	// ColorSlateGray:            0x708090,
	// ColorDimGray:              0x696969,
	// ColorDarkSlateGray:        0x2F4F4F,
	// ColorDarkGray:             0xA9A9A9,
	// ColorGray:                 0x808080,

	// ref all base16 colors with tcell strict/hard bases
	baseXtermAnsiTcellColors = []tcell.Color{tcell.ColorBlack, tcell.ColorMaroon, tcell.ColorGreen, tcell.ColorOlive, tcell.ColorNavy, tcell.ColorPurple, tcell.ColorTeal, tcell.ColorSilver, tcell.ColorGray, tcell.ColorRed, tcell.ColorLime, tcell.ColorYellow, tcell.ColorBlue, tcell.ColorFuchsia, tcell.ColorAqua, tcell.ColorWhite}
	levels                   []string
)

func init() {
	levels = make([]string, 0)
	levels = append(levels, "%sified")
	levels = append(levels, "%sish")
	// levels = append(levels, "%sicated")
	// levels = append(levels, "%s'd")
	// levels = append(levels, "%s-touched")
	levels = append(levels, "%s'ed")
	// CoolorClusterColors = getCoolorClusterColors()
	CoolorClusterColors = getBaseXtermClusterColors()
	// CoolorClusterColors = getNamedAnsiColors()
	// CoolorClusterColors = getBaseAnsiClusterColors()
	// colors := GenerateRandomColors(16)
	// colors = append(colors, tcell.GetColor("#7b7b7b"))
	// found := make([]Color, 0)
	// for _, v := range colors {
	// 	ccci := NewCoolorColorClusterInfo(NewIntCoolorColor(v.Hex()))
	// 	ccci.FindClusters()
	// 	ccci.Sort()
	// 	// fmt.Println(ccci.String())
	// 	minDistance := 0.0
	// 	minDistanceIndex := -1
	// 	for itc, tcol := range CoolorClusterColors {
	// 		vcol := MakeColorFromTcell(v)
	// 		if minDistance == 0.0 || minDistance > tcol.leadColor.DistanceLuv(vcol) {
	// 			minDistanceIndex = itc
	// 			minDistance = tcol.leadColor.DistanceLuv(vcol)
	// 		}
	// 	}
	// 	if minDistanceIndex != -1 {
	// 		// cc := MakeColorFromTcell(v)
	// 		// ClusterColors[minDistanceIndex].colors = append(ClusterColors[minDistanceIndex].colors, &cc)
	// 	}
	// }
	// for _, tcol := range ClusterColors {
	// 	colorss := lo.Map[*Color, string](tcol.colors, func(c *Color, i int) string {#e85c51
	// 		return c.GetCC().TerminalPreview()
	// 	})
	// 	fmt.Println(fmt.Sprintf("%s %v", tcol.name, strings.Join(colorss, " ")))
	// }
}
var t string = "#e85c51"
func (cci *CoolorColorClusterInfo) Debug() string {
	colorss := lo.Map[*CoolorColorDistanceInfo, string](cci.clusters, func(c *CoolorColorDistanceInfo, i int) string {
		return fmt.Sprintf("%s %0.2f", c.cluster.leadColor.GetCC().TerminalPreview(), c.distance)
	})
	return fmt.Sprintf("%s %s", cci.color.TerminalPreview(), strings.Join(colorss, " "))
}

func (cci *CoolorColorClusterInfo) String() string {
	rand.Seed(time.Now().UnixNano())
	main := cci.clusters[0].cluster.name
  fondler := cci.DemoteGray(cci.clusters[1].cluster.name, cci.clusters[2].cluster.name)
  adj := Generator().Adjectives(1)
  adjd := fmt.Sprintf("%s %s", adj[0], fondler)
	return fmt.Sprintf("%s [yellow:-:-]%s %s[-:-:-]", cci.color.TVPreview(), adjd, main)
}

func (cci *CoolorColorClusterInfo) DemoteGray(b, a string) string {
  if m, err := regexp.Match("silver|gray|grey", []byte(a)); m {
    return b
  } else if err != nil {
    panic(err)
  }
  return a
}

func (cci *CoolorColorClusterInfo) Sort() {
	sort.Sort(cci)
}

func (cci *CoolorColorClusterInfo) Swap(a, b int) {
	cci.clusters[a], cci.clusters[b] = cci.clusters[b], cci.clusters[a]
}

func (cci *CoolorColorClusterInfo) Less(a, b int) bool {
	return cci.clusters[a].distance < cci.clusters[b].distance
}

func (cci *CoolorColorClusterInfo) Len() int {
	return len(cci.clusters)
}

func (cci *CoolorColorClusterInfo) GenerateNeighbors(count int, maxDistance float64) {
	// for {
	//
	// }
}

func (cci *CoolorColorClusterInfo) FindClusters() {
	minDistance := 0.0
	minDistanceIndex := -1
	vcol := MakeColorFromTcell(*cci.color.Color)
	for itc, tcol := range CoolorClusterColors {
		distance := tcol.leadColor.DistanceCIEDE2000(vcol)
		ccdi := NewCoolorColorDistanceInfo(tcol)
		ccdi.distance = distance
		if minDistance == 0.0 || minDistance > distance {
			minDistanceIndex = itc
			minDistance = distance
		}
		cci.clusters = append(cci.clusters, ccdi)
		// if len(cci.clusters) > maxClusters {
		// 	break
		// }
	}
	if minDistanceIndex != -1 {
		cc := vcol
		CoolorClusterColors[minDistanceIndex].colors = append(CoolorClusterColors[minDistanceIndex].colors, &cc)
	}
}

var usefulColorNames []string = []string{"black", "green", "navy", "purple", "red", "lime", "yellow", "blue", "fuchsia", "aqua", "white", "azure", "chocolate", "crimson", "gold", "hotpink", "indigo", "lavender", "orange", "pink", "plum", "tomato", "turquoise", "violet", "grey", "darkgrey"}

func getNamedAnsiColors() ClusterPalettes {
	cps := make(ClusterPalettes, 0)
	for _, name := range usefulColorNames {
		tcol := tcell.ColorNames[name]
		col := MakeColorFromTcell(tcol)
		cps = append(cps, NewClusterFromCss(name, col.Hex()))
	}
	return cps
}

func getBaseXtermClusterColors() ClusterPalettes {
	cps := make(ClusterPalettes, 0)
	for i, c := range baseXtermAnsiColorNames {
		col := MakeColorFromTcell(baseXtermAnsiTcellColors[i])
		cps = append(cps, NewClusterFromCss(c, col.Hex()))
	}
	return cps
}

func getBaseAnsiClusterColors() ClusterPalettes {
	cps := make(ClusterPalettes, 0)
	for i, c := range baseAnsiNames {
		cps = append(cps, NewClusterFromCss(c, baseAnsi[i]))
	}
	for i, c := range baseAnsiNames {
		name := fmt.Sprintf("%s %s", brightAnsiPrefix, c)
		cps = append(cps, NewClusterFromCss(name, baseBrightAnsi[i]))
	}
	return cps
}

// func getCoolorClusterColors() ClusterPalettes {
// 	return ClusterPalettes{
// 		NewCluster("red", 255, 0, 0),
// 		NewCluster("orange", 255, 128, 0),
// 		NewCluster("yellow", 255, 255, 0),
// 		NewCluster("chartreuse", 128, 255, 0),
// 		NewCluster("puke", 128, 128, 0),
// 		NewCluster("green", 0, 255, 0),
// 		NewCluster("cerulean", 0, 128, 128),
// 		NewCluster("spring green", 0, 255, 128),
// 		NewCluster("navy green", 80, 190, 128),
// 		NewCluster("cyan", 0, 255, 255),
// 		NewCluster("azure", 0, 127, 255),
// 		NewCluster("blue", 0, 0, 255),
// 		NewCluster("violet", 127, 0, 255),
// 		NewCluster("magenta", 255, 0, 255),
// 		NewCluster("purple", 90, 0, 90),
// 		NewCluster("rose", 255, 0, 128),
// 		NewCluster("brown", 160, 40, 40),
// 		NewCluster("black", 0, 0, 0),
// 		NewCluster("grey", 200, 200, 200),
// 		// NewCluster("white", 255, 255, 255),
// 	}
// }

func NewClusterFromCss(s, css string) *CoolorColorCluster {
	col, err := Hex(css)
	if err != nil {
		return nil
	}
	ccc := &CoolorColorCluster{
		name:      s,
		leadColor: col,
		colors:    make([]*Color, 0),
	}

	return ccc
}

func NewCluster(s string, i1, i2, i3 uint8) *CoolorColorCluster {
	col := RGB255(i1, i2, i3)
	ccc := &CoolorColorCluster{
		name:      s,
		leadColor: col,
		colors:    make([]*Color, 0),
	}

	return ccc
}

func rgbToYIQ(r, g, b uint) float64 {
	return float64((r*299)+(g*587)+(b*114)) / 1000.0
}

// vim: ts=2 sw=2 et ft=go
