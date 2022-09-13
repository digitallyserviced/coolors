package coolor

// // Iter returns a copy of the Pallet in map form for iterating
// func (p *Pallet) Iter() map[string]*color.Color {
// 	return map[string]*color.Color{
// 		"bg":        p.bg,
// 		"bg_alt":    p.bg_alt,
// 		"fg":        p.fg,
// 		"fg_alt":    p.fg_alt,
// 		"pri":       p.pri,
// 		"sec":       p.sec,
// 		"primary":   p.pri,
// 		"secondary": p.sec,
// 		"alert":     p.alert,
// 		"cur":       p.cur,
// 		"cursor":    p.cur,
// 		"com":     p.com,
// 		"block":     p.block,
// 	}
// }
//
//
// // DefaultPallet fills a pallet with default values inspired by palenight
// func DefaultPallet() *Pallet {
// 	return &Pallet{
// 		bg:     color.NewColor("#292D3E"),
// 		bg_alt: color.NewColor("#697098"),
// 		pri:    color.NewColor("#c792ea"),
// 		sec:    color.NewColor("#C4E88D"),
// 		alert:  color.NewColor("#ff869a"),
// 		cur:    color.NewColor("#FFCB6B"),
// 		com:  color.NewColor("#82b1ff"),
// 		block:  color.NewColor("#939ede"),
// 		fg:     color.NewColor("#dde3eb"),
// 		fg_alt: color.NewColor("#C7D8FF"),
// 	}
// }
// // ApplyPallet reads a template from a reader, applies the given pallet to it
// // and then writes the filled in template to the writer.
// func ApplyPallet(r io.Reader, p *Pallet, w io.Writer) error {
//
// 	b, err := ioutil.ReadAll(r)
// 	if err != nil {
// 		return err
// 	}
//
// 	funcs := template.FuncMap{
// 		"hex":    color.HexString,
// 		"rgb":    color.RgbString,
// 		"rgb225": color.RgbString255,
// 		"hsv":    color.HsvString,
// 	}
//
// 	tmpl := template.Must(template.New("test").Funcs(funcs).Parse(string(b)))
//
// 	err = tmpl.Execute(w, p.Iter())
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
