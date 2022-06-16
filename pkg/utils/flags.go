package utils

import (
	"flag"

	"github.com/ntsd/cross-clipboard/pkg/config"
)

func ParseFlags() config.Config {
	c := config.Config{}

	flag.StringVar(&c.GroupName, "group", "default", "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&c.ProtocolID, "pid", "/cross-clipboard/0.0.1", "Sets a protocol id for stream headers")
	flag.StringVar(&c.ListenHost, "host", "0.0.0.0", "The bootstrap node host listen address\n")
	flag.IntVar(&c.ListenPort, "port", 4001, "node listen port")
	flag.IntVar(&c.MaxHistory, "max-history", 10, "max number of history clipboard")

	flag.Parse()
	return c
}
