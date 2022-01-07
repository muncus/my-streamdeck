// Simple app to let me control the Elgato Stream Deck.
package main

import (
	"os"
	"os/signal"

	streamdeck "github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/buttons"
	_ "github.com/magicmonkey/go-streamdeck/devices"
	"github.com/muncus/my-streamdeck/plugins"
	"github.com/muncus/my-streamdeck/plugins/googlemeet"
	"github.com/muncus/my-streamdeck/plugins/obswebsocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var deckDevice *streamdeck.StreamDeck

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	var err error
	deckDevice, err = streamdeck.New()
	if err != nil {
		log.Fatal().Msgf("Failed to open Stream Deck: %s", err)
		os.Exit(1)
	}
	log.Info().Msgf("Found streamdeck: %+v", deckDevice)

	// Meet Mutes
	meetPlugin, err := googlemeet.NewGoogleMeetPlugin(deckDevice)
	if err != nil {
		log.Error().Msgf("failed to initialize googlemeet plugin: %s", err)
	}
	deckDevice.AddButton(0, meetPlugin.VideoMuteButton)
	deckDevice.AddButton(5, meetPlugin.MuteButton)

	// OBS Plugin
	obsPlugin := obswebsocket.New(deckDevice)
	scene1 := obsPlugin.NewSceneButton("webcam")
	deckDevice.AddButton(4, scene1)
	scene2, err := buttons.NewImageFileButton("images/teapod-sad.png")
	if err != nil {
		log.Fatal().Msgf("Could not create Image button: %s", err)
	}
	scene2.SetActionHandler(obsPlugin.NewSceneChangeAction("sad-teapot"))
	deckDevice.AddButton(9, scene2)

	// Gracefully exit on interrupt, clearing buttons.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case <-c:
		cleanup()
	}
}

func cleanup() {
	log.Info().Msg("Cleaning up...")
	plugins.ClearButtons(deckDevice)
}
