package crossclipboard

import (
	"bufio"
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p/core/peer"
	"log"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multiaddr"
	"github.com/yqs112358/cross-clipboard/pkg/clipboard"
	"github.com/yqs112358/cross-clipboard/pkg/config"
	"github.com/yqs112358/cross-clipboard/pkg/crypto"
	"github.com/yqs112358/cross-clipboard/pkg/device"
	"github.com/yqs112358/cross-clipboard/pkg/devicemanager"
	"github.com/yqs112358/cross-clipboard/pkg/discovery"
	"github.com/yqs112358/cross-clipboard/pkg/stream"
	"github.com/yqs112358/cross-clipboard/pkg/xerror"
)

// CrossClipboard cross clipbaord struct
type CrossClipboard struct {
	Host   host.Host
	Config *config.Config

	ClipboardManager *clipboard.ClipboardManager
	DeviceManager    *devicemanager.DeviceManager

	streamHandler *stream.StreamHandler
	NewPeerChan   chan peer.AddrInfo

	LogChan   chan string
	ErrorChan chan error

	stopDiscovery chan struct{}
}

// NewCrossClipboard initial cross clipbaord
func NewCrossClipboard(cfg *config.Config) (*CrossClipboard, error) {
	cc := &CrossClipboard{
		Config:        cfg,
		NewPeerChan:   make(chan peer.AddrInfo),
		LogChan:       make(chan string),
		ErrorChan:     make(chan error),
		stopDiscovery: make(chan struct{}),
	}

	cc.ClipboardManager = clipboard.NewClipboardManager(cc.Config)
	cc.DeviceManager = devicemanager.NewDeviceManager(cc.Config)

	ctx := context.Background()

	// 0.0.0.0 will listen on any interface device.
	// TODO: change bad logic
	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cc.Config.Discovery.MDNS.ListenHost, cc.Config.Discovery.MDNS.ListenHost))
	if err != nil {
		return nil, xerror.NewFatalError("error to multiaddr.NewMultiaddr").Wrap(err)
	}

	// libp2p.New constructs a new libp2p Host.
	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(cc.Config.ID),
	)
	if err != nil {
		return nil, xerror.NewFatalError("error to libp2p.New").Wrap(err)
	}
	cc.Host = host

	pgpDecrypter, err := crypto.NewPGPDecrypter(cfg.PGPPrivateKey)
	if err != nil {
		return nil, xerror.NewFatalError("error to crypto.NewPGPDecrypter").Wrap(err)
	}

	go func() {
		err := cc.DeviceManager.Load()
		if err != nil {
			cc.ErrorChan <- xerror.NewFatalError("can not load device from setting").Wrap(err)
		}

		streamHandler := stream.NewStreamHandler(
			cc.Config,
			cc.ClipboardManager,
			cc.DeviceManager,
			cc.LogChan,
			cc.ErrorChan,
			pgpDecrypter,
		)
		cc.streamHandler = streamHandler

		// This function is called when a peer initiates a connection and starts a stream with this peer.
		cc.Host.SetStreamHandler(stream.PROTOCAL_ID, streamHandler.HandleStream)
		cc.LogChan <- fmt.Sprintf("[*] Your PeerID is: %s", host.ID().String())

		cc.startDiscoverers()
		cc.discoveryLoop(ctx)
	}()

	return cc, nil
}

func (cc *CrossClipboard) startDiscoverers() {
	mdnsDiscoverer := discovery.NewMdnsDiscoverer(cc.Config)
	err := mdnsDiscoverer.Init(cc.Host, cc.Config.GroupName, cc.NewPeerChan, cc.LogChan)
	if err != nil {
		cc.ErrorChan <- xerror.NewFatalError("error to discovery.InitMultiMDNS").Wrap(err)
	}
}

func (cc *CrossClipboard) discoveryLoop(ctx context.Context) {
	for {
		select {
		case peerInfo := <-cc.NewPeerChan: // when discover a peer
			dv := cc.DeviceManager.GetDevice(peerInfo.ID.String())
			if dv != nil && dv.Status == device.StatusBlocked {
				cc.ErrorChan <- xerror.NewRuntimeErrorf("device %s is blocked", peerInfo.ID.Loggable())
				continue
			}

			cc.LogChan <- fmt.Sprintf("connecting to peer: %s", peerInfo.ID.Loggable())

			retry := 1
			for ; retry < 5; retry++ { // retry to connect
				if err := cc.Host.Connect(ctx, peerInfo); err != nil {
					cc.ErrorChan <- xerror.NewRuntimeErrorf(
						"error to connect to peer %s, retrying %d",
						peerInfo.ID.Loggable(),
						retry,
					).Wrap(err)
					time.Sleep(time.Duration(retry*10) * time.Second)
					continue
				}
				break
			}
			if retry == 5 {
				cc.ErrorChan <- xerror.NewRuntimeErrorf("error to connect to peer %s", peerInfo.ID.Loggable())
				continue
			}

			// open a stream, this stream will be handled by handleStream other end
			stream, err := cc.Host.NewStream(ctx, peerInfo.ID, stream.PROTOCAL_ID)
			if err != nil {
				cc.ErrorChan <- xerror.NewRuntimeError("new stream error").Wrap(err)
				continue
			}

			if dv == nil {
				dv = device.NewDevice(peerInfo, stream)
			} else {
				dv.AddressInfo = peerInfo
				dv.Stream = stream
				dv.Reader = bufio.NewReader(stream)
				dv.Writer = bufio.NewWriter(stream)
			}

			cc.DeviceManager.UpdateDevice(dv)
			go cc.streamHandler.CreateReadData(dv.Reader, dv)

			cc.LogChan <- fmt.Sprintf("connected to peer host: %s", peerInfo)
		case <-cc.stopDiscovery: // when stop discovery
			cc.LogChan <- "stop discovery peer"
			return
		}
	}
}

func (cc *CrossClipboard) Stop() error {
	if cc.streamHandler != nil {
		for id, dv := range cc.DeviceManager.Devices {
			if dv.Status == device.StatusConnected {
				log.Printf("sending disconneced signal to peer %s \n", id)
				cc.streamHandler.SendSignal(dv, stream.SignalDisconnect)
			}
		}

		// sleep to wait sending disconnect signal
		time.Sleep(time.Second)

		for id, dv := range cc.DeviceManager.Devices {
			if dv.Status == device.StatusConnected {
				log.Printf("ending stream for peer %s \n", id)
				dv.Stream.Close()
			}
		}
	}

	cc.stopDiscovery <- struct{}{}

	err := cc.Host.Close()
	if err != nil {
		return xerror.NewFatalError("unable to close host").Wrap(err)
	}

	return nil
}
