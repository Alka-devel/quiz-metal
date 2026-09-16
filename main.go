package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

const (
	//RU-locale
	pressStart  = "Нажмите старт."
	appName     = "Квиз по металлам"
	chooseMetal = "Выберите металл!"
	hintStr     = "Подсказка"
	correct     = "Правильно"
	endGame     = "Конец игры."
	next        = "Идём дальше."
	youRight    = "Ты угадал. Это был"
)

var (
	gosans = font.Font{
		Typeface: "Google Sans",
		Style:    font.Regular,
		Weight:   100,
	}
	list = widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}

	ops     op.Ops
	input   widget.Editor
	theme   = material.NewTheme()
	ready   widget.Clickable
	buttons [57]widget.Clickable
	metals  []Metal

	ques                         = pressStart
	selected                     = -1
	currentQuestion, hint, score int
	started, failed, breac       bool
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title(appName),
			app.Size(600, 400),
			app.MinSize(600, 400),
		)
		if err := RunUI(w); err != nil {
			fmt.Println("Error:", err)
			os.Exit(0)
		}
	}()
	app.Main()
}

func RunUI(window *app.Window) error {
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			os.Exit(0)
			return e.Err
		case app.FrameEvent:
			ops.Reset()
			gtx := app.NewContext(&ops, e)
			paint.Fill(gtx.Ops, color.NRGBA{242, 231, 254, 255})
			interfaced(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func interfaced(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.Y = gtx.Dp(75)
			gtx.Constraints.Min.Y = gtx.Dp(75)
			return layout.UniformInset(5).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						r := gtx.Dp(unit.Dp(15))
						dl := gtx.Constraints
						saiz := image.Point{dl.Max.X, dl.Max.Y}
						rect := image.Rect(0, 0, saiz.X, saiz.Y)
						rr := clip.RRect{
							Rect: rect,
							NW:   r, NE: r,
							SW: r, SE: r}
						stack := rr.Push(gtx.Ops)
						paint.ColorOp{Color: color.NRGBA{219, 178, 255, 255}}.Add(gtx.Ops)
						paint.PaintOp{}.Add(gtx.Ops)
						stack.Pop()
						return layout.Dimensions{Size: saiz}
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(6).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								label := material.Body1(theme, ques)
								label.Font = gosans
								label.Color = color.NRGBA{0, 0, 0, 255}
								label.TextSize = 20
								return label.Layout(gtx)
							})
						})
					}),
				)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(5).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						r := gtx.Dp(unit.Dp(15))
						dl := gtx.Constraints
						saiz := image.Point{dl.Max.X, dl.Max.Y}
						rect := image.Rect(0, 0, saiz.X, saiz.Y)
						rr := clip.RRect{
							Rect: rect,
							NW:   r, NE: r,
							SW: r, SE: r}
						stack := rr.Push(gtx.Ops)
						paint.ColorOp{Color: color.NRGBA{219, 178, 255, 255}}.Add(gtx.Ops)
						paint.PaintOp{}.Add(gtx.Ops)
						stack.Pop()
						return layout.Dimensions{Size: saiz}
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(6).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							perRow := max((gtx.Constraints.Max.X-13)/100, 1)
							rows := (len(buttons) + perRow - 1) / perRow
							estim := (gtx.Constraints.Max.X - 13) / perRow
							return material.List(theme, &list).Layout(gtx, rows, func(gtx layout.Context, row int) layout.Dimensions {
								children := []layout.FlexChild{}
								start := row * perRow
								end := min(start+perRow, len(buttons))
								for i := start; i < end; i++ {
									index := i
									children = append(children,
										layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
											saiz := image.Point{
												X: estim,
												Y: estim,
											}
											gtx.Constraints.Min, gtx.Constraints.Max = saiz, saiz
											if buttons[index].Clicked(gtx) {
												if selected == index || failed {
													selected = -1
												} else {
													selected = index
												}
											}
											btn := material.Button(
												theme,
												&buttons[index],
												fmt.Sprintf("%s\n%s", Metals[i].Symbol, Metals[i].Name),
											)
											btn.TextSize = 13
											bg := color.NRGBA{187, 134, 252, 255}
											if selected == index && !failed {
												bg = color.NRGBA{120, 70, 220, 255}
											}
											btn.Background = bg
											return layout.UniformInset(unit.Dp(4)).Layout(
												gtx,
												btn.Layout,
											)
										}),
									)
								}
								return layout.Flex{
									Axis: layout.Horizontal,
								}.Layout(gtx, children...)
							},
							)
						})
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min = gtx.Constraints.Max
						return layout.SE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Bottom: 5,
								Right:  5,
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								if ready.Clicked(gtx) {
									if started && selected != -1 && !failed {
										tech(selected)
									} else if !started && !failed {
										quiz()
									} else if !failed && !strings.Contains(strings.ToLower(ques), strings.ToLower("Выберите металл")) {
										lastQues := ques
										ques = fmt.Sprintf("%s\n%s", lastQues, chooseMetal)
									}

								}
								m := op.Record(gtx.Ops)
								dims := layout.Stack{}.Layout(gtx,
									layout.Expanded(func(gtx layout.Context) layout.Dimensions {
										size := image.Pt(60, 60)
										r := 15
										defer clip.RRect{
											Rect: image.Rectangle{Max: size},
											NW:   r, NE: r, SW: r, SE: r,
										}.Push(gtx.Ops).Pop()
										paint.Fill(gtx.Ops, color.NRGBA{127, 57, 251, 255})
										for _, c := range ready.History() {
											drawInk(gtx, c)
										}
										return layout.Dimensions{Size: size}
									}),
									layout.Stacked(func(gtx layout.Context) layout.Dimensions {
										size := 60
										var ikonka *widget.Icon
										if started {
											ikonka, _ = widget.NewIcon(icons.ActionDone)
										} else {
											ikonka, _ = widget.NewIcon(icons.AVPlayArrow)
										}
										if breac {
											ikonka, _ = widget.NewIcon(icons.HardwareKeyboardArrowRight)
										}
										gtx.Constraints.Min = image.Pt(size, size)
										gtx.Constraints.Max = image.Pt(size, size)
										return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return ikonka.Layout(gtx, color.NRGBA{255, 255, 255, 255})
										})
									}),
								)
								c := m.Stop()
								r := 15
								defer clip.RRect{
									Rect: image.Rectangle{Max: dims.Size},
									NW:   r, NE: r, SW: r, SE: r,
								}.Push(gtx.Ops).Pop()
								c.Add(gtx.Ops)
								return material.Clickable(gtx, &ready, func(gtx layout.Context) layout.Dimensions {
									return dims
								})
							})
						})
					}),
				)
			})
		}),
	)
}
func drawInk(gtx layout.Context, c widget.Press) {
	const (
		expandDuration = float32(0.5)
		fadeDuration   = float32(0.9)
	)
	now := gtx.Now
	t := float32(now.Sub(c.Start).Seconds())
	end := c.End
	if end.IsZero() {
		end = now
	}
	endt := float32(end.Sub(c.Start).Seconds())
	var alphat float32
	{
		var haste float32
		if c.Cancelled {
			if h := 0.5 - endt/fadeDuration; h > 0 {
				haste = h
			}
		}
		half1 := t/fadeDuration + haste
		if half1 > 0.5 {
			half1 = 0.5
		}
		half2 := float32(now.Sub(end).Seconds())
		half2 /= fadeDuration
		half2 += haste
		if half2 > 0.5 {
			return
		}
		alphat = half1 + half2
	}
	sizet := t
	if c.Cancelled {
		sizet = endt
	}
	sizet /= expandDuration
	if !c.End.IsZero() || sizet <= 1.0 {
		gtx.Execute(op.InvalidateCmd{})
	}
	if sizet > 1.0 {
		sizet = 1.0
	}
	if alphat > .5 {
		alphat = 1.0 - alphat
	}
	t2 := alphat * 2
	alphaBezier := t2 * t2 * (3.0 - 2.0*t2)
	sizeBezier := sizet * sizet * (3.0 - 2.0*sizet)
	size := gtx.Constraints.Min.X
	if h := gtx.Constraints.Min.Y; h > size {
		size = h
	}
	size = int(float32(size) * 2 * float32(math.Sqrt(2)) * sizeBezier)
	alpha := 0.7 * alphaBezier
	const col = 0.8
	ba, bc := byte(alpha*0xff), byte(col*0xff)
	rgba := color.NRGBA{
		R: bc,
		G: bc,
		B: bc,
		A: ba,
	}
	ink := paint.ColorOp{Color: rgba}
	ink.Add(gtx.Ops)
	rr := size / 2
	defer op.Offset(c.Position.Add(image.Point{
		X: -rr,
		Y: -rr,
	})).Push(gtx.Ops).Pop()
	defer clip.UniformRRect(image.Rectangle{Max: image.Pt(size, size)}, rr).Push(gtx.Ops).Pop()
	paint.PaintOp{}.Add(gtx.Ops)
}

