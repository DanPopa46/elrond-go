package heartbeat_test

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"testing"
	"time"

<<<<<<< Updated upstream
	"github.com/ElrondNetwork/elrond-go/node/heartbeat"
	"github.com/ElrondNetwork/elrond-go/node/heartbeat/storage"
	"github.com/ElrondNetwork/elrond-go/node/mock"
	"github.com/ElrondNetwork/elrond-go/p2p"
=======
	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/node/heartbeat"
	"github.com/ElrondNetwork/elrond-go/node/mock"
>>>>>>> Stashed changes
	"github.com/stretchr/testify/assert"
)

//------- NewMonitor

<<<<<<< Updated upstream
func TestNewMonitor_NilMarshalizerShouldErr(t *testing.T) {
	t.Parallel()

	th := &mock.MockTimer{}
	mon, err := heartbeat.NewMonitor(
		nil,
		0,
		map[uint32][]string{0: {""}},
		time.Now(),
		&mock.MessageHandlerStub{},
		&mock.HeartbeatStorerStub{},
		th,
	)

	assert.Nil(t, mon)
	assert.Equal(t, heartbeat.ErrNilMarshalizer, err)
}

func TestNewMonitor_EmptyPublicKeyListShouldErr(t *testing.T) {
	t.Parallel()

	th := &mock.MockTimer{}
	mon, err := heartbeat.NewMonitor(
		&mock.MarshalizerMock{},
		0,
		make(map[uint32][]string),
		time.Now(),
		&mock.MessageHandlerStub{},
		&mock.HeartbeatStorerStub{},
		th,
	)

	assert.Nil(t, mon)
	assert.Equal(t, heartbeat.ErrEmptyPublicKeysMap, err)
}

func TestNewMonitor_NilMessageHandlerShouldErr(t *testing.T) {
	t.Parallel()

	th := &mock.MockTimer{}
	mon, err := heartbeat.NewMonitor(
		&mock.MarshalizerMock{},
		0,
		map[uint32][]string{0: {""}},
		time.Now(),
		nil,
		&mock.HeartbeatStorerStub{},
		th,
	)

	assert.Nil(t, mon)
	assert.Equal(t, heartbeat.ErrNilMessageHandler, err)
}

func TestNewMonitor_NilHeartbeatStorerShouldErr(t *testing.T) {
	t.Parallel()

	th := &mock.MockTimer{}
	mon, err := heartbeat.NewMonitor(
		&mock.MarshalizerMock{},
		0,
		map[uint32][]string{0: {""}},
		time.Now(),
		&mock.MessageHandlerStub{},
		nil,
		th,
	)

	assert.Nil(t, mon)
	assert.Equal(t, heartbeat.ErrNilHeartbeatStorer, err)
}

