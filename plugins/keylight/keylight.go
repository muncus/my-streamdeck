// package keylight contains methods for interacting with Elgato Key Lights from the stream deck.
package keylight

import (
	"context"
	"os"
	"time"

	kl "github.com/endocrimes/keylight-go"
	"github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/actionhandlers"
	"github.com/magicmonkey/go-streamdeck/buttons"
	"github.com/muncus/my-streamdeck/plugins"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger = log.Logger.With().Str("plugin", "keylight").Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})

type KeylightPlugin struct {
	d               *streamdeck.StreamDeck
	light           *kl.Device
	ready           bool
	brightIncrement int
	tempIncrement   int
	PowerToggle     plugins.ActionButton
	BrightnessInc   plugins.ActionButton
	BrightnessDec   plugins.ActionButton
}

func New(d *streamdeck.StreamDeck) *KeylightPlugin {
	plugin := &KeylightPlugin{
		ready:           false,
		d:               d,
		brightIncrement: 10,
		tempIncrement:   1,
	}
	// run discovery in a background goroutine.
	go plugin.prepare()

	// plugin.PowerToggle = buttons.NewTextButton("aziz, light!")
	plugin.PowerToggle, _ = buttons.NewImageFileButton("images/light_mode_bg.png")
	plugin.PowerToggle.SetActionHandler(plugin.LightAction(func(l *kl.Light) {
		if l.On > 0 {
			l.On = 0
		} else {
			l.On = 1
		}
	}))

	// TODO: add a min/max to these values.
	plugin.BrightnessInc, _ = buttons.NewImageFileButton("images/more_bright_bg.png")
	plugin.BrightnessInc.SetActionHandler(plugin.LightAction(func(l *kl.Light) {
		l.Brightness = l.Brightness + plugin.brightIncrement
	}))
	plugin.BrightnessDec, _ = buttons.NewImageFileButton("images/less_bright_bg.png")
	plugin.BrightnessDec.SetActionHandler(plugin.LightAction(func(l *kl.Light) {
		l.Brightness = l.Brightness - plugin.brightIncrement
	}))

	return plugin
}

func (p *KeylightPlugin) prepare() error {
	Logger.Debug().Msg("starting discovery")
	disc, err := kl.NewDiscovery()
	if err != nil {
		return err
	}
	dctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = disc.Run(dctx)
	if err != nil {
		return err
	}
	// NB: this only finds one light.
	p.light = <-disc.ResultsCh()
	Logger.Debug().Msgf("found: %s", p.light)
	p.ready = true
	return nil
}

// LightAction helps construct button actions that operate onall lights in a given LightGroup
func (p *KeylightPlugin) LightAction(lightfunc func(l *kl.Light)) streamdeck.ButtonActionHandler {
	return actionhandlers.NewCustomAction(func(streamdeck.Button) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		lg, err := p.light.FetchLightGroup(ctx)
		Logger.Debug().Msgf("LightGroup: %#v", lg)
		if err != nil {
			Logger.Error().Err(err).Msg("failed to fetch light group")
			return
		}
		for _, l := range lg.Lights {
			lightfunc(l)
		}
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		Logger.Debug().Msgf("Updating: %#v , %#v", lg, lg.Lights[0])
		_, err = p.light.UpdateLightGroup(ctx, lg)
		if err != nil {
			Logger.Error().Err(err).Msg("failed to update light group")
		}
	})
}
