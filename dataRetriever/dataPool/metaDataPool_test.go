package dataPool_test

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/dataRetriever/dataPool"
	"github.com/ElrondNetwork/elrond-go/dataRetriever/mock"
	"github.com/stretchr/testify/assert"
)

//------- NewDataPool

func TestNewMetaDataPool_NilMetaBlockShouldErr(t *testing.T) {
<<<<<<< Updated upstream
	t.Parallel()

	tdp, err := dataPool.NewMetaDataPool(
		nil,
		&mock.CacherStub{},
		&mock.CacherStub{},
		&mock.Uint64SyncMapCacherStub{},
		&mock.ShardedDataStub{},
		&mock.ShardedDataStub{},
=======
	tdp, err := dataPool.NewMetaDataPool(
		nil,
		&mock.ShardedDataStub{},
		&mock.CacherStub{},
		&mock.Uint64SyncMapCacherStub{},
>>>>>>> Stashed changes
	)

	assert.Equal(t, dataRetriever.ErrNilMetaBlockPool, err)
	assert.Nil(t, tdp)
}

func TestNewMetaDataPool_NilMiniBlockHeaderHashesShouldErr(t *testing.T) {
<<<<<<< Updated upstream
	t.Parallel()

=======
>>>>>>> Stashed changes
	tdp, err := dataPool.NewMetaDataPool(
		&mock.CacherStub{},
		nil,
		&mock.CacherStub{},
		&mock.Uint64SyncMapCacherStub{},
<<<<<<< Updated upstream
		&mock.ShardedDataStub{},
		&mock.ShardedDataStub{},
=======
>>>>>>> Stashed changes
	)

	assert.Equal(t, dataRetriever.ErrNilMiniBlockHashesPool, err)
	assert.Nil(t, tdp)
}

func TestNewMetaDataPool_NilShardHeaderShouldErr(t *testing.T) {
<<<<<<< Updated upstream
	t.Parallel()

	tdp, err := dataPool.NewMetaDataPool(
		&mock.CacherStub{},
		&mock.CacherStub{},
		nil,
		&mock.Uint64SyncMapCacherStub{},
		&mock.ShardedDataStub{},
		&mock.ShardedDataStub{},
=======
	tdp, err := dataPool.NewMetaDataPool(
		&mock.CacherStub{},
		&mock.ShardedDataStub{},
		nil,
		&mock.Uint64SyncMapCacherStub{},
>>>>>>> Stashed changes
	)

	assert.Equal(t, dataRetriever.ErrNilShardHeaderPool, err)
	assert.Nil(t, tdp)
}

func TestNewMetaDataPool_NilHeaderNoncesShouldErr(t *testing.T) {
<<<<<<< Updated upstream
	t.Parallel()

	tdp, err := dataPool.NewMetaDataPool(
		&mock.CacherStub{},
		&mock.CacherStub{},
		&mock.CacherStub{},
		nil,
		&mock.ShardedDataStub{},
		&mock.ShardedDataStub{},
	)

	assert.Equal(t, dataRetriever.ErrNilMetaBlockNoncesPool, err)
	assert.Nil(t, tdp)
}

func TestNewMetaDataPool_NilTxPoolShouldErr(t *testing.T) {
	t.Parallel()

	tdp, err := dataPool.NewMetaDataPool(
		&mock.CacherStub{},
		&mock.CacherStub{},
		&mock.CacherStub{},
		&mock.Uint64SyncMapCacherStub{},
		nil,
		&mock.ShardedDataStub{},
	)

	assert.Equal(t, dataRetriever.ErrNilTxDataPool, err)
	assert.Nil(t, tdp)
}

func TestNewMetaDataPool_NilUnsingedPoolNoncesShouldErr(t *testing.T) {
	t.Parallel()

	tdp, err := dataPool.NewMetaDataPool(
		&mock.CacherStub{},
		&mock.CacherStub{},
		&mock.CacherStub{},
		&mock.Uint64SyncMapCacherStub{},
		&mock.ShardedDataStub{},
		nil,
	)

	assert.Equal(t, dataRetriever.ErrNilUnsignedTransactionPool, err)
=======
	tdp, err := dataPool.NewMetaDataPool(
		&mock.CacherStub{},
		&mock.ShardedDataStub{},
		&mock.CacherStub{},
		nil,
	)

	assert.Equal(t, dataRetriever.ErrNilMetaBlockNoncesPool, err)
>>>>>>> Stashed changes
	assert.Nil(t, tdp)
}

func TestNewMetaDataPool_ConfigOk(t *testing.T) {
<<<<<<< Updated upstream
	t.Parallel()

	metaBlocks := &mock.CacherStub{}
	shardHeaders := &mock.CacherStub{}
	miniBlocks := &mock.CacherStub{}
	hdrsNonces := &mock.Uint64SyncMapCacherStub{}
	transactions := &mock.ShardedDataStub{}
	unsigned := &mock.ShardedDataStub{}

	tdp, err := dataPool.NewMetaDataPool(
		metaBlocks,
		miniBlocks,
		shardHeaders,
		hdrsNonces,
		transactions,
		unsigned,
=======
	metaChainBlocks := &mock.CacherStub{}
	shardHeaders := &mock.CacherStub{}
	miniBlockheaders := &mock.ShardedDataStub{}
	hdrsNonces := &mock.Uint64SyncMapCacherStub{}

	tdp, err := dataPool.NewMetaDataPool(
		metaChainBlocks,
		miniBlockheaders,
		shardHeaders,
		hdrsNonces,
>>>>>>> Stashed changes
	)

	assert.Nil(t, err)
	//pointer checking
<<<<<<< Updated upstream
	assert.True(t, metaBlocks == tdp.MetaBlocks())
	assert.True(t, shardHeaders == tdp.ShardHeaders())
	assert.True(t, miniBlocks == tdp.MiniBlocks())
	assert.True(t, hdrsNonces == tdp.HeadersNonces())
	assert.True(t, transactions == tdp.Transactions())
	assert.True(t, unsigned == tdp.UnsignedTransactions())
=======
	assert.True(t, metaChainBlocks == tdp.MetaChainBlocks())
	assert.True(t, shardHeaders == tdp.ShardHeaders())
	assert.True(t, miniBlockheaders == tdp.MiniBlockHashes())
	assert.True(t, hdrsNonces == tdp.HeadersNonces())
>>>>>>> Stashed changes
}
