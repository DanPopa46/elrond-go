package libp2p

import (
	"context"

<<<<<<< Updated upstream
	"github.com/ElrondNetwork/elrond-go/core/throttler"
=======
>>>>>>> Stashed changes
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/p2p/loadBalancer"
	"github.com/libp2p/go-libp2p-core/connmgr"
	libp2pCrypto "github.com/libp2p/go-libp2p-core/crypto"
<<<<<<< Updated upstream
=======
	"github.com/libp2p/go-libp2p-core/host"
>>>>>>> Stashed changes
	"github.com/libp2p/go-libp2p/p2p/net/mock"
)

// NewMemoryMessenger creates a new sandbox testable instance of libP2P messenger
// It should not open ports on current machine
// Should be used only in testing!
func NewMemoryMessenger(
	ctx context.Context,
	mockNet mocknet.Mocknet,
	peerDiscoverer p2p.PeerDiscoverer) (*networkMessenger, error) {

	if ctx == nil {
		return nil, p2p.ErrNilContext
	}
	if mockNet == nil {
		return nil, p2p.ErrNilMockNet
	}
	if peerDiscoverer == nil || peerDiscoverer.IsInterfaceNil() {
		return nil, p2p.ErrNilPeerDiscoverer
	}

	h, err := mockNet.GenPeer()
	if err != nil {
		return nil, err
	}

<<<<<<< Updated upstream
	lctx, err := NewLibp2pContext(ctx, NewConnectableHost(h))
=======
	lctx, err := NewLibp2pContext(ctx, NewConnectableHost(host.Host(h)))
>>>>>>> Stashed changes
	if err != nil {
		log.LogIfError(h.Close())
		return nil, err
	}

	mes, err := createMessenger(
		lctx,
		false,
		loadBalancer.NewOutgoingChannelLoadBalancer(),
		peerDiscoverer,
	)
	if err != nil {
		return nil, err
	}

<<<<<<< Updated upstream
	goRoutinesThrottler, err := throttler.NewNumGoRoutineThrottler(broadcastGoRoutines)
	if err != nil {
		log.LogIfError(h.Close())
		return nil, err
	}

	mes.goRoutinesThrottler = goRoutinesThrottler

=======
>>>>>>> Stashed changes
	return mes, err
}

// NewNetworkMessengerOnFreePort tries to create a new NetworkMessenger on a free port found in the system
// Should be used only in testing!
func NewNetworkMessengerOnFreePort(
	ctx context.Context,
	p2pPrivKey libp2pCrypto.PrivKey,
	conMgr connmgr.ConnManager,
	outgoingPLB p2p.ChannelLoadBalancer,
	peerDiscoverer p2p.PeerDiscoverer,
) (*networkMessenger, error) {
	return NewNetworkMessenger(
		ctx,
		0,
		p2pPrivKey,
		conMgr,
		outgoingPLB,
		peerDiscoverer,
		ListenLocalhostAddrWithIp4AndTcp,
	)
}
