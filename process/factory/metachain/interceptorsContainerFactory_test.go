package metachain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/process/factory"
	"github.com/ElrondNetwork/elrond-go/process/factory/metachain"
	"github.com/ElrondNetwork/elrond-go/process/mock"
	"github.com/ElrondNetwork/elrond-go/storage"
	"github.com/stretchr/testify/assert"
)

<<<<<<< Updated upstream
const maxTxNonceDeltaAllowed = 100

=======
>>>>>>> Stashed changes
var errExpected = errors.New("expected error")

func createStubTopicHandler(matchStrToErrOnCreate string, matchStrToErrOnRegister string) process.TopicHandler {
	return &mock.TopicHandlerStub{
		CreateTopicCalled: func(name string, createChannelForTopic bool) error {
			if matchStrToErrOnCreate == "" {
				return nil
			}

			if strings.Contains(name, matchStrToErrOnCreate) {
				return errExpected
			}

			return nil
		},
		RegisterMessageProcessorCalled: func(topic string, handler p2p.MessageProcessor) error {
			if matchStrToErrOnRegister == "" {
				return nil
			}

			if strings.Contains(topic, matchStrToErrOnRegister) {
				return errExpected
			}

			return nil
		},
	}
}

func createDataPools() dataRetriever.MetaPoolsHolder {
	pools := &mock.MetaPoolsHolderStub{
		ShardHeadersCalled: func() storage.Cacher {
			return &mock.CacherStub{}
		},
<<<<<<< Updated upstream
		MiniBlocksCalled: func() storage.Cacher {
			return &mock.CacherStub{}
		},
		MetaBlocksCalled: func() storage.Cacher {
=======
		MiniBlockHashesCalled: func() dataRetriever.ShardedDataCacherNotifier {
			return &mock.ShardedDataStub{}
		},
		MetaChainBlocksCalled: func() storage.Cacher {
>>>>>>> Stashed changes
			return &mock.CacherStub{}
		},
		HeadersNoncesCalled: func() dataRetriever.Uint64SyncMapCacher {
			return &mock.Uint64SyncMapCacherStub{}
		},
<<<<<<< Updated upstream
		TransactionsCalled: func() dataRetriever.ShardedDataCacherNotifier {
			return &mock.ShardedDataStub{}
		},
		UnsignedTransactionsCalled: func() dataRetriever.ShardedDataCacherNotifier {
			return &mock.ShardedDataStub{}
		},
=======
>>>>>>> Stashed changes
	}

	return pools
}

func createStore() *mock.ChainStorerMock {
	return &mock.ChainStorerMock{
		GetStorerCalled: func(unitType dataRetriever.UnitType) storage.Storer {
			return &mock.StorerStub{}
		},
	}
}

//------- NewInterceptorsContainerFactory

func TestNewInterceptorsContainerFactory_NilShardCoordinatorShouldErr(t *testing.T) {
	t.Parallel()

	icf, err := metachain.NewInterceptorsContainerFactory(
		nil,
<<<<<<< Updated upstream
		mock.NewNodesCoordinatorMock(),
=======
>>>>>>> Stashed changes
		&mock.TopicHandlerStub{},
		createStore(),
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		createDataPools(),
<<<<<<< Updated upstream
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
=======
		&mock.ChronologyValidatorStub{},
>>>>>>> Stashed changes
	)

	assert.Nil(t, icf)
	assert.Equal(t, process.ErrNilShardCoordinator, err)
}

<<<<<<< Updated upstream
func TestNewInterceptorsContainerFactory_NilNodesCoordinatorShouldErr(t *testing.T) {
	t.Parallel()

	icf, err := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
		nil,
		&mock.TopicHandlerStub{},
		createStore(),
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		createDataPools(),
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
	)

	assert.Nil(t, icf)
	assert.Equal(t, process.ErrNilNodesCoordinator, err)
}