func quiz() {
	metals = make([]Metal, len(Metals))
	copy(metals, Metals)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(metals), func(i, j int) {
		metals[i], metals[j] = metals[j], metals[i]
	})
	currentQuestion, hint, started, breac = 0, 0, true, false
	ques = fmt.Sprintf("%s %d: %s", hintStr, hint+1, metals[currentQuestion].Hints[hint])
}
func tech(index int) {
	metal := metals[currentQuestion]
	if Metals[index].Symbol == metal.Symbol {
		ques = fmt.Sprintf("%s: %s.\n", correct, metal.Name)
		currentQuestion++
		hint = -1
		score++
		if currentQuestion >= 5 {
			ques = fmt.Sprintf("%s %s (%s).\n%s", youRight, metal.Symbol, metal.Name, endGame)
			breac = false
			started = false
			failed = true
			selected = -1
		}
		return
	}
	hint++
	if hint < len(metal.Hints) {
		breac = false
		ques = fmt.Sprintf("%s %d: %s", hintStr, hint+1, metal.Hints[hint])
	}
	if hint >= len(metal.Hints) {
		ques = fmt.Sprintf("%s %s (%s). %s", youRight, metal.Symbol, metal.Name, next)
		breac = true
		currentQuestion++
		hint = -1
	}
	if currentQuestion >= 5 {
		ques = fmt.Sprintf("%s %s (%s).\n%s", youRight, metal.Symbol, metal.Name, endGame)
		breac = false
		started = false
		failed = true
		selected = -1
	}
}
