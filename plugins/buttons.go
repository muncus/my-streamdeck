package plugins

import (
	"fmt"
	"image"
	"os"

	"github.com/disintegration/gift"
	streamdeck "github.com/magicmonkey/go-streamdeck"
)

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
	b.active = active
	b.updateHandler(b)
}

// NewImageFromFile
func NewImageFromFile(fname string) (image.Image, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("could not open image file %s : %w", fname, err)
	}
	im, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("could not decode image file %s : %w", fname, err)
	}
	return im, nil
}
