package main

import (
	"strings"
	"fmt"

	"github.com/hugolgst/rich-go/client"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/unit"
)

type gui struct {
	id widget.Editor
	details widget.Editor
	state widget.Editor
	largeImage widget.Editor
	largeText widget.Editor
	smallImage widget.Editor
	smallText widget.Editor
	update widget.Clickable
}

func mainGui() {
	go func() {
		window := new(app.Window)

		window.Option(
			app.Title("GoRPC"),
			app.Size(unit.Dp(300), unit.Dp(320)),
		)

		if err := run(window); err != nil {
			panic(err)
		}
	}()

	app.Main()
}

func run(window *app.Window) error {
	var operations op.Ops

	theme := material.NewTheme()
	ui := gui{}

	for _, editor := range []*widget.Editor {
		&ui.id,
		&ui.details,
		&ui.state,
		&ui.largeImage,
		&ui.largeText,
		&ui.smallImage,
		&ui.smallText,
	} {
		editor.SingleLine = true
	}

	for {
		switch e := window.Event().(type) {
			case app.DestroyEvent:
				return e.Err

			case app.FrameEvent:
				digitsOnly(&ui.id)
				lowercase(&ui.largeImage)
				lowercase(&ui.smallImage)

				context := app.NewContext(&operations, e)

				for ui.update.Clicked(context) {
					activity := client.Activity{
						Details: ui.details.Text(),
						State: ui.state.Text(),
						LargeImage: ui.largeImage.Text(),
						LargeText: ui.largeText.Text(),
						SmallImage: ui.smallImage.Text(),
						SmallText: ui.smallText.Text(),
					}

					go func() {
						err := update(ui.id.Text(), activity, false, 5)

						if err != nil {
							fmt.Println(err)
						}
					}()
				}

				spacer := layout.Spacer{ Height: unit.Dp(10) }.Layout

				layout.UniformInset(16).Layout(context,
					func(context layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis: layout.Vertical,
						}.Layout(context,
							layout.Rigid(material.Editor(theme, &ui.id, "Client ID").Layout),
							layout.Rigid(spacer),
							layout.Rigid(material.Editor(theme, &ui.details, "Details").Layout),
							layout.Rigid(spacer),
							layout.Rigid(material.Editor(theme, &ui.state, "State").Layout),
							layout.Rigid(spacer),
							layout.Rigid(material.Editor(theme, &ui.largeImage, "Large Image").Layout),
							layout.Rigid(spacer),
							layout.Rigid(material.Editor(theme, &ui.largeText, "Large Hover Text").Layout),
							layout.Rigid(spacer),
							layout.Rigid(material.Editor(theme, &ui.smallImage, "Small Image").Layout),
							layout.Rigid(spacer),
							layout.Rigid(material.Editor(theme, &ui.smallText, "Small Hover Text").Layout),
							layout.Rigid(layout.Spacer{ Height: unit.Dp(20) }.Layout),
							layout.Rigid(material.Button(theme, &ui.update, "Update").Layout),
						)
					},
				)

				e.Frame(context.Ops)
			}
	}
}

func digitsOnly(editBox *widget.Editor) {
	input := editBox.Text()

	var builder strings.Builder

	for _, r := range input {
		if '0' <= r && r <= '9' {
			builder.WriteRune(r)
		}
	}

	if fixed := builder.String(); fixed != input {
		editBox.SetText(fixed)
	}
}

func lowercase(editBox *widget.Editor) {
	input := editBox.Text()
	lower := strings.ToLower(input)

	if lower != input {
		editBox.SetText(lower)
	}
}