package obswebsocket

import (
	"fmt"
	"image/color"
	"time"

	"github.com/christopher-dG/go-obs-websocket"
	"github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/actionhandlers"
	"github.com/magicmonkey/go-streamdeck/buttons"
	"github.com/magicmonkey/go-streamdeck/decorators"
	"github.com/muncus/my-streamdeck/plugins"

	_ "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var disabledButtonDecorator streamdeck.ButtonDecorator = decorators.NewBorder(15, color.RGBA{255, 0, 0, 150})

type OBSPlugin struct {
	d            *streamdeck.StreamDeck
	client       *obsws.Client
	ownedButtons []plugins.ActionButton // track all buttons we own, so they can be enabled/disabled
}

func New(d *streamdeck.StreamDeck) *OBSPlugin {
	plugin := &OBSPlugin{
		d: d,
		// TODO: make this settable
		client: &obsws.Client{
			Host: "localhost",
			Port: 4444,
		},
	}
	obsws.SetReceiveTimeout(5 * time.Second)
	go plugin.connect()
	return plugin
}

// connect to the obs websocket, and activate buttons.
func (p *OBSPlugin) connect() {
	p.client.Connect()
	fmt.Println("Connected.")
	p.setButtonsEnabled(p.client.Connected())
}

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

func (p *OBSPlugin) NewSceneButton(scenename string) plugins.ActionButton {
	btn := buttons.NewTextButton(scenename)
	btn.SetActionHandler(p.sceneChangeAction(scenename))
	return btn
}

func (p *OBSPlugin) sceneChangeAction(scene string) streamdeck.ButtonActionHandler {
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
