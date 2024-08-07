package settings

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/crypto-power/cryptopower/app"
	"github.com/crypto-power/cryptopower/ui/cryptomaterial"
	"github.com/crypto-power/cryptopower/ui/load"
	"github.com/crypto-power/cryptopower/ui/page/components"
	"github.com/crypto-power/cryptopower/ui/values"
)

const LicensePageID = "License"

const license = `ISC License

Copyright (c) 2018-2019 The Decred developers

Permission to use, copy, modify, and distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.`

type LicensePage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	pageContainer *widget.List
	backButton    cryptomaterial.IconButton
}

func NewLicensePage(l *load.Load) *LicensePage {
	pg := &LicensePage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(LicensePageID),
		pageContainer: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}
	pg.backButton = components.GetBackButton(l)

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *LicensePage) OnNavigatedTo() {}

// Layout draws the page UI components into the provided C
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *LicensePage) Layout(gtx C) D {
	if pg.Load.IsMobileView() {
		return pg.layoutMobile(gtx)
	}
	return pg.layoutDesktop(gtx)
}

func (pg *LicensePage) layoutDesktop(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(values.MarginPadding20).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.pageHeaderLayout),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding16, Bottom: values.MarginPadding20}.Layout(gtx, pg.pageContentLayout)
			}),
		)
	})
}

func (pg *LicensePage) layoutMobile(gtx layout.Context) layout.Dimensions {
	d := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      values.String(values.StrLicense),
			BackButton: pg.backButton,
			Back: func() {
				pg.ParentNavigator().CloseCurrentPage()
			},
			Body: func(gtx C) D {
				return pg.Theme.List(pg.pageContainer).Layout(gtx, 1, func(gtx C, _ int) D {
					return pg.Theme.Card().Layout(gtx, func(gtx C) D {
						return layout.UniformInset(values.MarginPadding25).Layout(gtx, func(gtx C) D {
							licenseText := pg.Theme.Body1(license)
							licenseText.Color = pg.Theme.Color.GrayText2
							return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, licenseText.Layout)
						})
					})
				})
			},
		}
		return sp.Layout(pg.ParentWindow(), gtx)
	}
	gtx.Constraints.Min.Y = gtx.Constraints.Max.Y

	return components.UniformMobile(gtx, false, false, d)
}

func (pg *LicensePage) pageHeaderLayout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.W.Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Right: values.MarginPadding16,
							Top:   values.MarginPaddingMinus2,
						}.Layout(gtx, pg.backButton.Layout)
					}),
					layout.Rigid(pg.Theme.Label(values.TextSize20, values.String(values.StrLicense)).Layout),
				)
			})
		}),
	)
}

func (pg *LicensePage) pageContentLayout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Dp(unit.Dp(560))
		gtx.Constraints.Max.X = gtx.Constraints.Min.X
		gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
		return pg.Theme.List(pg.pageContainer).Layout(gtx, 1, func(gtx C, _ int) D {
			return pg.Theme.Card().Layout(gtx, func(gtx C) D {
				return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
					licenseText := pg.Theme.Body1(license)
					licenseText.Color = pg.Theme.Color.GrayText1
					return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, licenseText.Layout)
				})
			})
		})
	})
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *LicensePage) HandleUserInteractions(gtx C) {
	if pg.backButton.Button.Clicked(gtx) {
		pg.ParentNavigator().CloseCurrentPage()
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *LicensePage) OnNavigatedFrom() {}