=======
>>>>>>> Stashed changes
func TestNewInterceptorsContainerFactory_NilTopicHandlerShouldErr(t *testing.T) {
	t.Parallel()

	icf, err := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
<<<<<<< Updated upstream
		mock.NewNodesCoordinatorMock(),
=======
>>>>>>> Stashed changes
		nil,
		createStore(),
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		createDataPools(),
<<<<<<< Updated upstream
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
=======
		&mock.ChronologyValidatorStub{},
>>>>>>> Stashed changes
	)

	assert.Nil(t, icf)
	assert.Equal(t, process.ErrNilMessenger, err)
}

func TestNewInterceptorsContainerFactory_NilBlockchainShouldErr(t *testing.T) {
	t.Parallel()

	icf, err := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
<<<<<<< Updated upstream
		mock.NewNodesCoordinatorMock(),
=======
>>>>>>> Stashed changes
		&mock.TopicHandlerStub{},
		nil,
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		createDataPools(),
<<<<<<< Updated upstream
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
=======
		&mock.ChronologyValidatorStub{},
>>>>>>> Stashed changes
	)

	assert.Nil(t, icf)
	assert.Equal(t, process.ErrNilStore, err)
}

func TestNewInterceptorsContainerFactory_NilMarshalizerShouldErr(t *testing.T) {
	t.Parallel()

	icf, err := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
<<<<<<< Updated upstream
		mock.NewNodesCoordinatorMock(),
=======
>>>>>>> Stashed changes
		&mock.TopicHandlerStub{},
		createStore(),
		nil,
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		createDataPools(),
<<<<<<< Updated upstream
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
=======
		&mock.ChronologyValidatorStub{},
>>>>>>> Stashed changes
	)

	assert.Nil(t, icf)
	assert.Equal(t, process.ErrNilMarshalizer, err)
}

func TestNewInterceptorsContainerFactory_NilHasherShouldErr(t *testing.T) {
	t.Parallel()

	icf, err := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
<<<<<<< Updated upstream
		mock.NewNodesCoordinatorMock(),
=======
>>>>>>> Stashed changes
		&mock.TopicHandlerStub{},
		createStore(),
		&mock.MarshalizerMock{},
		nil,
		mock.NewMultiSigner(),
		createDataPools(),
<<<<<<< Updated upstream
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
=======
		&mock.ChronologyValidatorStub{},
>>>>>>> Stashed changes
	)

	assert.Nil(t, icf)
	assert.Equal(t, process.ErrNilHasher, err)
}

func TestNewInterceptorsContainerFactory_NilMultiSignerShouldErr(t *testing.T) {
	t.Parallel()

	icf, err := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
<<<<<<< Updated upstream
		mock.NewNodesCoordinatorMock(),
=======
>>>>>>> Stashed changes
		&mock.TopicHandlerStub{},
		createStore(),
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		nil,
		createDataPools(),
<<<<<<< Updated upstream
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
=======
		&mock.ChronologyValidatorStub{},
>>>>>>> Stashed changes
	)

	assert.Nil(t, icf)
	assert.Equal(t, process.ErrNilMultiSigVerifier, err)
}

func TestNewInterceptorsContainerFactory_NilDataPoolShouldErr(t *testing.T) {
	t.Parallel()

	icf, err := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
<<<<<<< Updated upstream
		mock.NewNodesCoordinatorMock(),
=======
>>>>>>> Stashed changes
		&mock.TopicHandlerStub{},
		createStore(),
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		nil,
<<<<<<< Updated upstream
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
=======
		&mock.ChronologyValidatorStub{},
>>>>>>> Stashed changes
	)

	assert.Nil(t, icf)
	assert.Equal(t, process.ErrNilDataPoolHolder, err)
}

<<<<<<< Updated upstream
func TestNewInterceptorsContainerFactory_NilAccountsShouldErr(t *testing.T) {
	t.Parallel()

	icf, err := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
		mock.NewNodesCoordinatorMock(),
		&mock.TopicHandlerStub{},
		createStore(),
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		createDataPools(),
		nil,
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
	)

	assert.Nil(t, icf)
	assert.Equal(t, process.ErrNilAccountsAdapter, err)
}

