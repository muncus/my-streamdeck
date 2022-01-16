package googlemeet

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/actionhandlers"
	"github.com/magicmonkey/go-streamdeck/buttons"
	"github.com/muncus/my-streamdeck/plugins"
	"github.com/pelletier/go-toml"
	"github.com/rs/zerolog/log"
)

// Plugin that uses xdotool to send keypresses to chrome windows with google meet in them.
type GoogleMeetPlugin struct {
	windowCommand   string
	MuteButton      plugins.ActionButton
	VideoMuteButton plugins.ActionButton
	RaiseHandButton plugins.ActionButton
}

func NewGoogleMeetPlugin(d *streamdeck.StreamDeck, config *toml.Tree) (*GoogleMeetPlugin, error) {
	var err error
	p := &GoogleMeetPlugin{}
	if config.GetDefault("switch_to_window", false).(bool) {
		p.windowCommand = "windowactivate"
	} else {
		p.windowCommand = "windowfocus"
	}
	p.MuteButton, err = buttons.NewImageFileButton("images/microphone_bg.png")
	if err != nil {
		return &GoogleMeetPlugin{}, fmt.Errorf("failed to create image button: %s", err)
	}
	// Typical Meet window titles are "Meet - xxx-xxx-xxx".
	// The title remains when you close the meeting, so this cmd will bookmark it.
	// The Meet "home page" is titled "Google Meet", and will not match this pattern.
	p.MuteButton.SetActionHandler(actionhandlers.NewCustomAction(func(streamdeck.Button) {
		cmd := exec.Command("xdotool", "search", "--name", "Meet - *", p.windowCommand, "key", "ctrl+d")
		commandAction(cmd)
	}))

	// Video Mute.
	p.VideoMuteButton, err = buttons.NewImageFileButton("images/camera_toggle_bg.png")
	if err != nil {
		return &GoogleMeetPlugin{}, fmt.Errorf("failed to create image button: %s", err)
	}
	p.VideoMuteButton.SetActionHandler(actionhandlers.NewCustomAction(func(streamdeck.Button) {
		cmd := exec.Command("xdotool", "search", "--name", "Meet - *", p.windowCommand, "key", "ctrl+e")
		commandAction(cmd)
	}))

	// Raise hand
	p.RaiseHandButton, err = buttons.NewImageFileButton("images/hand_transparent_bg.png")
	if err != nil {
		return &GoogleMeetPlugin{}, fmt.Errorf("failed to create image button: %s", err)
	}
	p.RaiseHandButton.SetActionHandler(actionhandlers.NewCustomAction(func(streamdeck.Button) {
		cmd := exec.Command("xdotool", "search", "--name", "Meet - *", p.windowCommand, "key", "ctrl+alt+h")
		commandAction(cmd)
	}))

	return p, nil
}

func commandAction(cmd *exec.Cmd) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		exitErr := &exec.ExitError{}
		if errors.As(err, &exitErr) {
			log.Error().Err(err).Msgf("command exited %d: %s", exitErr.ExitCode(), exitErr.Stderr)
		}
		log.Error().Err(err).Msg(string(output))
	}
}
