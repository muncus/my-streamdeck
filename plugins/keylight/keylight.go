// package keylight contains methods for interacting with Elgato Key Lights from the stream deck.
package keylight

import (
	"context"
	"os"
	"time"

	kl "github.com/endocrimes/keylight-go"
	"github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/actionhandlers"
	"github.com/muncus/my-streamdeck/plugins"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger = log.Logger.With().Str("plugin", "keylight").Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})

type KeylightPlugin struct {
	light           *kl.Device
	quitter         chan (bool)
	ticker          time.Ticker
	brightIncrement int
	tempIncrement   int
	PowerToggle     *plugins.ImageButton
	BrightnessInc   *plugins.ImageButton
	BrightnessDec   *plugins.ImageButton
}

func New() *KeylightPlugin {
	plugin := &KeylightPlugin{
		brightIncrement: 10,
		tempIncrement:   1,
		ticker:          *time.NewTicker(5 * time.Second),
		quitter:         make(chan bool),
	}

	// plugin.PowerToggle = buttons.NewTextButton("aziz, light!")
	plugin.PowerToggle, _ = plugins.NewImageButtonFromFile("icons/light_mode_bg.png")
	plugin.PowerToggle.SetActionHandler(plugin.LightAction(func(l *kl.Light) {
		if l.On > 0 {
			l.On = 0
		} else {
			l.On = 1
		}
	}))

	// TODO: add a min/max to these values.
	plugin.BrightnessInc, _ = plugins.NewImageButtonFromFile("icons/more_bright_bg.png")
	plugin.BrightnessInc.SetActionHandler(plugin.LightAction(func(l *kl.Light) {
		l.Brightness = l.Brightness + plugin.brightIncrement
	}))
	plugin.BrightnessDec, _ = plugins.NewImageButtonFromFile("icons/less_bright_bg.png")
	plugin.BrightnessDec.SetActionHandler(plugin.LightAction(func(l *kl.Light) {
		l.Brightness = l.Brightness - plugin.brightIncrement
	}))

	plugin.setButtonsEnabled(false)
	// run discovery in a background goroutine.
	go plugin.watchButtonState()

	return plugin
}

// watchButtonState runs forever, updating owned buttons with a disabled decorator.
func (p *KeylightPlugin) watchButtonState() {
	for {
		select {
		case <-p.ticker.C:
			if p.light == nil {
				// no lights found yet. disable.
				p.setButtonsEnabled(false)
				p.discover()
			} else {
				p.setButtonsEnabled(true)
			}
		case <-p.quitter:
			return
		}
	}
}

func (p *KeylightPlugin) setButtonsEnabled(enabled bool) {
	p.PowerToggle.SetActive(enabled)
	p.BrightnessInc.SetActive(enabled)
	p.BrightnessDec.SetActive(enabled)
}

func (p *KeylightPlugin) discover() error {
	Logger.Debug().Msg("starting discovery")
	disc, err := kl.NewDiscovery()
	if err != nil {
		return err
	}
	dctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = disc.Run(dctx)
	if err != nil {
		return err
	}
	// NB: this only finds one light.
	p.light = <-disc.ResultsCh()
	Logger.Debug().Msgf("found: %s", p.light)
	return nil
}

// LightAction helps construct button actions that operate onall lights in a given LightGroup
func (p *KeylightPlugin) LightAction(lightfunc func(l *kl.Light)) streamdeck.ButtonActionHandler {
	return actionhandlers.NewCustomAction(func(streamdeck.Button) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		Logger.Debug().Msgf("Updating: %#v , %#v", lg, lg.Lights[0])
		_, err = p.light.UpdateLightGroup(ctx, lg)
		if err != nil {
			Logger.Error().Err(err).Msg("failed to update light group")
		}
	})
}
