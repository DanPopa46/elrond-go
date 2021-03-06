package shard

import (
	"github.com/ElrondNetwork/elrond-go/core/random"
	"github.com/ElrondNetwork/elrond-go/data/typeConverters"
	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/dataRetriever/factory/containers"
	"github.com/ElrondNetwork/elrond-go/dataRetriever/resolvers"
	"github.com/ElrondNetwork/elrond-go/dataRetriever/resolvers/topicResolverSender"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/process/factory"
	"github.com/ElrondNetwork/elrond-go/sharding"
)

const emptyExcludePeersOnTopic = ""

type resolversContainerFactory struct {
	shardCoordinator         sharding.Coordinator
	messenger                dataRetriever.TopicMessageHandler
	store                    dataRetriever.StorageService
	marshalizer              marshal.Marshalizer
	dataPools                dataRetriever.PoolsHolder
	uint64ByteSliceConverter typeConverters.Uint64ByteSliceConverter
	intRandomizer            dataRetriever.IntRandomizer
	dataPacker               dataRetriever.DataPacker
}

// NewResolversContainerFactory creates a new container filled with topic resolvers
func NewResolversContainerFactory(
	shardCoordinator sharding.Coordinator,
	messenger dataRetriever.TopicMessageHandler,
	store dataRetriever.StorageService,
	marshalizer marshal.Marshalizer,
	dataPools dataRetriever.PoolsHolder,
	uint64ByteSliceConverter typeConverters.Uint64ByteSliceConverter,
	dataPacker dataRetriever.DataPacker,
) (*resolversContainerFactory, error) {

	if shardCoordinator == nil || shardCoordinator.IsInterfaceNil() {
		return nil, dataRetriever.ErrNilShardCoordinator
	}
	if messenger == nil || messenger.IsInterfaceNil() {
		return nil, dataRetriever.ErrNilMessenger
	}
	if store == nil || store.IsInterfaceNil() {
		return nil, dataRetriever.ErrNilTxStorage
	}
	if marshalizer == nil || marshalizer.IsInterfaceNil() {
		return nil, dataRetriever.ErrNilMarshalizer
	}
	if dataPools == nil || dataPools.IsInterfaceNil() {
		return nil, dataRetriever.ErrNilDataPoolHolder
	}
	if uint64ByteSliceConverter == nil || uint64ByteSliceConverter.IsInterfaceNil() {
		return nil, dataRetriever.ErrNilUint64ByteSliceConverter
	}
	if dataPacker == nil || dataPacker.IsInterfaceNil() {
		return nil, dataRetriever.ErrNilDataPacker
	}

	return &resolversContainerFactory{
		shardCoordinator:         shardCoordinator,
		messenger:                messenger,
		store:                    store,
		marshalizer:              marshalizer,
		dataPools:                dataPools,
		uint64ByteSliceConverter: uint64ByteSliceConverter,
		intRandomizer:            &random.ConcurrentSafeIntRandomizer{},
		dataPacker:               dataPacker,
	}, nil
}

