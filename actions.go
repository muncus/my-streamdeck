package main

import (
	"os/exec"

	"github.com/rs/zerolog/log"
)

// Toggle microphone mute in active Google Meet window.
// Relies on xdotool.
func ToggleMeetMute() {
	cmd := exec.Command("xdotool", "search", "--name", "Meet *", "windowfocus", "key", "ctrl+d")
	output, err := cmd.CombinedOutput()
	log.Debug().Msg(string(output))
	log.Debug().Msgf("%#v", err)
}
