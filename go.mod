module github.com/muncus/my-streamdeck

go 1.16

require (
	github.com/christopher-dG/go-obs-websocket v0.0.0-20200720193653-c4fed10356a5
	github.com/endocrimes/keylight-go v0.0.0-20201110202118-a45c372ed336
	github.com/magicmonkey/go-streamdeck v0.0.5-0.20210523153817-6f7e604ec5b2
	github.com/rs/zerolog v1.22.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// replace github.com/magicmonkey/go-streamdeck => ../go-streamdeck/
