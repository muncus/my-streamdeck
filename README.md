# My Streamdeck Setup.

This is still very much in development, but i'd like to use my Elgato Streamdeck
more effectively with my typical (non-streamer) work.

## Plugins

Plugins bundle related functionality, can hold state or resources used by multiple Buttons.

### Google Meet

The `googlemeet` package sends keypresses to the browser window that is using
google meet. It does this by using the `xdotool` command, which must be
installed for the plugin to work.

It can toggle the Mute, Video Mute, or "Raise Hand" functionality in a running Google Meet window.
Optionally, it will switch to the Google Meet window on button press (see config in `streamdeck.toml`).

### OBS

The `obswebsocket` plugin communicates with OBS (obsproject.com) through the
[websocket plugin for
OBS](https://obsproject.com/forum/resources/obs-websocket-remote-control-obs-studio-from-websockets.466/).

Plugin functionality is fairly limited at the moment, but useful for switching
scenes, which is my primary use case.

### Keylight

The Keylight plugin works with Elgato Key Light Air (and probably the non-Air
version?). It can turn the light on/off, and adjust brightness. Color
temperature will come eventually, but is not a major need for me.

## Progress Notes

- use of Decorators in the streamdeck library leads to image corruption on the
deck. this may suggest we should switch to either the low-level interface, or a
different library for streamdeck communications.

### Button Actions
* To toggle mute in Google Meet:
    `xdotool search -name "Meet *" key ctrl+d`
    * There's also a way to do this with obs-websockets, If I'm using the obs virtual cam.

## Misc
- teapot icons are from the Noun Project [tea emoticon set](https://thenounproject.com/aomam/collection/teapot-emoticons-line)
- Other icons come from the Google [Material Design Icon collection](https://fonts.google.com/icons)