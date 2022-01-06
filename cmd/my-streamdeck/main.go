// Simple app to let me control the Elgato Stream Deck.
package main

import (
	"os"
	"os/signal"

	streamdeck "github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/actionhandlers"
	"github.com/magicmonkey/go-streamdeck/buttons"
	_ "github.com/magicmonkey/go-streamdeck/devices"
	"github.com/muncus/my-streamdeck/plugins"
	"github.com/muncus/my-streamdeck/plugins/googlemeet"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var deckDevice *streamdeck.StreamDeck

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	sd, err := streamdeck.New()
	if err != nil {
		log.Fatal().Msgf("Failed to open Stream Deck: %s", err)
		os.Exit(1)
	}
	deckDevice = sd
	log.Info().Msgf("Found streamdeck: %+v", sd)

	// define some actions.
	debugAction := actionhandlers.NewCustomAction(func(streamdeck.Button) {
		log.Debug().Msg("button was pressed!")
	})

	// // Set up a button to do something.
	// b1 := buttons.NewTextButton(" test ")
	// b1.SetActionHandler(debugAction)
	// sd.AddButton(0, b1)
	// sd.SetDecorator(b1.GetButtonIndex(), decorators.NewBorder(10, color.RGBA{255, 0, 0, 255}))

	// Meet Mutes
	meetPlugin, err := googlemeet.NewGoogleMeetPlugin(sd)
	if err != nil {
		log.Error().Msgf("failed to initialize googlemeet plugin: %s", err)
	}
	sd.AddButton(0, meetPlugin.VideoMuteButton)
	sd.AddButton(5, meetPlugin.MuteButton)

	// An image button.
	teapotButton, err := buttons.NewImageFileButton("images/teapod-sad.png")
	if err != nil {
		log.Fatal().Msgf("Could not create Image button: %s", err)
	}
	teapotButton.SetActionHandler(debugAction)
	sd.AddButton(2, teapotButton)

	// Gracefully exit, clearing buttons.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	select {
	case sig := <-c:
		_ = sig
		cleanup()
	}

	// // wait for us to be done.
	// var wg sync.WaitGroup
	// wg.Add(1)
	// wg.Wait()
}

func cleanup() {
	log.Info().Msg("Cleaning up...")
	plugins.ClearButtons(deckDevice)
}
