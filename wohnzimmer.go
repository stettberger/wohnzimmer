package wohnzimmer // import "github.com/stettberger/wohnzimmer"

import (
	"crypto/tls"

	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleopenal"
)

type Wohnzimmer struct {
	Config *gumble.Config
	Client *gumble.Client

	Address   string
	TLSConfig tls.Config

	Stream *gumbleopenal.Stream

	reconnect chan bool
}
