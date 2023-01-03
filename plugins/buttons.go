package plugins

import (
	"embed"
	"errors"
	"fmt"
	"image"
	"os"
	"os/exec"

	"github.com/disintegration/gift"
	streamdeck "github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/actionhandlers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Embedded filesystem with a collection of available icons.
//
//go:embed icons
var icons embed.FS

var Logger zerolog.Logger = log.Logger.With().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})

type ImageButton struct {
	img           image.Image
	action        streamdeck.ButtonActionHandler
	Size          *image.Point
	btnIndex      int
	updateHandler func(streamdeck.Button)
	active        bool
}

func NewImageButtonFromFile(fname string) (*ImageButton, error) {
	im, err := NewImageFromFile(fname)
	if err != nil {
		return &ImageButton{}, fmt.Errorf("failed to load image file %w", err)
	}
	return NewImageButton(im), nil
}

func NewImageButton(im image.Image) *ImageButton {
	ib := &ImageButton{
		Size:          &image.Point{X: 72, Y: 72},
		active:        true,
		updateHandler: func(_ streamdeck.Button) {}, // noop update handler.
	}
	ib.SetImage(im)
	return ib
}

func (b *ImageButton) SetImage(i image.Image) {
	// If we *always* resize, and draw onto a new RGBA, we will always have an RGBA.
	resizer := gift.New(
		gift.Resize(b.Size.X, b.Size.Y, gift.LanczosResampling),
	)
	dst := image.NewRGBA(resizer.Bounds(i.Bounds()))
	resizer.Draw(dst, i)
	b.img = dst
	b.updateHandler(b)
}

// ButtonDisplay interface methods.
func (b *ImageButton) GetImageForButton(btnSize int) image.Image {
	resizer := gift.New(
		gift.Resize(b.Size.X, b.Size.Y, gift.LanczosResampling),
	)
	// Gray out inactive buttons.
	if !b.IsActive() {
		resizer.Add(gift.Contrast(-70))
	}
	dst := image.NewRGBA(resizer.Bounds(b.img.Bounds()))
	resizer.Draw(dst, b.img)
	return dst
}

func (b *ImageButton) GetButtonIndex() int { return b.btnIndex }
func (b *ImageButton) SetButtonIndex(idx int) {
	b.btnIndex = idx
	b.updateHandler(b)
}

func (b *ImageButton) RegisterUpdateHandler(uh func(streamdeck.Button)) {
	b.updateHandler = uh
}

func (b *ImageButton) Pressed() {
	if b.action == nil {
		return
	}
	if !b.IsActive() {
		Logger.Debug().Msgf("button %d pressed, but is inactive", b.btnIndex)
		return
	}
	b.action.Pressed(b)
}

// ActionButton interface
func (b *ImageButton) SetActionHandler(act streamdeck.ButtonActionHandler) {
	b.action = act
}

func (b *ImageButton) IsActive() bool {
	return b.active
}
func (b *ImageButton) SetActive(active bool) {
	if b.active != active {
		b.active = active
		b.updateHandler(b)
	}
}

// NewImageFromFile
func NewImageFromFile(fname string) (image.Image, error) {
	f, err := icons.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("could not open image file %s : %w", fname, err)
	}
	im, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("could not decode image file %s : %w", fname, err)
	}
	return im, nil
}

// NewExecAction creates an action handler that executes the given command, logging errors.
func NewExecAction(c string, args ...string) streamdeck.ButtonActionHandler {
	return actionhandlers.NewCustomAction(func(b streamdeck.Button) {
		cmd := exec.Command(c, args...)
		log.Debug().Msgf("Running: %s", cmd.String())
		output, err := cmd.CombinedOutput()
		if err != nil {
			exitErr := &exec.ExitError{}
			if errors.As(err, &exitErr) {
				log.Error().Err(err).Msgf("command exited %d: %s", exitErr.ExitCode(), exitErr.Stderr)
			}
			log.Error().Err(err).Msg(string(output))
		}
		log.Debug().Msg(string(output))
	})
}

type MultiStateButton struct {
	// map state name to a thing to display.
	States       map[string]streamdeck.Button
	InitialState streamdeck.Button
	State        string
}

// AddState creates an available state for this button.
func (b *MultiStateButton) AddState(state string, button streamdeck.Button) {
	b.States[state] = button
}

// SetState updates the current state of the button.
// The state must exist. This method is commonly called in a button action.
func (b *MultiStateButton) SetState(s string) error {
	return nil
}