func TestNewMonitor_NilTimeHandlerShouldErr(t *testing.T) {
	t.Parallel()

	mon, err := heartbeat.NewMonitor(
		&mock.MarshalizerMock{},
		0,
		map[uint32][]string{0: {""}},
		time.Now(),
		&mock.MessageHandlerStub{},
		&mock.HeartbeatStorerStub{},
		nil,
	)

	assert.Nil(t, mon)
	assert.Equal(t, heartbeat.ErrNilTimer, err)
=======
func TestNewMonitor_NilSingleSignerShouldErr(t *testing.T) {
	t.Parallel()

	mon, err := heartbeat.NewMonitor(
		nil,
		&mock.KeyGenMock{},
		&mock.MarshalizerMock{},
		0,
		map[uint32][]string{0: {""}},
	)

	assert.Nil(t, mon)
	assert.Equal(t, heartbeat.ErrNilSingleSigner, err)
}

func TestNewMonitor_NilKeygenShouldErr(t *testing.T) {
	t.Parallel()

	mon, err := heartbeat.NewMonitor(
		&mock.SinglesignMock{},
		nil,
		&mock.MarshalizerMock{},
		0,
		map[uint32][]string{0: {""}},
	)

	assert.Nil(t, mon)
	assert.Equal(t, heartbeat.ErrNilKeyGenerator, err)
}

func TestNewMonitor_NilMarshalizerShouldErr(t *testing.T) {
	t.Parallel()

	mon, err := heartbeat.NewMonitor(
		&mock.SinglesignMock{},
		&mock.KeyGenMock{},
		nil,
		0,
		map[uint32][]string{0: {""}},
	)

	assert.Nil(t, mon)
	assert.Equal(t, heartbeat.ErrNilMarshalizer, err)
}

func TestNewMonitor_EmptyPublicKeyListShouldErr(t *testing.T) {
	t.Parallel()

	mon, err := heartbeat.NewMonitor(
		&mock.SinglesignMock{},
		&mock.KeyGenMock{},
		&mock.MarshalizerMock{},
		0,
		make(map[uint32][]string),
	)

	assert.Nil(t, mon)
	assert.Equal(t, heartbeat.ErrEmptyPublicKeysMap, err)
>>>>>>> Stashed changes
}

func TestNewMonitor_OkValsShouldCreatePubkeyMap(t *testing.T) {
	t.Parallel()

<<<<<<< Updated upstream
	th := &mock.MockTimer{}
	mon, err := heartbeat.NewMonitor(
		&mock.MarshalizerMock{},
		1,
		map[uint32][]string{0: {"pk1", "pk2"}},
		time.Now(),
		&mock.MessageHandlerStub{},
		&mock.HeartbeatStorerStub{
			UpdateGenesisTimeCalled: func(genesisTime time.Time) error {
				return nil
			},
			LoadHbmiDTOCalled: func(pubKey string) (*heartbeat.HeartbeatDTO, error) {
				return nil, errors.New("not found")
			},
			LoadKeysCalled: func() ([][]byte, error) {
				return nil, nil
			},
			SavePubkeyDataCalled: func(pubkey []byte, heartbeat *heartbeat.HeartbeatDTO) error {
				return nil
			},
		},
		th,
=======
	mon, err := heartbeat.NewMonitor(
		&mock.SinglesignMock{},
		&mock.KeyGenMock{},
		&mock.MarshalizerMock{},
		1,
		map[uint32][]string{0: {"pk1", "pk2"}},
>>>>>>> Stashed changes
	)

	assert.NotNil(t, mon)
	assert.Nil(t, err)
	hbStatus := mon.GetHeartbeats()
	assert.Equal(t, 2, len(hbStatus))
}

func TestNewMonitor_ShouldComputeShardId(t *testing.T) {
	t.Parallel()

<<<<<<< Updated upstream
	th := &mock.MockTimer{}
=======
>>>>>>> Stashed changes
	pksPerShards := map[uint32][]string{
		0: {"pk0"},
		1: {"pk1"},
	}

	maxDuration := time.Millisecond
	mon, err := heartbeat.NewMonitor(
<<<<<<< Updated upstream
		&mock.MarshalizerMock{},
		maxDuration,
		pksPerShards,
		time.Now(),
		&mock.MessageHandlerStub{},
		&mock.HeartbeatStorerStub{
			UpdateGenesisTimeCalled: func(genesisTime time.Time) error {
				return nil
			},
			LoadHbmiDTOCalled: func(pubKey string) (*heartbeat.HeartbeatDTO, error) {
				return nil, errors.New("not found")
			},
			LoadKeysCalled: func() ([][]byte, error) {
				return nil, nil
			},
			SavePubkeyDataCalled: func(pubkey []byte, heartbeat *heartbeat.HeartbeatDTO) error {
				return nil
			},
		},
		th,
=======
		&mock.SinglesignMock{},
		&mock.KeyGenMock{},
		&mock.MarshalizerMock{},
		maxDuration,
		pksPerShards,
>>>>>>> Stashed changes
	)

	assert.NotNil(t, mon)
	assert.Nil(t, err)
	hbStatus := mon.GetHeartbeats()

	assert.Equal(t, uint32(0), hbStatus[0].ComputedShardID)
	assert.Equal(t, uint32(1), hbStatus[1].ComputedShardID)
}

//------- ProcessReceivedMessage

<<<<<<< Updated upstream
func TestMonitor_ProcessReceivedMessageShouldWork(t *testing.T) {
	t.Parallel()

	pubKey := "pk1"

	th := &mock.MockTimer{}
	mon, _ := heartbeat.NewMonitor(
		&mock.MarshalizerMock{
			UnmarshalHandler: func(obj interface{}, buff []byte) error {
				(obj.(*heartbeat.Heartbeat)).Pubkey = []byte(pubKey)
				return nil
			},
		},
		time.Second*1000,
		map[uint32][]string{0: {pubKey}},
		time.Now(),
		&mock.MessageHandlerStub{
			CreateHeartbeatFromP2pMessageCalled: func(message p2p.MessageP2P) (*heartbeat.Heartbeat, error) {
				var rcvHb heartbeat.Heartbeat
				_ = json.Unmarshal(message.Data(), &rcvHb)
				return &rcvHb, nil
			},
		},
		&mock.HeartbeatStorerStub{
			UpdateGenesisTimeCalled: func(genesisTime time.Time) error {
				return nil
			},
			LoadHbmiDTOCalled: func(pubKey string) (*heartbeat.HeartbeatDTO, error) {
				return nil, errors.New("not found")
			},
			LoadKeysCalled: func() ([][]byte, error) {
				return nil, nil
			},
			SavePubkeyDataCalled: func(pubkey []byte, heartbeat *heartbeat.HeartbeatDTO) error {
				return nil
			},
			SaveKeysCalled: func(peersSlice [][]byte) error {
				return nil
			},
		},
		th,
	)

	hb := heartbeat.Heartbeat{
		Pubkey: []byte(pubKey),
	}
	hbBytes, _ := json.Marshal(hb)
	err := mon.ProcessReceivedMessage(&mock.P2PMessageStub{DataField: hbBytes}, nil)
=======
func TestMonitor_ProcessReceivedMessageNilMessageShouldErr(t *testing.T) {
	t.Parallel()

	mon, _ := heartbeat.NewMonitor(
		&mock.SinglesignMock{},
		&mock.KeyGenMock{},
		&mock.MarshalizerMock{},
		1,
		map[uint32][]string{0: {"pk1"}},
	)

	err := mon.ProcessReceivedMessage(nil)

	assert.Equal(t, heartbeat.ErrNilMessage, err)
}

func TestMonitor_ProcessReceivedMessageNilDataShouldErr(t *testing.T) {
	t.Parallel()

	mon, _ := heartbeat.NewMonitor(
		&mock.SinglesignMock{},
		&mock.KeyGenMock{},
		&mock.MarshalizerMock{},
		0,
		map[uint32][]string{0: {"pk1"}},
	)

	err := mon.ProcessReceivedMessage(&mock.P2PMessageStub{})

	assert.Equal(t, heartbeat.ErrNilDataToProcess, err)
}

func TestMonitor_ProcessReceivedMessageMarshalFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected err")

	mon, _ := heartbeat.NewMonitor(
		&mock.SinglesignMock{},
		&mock.KeyGenMock{},
		&mock.MarshalizerMock{
			UnmarshalHandler: func(obj interface{}, buff []byte) error {
				return errExpected
			},
		},
		1,
		map[uint32][]string{0: {"pk1"}},
	)

	err := mon.ProcessReceivedMessage(&mock.P2PMessageStub{DataField: []byte("")})

	assert.Equal(t, errExpected, err)
}

func TestMonitor_ProcessReceivedMessageWrongPubkeyShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected err")

	mon, _ := heartbeat.NewMonitor(
		&mock.SinglesignMock{},
		&mock.KeyGenMock{
			PublicKeyFromByteArrayMock: func(b []byte) (key crypto.PublicKey, e error) {
				return nil, errExpected
			},
		},
		&mock.MarshalizerMock{
			UnmarshalHandler: func(obj interface{}, buff []byte) error {
				return nil
			},
		},
		1,
		map[uint32][]string{0: {"pk1"}},
	)

	err := mon.ProcessReceivedMessage(&mock.P2PMessageStub{DataField: []byte("")})

	assert.Equal(t, errExpected, err)
}

func TestMonitor_ProcessReceivedMessageVerifyFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected err")

	mon, _ := heartbeat.NewMonitor(
		&mock.SinglesignStub{
			VerifyCalled: func(public crypto.PublicKey, msg []byte, sig []byte) error {
				return errExpected
			},
		},
		&mock.KeyGenMock{
			PublicKeyFromByteArrayMock: func(b []byte) (key crypto.PublicKey, e error) {
				return nil, nil
			},
		},
		&mock.MarshalizerMock{
			UnmarshalHandler: func(obj interface{}, buff []byte) error {
				return nil
			},
		},
		1,
		map[uint32][]string{0: {"pk1"}},
	)

	err := mon.ProcessReceivedMessage(&mock.P2PMessageStub{DataField: []byte("")})

	assert.Equal(t, errExpected, err)
}

