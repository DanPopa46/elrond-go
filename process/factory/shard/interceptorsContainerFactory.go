package shard

import (
<<<<<<< Updated upstream
	"github.com/ElrondNetwork/elrond-go/core/throttler"
=======
>>>>>>> Stashed changes
	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/hashing"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/process"
<<<<<<< Updated upstream
	"github.com/ElrondNetwork/elrond-go/process/dataValidators"
	"github.com/ElrondNetwork/elrond-go/process/factory"
	"github.com/ElrondNetwork/elrond-go/process/factory/containers"
	"github.com/ElrondNetwork/elrond-go/process/interceptors"
	interceptorFactory "github.com/ElrondNetwork/elrond-go/process/interceptors/factory"
	"github.com/ElrondNetwork/elrond-go/process/interceptors/processor"
	"github.com/ElrondNetwork/elrond-go/process/mock"
	"github.com/ElrondNetwork/elrond-go/process/rewardTransaction"
	"github.com/ElrondNetwork/elrond-go/sharding"
)

const numGoRoutines = 2000

type interceptorsContainerFactory struct {
	accounts               state.AccountsAdapter
	shardCoordinator       sharding.Coordinator
	messenger              process.TopicHandler
	store                  dataRetriever.StorageService
	marshalizer            marshal.Marshalizer
	hasher                 hashing.Hasher
	keyGen                 crypto.KeyGenerator
	singleSigner           crypto.SingleSigner
	multiSigner            crypto.MultiSigner
	dataPool               dataRetriever.PoolsHolder
	addrConverter          state.AddressConverter
	nodesCoordinator       sharding.NodesCoordinator
	argInterceptorFactory  *interceptorFactory.ArgInterceptedDataFactory
	globalTxThrottler      process.InterceptorThrottler
	maxTxNonceDeltaAllowed int
=======
	"github.com/ElrondNetwork/elrond-go/process/block/interceptors"
	"github.com/ElrondNetwork/elrond-go/process/dataValidators"
	"github.com/ElrondNetwork/elrond-go/process/factory"
	"github.com/ElrondNetwork/elrond-go/process/factory/containers"
	"github.com/ElrondNetwork/elrond-go/process/transaction"
	"github.com/ElrondNetwork/elrond-go/process/unsigned"
	"github.com/ElrondNetwork/elrond-go/sharding"
)

type interceptorsContainerFactory struct {
	accounts            state.AccountsAdapter
	shardCoordinator    sharding.Coordinator
	messenger           process.TopicHandler
	store               dataRetriever.StorageService
	marshalizer         marshal.Marshalizer
	hasher              hashing.Hasher
	keyGen              crypto.KeyGenerator
	singleSigner        crypto.SingleSigner
	multiSigner         crypto.MultiSigner
	dataPool            dataRetriever.PoolsHolder
	addrConverter       state.AddressConverter
	chronologyValidator process.ChronologyValidator
>>>>>>> Stashed changes
}

