package coolor

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/samber/lo"
)

type CoolorColorCluster struct {
	name      string
	leadColor Color
	colors    []*Color
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
    maxClusters int = 8
	ClusterPaletteTypes map[string]ClusterPalettes
	CoolorClusterColors ClusterPalettes
	baseAnsi            = []string{"#000000", "#c51e14", "#1dc121", "#c7c329", "#0a2fc4", "#c839c5", "#20c5c6", "#c7c7c7"}
	baseBrightAnsi      = []string{"#686868", "#fd6f6b", "#67f86f", "#fffa72", "#6a76fb", "#fd7cfc", "#68fdfe", "#ffffff"}
	baseAnsiNames       = []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "grey"}
	brightAnsiPrefix    = "bright"
	levels              []string
)

func init() {
	levels = make([]string, 0)
	levels = append(levels, "%sified")
	levels = append(levels, "%sish")
	levels = append(levels, "%sicated")
	levels = append(levels, "%s'd")
	levels = append(levels, "%s-touched")
	levels = append(levels, "%s'ed")
	// CoolorClusterColors = getCoolorClusterColors()
	CoolorClusterColors = getNamedAnsiColors()
	// CoolorClusterColors = getBaseAnsiClusterColors()
	colors := GenerateRandomColors(16)
	colors = append(colors, tcell.GetColor("#7b7b7b"))
	// found := make([]Color, 0)
	for _, v := range colors {
		ccci := NewCoolorColorClusterInfo(NewIntCoolorColor(v.Hex()))
		ccci.FindClusters()
		ccci.Sort()
		// fmt.Println(ccci.String())
		minDistance := 0.0
		minDistanceIndex := -1
		for itc, tcol := range CoolorClusterColors {
			vcol := MakeColorFromTcell(v)
			if minDistance == 0.0 || minDistance > tcol.leadColor.DistanceLuv(vcol) {
				minDistanceIndex = itc
				minDistance = tcol.leadColor.DistanceLuv(vcol)
			}
		}
		if minDistanceIndex != -1 {
			// cc := MakeColorFromTcell(v)
			// ClusterColors[minDistanceIndex].colors = append(ClusterColors[minDistanceIndex].colors, &cc)
		}
	}
	// for _, tcol := range ClusterColors {
	// 	colorss := lo.Map[*Color, string](tcol.colors, func(c *Color, i int) string {
	// 		return c.GetCC().TerminalPreview()
	// 	})
	// 	fmt.Println(fmt.Sprintf("%s %v", tcol.name, strings.Join(colorss, " ")))
	// }
}

func (cci *CoolorColorClusterInfo) Debug() string {
	colorss := lo.Map[*CoolorColorDistanceInfo, string](cci.clusters, func(c *CoolorColorDistanceInfo, i int) string {
		return fmt.Sprintf("%s %0.2f", c.cluster.leadColor.GetCC().TerminalPreview(), c.distance)
	})
	return fmt.Sprintf("%s %s", cci.color.TerminalPreview(), strings.Join(colorss, " "))
}

func (cci *CoolorColorClusterInfo) String() string {
	// colorss := lo.Map[*CoolorColorDistanceInfo, string](cci.clusters, func(c *CoolorColorDistanceInfo, i int) string {
	// 	return fmt.Sprintf("%s %0.2f", c.cluster.leadColor.GetCC().TerminalPreview(), c.distance)
	// })
	rand.Seed(time.Now().UnixNano())
	suff := lo.Sample[string](levels)
	rand.Seed(time.Now().UnixNano())
	suff2 := lo.Sample[string](levels)
	main := cci.clusters[0].cluster.name
	second := cci.clusters[1].cluster.name
	turd := cci.clusters[2].cluster.name
	return fmt.Sprintf("%s [yellow:-:-]%s %s %s[-:-:-]", cci.color.TerminalPreview(), fmt.Sprintf(suff2, turd), fmt.Sprintf(suff, second), main)
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
  for {
    
  }
}

func (cci *CoolorColorClusterInfo) FindClusters() {
	minDistance := 0.0
	minDistanceIndex := -1
	vcol := MakeColorFromTcell(*cci.color.color)
	for itc, tcol := range CoolorClusterColors {
		// distance := tcol.leadColor.DistanceCIEDE2000(vcol)
		distance := tcol.leadColor.DistanceRgb(vcol)
		ccdi := NewCoolorColorDistanceInfo(tcol)
		ccdi.distance = distance
		if minDistance == 0.0 || minDistance > distance {
			minDistanceIndex = itc
			minDistance = distance
		}
		cci.clusters = append(cci.clusters, ccdi)
    if len(cci.clusters) > maxClusters {
      break
    }
	}
	if minDistanceIndex != -1 {
		cc := vcol
		CoolorClusterColors[minDistanceIndex].colors = append(CoolorClusterColors[minDistanceIndex].colors, &cc)
	}
}

var (
  usefulColorNames []string = []string{"black", "green", "navy", "purple", "red", "lime", "yellow", "blue", "fuchsia", "aqua", "white", "azure", "chocolate", "crimson", "gold", "hotpink", "indigo", "lavender", "orange", "pink", "plum", "tomato", "turquoise", "violet", "grey", "darkgrey",}
)

