package redundancy

import (
	"sync"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
)

var log = logger.GetOrCreate("redundancy")

// maxRoundsOfInactivityAccepted defines the maximum rounds of inactivity accepted, after which the main or lower
// level redundancy machines will be considered inactive
const maxRoundsOfInactivityAccepted = 5

type nodeRedundancy struct {
	redundancyLevel     int64
	lastRoundIndexCheck int64
	roundsOfInactivity  uint64
	mutNodeRedundancy   sync.RWMutex
	messenger           P2PMessenger
}

// NewNodeRedundancy creates a node redundancy object which implements NodeRedundancyHandler interface
func NewNodeRedundancy(redundancyLevel int64, messenger P2PMessenger) (*nodeRedundancy, error) {
	if check.IfNil(messenger) {
		return nil, ErrNilMessenger
	}

	nr := &nodeRedundancy{
		redundancyLevel: redundancyLevel,
		messenger:       messenger,
	}

	return nr, nil
}

// IsRedundancyNode returns true if the current instance is used as a redundancy node
func (nr *nodeRedundancy) IsRedundancyNode() bool {
	return nr.redundancyLevel != 0
}

// IsMainMachineActive returns true if the main or lower level redundancy machines are active
func (nr *nodeRedundancy) IsMainMachineActive() bool {
	nr.mutNodeRedundancy.RLock()
	defer nr.mutNodeRedundancy.RUnlock()

	return nr.isMainMachineActive()
}

// AdjustInactivityIfNeeded increments rounds of inactivity for main or lower level redundancy machines if needed
func (nr *nodeRedundancy) AdjustInactivityIfNeeded(selfPubKey string, consensusPubKeys []string, roundIndex int64) {
	nr.mutNodeRedundancy.Lock()
	defer nr.mutNodeRedundancy.Unlock()

	if roundIndex <= nr.lastRoundIndexCheck {
		return
	}

	if nr.isMainMachineActive() {
		log.Debug("main or lower level redundancy machines are active", "node redundancy level", nr.redundancyLevel)
	} else {
		log.Warn("main or lower level redundancy machines are inactive", "node redundancy level", nr.redundancyLevel)
	}

	log.Debug("rounds of inactivity for main or lower level redundancy machines",
		"num", nr.roundsOfInactivity)

	for _, pubKey := range consensusPubKeys {
		if pubKey == selfPubKey {
			nr.roundsOfInactivity++
			break
		}
	}

	nr.lastRoundIndexCheck = roundIndex
}

// ResetInactivityIfNeeded resets rounds of inactivity for main or lower level redundancy machines if needed
func (nr *nodeRedundancy) ResetInactivityIfNeeded(selfPubKey string, consensusMsgPubKey string, consensusMsgPeerID core.PeerID) {
	if selfPubKey != consensusMsgPubKey {
		return
	}
	if consensusMsgPeerID == nr.messenger.ID() {
		return
	}

	nr.mutNodeRedundancy.Lock()
	nr.roundsOfInactivity = 0
	nr.mutNodeRedundancy.Unlock()
}

func (nr *nodeRedundancy) isMainMachineActive() bool {
	if nr.redundancyLevel < 0 {
		return true
	}

	return int64(nr.roundsOfInactivity) < maxRoundsOfInactivityAccepted*nr.redundancyLevel
}

// IsInterfaceNil returns true if there is no value under the interface
func (nr *nodeRedundancy) IsInterfaceNil() bool {
	return nr == nil
}
