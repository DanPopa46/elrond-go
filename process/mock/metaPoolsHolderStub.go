package mock

import (
	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/storage"
)

type MetaPoolsHolderStub struct {
<<<<<<< Updated upstream
	MetaBlocksCalled           func() storage.Cacher
	MiniBlocksCalled           func() storage.Cacher
	ShardHeadersCalled         func() storage.Cacher
	HeadersNoncesCalled        func() dataRetriever.Uint64SyncMapCacher
	TransactionsCalled         func() dataRetriever.ShardedDataCacherNotifier
	UnsignedTransactionsCalled func() dataRetriever.ShardedDataCacherNotifier
}

func (mphs *MetaPoolsHolderStub) Transactions() dataRetriever.ShardedDataCacherNotifier {
	return mphs.TransactionsCalled()
}

func (mphs *MetaPoolsHolderStub) UnsignedTransactions() dataRetriever.ShardedDataCacherNotifier {
	return mphs.UnsignedTransactionsCalled()
}

func (mphs *MetaPoolsHolderStub) MetaBlocks() storage.Cacher {
	return mphs.MetaBlocksCalled()
}

func (mphs *MetaPoolsHolderStub) MiniBlocks() storage.Cacher {
	return mphs.MiniBlocksCalled()
=======
	MetaChainBlocksCalled func() storage.Cacher
	MiniBlockHashesCalled func() dataRetriever.ShardedDataCacherNotifier
	ShardHeadersCalled    func() storage.Cacher
	HeadersNoncesCalled   func() dataRetriever.Uint64SyncMapCacher
}

func (mphs *MetaPoolsHolderStub) MetaChainBlocks() storage.Cacher {
	return mphs.MetaChainBlocksCalled()
}

func (mphs *MetaPoolsHolderStub) MiniBlockHashes() dataRetriever.ShardedDataCacherNotifier {
	return mphs.MiniBlockHashesCalled()
>>>>>>> Stashed changes
}

func (mphs *MetaPoolsHolderStub) ShardHeaders() storage.Cacher {
	return mphs.ShardHeadersCalled()
}

func (mphs *MetaPoolsHolderStub) HeadersNonces() dataRetriever.Uint64SyncMapCacher {
	return mphs.HeadersNoncesCalled()
}

// IsInterfaceNil returns true if there is no value under the interface
func (mphs *MetaPoolsHolderStub) IsInterfaceNil() bool {
	if mphs == nil {
		return true
	}
	return false
}