func TestNewInterceptorsContainerFactory_NilAddrConvShouldErr(t *testing.T) {
	t.Parallel()

	icf, err := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
		mock.NewNodesCoordinatorMock(),
		&mock.TopicHandlerStub{},
		createStore(),
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		createDataPools(),
		&mock.AccountsStub{},
		nil,
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
	)

	assert.Nil(t, icf)
	assert.Equal(t, process.ErrNilAddressConverter, err)
}

func TestNewInterceptorsContainerFactory_NilSingleSignerShouldErr(t *testing.T) {
	t.Parallel()

	icf, err := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
		mock.NewNodesCoordinatorMock(),
		&mock.TopicHandlerStub{},
		createStore(),
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		createDataPools(),
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		nil,
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
	)

	assert.Nil(t, icf)
	assert.Equal(t, process.ErrNilSingleSigner, err)
}

func TestNewInterceptorsContainerFactory_NilKeyGenShouldErr(t *testing.T) {
	t.Parallel()

	icf, err := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
		mock.NewNodesCoordinatorMock(),
		&mock.TopicHandlerStub{},
		createStore(),
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		createDataPools(),
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		nil,
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
	)

	assert.Nil(t, icf)
	assert.Equal(t, process.ErrNilKeyGen, err)
}

func TestNewInterceptorsContainerFactory_NilFeeHandlerShouldErr(t *testing.T) {
	t.Parallel()

	icf, err := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
		mock.NewNodesCoordinatorMock(),
		&mock.TopicHandlerStub{},
		createStore(),
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		createDataPools(),
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		nil,
	)

	assert.Nil(t, icf)
	assert.Equal(t, process.ErrNilEconomicsFeeHandler, err)
}

=======
>>>>>>> Stashed changes
func TestNewInterceptorsContainerFactory_ShouldWork(t *testing.T) {
	t.Parallel()

	icf, err := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
<<<<<<< Updated upstream
		mock.NewNodesCoordinatorMock(),
=======
>>>>>>> Stashed changes
		&mock.TopicHandlerStub{},
		createStore(),
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		createDataPools(),
<<<<<<< Updated upstream
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
=======
		&mock.ChronologyValidatorStub{},
>>>>>>> Stashed changes
	)

	assert.NotNil(t, icf)
	assert.Nil(t, err)
}

//------- Create

func TestInterceptorsContainerFactory_CreateTopicMetablocksFailsShouldErr(t *testing.T) {
	t.Parallel()

	icf, _ := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
<<<<<<< Updated upstream
		mock.NewNodesCoordinatorMock(),
=======
>>>>>>> Stashed changes
		createStubTopicHandler(factory.MetachainBlocksTopic, ""),
		createStore(),
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		createDataPools(),
<<<<<<< Updated upstream
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
=======
		&mock.ChronologyValidatorStub{},
>>>>>>> Stashed changes
	)

	container, err := icf.Create()

	assert.Nil(t, container)
	assert.Equal(t, errExpected, err)
}

func TestInterceptorsContainerFactory_CreateTopicShardHeadersForMetachainFailsShouldErr(t *testing.T) {
	t.Parallel()

	icf, _ := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
<<<<<<< Updated upstream
		mock.NewNodesCoordinatorMock(),
=======
>>>>>>> Stashed changes
		createStubTopicHandler(factory.ShardHeadersForMetachainTopic, ""),
		createStore(),
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		createDataPools(),
<<<<<<< Updated upstream
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
=======
		&mock.ChronologyValidatorStub{},
>>>>>>> Stashed changes
	)

	container, err := icf.Create()

	assert.Nil(t, container)
	assert.Equal(t, errExpected, err)
}

func TestInterceptorsContainerFactory_CreateRegisterForMetablocksFailsShouldErr(t *testing.T) {
	t.Parallel()

	icf, _ := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
<<<<<<< Updated upstream
		mock.NewNodesCoordinatorMock(),
=======
>>>>>>> Stashed changes
		createStubTopicHandler("", factory.MetachainBlocksTopic),
		createStore(),
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		createDataPools(),
<<<<<<< Updated upstream
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
=======
		&mock.ChronologyValidatorStub{},
>>>>>>> Stashed changes
	)

	container, err := icf.Create()

	assert.Nil(t, container)
	assert.Equal(t, errExpected, err)
}

