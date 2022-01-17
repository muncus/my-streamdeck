// Simple app to let me control the Elgato Stream Deck.
package main

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	streamdeck "github.com/magicmonkey/go-streamdeck"
	_ "github.com/magicmonkey/go-streamdeck/devices"
	"github.com/muncus/my-streamdeck/plugins"
	"github.com/muncus/my-streamdeck/plugins/googlemeet"
	"github.com/muncus/my-streamdeck/plugins/keylight"
	"github.com/muncus/my-streamdeck/plugins/obswebsocket"
	"github.com/pelletier/go-toml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var deckDevice *streamdeck.StreamDeck
var configFile = flag.String("config", "", "Config file, in yaml format")

func init() {
	// set up logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	var err error

	flag.Parse()

	// load a config file.
	abspath, _ := filepath.Abs(*configFile)
	config, err := toml.LoadFile(abspath)
	if err != nil {
		log.Fatal().Msgf("failed to read config file (%s): %s", abspath, err)
	}
	l, err := zerolog.ParseLevel(config.GetDefault("log_level", "debug").(string))
	if err != nil {
		log.Fatal().Msgf("invalid log_level (%s): %s", config.Get("log_level"), err)
	}
	zerolog.SetGlobalLevel(l)

	deckDevice, err = streamdeck.New()
	if err != nil {
		log.Fatal().Msgf("Failed to open Stream Deck: %s", err)
		os.Exit(1)
	}
	log.Debug().Msgf("Found streamdeck: %+v", deckDevice)

	// Meet Mutes
	meetPlugin, err := googlemeet.NewGoogleMeetPlugin(
		deckDevice,
		config.GetDefault("googlemeet", &toml.Tree{}).(*toml.Tree))
	if err != nil {
		log.Fatal().Msgf("failed to initialize googlemeet plugin: %s", err)
	}
	deckDevice.AddButton(0, meetPlugin.VideoMuteButton)
	deckDevice.AddButton(5, meetPlugin.MuteButton)
	deckDevice.AddButton(10, meetPlugin.RaiseHandButton)

	// OBS Plugin
	obsPlugin, err := obswebsocket.New(
		deckDevice,
		config.GetDefault("obswebsocket", &toml.Tree{}).(*toml.Tree))
	if err != nil {
		log.Fatal().Msgf("failed to initialize obswebsocket plugin: %s", err)
	}
	defer obsPlugin.Close()
	scene1, err := plugins.NewImageButtonFromFile("images/webcam_bg.png")
	if err != nil {
		log.Fatal().Msgf("Could not create Image button: %s", err)
	}
	scene1.SetActionHandler(obsPlugin.NewSceneChangeAction("webcam"))
	obsPlugin.ManageButton(scene1)
	deckDevice.AddButton(4, scene1)

	scene2, err := plugins.NewImageButtonFromFile("images/teapod-sad.png")
	if err != nil {
		log.Fatal().Msgf("Could not create Image button: %s", err)
	}
	scene2.SetActionHandler(obsPlugin.NewSceneChangeAction("sad-teapot"))
	obsPlugin.ManageButton(scene2)
	deckDevice.AddButton(9, scene2)

	//NB: this button is not managed by the OBSPlugin, because it should not be disabled when obs is inactive.
	obsbtn, err := plugins.NewImageButtonFromFile("images/obs.png")
	if err != nil {
		log.Fatal().Msgf("Could not create Image button: %s", err)
	}
	obsbtn.SetActionHandler(obsPlugin.LaunchOBSAction())
	deckDevice.AddButton(14, obsbtn)

	// keylights
	lightPlugin := keylight.New(deckDevice)
	deckDevice.AddButton(2, lightPlugin.PowerToggle)
	deckDevice.AddButton(7, lightPlugin.BrightnessInc)
	deckDevice.AddButton(12, lightPlugin.BrightnessDec)

	// Gracefully exit on interrupt, clearing buttons.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	select {
	case <-c:
		cleanup()
	}
}

func cleanup() {
	log.Debug().Msg("Cleaning up...")
	plugins.ClearButtons(deckDevice)
}
