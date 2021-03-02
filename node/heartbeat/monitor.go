package heartbeat

import (
<<<<<<< Updated upstream
	"bytes"
	"encoding/hex"
	"fmt"
=======
	"encoding/hex"
>>>>>>> Stashed changes
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/logger"
<<<<<<< Updated upstream
=======
	"github.com/ElrondNetwork/elrond-go/crypto"
>>>>>>> Stashed changes
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/statusHandler"
)

var log = logger.DefaultLogger()

// Monitor represents the heartbeat component that processes received heartbeat messages
type Monitor struct {
<<<<<<< Updated upstream
	maxDurationPeerUnresponsive time.Duration
	marshalizer                 marshal.Marshalizer
	mutHeartbeatMessages        sync.RWMutex
	heartbeatMessages           map[string]*heartbeatMessageInfo
	mutPubKeysMap               sync.RWMutex
	pubKeysMap                  map[uint32][]string
	mutFullPeersSlice           sync.RWMutex
	fullPeersSlice              [][]byte
	appStatusHandler            core.AppStatusHandler
	genesisTime                 time.Time
	messageHandler              MessageHandler
	storer                      HeartbeatStorageHandler
	timer                       Timer
=======
	singleSigner                crypto.SingleSigner
	maxDurationPeerUnresponsive time.Duration
	keygen                      crypto.KeyGenerator
	marshalizer                 marshal.Marshalizer
	heartbeatMessages           map[string]*heartbeatMessageInfo
	mutHeartbeatMessages        sync.RWMutex
	pubKeysMap                  map[uint32][]string
	mutPubKeysMap               sync.RWMutex
	appStatusHandler            core.AppStatusHandler
>>>>>>> Stashed changes
}

