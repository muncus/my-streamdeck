package plugins

import streamdeck "github.com/magicmonkey/go-streamdeck"

// ActionButton adds a SetActionHandler() method to the streamdeck.Button interface
// All existing button types satisfy this interface already.
type ActionButton interface {
	streamdeck.Button
	SetActionHandler(streamdeck.ButtonActionHandler)
}
