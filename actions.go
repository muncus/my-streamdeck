package main

import (
	"os/exec"

	"github.com/rs/zerolog/log"
)

// Toggle microphone mute in active Google Meet window.
// Relies on xdotool.
func ToggleMeetMute() {
	// Typical Meet window titles are "Meet - xxx-xxx-xxx".
	// The title remains when you close the meeting, so this cmd will bookmark it.
	// The Meet "home page" is titled "Google Meet", and will not match this pattern.
	cmd := exec.Command("xdotool", "search", "--name", "Meet - *", "windowfocus", "key", "ctrl+d")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Debug().Msg(string(output))
		log.Debug().Msgf("%#v", err)
	}
}