// NewMonitor returns a new monitor instance
func NewMonitor(
<<<<<<< Updated upstream
	marshalizer marshal.Marshalizer,
	maxDurationPeerUnresponsive time.Duration,
	pubKeysMap map[uint32][]string,
	genesisTime time.Time,
	messageHandler MessageHandler,
	storer HeartbeatStorageHandler,
	timer Timer,
) (*Monitor, error) {

=======
	singleSigner crypto.SingleSigner,
	keygen crypto.KeyGenerator,
	marshalizer marshal.Marshalizer,
	maxDurationPeerUnresponsive time.Duration,
	pubKeysMap map[uint32][]string,
) (*Monitor, error) {

	if singleSigner == nil || singleSigner.IsInterfaceNil() {
		return nil, ErrNilSingleSigner
	}
	if keygen == nil || keygen.IsInterfaceNil() {
		return nil, ErrNilKeyGenerator
	}
>>>>>>> Stashed changes
	if marshalizer == nil || marshalizer.IsInterfaceNil() {
		return nil, ErrNilMarshalizer
	}
	if len(pubKeysMap) == 0 {
		return nil, ErrEmptyPublicKeysMap
	}
<<<<<<< Updated upstream
	if messageHandler == nil || messageHandler.IsInterfaceNil() {
		return nil, ErrNilMessageHandler
	}
	if storer == nil || storer.IsInterfaceNil() {
		return nil, ErrNilHeartbeatStorer
	}
	if timer == nil || timer.IsInterfaceNil() {
		return nil, ErrNilTimer
	}

	mon := &Monitor{
=======

	pubKeysMapCopy := make(map[uint32][]string, 0)

	mon := &Monitor{
		singleSigner:                singleSigner,
		keygen:                      keygen,
>>>>>>> Stashed changes
		marshalizer:                 marshalizer,
		heartbeatMessages:           make(map[string]*heartbeatMessageInfo),
		maxDurationPeerUnresponsive: maxDurationPeerUnresponsive,
		appStatusHandler:            &statusHandler.NilStatusHandler{},
<<<<<<< Updated upstream
		genesisTime:                 genesisTime,
		messageHandler:              messageHandler,
		storer:                      storer,
		timer:                       timer,
	}

	err := mon.storer.UpdateGenesisTime(genesisTime)
	if err != nil {
		return nil, err
	}

	err = mon.initializeHeartbeatMessagesInfo(pubKeysMap)
	if err != nil {
		return nil, err
	}

	err = mon.loadRestOfPubKeysFromStorage()
	if err != nil {
		log.Debug(fmt.Sprintf("heartbeat can't load public keys from storage: %s", err.Error()))
	}

	return mon, nil
}

func (m *Monitor) initializeHeartbeatMessagesInfo(pubKeysMap map[uint32][]string) error {
	pubKeysMapCopy := make(map[uint32][]string, 0)
	for shardId, pubKeys := range pubKeysMap {
		for _, pubkey := range pubKeys {
			err := m.loadHbmiFromStorer(pubkey)
			if err != nil { // if pubKey not found in DB, create a new instance
				mhbi, errNewHbmi := newHeartbeatMessageInfo(m.maxDurationPeerUnresponsive, true, m.genesisTime, m.timer)
				if errNewHbmi != nil {
					return errNewHbmi
				}

				mhbi.genesisTime = m.genesisTime
				mhbi.computedShardID = shardId
				m.heartbeatMessages[pubkey] = mhbi
			}
			pubKeysMapCopy[shardId] = append(pubKeysMapCopy[shardId], pubkey)
		}
	}

	m.pubKeysMap = pubKeysMapCopy
	return nil
}

func (m *Monitor) loadRestOfPubKeysFromStorage() error {
	peersSlice, err := m.storer.LoadKeys()
	if err != nil {
		return err
	}

	for _, peer := range peersSlice {
		_, ok := m.heartbeatMessages[string(peer)]
		if !ok { // peer not in nodes map
			err = m.loadHbmiFromStorer(string(peer))
			if err != nil {
				continue
			}
		}
	}

	return nil
}

func (m *Monitor) loadHbmiFromStorer(pubKey string) error {
	hbmiDTO, err := m.storer.LoadHbmiDTO(pubKey)
	if err != nil {
		return err
	}

	receivedHbmi := m.convertFromExportedStruct(*hbmiDTO, m.maxDurationPeerUnresponsive)
	receivedHbmi.getTimeHandler = m.timer.Now
	receivedHbmi.isActive = m.timer.Now().Sub(receivedHbmi.lastUptimeDowntime) <= m.maxDurationPeerUnresponsive
	receivedHbmi.lastUptimeDowntime = m.timer.Now()
	receivedHbmi.genesisTime = m.genesisTime

	m.heartbeatMessages[pubKey] = &receivedHbmi

	return nil
=======
	}

	for shardId, pubKeys := range pubKeysMap {
		for _, pubkey := range pubKeys {
			pubKeysMapCopy[shardId] = append(pubKeysMapCopy[shardId], pubkey)
			mhbi, err := newHeartbeatMessageInfo(maxDurationPeerUnresponsive, true)
			if err != nil {
				return nil, err
			}

			mhbi.computedShardID = shardId
			mon.heartbeatMessages[pubkey] = mhbi
		}
	}
	mon.pubKeysMap = pubKeysMapCopy
	return mon, nil
>>>>>>> Stashed changes
}

// SetAppStatusHandler will set the AppStatusHandler which will be used for monitoring
func (m *Monitor) SetAppStatusHandler(ash core.AppStatusHandler) error {
	if ash == nil || ash.IsInterfaceNil() {
		return ErrNilAppStatusHandler
	}

	m.appStatusHandler = ash
	return nil
}

// ProcessReceivedMessage satisfies the p2p.MessageProcessor interface so it can be called
// by the p2p subsystem each time a new heartbeat message arrives
<<<<<<< Updated upstream
func (m *Monitor) ProcessReceivedMessage(message p2p.MessageP2P, _ func(buffToSend []byte)) error {
	hbRecv, err := m.messageHandler.CreateHeartbeatFromP2pMessage(message)
	if err != nil {
		return err
	}

	//message is validated, process should be done async, method can return nil
	go m.addHeartbeatMessageToMap(hbRecv)

	go m.computeAllHeartbeatMessages()

	return nil
}

func (m *Monitor) addHeartbeatMessageToMap(hb *Heartbeat) {
	pubKeyStr := string(hb.Pubkey)
	m.mutHeartbeatMessages.Lock()
	hbmi, ok := m.heartbeatMessages[pubKeyStr]
	if hbmi == nil || !ok {
		var err error
		hbmi, err = newHeartbeatMessageInfo(m.maxDurationPeerUnresponsive, false, m.genesisTime, m.timer)
		if err != nil {
			log.Error(err.Error())
			m.mutHeartbeatMessages.Unlock()
			return
		}
		m.heartbeatMessages[pubKeyStr] = hbmi
	}
	m.mutHeartbeatMessages.Unlock()

	computedShardID := m.computeShardID(pubKeyStr)

	hbmi.updateMutex.Lock()
	hbmi.HeartbeatReceived(computedShardID, hb.ShardID, hb.VersionNumber, hb.NodeDisplayName)
	hbDTO := m.convertToExportedStruct(hbmi)
	hbmi.updateMutex.Unlock()

	err := m.storer.SavePubkeyData(hb.Pubkey, &hbDTO)
	if err != nil {
		log.Error(fmt.Sprintf("cannot save heartbeat to db: %s", err.Error()))
	}
	m.addPeerToFullPeersSlice(hb.Pubkey)
}

