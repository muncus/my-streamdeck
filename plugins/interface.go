package plugins

import (
	"image/color"

	streamdeck "github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/buttons"
)

// ActionButton adds a SetActionHandler() method to the streamdeck.Button interface
// All existing button types satisfy this interface already.
type ActionButton interface {
	streamdeck.Button
	SetActionHandler(streamdeck.ButtonActionHandler)
}

// Black out all the buttons on the deck.
// Currently only works for original and smaller.
func ClearButtons(sd *streamdeck.StreamDeck) {
	nullbutton := buttons.NewColourButton(color.Black)
	// FIXME: this is hardcoded max button index.
	for i := 0; i < 15; i++ {
		sd.AddButton(i, nullbutton)
	}

}
