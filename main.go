// Simple app to let me control the Elgato Stream Deck.
package main

import (
	"image/color"
	"os"
	"sync"

	streamdeck "github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/actionhandlers"
	"github.com/magicmonkey/go-streamdeck/buttons"
	_ "github.com/magicmonkey/go-streamdeck/devices"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	sd, err := streamdeck.New()
	if err != nil {
		log.Fatal().Msgf("Failed to open Stream Deck: %s", err)
		os.Exit(1)
	}
	log.Info().Msgf("Found streamdeck: %+v", sd)

	// Set up a button to do something.
	b1 := buttons.NewTextButton("   test   ")
	b1action := actionhandlers.NewCustomAction(func(streamdeck.Button) {
		log.Debug().Msg("button was pressed!")
	})
	b1.SetActionHandler(b1action)
	sd.AddButton(0, b1)

	// An image button.
	b2, err := buttons.NewImageFileButton("images/mic.png")
	if err != nil {
		log.Fatal().Msgf("Could not create Image button: %s", err)
	}
	b2.SetActionHandler(b1action)
	sd.AddButton(1, b2)

	cButton := buttons.NewColourButton(color.RGBA{255, 255, 0, 255})
	sd.AddButton(2, cButton)

	// An image button.
	b3, err := buttons.NewImageFileButton("images/teapod-sad.png")
	if err != nil {
		log.Fatal().Msgf("Could not create Image button: %s", err)
	}
	b3.SetActionHandler(b1action)
	sd.AddButton(3, b3)

	// wait for us to be done.
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