func TestMonitor_ProcessReceivedMessageShouldWork(t *testing.T) {
	t.Parallel()

	pubKey := "pk1"

	mon, _ := heartbeat.NewMonitor(
		&mock.SinglesignStub{
			VerifyCalled: func(public crypto.PublicKey, msg []byte, sig []byte) error {
				return nil
			},
		},
		&mock.KeyGenMock{
			PublicKeyFromByteArrayMock: func(b []byte) (key crypto.PublicKey, e error) {
				return nil, nil
			},
		},
		&mock.MarshalizerMock{
			UnmarshalHandler: func(obj interface{}, buff []byte) error {
				(obj.(*heartbeat.Heartbeat)).Pubkey = []byte(pubKey)
				return nil
			},
		},
		time.Second*1000,
		map[uint32][]string{0: {pubKey}},
	)

	err := mon.ProcessReceivedMessage(&mock.P2PMessageStub{DataField: []byte("")})
>>>>>>> Stashed changes
	assert.Nil(t, err)

	//a delay is mandatory for the go routine to finish its job
	time.Sleep(time.Second)

	hbStatus := mon.GetHeartbeats()
	assert.Equal(t, 1, len(hbStatus))
	assert.Equal(t, hex.EncodeToString([]byte(pubKey)), hbStatus[0].HexPublicKey)
}

