package discovery

import (
	"fmt"
	"github.com/yqs112358/cross-clipboard/pkg/config"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// DiscoveryNotifee noti struct when discover a new peer
type DiscoveryNotifee struct {
	PeerHost host.Host
	PeerChan chan peer.AddrInfo
	LogChan  chan string
}

// HandlePeerFound interface to be called when new  peer is found
func (n *DiscoveryNotifee) HandlePeerFound(peerInfo peer.AddrInfo) {
	n.LogChan <- fmt.Sprintf("discovered peer: %s", peerInfo)
	if n.PeerHost.ID() != peerInfo.ID {
		n.PeerChan <- peerInfo
	}
}

type MulticastDNS struct {
	cfg *config.Config
}

func NewMdnsDiscoverer(c *config.Config) *MulticastDNS {
	return &MulticastDNS{cfg: c}
}

func (m *MulticastDNS) Init(peerHost host.Host, serviceName string, peerChan chan peer.AddrInfo, logChan chan string) error {
	// register with service so that we get notified about peer discovery
	n := &DiscoveryNotifee{
		PeerHost: peerHost,
		PeerChan: peerChan,
		LogChan:  logChan,
	}

	// An hour might be a long long period in practical applications. But this is fine for us
	ser := mdns.NewMdnsService(peerHost, serviceName, n)
	if err := ser.Start(); err != nil {
		return err
	}

	return nil
}
