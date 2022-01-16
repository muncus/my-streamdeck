// package obswebsocket contains the plugin for communicating with OBS Studio, using the obs-websocket addon for OBS.
// The protocol for this socket is described at https://github.com/obsproject/obs-websocket/blob/master/docs/generated/protocol.md
package obswebsocket

import (
	"fmt"
	"image/color"
	"io"
	"os"
	"time"

	obsws "github.com/christopher-dG/go-obs-websocket"
	"github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/actionhandlers"
	"github.com/magicmonkey/go-streamdeck/buttons"
	"github.com/magicmonkey/go-streamdeck/decorators"
	"github.com/muncus/my-streamdeck/plugins"
	"github.com/pelletier/go-toml"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// A shared ButtonDecorator used to indicate that the button will not function.
var disabledButtonDecorator streamdeck.ButtonDecorator = decorators.NewBorder(15, color.RGBA{255, 0, 0, 150})
var Logger zerolog.Logger = log.Logger.With().Str("plugin", "obswebsocket").Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})

// OBSPluginConfig describes valid config options that can be specified for this plugin
type OBSPluginConfig struct {
	Host     string
	Port     int
	Password string
}

type OBSPlugin struct {
	d            *streamdeck.StreamDeck
	client       *obsws.Client
	ownedButtons []plugins.ActionButton // track all buttons we own, so they can be enabled/disabled
	ticker       *time.Ticker
	quitter      chan bool
}

// New creates a new instance of the OBS plugin, to display on the given streamdeck.
// config may contain fields from OBSPluginConfig.
func New(d *streamdeck.StreamDeck, config *toml.Tree) (*OBSPlugin, error) {
	// obsws logs are considered at "debug" level, so disable them if we've asked for no debug logs.
	// the standard logging library does not have a concept of log levels, so we set the output stream manually
	if e := log.Debug(); e.Enabled() {
		obsws.Logger.SetOutput(Logger.With().Str("level", "debug").Logger())
	} else {
		obsws.Logger.SetOutput(io.Discard)
	}
	configstruct := &OBSPluginConfig{}
	err := toml.Unmarshal([]byte(config.String()), configstruct)
	if err != nil {
		return &OBSPlugin{}, fmt.Errorf("failed to parse OBS config: %w", err)
	}
	plugin := &OBSPlugin{
		d: d,
		client: &obsws.Client{
			Host:     configstruct.Host,
			Port:     configstruct.Port,
			Password: configstruct.Password,
		},
		ticker:  time.NewTicker(5 * time.Second),
		quitter: make(chan bool),
	}
	obsws.SetReceiveTimeout(5 * time.Second)
	plugin.connect()
	go plugin.watchConnectionState()
	return plugin, nil
}

// watchConnectionState updates button appearance when we connect/disconnect from OBS
func (p *OBSPlugin) watchConnectionState() {
	for {
		select {
		case <-p.ticker.C:
			p.connect()
		case <-p.quitter:
			return
		}
	}
}

// connect to the obs websocket, and activate buttons.
func (p *OBSPlugin) connect() {
	p.setButtonsEnabled(p.client.Connected())
	if !p.client.Connected() {
		p.client.Connect()
		p.setButtonsEnabled(p.client.Connected())
	}
}

// setButtonsEnabled marks buttons as disabled when not connected to an OBS instance.
// FIXME: decorator use is currently disabled, as it causes some visual artifacts.
func (p *OBSPlugin) setButtonsEnabled(enabled bool) {
	if !enabled {
		for _, b := range p.ownedButtons {
			_ = b
			// p.d.SetDecorator(b.GetButtonIndex(), disabledButtonDecorator)
		}
	} else {
		for _, b := range p.ownedButtons {
			_ = b
			// p.d.UnsetDecorator(b.GetButtonIndex())
		}
	}
}

// NewSceneButton creates a button that will change to the named scene when pressed
// the button appearance will be the name of the desired scene.
func (p *OBSPlugin) NewSceneButton(scenename string) plugins.ActionButton {
	btn := buttons.NewTextButton(scenename)
	btn.SetActionHandler(p.NewSceneChangeAction(scenename))
	p.ownedButtons = append(p.ownedButtons, btn)
	return btn
}

// NewSceneChangeAction returns a handler that changes scene to the named scene in OBS.
func (p *OBSPlugin) NewSceneChangeAction(scene string) streamdeck.ButtonActionHandler {
	a := actionhandlers.NewCustomAction(func(streamdeck.Button) {
		req := obsws.NewSetCurrentSceneRequest(scene)
		resp, err := req.SendReceive(*p.client)
		if err != nil {
			Logger.Warn().Err(err)
			return
		}
		Logger.Info().Msg(resp.Status())
	})
	return a
}

// ManageButton lets this plugin decorate the given button on connect/disconnect.
// It is intended for buttons whose actions depend on OBSPlugin being connected to OBS, but were not constructed with New*Button methods.
func (p *OBSPlugin) ManageButton(b plugins.ActionButton) {
	p.ownedButtons = append(p.ownedButtons, b)
}

// Close should be called when exiting. it cleans up background goroutines.
func (p *OBSPlugin) Close() {
	close(p.quitter)
}