<<<<<<< Updated upstream
// Create returns a resolver container that will hold all resolvers in the system
func (rcf *resolversContainerFactory) Create() (dataRetriever.ResolversContainer, error) {
	container := containers.NewResolversContainer()

	keys, resolverSlice, err := rcf.generateTxResolvers(
		factory.TransactionTopic,
		dataRetriever.TransactionUnit,
		rcf.dataPools.Transactions(),
	)
=======
// Create returns an interceptor container that will hold all interceptors in the system
func (rcf *resolversContainerFactory) Create() (dataRetriever.ResolversContainer, error) {
	container := containers.NewResolversContainer()

	keys, resolverSlice, err := rcf.generateTxResolvers(factory.TransactionTopic, dataRetriever.TransactionUnit, rcf.dataPools.Transactions())
>>>>>>> Stashed changes
	if err != nil {
		return nil, err
	}
	err = container.AddMultiple(keys, resolverSlice)
	if err != nil {
		return nil, err
	}

	keys, resolverSlice, err = rcf.generateTxResolvers(
		factory.UnsignedTransactionTopic,
		dataRetriever.UnsignedTransactionUnit,
		rcf.dataPools.UnsignedTransactions(),
	)
	if err != nil {
		return nil, err
	}
	err = container.AddMultiple(keys, resolverSlice)
	if err != nil {
		return nil, err
	}

<<<<<<< Updated upstream
	keys, resolverSlice, err = rcf.generateTxResolvers(
		factory.RewardsTransactionTopic,
		dataRetriever.RewardTransactionUnit,
		rcf.dataPools.RewardTransactions(),
	)
	if err != nil {
		return nil, err
	}

	err = container.AddMultiple(keys, resolverSlice)
	if err != nil {
		return nil, err
	}

=======
>>>>>>> Stashed changes
	keys, resolverSlice, err = rcf.generateHdrResolver()
	if err != nil {
		return nil, err
	}
	err = container.AddMultiple(keys, resolverSlice)
	if err != nil {
		return nil, err
	}

	keys, resolverSlice, err = rcf.generateMiniBlocksResolvers()
	if err != nil {
		return nil, err
	}
	err = container.AddMultiple(keys, resolverSlice)
	if err != nil {
		return nil, err
	}

	keys, resolverSlice, err = rcf.generatePeerChBlockBodyResolver()
	if err != nil {
		return nil, err
	}
	err = container.AddMultiple(keys, resolverSlice)
	if err != nil {
		return nil, err
	}

	keys, resolverSlice, err = rcf.generateMetachainShardHeaderResolver()
	if err != nil {
		return nil, err
	}
	err = container.AddMultiple(keys, resolverSlice)
	if err != nil {
		return nil, err
	}

	keys, resolverSlice, err = rcf.generateMetablockHeaderResolver()
	if err != nil {
		return nil, err
	}
	err = container.AddMultiple(keys, resolverSlice)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func (rcf *resolversContainerFactory) createTopicAndAssignHandler(
	topicName string,
	resolver dataRetriever.Resolver,
	createChannel bool,
) (dataRetriever.Resolver, error) {

	err := rcf.messenger.CreateTopic(topicName, createChannel)
	if err != nil {
		return nil, err
	}

	return resolver, rcf.messenger.RegisterMessageProcessor(topicName, resolver)
}

//------- Tx resolvers

func (rcf *resolversContainerFactory) generateTxResolvers(
	topic string,
	unit dataRetriever.UnitType,
	dataPool dataRetriever.ShardedDataCacherNotifier,
) ([]string, []dataRetriever.Resolver, error) {

	shardC := rcf.shardCoordinator

	noOfShards := shardC.NumberOfShards()

<<<<<<< Updated upstream
	keys := make([]string, noOfShards+1)
	resolverSlice := make([]dataRetriever.Resolver, noOfShards+1)
=======
	keys := make([]string, noOfShards)
	resolverSlice := make([]dataRetriever.Resolver, noOfShards)
>>>>>>> Stashed changes

	for idx := uint32(0); idx < noOfShards; idx++ {
		identifierTx := topic + shardC.CommunicationIdentifier(idx)
		excludePeersFromTopic := topic + shardC.CommunicationIdentifier(shardC.SelfId())

		resolver, err := rcf.createTxResolver(identifierTx, excludePeersFromTopic, unit, dataPool)
		if err != nil {
			return nil, nil, err
		}

		resolverSlice[idx] = resolver
		keys[idx] = identifierTx
	}

<<<<<<< Updated upstream
	identifierTx := topic + shardC.CommunicationIdentifier(sharding.MetachainShardId)
	excludePeersFromTopic := topic + shardC.CommunicationIdentifier(shardC.SelfId())

	resolver, err := rcf.createTxResolver(identifierTx, excludePeersFromTopic, unit, dataPool)
	if err != nil {
		return nil, nil, err
	}

	resolverSlice[noOfShards] = resolver
	keys[noOfShards] = identifierTx

=======
>>>>>>> Stashed changes
	return keys, resolverSlice, nil
}

func (rcf *resolversContainerFactory) createTxResolver(
	topic string,
	excludedTopic string,
	unit dataRetriever.UnitType,
	dataPool dataRetriever.ShardedDataCacherNotifier,
) (dataRetriever.Resolver, error) {

	txStorer := rcf.store.GetStorer(unit)

<<<<<<< Updated upstream
	resolverSender, err := rcf.createOneResolverSender(topic, excludedTopic)
=======
	peerListCreator, err := topicResolverSender.NewDiffPeerListCreator(rcf.messenger, topic, excludedTopic)
	if err != nil {
		return nil, err
	}

	//TODO instantiate topic sender resolver with the shard IDs for which this resolver is supposed to serve the data
	// this will improve the serving of transactions as the searching will be done only on 2 sharded data units
	resolverSender, err := topicResolverSender.NewTopicResolverSender(
		rcf.messenger,
		topic,
		peerListCreator,
		rcf.marshalizer,
		rcf.intRandomizer,
		uint32(0),
	)
>>>>>>> Stashed changes
	if err != nil {
		return nil, err
	}

	resolver, err := resolvers.NewTxResolver(
		resolverSender,
		dataPool,
		txStorer,
		rcf.marshalizer,
		rcf.dataPacker,
	)
	if err != nil {
		return nil, err
	}

	//add on the request topic
	return rcf.createTopicAndAssignHandler(
		topic+resolverSender.TopicRequestSuffix(),
		resolver,
		false)
}

//------- Hdr resolver

func (rcf *resolversContainerFactory) generateHdrResolver() ([]string, []dataRetriever.Resolver, error) {
	shardC := rcf.shardCoordinator

	//only one intrashard header topic
	identifierHdr := factory.HeadersTopic + shardC.CommunicationIdentifier(shardC.SelfId())

	peerListCreator, err := topicResolverSender.NewDiffPeerListCreator(rcf.messenger, identifierHdr, emptyExcludePeersOnTopic)
	if err != nil {
		return nil, nil, err
	}

	hdrStorer := rcf.store.GetStorer(dataRetriever.BlockHeaderUnit)
	resolverSender, err := topicResolverSender.NewTopicResolverSender(
		rcf.messenger,
		identifierHdr,
		peerListCreator,
		rcf.marshalizer,
		rcf.intRandomizer,
		shardC.SelfId(),
	)
	if err != nil {
		return nil, nil, err
	}

	hdrNonceHashDataUnit := dataRetriever.ShardHdrNonceHashDataUnit + dataRetriever.UnitType(shardC.SelfId())
	hdrNonceStore := rcf.store.GetStorer(hdrNonceHashDataUnit)
	resolver, err := resolvers.NewHeaderResolver(
		resolverSender,
		rcf.dataPools.Headers(),
		rcf.dataPools.HeadersNonces(),
		hdrStorer,
		hdrNonceStore,
		rcf.marshalizer,
		rcf.uint64ByteSliceConverter,
	)
	if err != nil {
		return nil, nil, err
	}
	//add on the request topic
	_, err = rcf.createTopicAndAssignHandler(
		identifierHdr+resolverSender.TopicRequestSuffix(),
		resolver,
		false)
	if err != nil {
		return nil, nil, err
	}

	err = rcf.createTopicHeadersForMetachain()
	if err != nil {
		return nil, nil, err
	}

	return []string{identifierHdr}, []dataRetriever.Resolver{resolver}, nil
}

func (rcf *resolversContainerFactory) createTopicHeadersForMetachain() error {
	shardC := rcf.shardCoordinator
	identifierHdr := factory.ShardHeadersForMetachainTopic + shardC.CommunicationIdentifier(sharding.MetachainShardId)

	return rcf.messenger.CreateTopic(identifierHdr, true)
}

//------- MiniBlocks resolvers

func (rcf *resolversContainerFactory) generateMiniBlocksResolvers() ([]string, []dataRetriever.Resolver, error) {
	shardC := rcf.shardCoordinator
	noOfShards := shardC.NumberOfShards()
<<<<<<< Updated upstream
	keys := make([]string, noOfShards+1)
	resolverSlice := make([]dataRetriever.Resolver, noOfShards+1)
=======
	keys := make([]string, noOfShards)
	resolverSlice := make([]dataRetriever.Resolver, noOfShards)
>>>>>>> Stashed changes

	for idx := uint32(0); idx < noOfShards; idx++ {
		identifierMiniBlocks := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(idx)
		excludePeersFromTopic := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(shardC.SelfId())

		resolver, err := rcf.createMiniBlocksResolver(identifierMiniBlocks, excludePeersFromTopic)
		if err != nil {
			return nil, nil, err
		}

		resolverSlice[idx] = resolver
		keys[idx] = identifierMiniBlocks
	}

<<<<<<< Updated upstream
	identifierMiniBlocks := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(sharding.MetachainShardId)
	excludePeersFromTopic := factory.MiniBlocksTopic + shardC.CommunicationIdentifier(shardC.SelfId())

	resolver, err := rcf.createMiniBlocksResolver(identifierMiniBlocks, excludePeersFromTopic)
	if err != nil {
		return nil, nil, err
	}

	resolverSlice[noOfShards] = resolver
	keys[noOfShards] = identifierMiniBlocks

=======
>>>>>>> Stashed changes
	return keys, resolverSlice, nil
}

func (rcf *resolversContainerFactory) createMiniBlocksResolver(topic string, excludedTopic string) (dataRetriever.Resolver, error) {
	miniBlocksStorer := rcf.store.GetStorer(dataRetriever.MiniBlockUnit)

<<<<<<< Updated upstream
	resolverSender, err := rcf.createOneResolverSender(topic, excludedTopic)
=======
	peerListCreator, err := topicResolverSender.NewDiffPeerListCreator(rcf.messenger, topic, excludedTopic)
	if err != nil {
		return nil, err
	}

	resolverSender, err := topicResolverSender.NewTopicResolverSender(
		rcf.messenger,
		topic,
		peerListCreator,
		rcf.marshalizer,
		rcf.intRandomizer,
		uint32(0),
	)
>>>>>>> Stashed changes
	if err != nil {
		return nil, err
	}

	txBlkResolver, err := resolvers.NewGenericBlockBodyResolver(
		resolverSender,
		rcf.dataPools.MiniBlocks(),
		miniBlocksStorer,
		rcf.marshalizer,
	)
	if err != nil {
		return nil, err
	}

	//add on the request topic
	return rcf.createTopicAndAssignHandler(
		topic+resolverSender.TopicRequestSuffix(),
		txBlkResolver,
		false)
}

//------- PeerChBlocks resolvers

func (rcf *resolversContainerFactory) generatePeerChBlockBodyResolver() ([]string, []dataRetriever.Resolver, error) {
	shardC := rcf.shardCoordinator

	//only one intrashard peer change blocks topic
	identifierPeerCh := factory.PeerChBodyTopic + shardC.CommunicationIdentifier(shardC.SelfId())
	peerBlockBodyStorer := rcf.store.GetStorer(dataRetriever.PeerChangesUnit)

	peerListCreator, err := topicResolverSender.NewDiffPeerListCreator(rcf.messenger, identifierPeerCh, emptyExcludePeersOnTopic)
	if err != nil {
		return nil, nil, err
	}

	resolverSender, err := topicResolverSender.NewTopicResolverSender(
		rcf.messenger,
		identifierPeerCh,
		peerListCreator,
		rcf.marshalizer,
		rcf.intRandomizer,
		shardC.SelfId(),
	)
	if err != nil {
		return nil, nil, err
	}

	resolver, err := resolvers.NewGenericBlockBodyResolver(
		resolverSender,
		rcf.dataPools.MiniBlocks(),
		peerBlockBodyStorer,
		rcf.marshalizer,
	)
	if err != nil {
		return nil, nil, err
	}
	//add on the request topic
	_, err = rcf.createTopicAndAssignHandler(
		identifierPeerCh+resolverSender.TopicRequestSuffix(),
		resolver,
		false)
	if err != nil {
		return nil, nil, err
	}

	return []string{identifierPeerCh}, []dataRetriever.Resolver{resolver}, nil
}

//------- MetachainShardHeaderResolvers

func (rcf *resolversContainerFactory) generateMetachainShardHeaderResolver() ([]string, []dataRetriever.Resolver, error) {
	shardC := rcf.shardCoordinator

	//only one metachain header topic
	//example: shardHeadersForMetachain_0_META
	identifierHdr := factory.ShardHeadersForMetachainTopic + shardC.CommunicationIdentifier(sharding.MetachainShardId)
	peerListCreator, err := topicResolverSender.NewDiffPeerListCreator(rcf.messenger, identifierHdr, emptyExcludePeersOnTopic)
	if err != nil {
		return nil, nil, err
	}

	hdrStorer := rcf.store.GetStorer(dataRetriever.BlockHeaderUnit)
	resolverSender, err := topicResolverSender.NewTopicResolverSender(
		rcf.messenger,
		identifierHdr,
		peerListCreator,
		rcf.marshalizer,
		rcf.intRandomizer,
		shardC.SelfId(),
	)
	if err != nil {
		return nil, nil, err
	}

	hdrNonceHashDataUnit := dataRetriever.ShardHdrNonceHashDataUnit + dataRetriever.UnitType(shardC.SelfId())
	hdrNonceStore := rcf.store.GetStorer(hdrNonceHashDataUnit)
	resolver, err := resolvers.NewHeaderResolver(
		resolverSender,
		rcf.dataPools.Headers(),
		rcf.dataPools.HeadersNonces(),
		hdrStorer,
		hdrNonceStore,
		rcf.marshalizer,
		rcf.uint64ByteSliceConverter,
	)
	if err != nil {
		return nil, nil, err
	}

	//add on the request topic
	_, err = rcf.createTopicAndAssignHandler(
		identifierHdr+resolverSender.TopicRequestSuffix(),
		resolver,
		false)
	if err != nil {
		return nil, nil, err
	}

	return []string{identifierHdr}, []dataRetriever.Resolver{resolver}, nil
}

//------- MetaBlockHeaderResolvers

func (rcf *resolversContainerFactory) generateMetablockHeaderResolver() ([]string, []dataRetriever.Resolver, error) {
	shardC := rcf.shardCoordinator

	//only one metachain header block topic
	//this is: metachainBlocks
	identifierHdr := factory.MetachainBlocksTopic
	hdrStorer := rcf.store.GetStorer(dataRetriever.MetaBlockUnit)

	metaAndCrtShardTopic := factory.ShardHeadersForMetachainTopic + shardC.CommunicationIdentifier(sharding.MetachainShardId)
	excludedPeersOnTopic := factory.TransactionTopic + shardC.CommunicationIdentifier(shardC.SelfId())

	peerListCreator, err := topicResolverSender.NewDiffPeerListCreator(rcf.messenger, metaAndCrtShardTopic, excludedPeersOnTopic)
	if err != nil {
		return nil, nil, err
	}

	resolverSender, err := topicResolverSender.NewTopicResolverSender(
		rcf.messenger,
		identifierHdr,
		peerListCreator,
		rcf.marshalizer,
		rcf.intRandomizer,
		sharding.MetachainShardId,
	)
	if err != nil {
		return nil, nil, err
	}

	hdrNonceStore := rcf.store.GetStorer(dataRetriever.MetaHdrNonceHashDataUnit)
	resolver, err := resolvers.NewHeaderResolver(
		resolverSender,
		rcf.dataPools.MetaBlocks(),
		rcf.dataPools.HeadersNonces(),
		hdrStorer,
		hdrNonceStore,
		rcf.marshalizer,
		rcf.uint64ByteSliceConverter,
	)
	if err != nil {
		return nil, nil, err
	}

	//add on the request topic
	_, err = rcf.createTopicAndAssignHandler(
		identifierHdr+resolverSender.TopicRequestSuffix(),
		resolver,
		false)
	if err != nil {
		return nil, nil, err
	}

	return []string{identifierHdr}, []dataRetriever.Resolver{resolver}, nil
}

<<<<<<< Updated upstream
func (rcf *resolversContainerFactory) createOneResolverSender(
	topic string,
	excludedTopic string,
) (dataRetriever.TopicResolverSender, error) {

	peerListCreator, err := topicResolverSender.NewDiffPeerListCreator(rcf.messenger, topic, excludedTopic)
	if err != nil {
		return nil, err
	}

	//TODO instantiate topic sender resolver with the shard IDs for which this resolver is supposed to serve the data
	// this will improve the serving of transactions as the searching will be done only on 2 sharded data units
	resolverSender, err := topicResolverSender.NewTopicResolverSender(
		rcf.messenger,
		topic,
		peerListCreator,
		rcf.marshalizer,
		rcf.intRandomizer,
		uint32(0),
	)
	if err != nil {
		return nil, err
	}

	return resolverSender, nil
}

=======
>>>>>>> Stashed changes
// IsInterfaceNil returns true if there is no value under the interface
func (rcf *resolversContainerFactory) IsInterfaceNil() bool {
	if rcf == nil {
		return true
	}
	return false
}
