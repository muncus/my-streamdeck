package googlemeet

import (
	"fmt"
	"os/exec"

	"github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/actionhandlers"
	"github.com/magicmonkey/go-streamdeck/buttons"
	"github.com/muncus/my-streamdeck/plugins"
	"github.com/rs/zerolog/log"
)

// Plugin that uses xdotool to send keypresses to chrome windows with google meet in them.
type GoogleMeetPlugin struct {
	MuteButton      plugins.ActionButton
	VideoMuteButton plugins.ActionButton
	RaiseHandButton plugins.ActionButton
}

func NewGoogleMeetPlugin(d *streamdeck.StreamDeck) (*GoogleMeetPlugin, error) {
	var err error
	p := &GoogleMeetPlugin{}
	p.MuteButton, err = buttons.NewImageFileButton("images/mic.png")
	if err != nil {
		return &GoogleMeetPlugin{}, fmt.Errorf("failed to create image button: %s", err)
	}
	// Typical Meet window titles are "Meet - xxx-xxx-xxx".
	// The title remains when you close the meeting, so this cmd will bookmark it.
	// The Meet "home page" is titled "Google Meet", and will not match this pattern.
	p.MuteButton.SetActionHandler(actionhandlers.NewCustomAction(func(streamdeck.Button) {
		cmd := exec.Command("xdotool", "search", "--name", "Meet - *", "windowfocus", "key", "ctrl+d")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Debug().Msg(string(output))
			log.Debug().Msgf("%#v", err)
		}
	}))

	// Video Mute.
	p.VideoMuteButton, err = buttons.NewImageFileButton("images/camera.png")
	if err != nil {
		return &GoogleMeetPlugin{}, fmt.Errorf("failed to create image button: %s", err)
	}
	p.VideoMuteButton.SetActionHandler(actionhandlers.NewCustomAction(func(streamdeck.Button) {
		cmd := exec.Command("xdotool", "search", "--name", "Meet - *", "windowfocus", "key", "ctrl+e")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Debug().Msg(string(output))
			log.Debug().Msgf("%#v", err)
		}
	}))

	// Raise hand
	p.RaiseHandButton, err = buttons.NewImageFileButton("images/hand.png")
	if err != nil {
		return &GoogleMeetPlugin{}, fmt.Errorf("failed to create image button: %s", err)
	}
	p.RaiseHandButton.SetActionHandler(actionhandlers.NewCustomAction(func(streamdeck.Button) {
		cmd := exec.Command("xdotool", "search", "--name", "Meet - *", "windowfocus", "key", "ctrl+alt+h")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Debug().Msg(string(output))
			log.Debug().Msgf("%#v", err)
		}
	}))

	return p, nil
}
