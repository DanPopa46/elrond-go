package resolvers

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go/data/block"
	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/storage"
)

<<<<<<< Updated upstream
// genericBlockBodyResolver is a wrapper over Resolver that is specialized in resolving block body requests
type genericBlockBodyResolver struct {
=======
// GenericBlockBodyResolver is a wrapper over Resolver that is specialized in resolving block body requests
type GenericBlockBodyResolver struct {
>>>>>>> Stashed changes
	dataRetriever.TopicResolverSender
	miniBlockPool    storage.Cacher
	miniBlockStorage storage.Storer
	marshalizer      marshal.Marshalizer
}

// NewGenericBlockBodyResolver creates a new block body resolver
func NewGenericBlockBodyResolver(
	senderResolver dataRetriever.TopicResolverSender,
	miniBlockPool storage.Cacher,
	miniBlockStorage storage.Storer,
<<<<<<< Updated upstream
	marshalizer marshal.Marshalizer,
) (*genericBlockBodyResolver, error) {
=======
	marshalizer marshal.Marshalizer) (*GenericBlockBodyResolver, error) {
>>>>>>> Stashed changes

	if senderResolver == nil || senderResolver.IsInterfaceNil() {
		return nil, dataRetriever.ErrNilResolverSender
	}

	if miniBlockPool == nil || miniBlockPool.IsInterfaceNil() {
		return nil, dataRetriever.ErrNilBlockBodyPool
	}

	if miniBlockStorage == nil || miniBlockStorage.IsInterfaceNil() {
		return nil, dataRetriever.ErrNilBlockBodyStorage
	}

	if marshalizer == nil || marshalizer.IsInterfaceNil() {
		return nil, dataRetriever.ErrNilMarshalizer
	}

<<<<<<< Updated upstream
	bbResolver := &genericBlockBodyResolver{
=======
	bbResolver := &GenericBlockBodyResolver{
>>>>>>> Stashed changes
		TopicResolverSender: senderResolver,
		miniBlockPool:       miniBlockPool,
		miniBlockStorage:    miniBlockStorage,
		marshalizer:         marshalizer,
	}

	return bbResolver, nil
}

// ProcessReceivedMessage will be the callback func from the p2p.Messenger and will be called each time a new message was received
// (for the topic this validator was registered to, usually a request topic)
<<<<<<< Updated upstream
func (gbbRes *genericBlockBodyResolver) ProcessReceivedMessage(message p2p.MessageP2P, _ func(buffToSend []byte)) error {
=======
func (gbbRes *GenericBlockBodyResolver) ProcessReceivedMessage(message p2p.MessageP2P) error {
>>>>>>> Stashed changes
	rd := &dataRetriever.RequestData{}
	err := rd.Unmarshal(gbbRes.marshalizer, message)
	if err != nil {
		return err
	}

	buff, err := gbbRes.resolveBlockBodyRequest(rd)
	if err != nil {
		return err
	}

	if buff == nil {
		log.Debug(fmt.Sprintf("missing data: %v", rd))
		return nil
	}

	return gbbRes.Send(buff, message.Peer())
}

<<<<<<< Updated upstream
func (gbbRes *genericBlockBodyResolver) resolveBlockBodyRequest(rd *dataRetriever.RequestData) ([]byte, error) {
=======
func (gbbRes *GenericBlockBodyResolver) resolveBlockBodyRequest(rd *dataRetriever.RequestData) ([]byte, error) {

>>>>>>> Stashed changes
	if rd.Value == nil {
		return nil, dataRetriever.ErrNilValue
	}

<<<<<<< Updated upstream
	hashes, err := gbbRes.miniBlockHashesFromRequestType(rd)
	if err != nil {
		return nil, err
	}

	miniBlocks, _ := gbbRes.GetMiniBlocks(hashes)
	if len(miniBlocks) == 0 {
		return nil, dataRetriever.ErrEmptyMiniBlockSlice
=======
	miniBlockHashes, err := gbbRes.miniBlockHashesFromRequestType(rd)
	if err != nil {
		return nil, err
	}
	miniBlocks := gbbRes.GetMiniBlocks(miniBlockHashes)

	if miniBlocks == nil {
		return nil, dataRetriever.ErrNilMiniBlocks
>>>>>>> Stashed changes
	}

	buff, err := gbbRes.marshalizer.Marshal(miniBlocks)
	if err != nil {
		return nil, err
	}

	return buff, nil
}

<<<<<<< Updated upstream
func (gbbRes *genericBlockBodyResolver) miniBlockHashesFromRequestType(requestData *dataRetriever.RequestData) ([][]byte, error) {
=======
func (gbbRes *GenericBlockBodyResolver) miniBlockHashesFromRequestType(requestData *dataRetriever.RequestData) ([][]byte, error) {
>>>>>>> Stashed changes
	miniBlockHashes := make([][]byte, 0)

	switch requestData.Type {
	case dataRetriever.HashType:
		miniBlockHashes = append(miniBlockHashes, requestData.Value)

	case dataRetriever.HashArrayType:
		err := gbbRes.marshalizer.Unmarshal(&miniBlockHashes, requestData.Value)

		if err != nil {
			return nil, dataRetriever.ErrUnmarshalMBHashes
		}

	default:
		return nil, dataRetriever.ErrInvalidRequestType
	}

	return miniBlockHashes, nil
}

// RequestDataFromHash requests a block body from other peers having input the block body hash
<<<<<<< Updated upstream
func (gbbRes *genericBlockBodyResolver) RequestDataFromHash(hash []byte) error {
=======
func (gbbRes *GenericBlockBodyResolver) RequestDataFromHash(hash []byte) error {
>>>>>>> Stashed changes
	return gbbRes.SendOnRequestTopic(&dataRetriever.RequestData{
		Type:  dataRetriever.HashType,
		Value: hash,
	})
}

// RequestDataFromHashArray requests a block body from other peers having input the block body hash
<<<<<<< Updated upstream
func (gbbRes *genericBlockBodyResolver) RequestDataFromHashArray(hashes [][]byte) error {
=======
func (gbbRes *GenericBlockBodyResolver) RequestDataFromHashArray(hashes [][]byte) error {
>>>>>>> Stashed changes
	hash, err := gbbRes.marshalizer.Marshal(hashes)

	if err != nil {
		return err
	}

	return gbbRes.SendOnRequestTopic(&dataRetriever.RequestData{
		Type:  dataRetriever.HashArrayType,
		Value: hash,
	})
}

// GetMiniBlocks method returns a list of deserialized mini blocks from a given hash list either from data pool or from storage
<<<<<<< Updated upstream
func (gbbRes *genericBlockBodyResolver) GetMiniBlocks(hashes [][]byte) (block.MiniBlockSlice, [][]byte) {
	miniBlocks, missingMiniBlocksHashes := gbbRes.GetMiniBlocksFromPool(hashes)
	if len(missingMiniBlocksHashes) == 0 {
		return miniBlocks, missingMiniBlocksHashes
	}

	miniBlocksFromStorer, missingMiniBlocksHashes := gbbRes.getMiniBlocksFromStorer(missingMiniBlocksHashes)
	miniBlocks = append(miniBlocks, miniBlocksFromStorer...)

	return miniBlocks, missingMiniBlocksHashes
}

// GetMiniBlocksFromPool method returns a list of deserialized mini blocks from a given hash list from data pool
func (gbbRes *genericBlockBodyResolver) GetMiniBlocksFromPool(hashes [][]byte) (block.MiniBlockSlice, [][]byte) {
	miniBlocks := make(block.MiniBlockSlice, 0)
	missingMiniBlocksHashes := make([][]byte, 0)

	for i := 0; i < len(hashes); i++ {
		obj, ok := gbbRes.miniBlockPool.Peek(hashes[i])
		if !ok {
			missingMiniBlocksHashes = append(missingMiniBlocksHashes, hashes[i])
			continue
		}

		miniBlock, ok := obj.(*block.MiniBlock)
		if !ok {
			missingMiniBlocksHashes = append(missingMiniBlocksHashes, hashes[i])
			continue
		}

		miniBlocks = append(miniBlocks, miniBlock)
	}

	return miniBlocks, missingMiniBlocksHashes
}

// getMiniBlocksFromStorer returns a list of mini blocks from storage and a list of missing hashes
func (gbbRes *genericBlockBodyResolver) getMiniBlocksFromStorer(hashes [][]byte) (block.MiniBlockSlice, [][]byte) {
	miniBlocks := make(block.MiniBlockSlice, 0)
	missingMiniBlocksHashes := make([][]byte, 0)

	for i := 0; i < len(hashes); i++ {
		buff, err := gbbRes.miniBlockStorage.Get(hashes[i])
		if err != nil {
			log.Debug(err.Error())
			missingMiniBlocksHashes = append(missingMiniBlocksHashes, hashes[i])
			continue
		}

		miniBlock := &block.MiniBlock{}
		err = gbbRes.marshalizer.Unmarshal(miniBlock, buff)
		if err != nil {
			log.Debug(err.Error())
			gbbRes.miniBlockPool.Remove([]byte(hashes[i]))
			err = gbbRes.miniBlockStorage.Remove([]byte(hashes[i]))
			if err != nil {
				log.Debug(err.Error())
			}

			missingMiniBlocksHashes = append(missingMiniBlocksHashes, hashes[i])
			continue
		}

		miniBlocks = append(miniBlocks, miniBlock)
	}

	return miniBlocks, missingMiniBlocksHashes
}

// IsInterfaceNil returns true if there is no value under the interface
func (gbbRes *genericBlockBodyResolver) IsInterfaceNil() bool {
=======
func (gbbRes *GenericBlockBodyResolver) GetMiniBlocks(hashes [][]byte) block.MiniBlockSlice {
	miniBlocks := gbbRes.getMiniBlocks(hashes)
	if miniBlocks == nil {
		return nil
	}

	mbLength := len(hashes)
	expandedMiniBlocks := make(block.MiniBlockSlice, mbLength)

	for i := 0; i < mbLength; i++ {
		mb := &block.MiniBlock{}
		err := gbbRes.marshalizer.Unmarshal(mb, miniBlocks[i])

		if err != nil {
			log.Debug(err.Error())
			gbbRes.miniBlockPool.Remove(hashes[i])
			err = gbbRes.miniBlockStorage.Remove(hashes[i])
			if err != nil {
				log.Debug(err.Error())
			}

			return nil
		}

		expandedMiniBlocks[i] = mb
	}

	return expandedMiniBlocks
}

// getMiniBlocks method returns a list of serialized mini blocks from a given hash list either from data pool or from storage
func (gbbRes *GenericBlockBodyResolver) getMiniBlocks(hashes [][]byte) [][]byte {
	miniBlocks := gbbRes.getMiniBlocksFromCache(hashes)

	if miniBlocks != nil {
		return miniBlocks
	}

	return gbbRes.getMiniBlocksFromStorer(hashes)
}

// getMiniBlocksFromCache returns a full list of miniblocks from cache.
// If any of the miniblocks is missing the function returns nil
func (gbbRes *GenericBlockBodyResolver) getMiniBlocksFromCache(hashes [][]byte) [][]byte {
	miniBlocksLen := len(hashes)
	miniBlocks := make([][]byte, miniBlocksLen)

	for i := 0; i < miniBlocksLen; i++ {
		cachedMB, _ := gbbRes.miniBlockPool.Peek(hashes[i])

		if cachedMB == nil {
			return nil
		}

		buff, err := gbbRes.marshalizer.Marshal(cachedMB)
		if err != nil {
			log.Debug(err.Error())
			return nil
		}

		miniBlocks[i] = buff
	}

	return miniBlocks
}

// getMiniBlocksFromStorer returns a full list of MiniBlocks from the storage unit.
// If any MiniBlock is missing or is invalid, it is removed and the function returns nil
func (gbbRes *GenericBlockBodyResolver) getMiniBlocksFromStorer(hashes [][]byte) [][]byte {
	miniBlocksLen := len(hashes)
	miniBlocks := make([][]byte, miniBlocksLen)

	for i := 0; i < miniBlocksLen; i++ {
		buff, err := gbbRes.miniBlockStorage.Get(hashes[i])
		if err != nil {
			log.Debug(err.Error())
			return nil
		}

		miniBlocks[i] = buff
	}

	return miniBlocks
}

// IsInterfaceNil returns true if there is no value under the interface
func (gbbRes *GenericBlockBodyResolver) IsInterfaceNil() bool {
>>>>>>> Stashed changes
	if gbbRes == nil {
		return true
	}
	return false
}
