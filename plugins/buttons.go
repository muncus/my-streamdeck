package plugins

import (
	"fmt"
	"image"
	"os"

	"github.com/disintegration/gift"
	streamdeck "github.com/magicmonkey/go-streamdeck"
)

type ImageButton struct {
	img      image.Image
	action   streamdeck.ButtonActionHandler
	Size     *image.Point
	btnIndex int
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
		Size: &image.Point{X: 72, Y: 72},
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
}

// ButtonDisplay interface methods.
func (b *ImageButton) GetImageForButton(btnSize int) image.Image {
	if btnSize == b.Size.X {
		return b.img
	}
	resizer := gift.New(
		gift.Resize(b.Size.X, b.Size.Y, gift.LanczosResampling),
	)
	dst := image.NewRGBA(resizer.Bounds(b.img.Bounds()))
	resizer.Draw(dst, b.img)
	return dst
}

func (b *ImageButton) GetButtonIndex() int    { return b.btnIndex }
func (b *ImageButton) SetButtonIndex(idx int) { b.btnIndex = idx }

func (b *ImageButton) RegisterUpdateHandler(func(streamdeck.Button)) {}
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
