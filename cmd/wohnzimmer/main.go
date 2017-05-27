package main // import "dokucode.de/wohnzimmer/cmd/wohnzimmer"

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"

	"github.com/stettberger/wohnzimmer"
	"layeh.com/gumble/gumble"
	_ "layeh.com/gumble/opus"
)

func main() {
	// Command line flags
	server := flag.String("server", "localhost:64738", "the server to connect to")
	username := flag.String("username", "", "the username of the client")
	password := flag.String("password", "", "the password of the server")
	insecure := flag.Bool("insecure", false, "skip server certificate verification")
	certificate := flag.String("certificate", "", "PEM encoded certificate and private key")

	flag.Parse()

	// Initialize
	w := wohnzimmer.Wohnzimmer{
		Config: gumble.NewConfig(),
		Address: *server,
	}

	w.Config.Username = *username
	w.Config.Password = *password

	if *insecure {
		w.TLSConfig.InsecureSkipVerify = true
	}
	if *certificate != "" {
		cert, err := tls.LoadX509KeyPair(*certificate, *certificate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		w.TLSConfig.Certificates = append(w.TLSConfig.Certificates, cert)
	}

	w.Start();

	keepAlive := make(chan bool)
	<-keepAlive
}
