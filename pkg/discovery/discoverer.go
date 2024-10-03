package discovery

import (
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Discoverer interface {
	Init(host host.Host, serviceName string, logChan chan string) (chan peer.AddrInfo, error)
}
