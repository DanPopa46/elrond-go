package mock

import (
<<<<<<< Updated upstream
	"github.com/ElrondNetwork/elrond-go/core/indexer"
=======
>>>>>>> Stashed changes
	"github.com/ElrondNetwork/elrond-go/core/statistics"
	"github.com/ElrondNetwork/elrond-go/data"
	"github.com/ElrondNetwork/elrond-go/data/block"
)

// IndexerMock is a mock implementation fot the Indexer interface
type IndexerMock struct {
	SaveBlockCalled func(body block.Body, header *block.Header)
}

<<<<<<< Updated upstream
func (im *IndexerMock) SaveBlock(body data.BodyHandler, header data.HeaderHandler, txPool map[string]data.TransactionHandler, signersIndexes []uint64) {
	panic("implement me")
}

func (im *IndexerMock) SaveMetaBlock(header data.HeaderHandler, signersIndexes []uint64) {
	return
}

=======
func (im *IndexerMock) SaveBlock(body data.BodyHandler, header data.HeaderHandler, txPool map[string]data.TransactionHandler) {
	panic("implement me")
}

>>>>>>> Stashed changes
func (im *IndexerMock) UpdateTPS(tpsBenchmark statistics.TPSBenchmark) {
	panic("implement me")
}

<<<<<<< Updated upstream
func (im *IndexerMock) SaveRoundInfo(roundInfo indexer.RoundInfo) {
	panic("implement me")
}

func (im *IndexerMock) SaveValidatorsPubKeys(validatorsPubKeys map[uint32][][]byte) {
	panic("implement me")
}

=======
>>>>>>> Stashed changes
// IsInterfaceNil returns true if there is no value under the interface
func (im *IndexerMock) IsInterfaceNil() bool {
	if im == nil {
		return true
	}
	return false
}
<<<<<<< Updated upstream

func (im *IndexerMock) IsNilIndexer() bool {
	return false
}
=======
>>>>>>> Stashed changes