func TestInterceptorsContainerFactory_CreateRegisterShardHeadersForMetachainFailsShouldErr(t *testing.T) {
	t.Parallel()

	icf, _ := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
<<<<<<< Updated upstream
		mock.NewNodesCoordinatorMock(),
=======
>>>>>>> Stashed changes
		createStubTopicHandler("", factory.ShardHeadersForMetachainTopic),
		createStore(),
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		createDataPools(),
<<<<<<< Updated upstream
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
=======
		&mock.ChronologyValidatorStub{},
>>>>>>> Stashed changes
	)

	container, err := icf.Create()

	assert.Nil(t, container)
	assert.Equal(t, errExpected, err)
}

func TestInterceptorsContainerFactory_CreateShouldWork(t *testing.T) {
	t.Parallel()

	icf, _ := metachain.NewInterceptorsContainerFactory(
		mock.NewOneShardCoordinatorMock(),
<<<<<<< Updated upstream
		mock.NewNodesCoordinatorMock(),
=======
>>>>>>> Stashed changes
		&mock.TopicHandlerStub{
			CreateTopicCalled: func(name string, createChannelForTopic bool) error {
				return nil
			},
			RegisterMessageProcessorCalled: func(topic string, handler p2p.MessageProcessor) error {
				return nil
			},
		},
		createStore(),
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		createDataPools(),
<<<<<<< Updated upstream
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
=======
		&mock.ChronologyValidatorStub{},
>>>>>>> Stashed changes
	)

	container, err := icf.Create()

	assert.NotNil(t, container)
	assert.Nil(t, err)
}

func TestInterceptorsContainerFactory_With4ShardsShouldWork(t *testing.T) {
	t.Parallel()

	noOfShards := 4

	shardCoordinator := mock.NewMultipleShardsCoordinatorMock()
	shardCoordinator.SetNoShards(uint32(noOfShards))
	shardCoordinator.CurrentShard = 1

<<<<<<< Updated upstream
	nodesCoordinator := &mock.NodesCoordinatorMock{
		ShardConsensusSize: 1,
		MetaConsensusSize:  1,
		NbShards:           uint32(noOfShards),
		ShardId:            1,
	}

	icf, _ := metachain.NewInterceptorsContainerFactory(
		shardCoordinator,
		nodesCoordinator,
=======
	icf, _ := metachain.NewInterceptorsContainerFactory(
		shardCoordinator,
>>>>>>> Stashed changes
		&mock.TopicHandlerStub{
			CreateTopicCalled: func(name string, createChannelForTopic bool) error {
				return nil
			},
			RegisterMessageProcessorCalled: func(topic string, handler p2p.MessageProcessor) error {
				return nil
			},
		},
		createStore(),
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewMultiSigner(),
		createDataPools(),
<<<<<<< Updated upstream
		&mock.AccountsStub{},
		&mock.AddressConverterMock{},
		&mock.SignerMock{},
		&mock.SingleSignKeyGenMock{},
		maxTxNonceDeltaAllowed,
		&mock.FeeHandlerStub{},
	)

	container, err := icf.Create()

	numInterceptorsMetablock := 1
	numInterceptorsShardHeadersForMetachain := noOfShards
	numInterceptorsTransactionsForMetachain := noOfShards + 1
	numInterceptorsUnsignedTxsForMetachain := noOfShards + 1
	totalInterceptors := numInterceptorsMetablock + numInterceptorsShardHeadersForMetachain +
		numInterceptorsTransactionsForMetachain + numInterceptorsUnsignedTxsForMetachain

	assert.Nil(t, err)
=======
		&mock.ChronologyValidatorStub{},
	)

	container, _ := icf.Create()

	numInterceptorsMetablock := 1
	numInterceptorsShardHeadersForMetachain := noOfShards
	totalInterceptors := numInterceptorsMetablock + numInterceptorsShardHeadersForMetachain

>>>>>>> Stashed changes
	assert.Equal(t, totalInterceptors, container.Len())
}