func TestMonitor_ProcessReceivedMessageWithNewPublicKey(t *testing.T) {
	t.Parallel()

	pubKey := "pk1"

<<<<<<< Updated upstream
	th := &mock.MockTimer{}
	mon, _ := heartbeat.NewMonitor(
		&mock.MarshalizerMock{
			UnmarshalHandler: func(obj interface{}, buff []byte) error {
				(obj.(*heartbeat.Heartbeat)).Pubkey = []byte(pubKey)
				return nil
			},
		},
		time.Second*1000,
		map[uint32][]string{0: {"pk2"}},
		time.Now(),
		&mock.MessageHandlerStub{
			CreateHeartbeatFromP2pMessageCalled: func(message p2p.MessageP2P) (*heartbeat.Heartbeat, error) {
				var rcvHb heartbeat.Heartbeat
				_ = json.Unmarshal(message.Data(), &rcvHb)
				return &rcvHb, nil
			},
		},
		&mock.HeartbeatStorerStub{
			UpdateGenesisTimeCalled: func(genesisTime time.Time) error {
				return nil
			},
			LoadHbmiDTOCalled: func(pubKey string) (*heartbeat.HeartbeatDTO, error) {
				return nil, errors.New("not found")
			},
			LoadKeysCalled: func() ([][]byte, error) {
				return nil, nil
			},
			SavePubkeyDataCalled: func(pubkey []byte, heartbeat *heartbeat.HeartbeatDTO) error {
				return nil
			},
			SaveKeysCalled: func(peersSlice [][]byte) error {
				return nil
			},
		},
		th,
	)

	hb := heartbeat.Heartbeat{
		Pubkey: []byte(pubKey),
	}
	hbBytes, _ := json.Marshal(hb)
	err := mon.ProcessReceivedMessage(&mock.P2PMessageStub{DataField: hbBytes}, nil)
=======
	mon, _ := heartbeat.NewMonitor(
		&mock.SinglesignStub{
			VerifyCalled: func(public crypto.PublicKey, msg []byte, sig []byte) error {
				return nil
			},
		},
		&mock.KeyGenMock{
			PublicKeyFromByteArrayMock: func(b []byte) (key crypto.PublicKey, e error) {
				return nil, nil
			},
		},
		&mock.MarshalizerMock{
			UnmarshalHandler: func(obj interface{}, buff []byte) error {
				(obj.(*heartbeat.Heartbeat)).Pubkey = []byte(pubKey)
				return nil
			},
		},
		time.Second*1000,
		map[uint32][]string{0: {"pk2"}},
	)

	err := mon.ProcessReceivedMessage(&mock.P2PMessageStub{DataField: []byte("")})
>>>>>>> Stashed changes
	assert.Nil(t, err)

	//a delay is mandatory for the go routine to finish its job
	time.Sleep(time.Second)

	//there should be 2 heartbeats, because a new one should have been added with pk2
	hbStatus := mon.GetHeartbeats()
	assert.Equal(t, 2, len(hbStatus))
	assert.Equal(t, hex.EncodeToString([]byte(pubKey)), hbStatus[0].HexPublicKey)
}

