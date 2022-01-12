// package obswebsocket contains the plugin for communicating with OBS Studio, using the obs-websocket addon for OBS.
// The protocol for this socket is described at https://github.com/obsproject/obs-websocket/blob/master/docs/generated/protocol.md
package obswebsocket

import (
	"fmt"
	"image/color"
	"time"

	obsws "github.com/christopher-dG/go-obs-websocket"
	"github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/actionhandlers"
	"github.com/magicmonkey/go-streamdeck/buttons"
	"github.com/magicmonkey/go-streamdeck/decorators"
	"github.com/muncus/my-streamdeck/plugins"
	"github.com/pelletier/go-toml"

	_ "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var disabledButtonDecorator streamdeck.ButtonDecorator = decorators.NewBorder(15, color.RGBA{255, 0, 0, 150})

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
}

func New(d *streamdeck.StreamDeck, config *toml.Tree) (*OBSPlugin, error) {
	configstruct := &OBSPluginConfig{}
	err := toml.Unmarshal([]byte(config.Get("obswebsocket").(*toml.Tree).String()), configstruct)
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
	}
	obsws.SetReceiveTimeout(5 * time.Second)
	go plugin.connect()
	return plugin, nil
}

// connect to the obs websocket, and activate buttons.
func (p *OBSPlugin) connect() {
	p.client.Connect()
	fmt.Println("Connected.")
	p.setButtonsEnabled(p.client.Connected())
}

// TODO: use a mechanic like this to disable OBS buttons when we're not connected to obs.
func (p *OBSPlugin) setButtonsEnabled(enabled bool) {
	if enabled {
		for _, b := range p.ownedButtons {
			p.d.SetDecorator(b.GetButtonIndex(), disabledButtonDecorator)
		}
	} else {
		for _, b := range p.ownedButtons {
			p.d.UnsetDecorator(b.GetButtonIndex())
		}
	}
}

// NewSceneButton creates a button that will change to the named scene when pressed
// the button appearance will be the name of the desired scene.
func (p *OBSPlugin) NewSceneButton(scenename string) plugins.ActionButton {
	btn := buttons.NewTextButton(scenename)
	btn.SetActionHandler(p.NewSceneChangeAction(scenename))
	return btn
}

// NewSceneChangeAction returns a handler that changes scene to the named scene in OBS.
func (p *OBSPlugin) NewSceneChangeAction(scene string) streamdeck.ButtonActionHandler {
	a := actionhandlers.NewCustomAction(func(streamdeck.Button) {
		req := obsws.NewSetCurrentSceneRequest(scene)
		resp, err := req.SendReceive(*p.client)
		if err != nil {
			log.Warn().Err(err)
			return
		}
		log.Info().Msg(resp.Status())
	})
	return a
}
