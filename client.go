package wohnzimmer // import "dokucode.de/wohnzimmer"

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/kennygrant/sanitize"
	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleopenal"
	"layeh.com/gumble/gumbleutil"
)

func (b *Wohnzimmer) Start() {
	b.Config.Attach(gumbleutil.AutoBitrate)
	b.Config.Attach(b)

	// Audio
	if os.Getenv("ALSOFT_LOGLEVEL") == "" {
		os.Setenv("ALSOFT_LOGLEVEL", "0")
	}

	b.reconnect = make(chan bool)

	for {
		go b.connect();
		<- b.reconnect;
		b.Print("Reconnect in 5 seconds...")
		time.Sleep(time.Second * 5)
	}
}


func (b *Wohnzimmer) connect() {
	// connect
	var err error
	_, err = gumble.DialWithDialer(new(net.Dialer), b.Address, b.Config, &b.TLSConfig)
	if err != nil {
		b.Print(fmt.Sprintf("dial failed: %s", err))
		b.reconnect <- true
		return
	}

	if stream, err := gumbleopenal.New(b.Client); err != nil {
		fmt.Fprintf(os.Stderr, "openal failed %s\n", err)
		b.reconnect <- true
		return
	} else {
		b.Stream = stream
	}

	b.Stream.StartSource()
}

func (b *Wohnzimmer) Print(line string) {
	now := time.Now()
	fmt.Printf("[%02d:%02d:%02d] %s\n", now.Hour(), now.Minute(), now.Second(), line)
}

func esc(str string) string {
	return sanitize.HTML(str)
}

func (b *Wohnzimmer) OnConnect(e *gumble.ConnectEvent) {
	b.Client = e.Client

	b.Print(fmt.Sprintf("To: %s", e.Client.Self.Channel.Name))
	b.Print(fmt.Sprintf("Connected to %s", b.Client.Conn.RemoteAddr()))
	if e.WelcomeMessage != nil {
		b.Print(fmt.Sprintf("Welcome message: %s", esc(*e.WelcomeMessage)))
	}
}

func (b *Wohnzimmer) OnDisconnect(e *gumble.DisconnectEvent) {
	var reason string
	switch e.Type {
	case gumble.DisconnectError:
		reason = "connection error"
	}
	if reason == "" {
		b.Print("Disconnected")
	} else {
		b.Print("Disconnected: " + reason)
	}

	b.reconnect <- true
}

func (b *Wohnzimmer) OnTextMessage(e *gumble.TextMessageEvent) {
	b.Print(fmt.Sprintf("<%s>: %s", e.Sender, e.Message))
}

func (b *Wohnzimmer) OnUserChange(e *gumble.UserChangeEvent) {
	// if e.Type.Has(gumble.UserChangeChannel) && e.User == b.Client.Self {
	// 	b.UpdateInputStatus(fmt.Sprintf("To: %s", e.User.Channel.Name))
	// }

}

func (b *Wohnzimmer) OnChannelChange(e *gumble.ChannelChangeEvent) {

}

func (b *Wohnzimmer) OnPermissionDenied(e *gumble.PermissionDeniedEvent) {
	var info string
	switch e.Type {
	case gumble.PermissionDeniedOther:
		info = e.String
	case gumble.PermissionDeniedPermission:
		info = "insufficient permissions"
	case gumble.PermissionDeniedSuperUser:
		info = "cannot modify SuperUser"
	case gumble.PermissionDeniedInvalidChannelName:
		info = "invalid channel name"
	case gumble.PermissionDeniedTextTooLong:
		info = "text too long"
	case gumble.PermissionDeniedTemporaryChannel:
		info = "temporary channel"
	case gumble.PermissionDeniedMissingCertificate:
		info = "missing certificate"
	case gumble.PermissionDeniedInvalidUserName:
		info = "invalid user name"
	case gumble.PermissionDeniedChannelFull:
		info = "channel full"
	case gumble.PermissionDeniedNestingLimit:
		info = "nesting limit"
	}
	b.Print(fmt.Sprintf("Permission denied: %s", info))
}

func (b *Wohnzimmer) OnUserList(e *gumble.UserListEvent) {
}

func (b *Wohnzimmer) OnACL(e *gumble.ACLEvent) {
}

func (b *Wohnzimmer) OnBanList(e *gumble.BanListEvent) {
}

func (b *Wohnzimmer) OnContextActionChange(e *gumble.ContextActionChangeEvent) {
}

func (b *Wohnzimmer) OnServerConfig(e *gumble.ServerConfigEvent) {
}
