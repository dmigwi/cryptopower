package cryptomaterial

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/crypto-power/cryptopower/app"
	"github.com/crypto-power/cryptopower/ui/values"
)

func drawInk(gtx layout.Context, c widget.Press, highlightColor color.NRGBA) {
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
	// Cover the entire constraints min rectangle.
	size *= 2 * int(math.Sqrt(2))
	// Apply curve values to size and color.
	size *= int(sizeBezier)
	alpha := 0.7 * alphaBezier
	const col = 0.8
	ba, _ := byte(alpha*0xff), byte(col*0xff)
	rgba := mulAlpha(highlightColor, ba)
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

func GenerateRandomNumber() int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int()
}

func UniformPaddingWithTopInset(topInset unit.Dp, gtx layout.Context, body layout.Widget, isMobileView ...bool) D {
	_isMobileView := len(isMobileView) > 0 && isMobileView[0]
	width := gtx.Constraints.Max.X
	paddingHorizontal := values.MarginPadding24
	bottomPadding := values.MarginPadding24

	if (width - 2*gtx.Dp(paddingHorizontal)) > gtx.Dp(values.AppWidth) {
		paddingValue := float32(width-gtx.Dp(values.AppWidth)) / 4
		paddingHorizontal = unit.Dp(paddingValue)
	}

	if _isMobileView {
		paddingHorizontal = values.MarginPadding16
		bottomPadding = values.MarginPadding0
	}

	return layout.Inset{
		Top:    topInset,
		Right:  paddingHorizontal,
		Bottom: bottomPadding,
		Left:   paddingHorizontal,
	}.Layout(gtx, body)
}

func UniformPadding(gtx layout.Context, body layout.Widget, isMobileView ...bool) D {
	return UniformPaddingWithTopInset(values.MarginPadding5, gtx, body, isMobileView...)
}

func DisableLayout(currentPage app.Page, gtx C, titleLayout, subtitleLayout func(gtx C) D, transparency uint8, color color.NRGBA, actionButton *Button) D {
	return layout.Stack{Alignment: layout.N}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			if currentPage == nil {
				return D{}
			}
			mgtx := gtx.Disabled()
			return currentPage.Layout(mgtx)
		}),
		layout.Stacked(func(gtx C) D {
			overlayColor := color
			overlayColor.A = transparency
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			FillMax(gtx, overlayColor, 10)
			if titleLayout == nil && subtitleLayout == nil && actionButton == nil {
				return D{}
			}

			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if titleLayout == nil {
							return D{}
						}
						return titleLayout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						if subtitleLayout == nil {
							return D{}
						}
						return subtitleLayout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						if actionButton == nil {
							return D{}
						}
						actionButton.TextSize = values.TextSize14
						return actionButton.Layout(gtx)
					}),
				)
			})
		}),
	)
}