func (m *Monitor) addPeerToFullPeersSlice(pubKey []byte) {
	m.mutFullPeersSlice.Lock()
	defer m.mutFullPeersSlice.Unlock()
	if !m.isPeerInFullPeersSlice(pubKey) {
		m.fullPeersSlice = append(m.fullPeersSlice, pubKey)
		err := m.storer.SaveKeys(m.fullPeersSlice)
		if err != nil {
			log.Error(fmt.Sprintf("can't store the keys slice: %s", err.Error()))
		}
	}
}

func (m *Monitor) isPeerInFullPeersSlice(pubKey []byte) bool {
	for _, peer := range m.fullPeersSlice {
		if bytes.Equal(peer, pubKey) {
			return true
		}
	}

	return false
=======
func (m *Monitor) ProcessReceivedMessage(message p2p.MessageP2P) error {
	if message == nil || message.IsInterfaceNil() {
		return ErrNilMessage
	}
	if message.Data() == nil {
		return ErrNilDataToProcess
	}

	hbRecv := &Heartbeat{}

	err := m.marshalizer.Unmarshal(hbRecv, message.Data())
	if err != nil {
		return err
	}

	err = m.verifySignature(hbRecv)
	if err != nil {
		return err
	}

	//message is validated, process should be done async, method can return nil
	go func(msg p2p.MessageP2P, hb *Heartbeat) {
		m.mutHeartbeatMessages.Lock()
		defer m.mutHeartbeatMessages.Unlock()

		pe := m.heartbeatMessages[string(hb.Pubkey)]
		if pe == nil {
			pe, err = newHeartbeatMessageInfo(m.maxDurationPeerUnresponsive, false)
			if err != nil {
				log.Error(err.Error())
				return
			}
			m.heartbeatMessages[string(hb.Pubkey)] = pe
		}

		computedShardID := m.computeShardID(string(hb.Pubkey))
		pe.HeartbeatReceived(computedShardID, hb.ShardID, hb.VersionNumber, hb.NodeDisplayName)
		m.updateAllHeartbeatMessages()
	}(message, hbRecv)

	return nil
>>>>>>> Stashed changes
}

func (m *Monitor) computeShardID(pubkey string) uint32 {
	// TODO : the shard ID will be recomputed at the end of an epoch / beginning of a new one.
	//  For the moment, just find the shard ID from a copy of the initial pub keys map
	m.mutPubKeysMap.RLock()
	defer m.mutPubKeysMap.RUnlock()
	for shardID, pubKeysSlice := range m.pubKeysMap {
		for _, pKey := range pubKeysSlice {
			if pKey == pubkey {
				return shardID
			}
		}
	}

	// if not found, return the latest known computed shard ID
	return m.heartbeatMessages[pubkey].computedShardID
}

<<<<<<< Updated upstream
func (m *Monitor) computeAllHeartbeatMessages() {
	m.mutHeartbeatMessages.Lock()
	counterActiveValidators := 0
	counterConnectedNodes := 0
	for _, v := range m.heartbeatMessages {
		v.computeActive(m.timer.Now())
=======
func (m *Monitor) verifySignature(hbRecv *Heartbeat) error {
	senderPubKey, err := m.keygen.PublicKeyFromByteArray(hbRecv.Pubkey)
	if err != nil {
		return err
	}

	copiedHeartbeat := *hbRecv
	copiedHeartbeat.Signature = nil
	buffCopiedHeartbeat, err := m.marshalizer.Marshal(copiedHeartbeat)
	if err != nil {
		return err
	}

	return m.singleSigner.Verify(senderPubKey, buffCopiedHeartbeat, hbRecv.Signature)
}

func (m *Monitor) updateAllHeartbeatMessages() {
	counterActiveValidators := 0
	counterConnectedNodes := 0
	for _, v := range m.heartbeatMessages {
		v.updateFields()

>>>>>>> Stashed changes
		if v.isActive {
			counterConnectedNodes++

			if v.isValidator {
				counterActiveValidators++
			}
		}
	}
<<<<<<< Updated upstream
	m.mutHeartbeatMessages.Unlock()
=======
>>>>>>> Stashed changes

	m.appStatusHandler.SetUInt64Value(core.MetricLiveValidatorNodes, uint64(counterActiveValidators))
	m.appStatusHandler.SetUInt64Value(core.MetricConnectedNodes, uint64(counterConnectedNodes))
}

// GetHeartbeats returns the heartbeat status
func (m *Monitor) GetHeartbeats() []PubKeyHeartbeat {
<<<<<<< Updated upstream
	m.computeAllHeartbeatMessages()

	m.mutHeartbeatMessages.Lock()
	status := make([]PubKeyHeartbeat, len(m.heartbeatMessages))
=======
	m.mutHeartbeatMessages.RLock()
	status := make([]PubKeyHeartbeat, len(m.heartbeatMessages))

>>>>>>> Stashed changes
	idx := 0
	for k, v := range m.heartbeatMessages {
		status[idx] = PubKeyHeartbeat{
			HexPublicKey:    hex.EncodeToString([]byte(k)),
			TimeStamp:       v.timeStamp,
			MaxInactiveTime: v.maxInactiveTime,
			IsActive:        v.isActive,
			ReceivedShardID: v.receivedShardID,
			ComputedShardID: v.computedShardID,
<<<<<<< Updated upstream
			TotalUpTime:     int(v.totalUpTime.Seconds()),
			TotalDownTime:   int(v.totalDownTime.Seconds()),
=======
			TotalUpTime:     v.totalUpTime,
			TotalDownTime:   v.totalDownTime,
>>>>>>> Stashed changes
			VersionNumber:   v.versionNumber,
			IsValidator:     v.isValidator,
			NodeDisplayName: v.nodeDisplayName,
		}
		idx++
<<<<<<< Updated upstream
	}
	m.mutHeartbeatMessages.Unlock()
=======

	}
	m.mutHeartbeatMessages.RUnlock()
>>>>>>> Stashed changes

	sort.Slice(status, func(i, j int) bool {
		return strings.Compare(status[i].HexPublicKey, status[j].HexPublicKey) < 0
	})

	return status
}

// IsInterfaceNil returns true if there is no value under the interface
func (m *Monitor) IsInterfaceNil() bool {
	if m == nil {
		return true
	}
	return false
}
<<<<<<< Updated upstream

func (m *Monitor) convertToExportedStruct(v *heartbeatMessageInfo) HeartbeatDTO {
	return HeartbeatDTO{
		TimeStamp:          v.timeStamp,
		MaxInactiveTime:    v.maxInactiveTime,
		IsActive:           v.isActive,
		ReceivedShardID:    v.receivedShardID,
		ComputedShardID:    v.computedShardID,
		TotalUpTime:        v.totalUpTime,
		TotalDownTime:      v.totalDownTime,
		VersionNumber:      v.versionNumber,
		IsValidator:        v.isValidator,
		NodeDisplayName:    v.nodeDisplayName,
		LastUptimeDowntime: v.lastUptimeDowntime,
		GenesisTime:        v.genesisTime,
	}
}

func (m *Monitor) convertFromExportedStruct(hbDTO HeartbeatDTO, maxDuration time.Duration) heartbeatMessageInfo {
	hbmi := heartbeatMessageInfo{
		maxDurationPeerUnresponsive: maxDuration,
		maxInactiveTime:             hbDTO.MaxInactiveTime,
		timeStamp:                   hbDTO.TimeStamp,
		isActive:                    hbDTO.IsActive,
		totalUpTime:                 hbDTO.TotalUpTime,
		totalDownTime:               hbDTO.TotalDownTime,
		receivedShardID:             hbDTO.ReceivedShardID,
		computedShardID:             hbDTO.ComputedShardID,
		versionNumber:               hbDTO.VersionNumber,
		nodeDisplayName:             hbDTO.NodeDisplayName,
		isValidator:                 hbDTO.IsValidator,
		lastUptimeDowntime:          hbDTO.LastUptimeDowntime,
		genesisTime:                 hbDTO.GenesisTime,
	}

	return hbmi
}
=======
>>>>>>> Stashed changes