func TestMonitor_ProcessReceivedMessageWithNewShardID(t *testing.T) {
	t.Parallel()

	pubKey := []byte("pk1")

<<<<<<< Updated upstream
	th := &mock.MockTimer{}
	mon, _ := heartbeat.NewMonitor(
=======
	mon, _ := heartbeat.NewMonitor(
		&mock.SinglesignStub{
			VerifyCalled: func(public crypto.PublicKey, msg []byte, sig []byte) error {
				return nil
			},
		},
		&mock.KeyGenMock{
			PublicKeyFromByteArrayMock: func(b []byte) (key crypto.PublicKey, e error) {
				return nil, nil
			},
		},
>>>>>>> Stashed changes
		&mock.MarshalizerMock{
			UnmarshalHandler: func(obj interface{}, buff []byte) error {
				var rcvdHb heartbeat.Heartbeat
				_ = json.Unmarshal(buff, &rcvdHb)
				(obj.(*heartbeat.Heartbeat)).Pubkey = rcvdHb.Pubkey
				(obj.(*heartbeat.Heartbeat)).ShardID = rcvdHb.ShardID
				return nil
			},
		},
		time.Second*1000,
		map[uint32][]string{0: {"pk1"}},
<<<<<<< Updated upstream
		time.Now(),
		&mock.MessageHandlerStub{
			CreateHeartbeatFromP2pMessageCalled: func(message p2p.MessageP2P) (*heartbeat.Heartbeat, error) {
				var rcvHb heartbeat.Heartbeat
				_ = json.Unmarshal(message.Data(), &rcvHb)
				return &rcvHb, nil
			},
		},
		&mock.HeartbeatStorerStub{
			UpdateGenesisTimeCalled: func(genesisTime time.Time) error {
				return nil
			},
			LoadHbmiDTOCalled: func(pubKey string) (*heartbeat.HeartbeatDTO, error) {
				return nil, errors.New("not found")
			},
			LoadKeysCalled: func() ([][]byte, error) {
				return nil, nil
			},
			SavePubkeyDataCalled: func(pubkey []byte, heartbeat *heartbeat.HeartbeatDTO) error {
				return nil
			},
			SaveKeysCalled: func(peersSlice [][]byte) error {
				return nil
			},
		},
		th,
=======
>>>>>>> Stashed changes
	)

	// First send from pk1 from shard 0
	hb := &heartbeat.Heartbeat{
		Pubkey:  pubKey,
		ShardID: uint32(0),
	}

	buffToSend, err := json.Marshal(hb)
	assert.Nil(t, err)

<<<<<<< Updated upstream
	err = mon.ProcessReceivedMessage(&mock.P2PMessageStub{DataField: buffToSend}, nil)
=======
	err = mon.ProcessReceivedMessage(&mock.P2PMessageStub{DataField: buffToSend})
>>>>>>> Stashed changes
	assert.Nil(t, err)

	//a delay is mandatory for the go routine to finish its job
	time.Sleep(time.Second)

	hbStatus := mon.GetHeartbeats()

	assert.Equal(t, uint32(0), hbStatus[0].ReceivedShardID)

	// now we send a new heartbeat which will contain a new shard id
	hb = &heartbeat.Heartbeat{
		Pubkey:  pubKey,
		ShardID: uint32(1),
	}

	buffToSend, err = json.Marshal(hb)

	assert.Nil(t, err)

<<<<<<< Updated upstream
	err = mon.ProcessReceivedMessage(&mock.P2PMessageStub{DataField: buffToSend}, nil)
=======
	err = mon.ProcessReceivedMessage(&mock.P2PMessageStub{DataField: buffToSend})
>>>>>>> Stashed changes

	time.Sleep(1 * time.Second)

	hbStatus = mon.GetHeartbeats()

	// check if shard ID is changed at the same status
	assert.Equal(t, uint32(1), hbStatus[0].ReceivedShardID)
	assert.Equal(t, 1, len(hbStatus))
}

func TestMonitor_ProcessReceivedMessageShouldSetPeerInactive(t *testing.T) {
	t.Parallel()

	pubKey1 := "pk1-should-stay-online"
	pubKey2 := "pk2-should-go-offline"
<<<<<<< Updated upstream
	storer, _ := storage.NewHeartbeatDbStorer(mock.NewStorerMock(), &mock.MarshalizerFake{})
	th := &mock.MockTimer{}
	mon, _ := heartbeat.NewMonitor(
=======

	mon, _ := heartbeat.NewMonitor(
		&mock.SinglesignStub{
			VerifyCalled: func(public crypto.PublicKey, msg []byte, sig []byte) error {
				return nil
			},
		},
		&mock.KeyGenMock{
			PublicKeyFromByteArrayMock: func(b []byte) (key crypto.PublicKey, e error) {
				return nil, nil
			},
		},
>>>>>>> Stashed changes
		&mock.MarshalizerMock{
			UnmarshalHandler: func(obj interface{}, buff []byte) error {
				var rcvdHb heartbeat.Heartbeat
				_ = json.Unmarshal(buff, &rcvdHb)
				(obj.(*heartbeat.Heartbeat)).Pubkey = rcvdHb.Pubkey
				(obj.(*heartbeat.Heartbeat)).ShardID = rcvdHb.ShardID
				return nil
			},
		},
<<<<<<< Updated upstream
		time.Second*5,
		map[uint32][]string{0: {pubKey1, pubKey2}},
		th.Now(),
		&mock.MessageHandlerStub{
			CreateHeartbeatFromP2pMessageCalled: func(message p2p.MessageP2P) (*heartbeat.Heartbeat, error) {
				var rcvHb heartbeat.Heartbeat
				_ = json.Unmarshal(message.Data(), &rcvHb)
				return &rcvHb, nil
			},
		},
		storer,
		th,
=======
		time.Millisecond*5,
		map[uint32][]string{0: {pubKey1, pubKey2}},
>>>>>>> Stashed changes
	)

	// First send from pk1
	err := sendHbMessageFromPubKey(pubKey1, mon)
	assert.Nil(t, err)

	// Send from pk2
	err = sendHbMessageFromPubKey(pubKey2, mon)
	assert.Nil(t, err)

	// set pk2 to inactive as max inactive time is lower
<<<<<<< Updated upstream
	time.Sleep(10 * time.Millisecond)
	th.IncrementSeconds(6)
=======
	time.Sleep(6 * time.Millisecond)
>>>>>>> Stashed changes

	// Check that both are added
	hbStatus := mon.GetHeartbeats()
	assert.Equal(t, 2, len(hbStatus))
<<<<<<< Updated upstream
	//assert.False(t, hbStatus[1].IsActive)

	// Now send a message from pk1 in order to see that pk2 is not active anymore
	err = sendHbMessageFromPubKey(pubKey1, mon)
	time.Sleep(5 * time.Millisecond)
	assert.Nil(t, err)

	th.IncrementSeconds(4)
=======

	// Now send a message from pk1 in order to see that pk2 is not active anymore
	err = sendHbMessageFromPubKey(pubKey1, mon)
	assert.Nil(t, err)

	time.Sleep(5 * time.Millisecond)
>>>>>>> Stashed changes

	hbStatus = mon.GetHeartbeats()

	// check if pk1 is still on
	assert.True(t, hbStatus[0].IsActive)
	// check if pk2 was set to offline by pk1
	assert.False(t, hbStatus[1].IsActive)
}

func sendHbMessageFromPubKey(pubKey string, mon *heartbeat.Monitor) error {
	hb := &heartbeat.Heartbeat{
		Pubkey: []byte(pubKey),
	}
	buffToSend, _ := json.Marshal(hb)
<<<<<<< Updated upstream
	err := mon.ProcessReceivedMessage(&mock.P2PMessageStub{DataField: buffToSend}, nil)
=======
	err := mon.ProcessReceivedMessage(&mock.P2PMessageStub{DataField: buffToSend})
>>>>>>> Stashed changes
	return err
}