func getNamedAnsiColors() ClusterPalettes {
  cps := make(ClusterPalettes, 0)
  for _, name  := range usefulColorNames {
    tcol := tcell.ColorNames[name]
    col := MakeColorFromTcell(tcol)
    cps = append(cps, NewClusterFromCss(name, col.Hex()))
  }
	return cps
}
func getBaseAnsiClusterColors() ClusterPalettes {
  cps := make(ClusterPalettes, 0)
  for i, c  := range baseAnsiNames {
    cps = append(cps, NewClusterFromCss(c, baseAnsi[i]))
  }
  for i, c  := range baseAnsiNames {
    name := fmt.Sprintf("%s %s", brightAnsiPrefix, c)
    cps = append(cps, NewClusterFromCss(name, baseBrightAnsi[i]))
  }
	return cps
}

func getCoolorClusterColors() ClusterPalettes {
	return ClusterPalettes{
		NewCluster("red", 255, 0, 0),
		NewCluster("orange", 255, 128, 0),
		NewCluster("yellow", 255, 255, 0),
		NewCluster("chartreuse", 128, 255, 0),
		NewCluster("puke", 128, 128, 0),
		NewCluster("green", 0, 255, 0),
		NewCluster("cerulean", 0, 128, 128),
		NewCluster("spring green", 0, 255, 128),
		NewCluster("navy green", 80, 190, 128),
		NewCluster("cyan", 0, 255, 255),
		NewCluster("azure", 0, 127, 255),
		NewCluster("blue", 0, 0, 255),
		NewCluster("violet", 127, 0, 255),
		NewCluster("magenta", 255, 0, 255),
		NewCluster("purple", 90, 0, 90),
		NewCluster("rose", 255, 0, 128),
		NewCluster("brown", 160, 40, 40),
		NewCluster("black", 0, 0, 0),
		NewCluster("grey", 200, 200, 200),
		// NewCluster("white", 255, 255, 255),
	}
}

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

// func ClusterCoolorPalette(cp *CoolorPalette)

/*
// import "github.com/digitallyserviced/coolors/color"

/*

function colorDistance(color1, color2) {
  const x =
    Math.pow(color1[0] - color2[0], 2) +
    Math.pow(color1[1] - color2[1], 2) +
    Math.pow(color1[2] - color2[2], 2);
  return Math.sqrt(x);
}
*/

// if tcol.leadColor.AlmostEqualRgb(vcol) {
// 	// found = append(found, vcol)
// 	// tcol.colors = append(tcol.colors, &vcol)
// }
// fmt.Printf("%v %f %s %s\n", v.leadColor.AlmostEqualRgb(tcol), v.leadColor.DistanceRgb(tcol), v.leadColor.GetCC().TerminalPreview(), tcol.GetCC().TerminalPreview())
// c, err := Hex(fmt.Sprintf("#%s",v.leadColor.Hex()))
// if err != nil {
//   panic(err)
// }
// col, ok := MakeColor(c)
// if !ok {
//   panic(fmt.Errorf("%v not a color", col))
// }

// if tcol.leadColor.AlmostEqualRgb(vcol) {
// found = append(found, vcol)
// tcol.colors = append(tcol.colors, &vcol)
// }
// fmt.Printf("%v %f %s %s\n", v.leadColor.AlmostEqualRgb(tcol), v.leadColor.DistanceRgb(tcol), v.leadColor.GetCC().TerminalPreview(), tcol.GetCC().TerminalPreview())
// foundcc := lo.Uniq[string](lo.Map[Color, string](found, func(c Color, i int) string {
//   return c.GetCC().TerminalPreview()
// }))
// fmt.Println(fmt.Sprintf("%s %v", tcol.name, strings.Join(foundcc, " ")))
/*

function oneDimensionSorting(colors, dim) {
  return colors
    .sort((colorA, colorB) => {
      if (colorA.hsl[dim] < colorB.hsl[dim]) {
        return -1;
      } else if (colorA.hsl[dim] > colorB.hsl[dim]) {
        return 1;
      } else {
        return 0;
      }
    });
}

function sortWithClusters(colorsToSort) {
  const mappedColors = colorsToSort
    .map((color) => {
      const isRgba = color.includes('rgba');
      if (isRgba) {
        return blendRgbaWithWhite(color);
      } else {
        return color;
      }
    })
    .map(colorUtil.color);

  mappedColors.forEach((color) => {
    let minDistance;
    let minDistanceClusterIndex;

    clusters.forEach((cluster, clusterIndex) => {
      const colorRgbArr = [color.rgb.r, color.rgb.g, color.rgb.b];
      const distance = colorDistance(colorRgbArr, cluster.leadColor);
      if (typeof minDistance === 'undefined' || minDistance > distance) {
        minDistance = distance;
        minDistanceClusterIndex = clusterIndex;
      }
    });

    clusters[minDistanceClusterIndex].colors.push(color);
  });

  clusters.forEach((cluster) => {
    const dim = ['white', 'grey', 'black'].includes(cluster.name) ? 'l' : 's';
    cluster.colors = oneDimensionSorting(cluster.colors, dim)
  });

  return clusters;
}

const sortedClusters = sortWithClusters(colors);
const sortedColors = sortedClusters.reduce((acc, curr) => {
  const colors = curr.colors.map((color) => color.hex);
  return [...acc, ...colors];
}, []);
renderColors(sortedColors, '#sorted');
*/
// vim: ts=2 sw=2 et ft=go
