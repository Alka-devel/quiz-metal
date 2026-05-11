package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"os"
	pusy "quiz/main/res"
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

type Metal struct {
	Number int
	Symbol string
	Name   string
	Hints  []string
}

var gosans = font.Font{
	Typeface: "Google Sans",
	Style:    font.Regular,
	Weight:   100,
}
var list = widget.List{
	List: layout.List{
		Axis: layout.Vertical,
	},
}

var ops op.Ops
var input widget.Editor
var theme = material.NewTheme()
var ready widget.Clickable
var buttons [57]widget.Clickable
var metals []pusy.Metal

var ques = "Нажмите старт."
var selected = -1
var currentQuestion, hint, score int
var started, failed bool

func main() {
	go func() {
		w := new(app.Window)

		w.Option(
			app.Title("Квиз по металлам"),
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
			interfaced(gtx, window)
			e.Frame(gtx.Ops)
		}
	}
}

func interfaced(gtx layout.Context, win *app.Window) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
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
							SW: r, SE: r,
						}
						stack := rr.Push(gtx.Ops)
						paint.ColorOp{Color: color.NRGBA{219, 178, 255, 255}}.Add(gtx.Ops)
						//paint.ColorOp{Color: color.NRGBA{111, 11, 111, 255}}.Add(gtx.Ops)
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
		layout.Flexed(0.8, func(gtx layout.Context) layout.Dimensions {

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
							SW: r, SE: r,
						}
						stack := rr.Push(gtx.Ops)
						paint.ColorOp{Color: color.NRGBA{219, 178, 255, 255}}.Add(gtx.Ops)
						//paint.ColorOp{Color: color.NRGBA{111, 11, 111, 255}}.Add(gtx.Ops)
						paint.PaintOp{}.Add(gtx.Ops)
						stack.Pop()
						return layout.Dimensions{Size: saiz}
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(6).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							perRow := max((gtx.Constraints.Max.X-13)/100, 1)

							rows := (len(buttons) + perRow - 1) / perRow
							estim := (gtx.Constraints.Max.X - 13) / perRow
							//fmt.Printf("perRow: %d rows: %d max: %d estim. size: %d \n", perRow, rows, gtx.Constraints.Max, estim)

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
													selected = -1 // выключить
												} else {
													selected = index // выбрать новую
												}
											}
											btn := material.Button(
												theme,
												&buttons[index],
												fmt.Sprintf("%s\n%s", pusy.Metals[i].Symbol, pusy.Metals[i].Name ),
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
										posay(selected)
									} else if !started && !failed {
										quiz()
									} else if !failed {
										lastQues := ques
										ques = fmt.Sprintf("%s\nВыберите металл!", lastQues)
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

								// click hitbox
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
	// duration is the number of seconds for the
	// completed animation: expand while fading in, then
	// out.
	const (
		expandDuration = float32(0.5)
		fadeDuration   = float32(0.9)
	)

	now := gtx.Now

	t := float32(now.Sub(c.Start).Seconds())

	end := c.End
	if end.IsZero() {
		// If the press hasn't ended, don't fade-out.
		end = now
	}

	endt := float32(end.Sub(c.Start).Seconds())

	// Compute the fade-in/out position in [0;1].
	var alphat float32
	{
		var haste float32
		if c.Cancelled {
			// If the press was cancelled before the inkwell
			// was fully faded in, fast forward the animation
			// to match the fade-out.
			if h := 0.5 - endt/fadeDuration; h > 0 {
				haste = h
			}
		}
		// Fade in.
		half1 := t/fadeDuration + haste
		if half1 > 0.5 {
			half1 = 0.5
		}

		// Fade out.
		half2 := float32(now.Sub(end).Seconds())
		half2 /= fadeDuration
		half2 += haste
		if half2 > 0.5 {
			// Too old.
			return
		}

		alphat = half1 + half2
	}

	// Compute the expand position in [0;1].
	sizet := t
	if c.Cancelled {
		// Freeze expansion of cancelled presses.
		sizet = endt
	}
	sizet /= expandDuration

	// Animate only ended presses, and presses that are fading in.
	if !c.End.IsZero() || sizet <= 1.0 {
		gtx.Execute(op.InvalidateCmd{})
	}

	if sizet > 1.0 {
		sizet = 1.0
	}

	if alphat > .5 {
		// Start fadeout after half the animation.
		alphat = 1.0 - alphat
	}
	// Twice the speed to attain fully faded in at 0.5.
	t2 := alphat * 2
	// Beziér ease-in curve.
	alphaBezier := t2 * t2 * (3.0 - 2.0*t2)
	sizeBezier := sizet * sizet * (3.0 - 2.0*sizet)
	size := gtx.Constraints.Min.X
	if h := gtx.Constraints.Min.Y; h > size {
		size = h
	}
	// Cover the entire constraints min rectangle and
	// apply curve values to size and color.
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
	metals = make([]pusy.Metal, len(pusy.Metals))
	copy(metals, pusy.Metals)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(metals), func(i, j int) {
		metals[i], metals[j] = metals[j], metals[i]
	})
	//score := 0
	currentQuestion, hint, started = 0, 0, true
	ques = metals[currentQuestion].Hints[hint]
}
func posay(index int) {
	metal := metals[currentQuestion]
	if pusy.Metals[index].Symbol == metal.Symbol {
		ques = fmt.Sprintf("Правильно: %s.\n", metal.Name)

		currentQuestion++
		hint = -1
		score++
		if currentQuestion >= 5 {
			ques = fmt.Sprintf("Ты угадал. Это был %s (%s).\nКонец игры", metal.Symbol, metal.Name)
			started = false
			failed = true
			selected = -1
		}
		return
	}
	hint++
	if hint < len(metal.Hints) {
		ques = metal.Hints[hint]
	}
	if hint >= len(metal.Hints) {
		ques = fmt.Sprintf("Ты не угадал. Это был %s (%s). Идём дальше.", metal.Symbol, metal.Name)

		currentQuestion++
		hint = -1
	}
	if currentQuestion >= 5 {
		ques = fmt.Sprintf("Ты не угадал. Это был %s (%s).\nКонец игры", metal.Symbol, metal.Name)
		started = false
		failed = true
		selected = -1
	}

}