// NewInterceptorsContainerFactory is responsible for creating a new interceptors factory object
func NewInterceptorsContainerFactory(
	accounts state.AccountsAdapter,
	shardCoordinator sharding.Coordinator,
<<<<<<< Updated upstream
	nodesCoordinator sharding.NodesCoordinator,
=======
>>>>>>> Stashed changes
	messenger process.TopicHandler,
	store dataRetriever.StorageService,
	marshalizer marshal.Marshalizer,
	hasher hashing.Hasher,
	keyGen crypto.KeyGenerator,
	singleSigner crypto.SingleSigner,
	multiSigner crypto.MultiSigner,
	dataPool dataRetriever.PoolsHolder,
	addrConverter state.AddressConverter,
<<<<<<< Updated upstream
	maxTxNonceDeltaAllowed int,
	txFeeHandler process.FeeHandler,
=======
	chronologyValidator process.ChronologyValidator,
>>>>>>> Stashed changes
) (*interceptorsContainerFactory, error) {
	if accounts == nil || accounts.IsInterfaceNil() {
		return nil, process.ErrNilAccountsAdapter
	}
	if shardCoordinator == nil || shardCoordinator.IsInterfaceNil() {
		return nil, process.ErrNilShardCoordinator
	}
<<<<<<< Updated upstream
	if messenger == nil || messenger.IsInterfaceNil() {
=======
	if messenger == nil {
>>>>>>> Stashed changes
		return nil, process.ErrNilMessenger
	}
	if store == nil || store.IsInterfaceNil() {
		return nil, process.ErrNilBlockChain
	}
	if marshalizer == nil || marshalizer.IsInterfaceNil() {
		return nil, process.ErrNilMarshalizer
	}
	if hasher == nil || hasher.IsInterfaceNil() {
		return nil, process.ErrNilHasher
	}
	if keyGen == nil || keyGen.IsInterfaceNil() {
		return nil, process.ErrNilKeyGen
	}
	if singleSigner == nil || singleSigner.IsInterfaceNil() {
		return nil, process.ErrNilSingleSigner
	}
	if multiSigner == nil || multiSigner.IsInterfaceNil() {
		return nil, process.ErrNilMultiSigVerifier
	}
	if dataPool == nil || dataPool.IsInterfaceNil() {
		return nil, process.ErrNilDataPoolHolder
	}
	if addrConverter == nil || addrConverter.IsInterfaceNil() {
		return nil, process.ErrNilAddressConverter
	}
<<<<<<< Updated upstream
	if nodesCoordinator == nil || nodesCoordinator.IsInterfaceNil() {
		return nil, process.ErrNilNodesCoordinator
	}
	if txFeeHandler == nil || txFeeHandler.IsInterfaceNil() {
		return nil, process.ErrNilEconomicsFeeHandler
	}

	argInterceptorFactory := &interceptorFactory.ArgInterceptedDataFactory{
		Marshalizer:      marshalizer,
		Hasher:           hasher,
		ShardCoordinator: shardCoordinator,
		MultiSigVerifier: multiSigner,
		NodesCoordinator: nodesCoordinator,
		KeyGen:           keyGen,
		Signer:           singleSigner,
		AddrConv:         addrConverter,
		FeeHandler:       txFeeHandler,
	}

	icf := &interceptorsContainerFactory{
		accounts:               accounts,
		shardCoordinator:       shardCoordinator,
		messenger:              messenger,
		store:                  store,
		marshalizer:            marshalizer,
		hasher:                 hasher,
		keyGen:                 keyGen,
		singleSigner:           singleSigner,
		multiSigner:            multiSigner,
		dataPool:               dataPool,
		addrConverter:          addrConverter,
		nodesCoordinator:       nodesCoordinator,
		argInterceptorFactory:  argInterceptorFactory,
		maxTxNonceDeltaAllowed: maxTxNonceDeltaAllowed,
	}

	var err error
	icf.globalTxThrottler, err = throttler.NewNumGoRoutineThrottler(numGoRoutines)
	if err != nil {
		return nil, err
	}

	return icf, nil
=======
	if chronologyValidator == nil || chronologyValidator.IsInterfaceNil() {
		return nil, process.ErrNilChronologyValidator
	}

	return &interceptorsContainerFactory{
		accounts:            accounts,
		shardCoordinator:    shardCoordinator,
		messenger:           messenger,
		store:               store,
		marshalizer:         marshalizer,
		hasher:              hasher,
		keyGen:              keyGen,
		singleSigner:        singleSigner,
		multiSigner:         multiSigner,
		dataPool:            dataPool,
		addrConverter:       addrConverter,
		chronologyValidator: chronologyValidator,
	}, nil
>>>>>>> Stashed changes
}

// Create returns an interceptor container that will hold all interceptors in the system
func (icf *interceptorsContainerFactory) Create() (process.InterceptorsContainer, error) {
	container := containers.NewInterceptorsContainer()

	keys, interceptorSlice, err := icf.generateTxInterceptors()
	if err != nil {
		return nil, err
	}

	err = container.AddMultiple(keys, interceptorSlice)
	if err != nil {
		return nil, err
	}

	keys, interceptorSlice, err = icf.generateUnsignedTxsInterceptors()
	if err != nil {
		return nil, err
	}

	err = container.AddMultiple(keys, interceptorSlice)
	if err != nil {
		return nil, err
	}

<<<<<<< Updated upstream
	keys, interceptorSlice, err = icf.generateRewardTxInterceptors()
=======
	keys, interceptorSlice, err = icf.generateHdrInterceptor()
>>>>>>> Stashed changes
	if err != nil {
		return nil, err
	}

	err = container.AddMultiple(keys, interceptorSlice)
	if err != nil {
		return nil, err
	}

<<<<<<< Updated upstream
	keys, interceptorSlice, err = icf.generateHdrInterceptor()
=======
	keys, interceptorSlice, err = icf.generateMiniBlocksInterceptors()
>>>>>>> Stashed changes
	if err != nil {
		return nil, err
	}

	err = container.AddMultiple(keys, interceptorSlice)
	if err != nil {
		return nil, err
	}

<<<<<<< Updated upstream
	keys, interceptorSlice, err = icf.generateMiniBlocksInterceptors()
=======
	keys, interceptorSlice, err = icf.generatePeerChBlockBodyInterceptor()
>>>>>>> Stashed changes
	if err != nil {
		return nil, err
	}

	err = container.AddMultiple(keys, interceptorSlice)
	if err != nil {
		return nil, err
	}

	keys, interceptorSlice, err = icf.generateMetachainHeaderInterceptor()
	if err != nil {
		return nil, err
	}

	err = container.AddMultiple(keys, interceptorSlice)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func (icf *interceptorsContainerFactory) createTopicAndAssignHandler(
	topic string,
	interceptor process.Interceptor,
	createChannel bool,
) (process.Interceptor, error) {

	err := icf.messenger.CreateTopic(topic, createChannel)
	if err != nil {
		return nil, err
	}

	return interceptor, icf.messenger.RegisterMessageProcessor(topic, interceptor)
}

//------- Tx interceptors

func (icf *interceptorsContainerFactory) generateTxInterceptors() ([]string, []process.Interceptor, error) {
	shardC := icf.shardCoordinator

	noOfShards := shardC.NumberOfShards()

	keys := make([]string, noOfShards)
	interceptorSlice := make([]process.Interceptor, noOfShards)

	for idx := uint32(0); idx < noOfShards; idx++ {
		identifierTx := factory.TransactionTopic + shardC.CommunicationIdentifier(idx)

		interceptor, err := icf.createOneTxInterceptor(identifierTx)
		if err != nil {
			return nil, nil, err
		}

		keys[int(idx)] = identifierTx
		interceptorSlice[int(idx)] = interceptor
	}

	//tx interceptor for metachain topic
	identifierTx := factory.TransactionTopic + shardC.CommunicationIdentifier(sharding.MetachainShardId)

	interceptor, err := icf.createOneTxInterceptor(identifierTx)
	if err != nil {
		return nil, nil, err
	}

	keys = append(keys, identifierTx)
	interceptorSlice = append(interceptorSlice, interceptor)
	return keys, interceptorSlice, nil
}

<<<<<<< Updated upstream
func (icf *interceptorsContainerFactory) createOneTxInterceptor(topic string) (process.Interceptor, error) {
	txValidator, err := dataValidators.NewTxValidator(icf.accounts, icf.shardCoordinator, icf.maxTxNonceDeltaAllowed)
=======
func (icf *interceptorsContainerFactory) createOneTxInterceptor(identifier string) (process.Interceptor, error) {
	//TODO implement other TxHandlerProcessValidator that will check the tx nonce against account's nonce
	txValidator, err := dataValidators.NewStorageTxValidator(icf.hasher, icf.marshalizer, icf.store.GetStorer(dataRetriever.TransactionUnit))
>>>>>>> Stashed changes
	if err != nil {
		return nil, err
	}

<<<<<<< Updated upstream
	argProcessor := &processor.ArgTxInterceptorProcessor{
		ShardedDataCache: icf.dataPool.Transactions(),
		TxValidator:      txValidator,
	}
	txProcessor, err := processor.NewTxInterceptorProcessor(argProcessor)
	if err != nil {
		return nil, err
	}

	txFactory, err := interceptorFactory.NewShardInterceptedDataFactory(
		icf.argInterceptorFactory,
		interceptorFactory.InterceptedTx,
	)
	if err != nil {
		return nil, err
	}

	interceptor, err := interceptors.NewMultiDataInterceptor(
		icf.marshalizer,
		txFactory,
		txProcessor,
		icf.globalTxThrottler,
	)
=======
	interceptor, err := transaction.NewTxInterceptor(
		icf.marshalizer,
		icf.dataPool.Transactions(),
		txValidator,
		icf.addrConverter,
		icf.hasher,
		icf.singleSigner,
		icf.keyGen,
		icf.shardCoordinator)

>>>>>>> Stashed changes
	if err != nil {
		return nil, err
	}

<<<<<<< Updated upstream
	return icf.createTopicAndAssignHandler(topic, interceptor, true)
=======
	return icf.createTopicAndAssignHandler(identifier, interceptor, true)
>>>>>>> Stashed changes
}

//------- Unsigned transactions interceptors

func (icf *interceptorsContainerFactory) generateUnsignedTxsInterceptors() ([]string, []process.Interceptor, error) {
	shardC := icf.shardCoordinator

	noOfShards := shardC.NumberOfShards()

	keys := make([]string, noOfShards)
	interceptorSlice := make([]process.Interceptor, noOfShards)

	for idx := uint32(0); idx < noOfShards; idx++ {
		identifierScr := factory.UnsignedTransactionTopic + shardC.CommunicationIdentifier(idx)

		interceptor, err := icf.createOneUnsignedTxInterceptor(identifierScr)
		if err != nil {
			return nil, nil, err
		}

		keys[int(idx)] = identifierScr
		interceptorSlice[int(idx)] = interceptor
	}

	identifierTx := factory.UnsignedTransactionTopic + shardC.CommunicationIdentifier(sharding.MetachainShardId)

	interceptor, err := icf.createOneUnsignedTxInterceptor(identifierTx)
	if err != nil {
		return nil, nil, err
	}

	keys = append(keys, identifierTx)
	interceptorSlice = append(interceptorSlice, interceptor)
	return keys, interceptorSlice, nil
}

<<<<<<< Updated upstream
func (icf *interceptorsContainerFactory) createOneUnsignedTxInterceptor(topic string) (process.Interceptor, error) {
	//TODO replace the nil tx validator with white list validator
	txValidator, err := mock.NewNilTxValidator()
	if err != nil {
		return nil, err
	}

	argProcessor := &processor.ArgTxInterceptorProcessor{
		ShardedDataCache: icf.dataPool.UnsignedTransactions(),
		TxValidator:      txValidator,
	}
	txProcessor, err := processor.NewTxInterceptorProcessor(argProcessor)
	if err != nil {
		return nil, err
	}

	txFactory, err := interceptorFactory.NewShardInterceptedDataFactory(
		icf.argInterceptorFactory,
		interceptorFactory.InterceptedUnsignedTx,
	)
	if err != nil {
		return nil, err
	}

	interceptor, err := interceptors.NewMultiDataInterceptor(
		icf.marshalizer,
		txFactory,
		txProcessor,
		icf.globalTxThrottler,
	)
	if err != nil {
		return nil, err
	}

	return icf.createTopicAndAssignHandler(topic, interceptor, true)
}

//------- Reward transactions interceptors

func (icf *interceptorsContainerFactory) generateRewardTxInterceptors() ([]string, []process.Interceptor, error) {
	shardC := icf.shardCoordinator

	noOfShards := shardC.NumberOfShards()

	keys := make([]string, noOfShards)
	interceptorSlice := make([]process.Interceptor, noOfShards)

	for idx := uint32(0); idx < noOfShards; idx++ {
		identifierScr := factory.RewardsTransactionTopic + shardC.CommunicationIdentifier(idx)

		interceptor, err := icf.createOneRewardTxInterceptor(identifierScr)
		if err != nil {
			return nil, nil, err
		}

		keys[int(idx)] = identifierScr
		interceptorSlice[int(idx)] = interceptor
	}

	identifierTx := factory.RewardsTransactionTopic + shardC.CommunicationIdentifier(sharding.MetachainShardId)

	interceptor, err := icf.createOneRewardTxInterceptor(identifierTx)
	if err != nil {
		return nil, nil, err
	}

	keys = append(keys, identifierTx)
	interceptorSlice = append(interceptorSlice, interceptor)

	return keys, interceptorSlice, nil
}

func (icf *interceptorsContainerFactory) createOneRewardTxInterceptor(identifier string) (process.Interceptor, error) {
	rewardTxStorer := icf.store.GetStorer(dataRetriever.RewardTransactionUnit)

	interceptor, err := rewardTransaction.NewRewardTxInterceptor(
		icf.marshalizer,
		icf.dataPool.RewardTransactions(),
		rewardTxStorer,
		icf.addrConverter,
		icf.hasher,
		icf.shardCoordinator,
	)
=======
func (icf *interceptorsContainerFactory) createOneUnsignedTxInterceptor(identifier string) (process.Interceptor, error) {
	uTxStorer := icf.store.GetStorer(dataRetriever.UnsignedTransactionUnit)

	interceptor, err := unsigned.NewUnsignedTxInterceptor(
		icf.marshalizer,
		icf.dataPool.UnsignedTransactions(),
		uTxStorer,
		icf.addrConverter,
		icf.hasher,
		icf.shardCoordinator)

>>>>>>> Stashed changes
	if err != nil {
		return nil, err
	}

	return icf.createTopicAndAssignHandler(identifier, interceptor, true)
}

//------- Hdr interceptor

func (icf *interceptorsContainerFactory) generateHdrInterceptor() ([]string, []process.Interceptor, error) {
	shardC := icf.shardCoordinator
	//TODO implement other HeaderHandlerProcessValidator that will check the header's nonce
	// against blockchain's latest nonce - k finality
	hdrValidator, err := dataValidators.NewNilHeaderValidator()
	if err != nil {
		return nil, nil, err
	}

<<<<<<< Updated upstream
	hdrFactory, err := interceptorFactory.NewShardInterceptedDataFactory(
		icf.argInterceptorFactory,
		interceptorFactory.InterceptedShardHeader,
	)
	if err != nil {
		return nil, nil, err
	}

	argProcessor := &processor.ArgHdrInterceptorProcessor{
		Headers:       icf.dataPool.Headers(),
		HeadersNonces: icf.dataPool.HeadersNonces(),
		HdrValidator:  hdrValidator,
	}
	hdrProcessor, err := processor.NewHdrInterceptorProcessor(argProcessor)
	if err != nil {
		return nil, nil, err
	}

	//only one intrashard header topic
	interceptor, err := interceptors.NewSingleDataInterceptor(
		hdrFactory,
		hdrProcessor,
		icf.globalTxThrottler,
=======
	//only one intrashard header topic
	identifierHdr := factory.HeadersTopic + shardC.CommunicationIdentifier(shardC.SelfId())
	interceptor, err := interceptors.NewHeaderInterceptor(
		icf.marshalizer,
		icf.dataPool.Headers(),
		icf.dataPool.HeadersNonces(),
		hdrValidator,
		icf.multiSigner,
		icf.hasher,
		icf.shardCoordinator,
		icf.chronologyValidator,
>>>>>>> Stashed changes
	)
	if err != nil {
		return nil, nil, err
	}
<<<<<<< Updated upstream

	identifierHdr := factory.HeadersTopic + shardC.CommunicationIdentifier(shardC.SelfId())
=======
>>>>>>> Stashed changes
	_, err = icf.createTopicAndAssignHandler(identifierHdr, interceptor, true)
	if err != nil {
		return nil, nil, err
	}

	return []string{identifierHdr}, []process.Interceptor{interceptor}, nil
}

//------- MiniBlocks interceptors

func (icf *interceptorsContainerFactory) generateMiniBlocksInterceptors() ([]string, []process.Interceptor, error) {
	shardC := icf.shardCoordinator
	noOfShards := shardC.NumberOfShards()
<<<<<<< Updated upstream
	keys := make([]string, noOfShards+1)
	interceptorsSlice := make([]process.Interceptor, noOfShards+1)
=======
	keys := make([]string, noOfShards)
	interceptorSlice := make([]process.Interceptor, noOfShards)
>>>>>>> Stashed changes

	for idx := uint32(0); idx < noOfShards; idx++ {
		identifierMiniBlocks := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(idx)

		interceptor, err := icf.createOneMiniBlocksInterceptor(identifierMiniBlocks)
		if err != nil {
			return nil, nil, err
		}

		keys[int(idx)] = identifierMiniBlocks
<<<<<<< Updated upstream
		interceptorsSlice[int(idx)] = interceptor
	}

	identifierMiniBlocks := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(sharding.MetachainShardId)

	interceptor, err := icf.createOneMiniBlocksInterceptor(identifierMiniBlocks)
	if err != nil {
		return nil, nil, err
	}

	keys[noOfShards] = identifierMiniBlocks
	interceptorsSlice[noOfShards] = interceptor

	return keys, interceptorsSlice, nil
}

func (icf *interceptorsContainerFactory) createOneMiniBlocksInterceptor(topic string) (process.Interceptor, error) {
	argProcessor := &processor.ArgTxBodyInterceptorProcessor{
		MiniblockCache:   icf.dataPool.MiniBlocks(),
		Marshalizer:      icf.marshalizer,
		Hasher:           icf.hasher,
		ShardCoordinator: icf.shardCoordinator,
	}
	txBlockBodyProcessor, err := processor.NewTxBodyInterceptorProcessor(argProcessor)
=======
		interceptorSlice[int(idx)] = interceptor
	}

	return keys, interceptorSlice, nil
}

func (icf *interceptorsContainerFactory) createOneMiniBlocksInterceptor(identifier string) (process.Interceptor, error) {
	txBlockBodyStorer := icf.store.GetStorer(dataRetriever.MiniBlockUnit)

	interceptor, err := interceptors.NewTxBlockBodyInterceptor(
		icf.marshalizer,
		icf.dataPool.MiniBlocks(),
		txBlockBodyStorer,
		icf.hasher,
		icf.shardCoordinator,
	)

>>>>>>> Stashed changes
	if err != nil {
		return nil, err
	}

<<<<<<< Updated upstream
	txFactory, err := interceptorFactory.NewShardInterceptedDataFactory(
		icf.argInterceptorFactory,
		interceptorFactory.InterceptedTxBlockBody,
	)
	if err != nil {
		return nil, err
	}

	interceptor, err := interceptors.NewSingleDataInterceptor(
		txFactory,
		txBlockBodyProcessor,
		icf.globalTxThrottler,
	)
	if err != nil {
		return nil, err
	}

	return icf.createTopicAndAssignHandler(topic, interceptor, true)
=======
	return icf.createTopicAndAssignHandler(identifier, interceptor, true)
}

//------- PeerChBlocks interceptor

func (icf *interceptorsContainerFactory) generatePeerChBlockBodyInterceptor() ([]string, []process.Interceptor, error) {
	shardC := icf.shardCoordinator

	//only one intrashard peer change blocks topic
	identifierPeerCh := factory.PeerChBodyTopic + shardC.CommunicationIdentifier(shardC.SelfId())
	peerBlockBodyStorer := icf.store.GetStorer(dataRetriever.PeerChangesUnit)

	interceptor, err := interceptors.NewPeerBlockBodyInterceptor(
		icf.marshalizer,
		icf.dataPool.PeerChangesBlocks(),
		peerBlockBodyStorer,
		icf.hasher,
		shardC,
	)
	if err != nil {
		return nil, nil, err
	}
	_, err = icf.createTopicAndAssignHandler(identifierPeerCh, interceptor, true)
	if err != nil {
		return nil, nil, err
	}

	return []string{identifierPeerCh}, []process.Interceptor{interceptor}, nil
>>>>>>> Stashed changes
}

//------- MetachainHeader interceptors

func (icf *interceptorsContainerFactory) generateMetachainHeaderInterceptor() ([]string, []process.Interceptor, error) {
	identifierHdr := factory.MetachainBlocksTopic
	//TODO implement other HeaderHandlerProcessValidator that will check the header's nonce
	// against blockchain's latest nonce - k finality
	hdrValidator, err := dataValidators.NewNilHeaderValidator()
	if err != nil {
		return nil, nil, err
	}

<<<<<<< Updated upstream
	hdrFactory, err := interceptorFactory.NewShardInterceptedDataFactory(
		icf.argInterceptorFactory,
		interceptorFactory.InterceptedMetaHeader,
	)
	if err != nil {
		return nil, nil, err
	}

	argProcessor := &processor.ArgHdrInterceptorProcessor{
		Headers:       icf.dataPool.MetaBlocks(),
		HeadersNonces: icf.dataPool.HeadersNonces(),
		HdrValidator:  hdrValidator,
	}
	hdrProcessor, err := processor.NewHdrInterceptorProcessor(argProcessor)
	if err != nil {
		return nil, nil, err
	}

	//only one metachain header topic
	interceptor, err := interceptors.NewSingleDataInterceptor(
		hdrFactory,
		hdrProcessor,
		icf.globalTxThrottler,
=======
	interceptor, err := interceptors.NewMetachainHeaderInterceptor(
		icf.marshalizer,
		icf.dataPool.MetaBlocks(),
		icf.dataPool.HeadersNonces(),
		hdrValidator,
		icf.multiSigner,
		icf.hasher,
		icf.shardCoordinator,
		icf.chronologyValidator,
>>>>>>> Stashed changes
	)
	if err != nil {
		return nil, nil, err
	}
<<<<<<< Updated upstream

=======
>>>>>>> Stashed changes
	_, err = icf.createTopicAndAssignHandler(identifierHdr, interceptor, true)
	if err != nil {
		return nil, nil, err
	}

	return []string{identifierHdr}, []process.Interceptor{interceptor}, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (icf *interceptorsContainerFactory) IsInterfaceNil() bool {
	if icf == nil {
		return true
	}
	return false
}
