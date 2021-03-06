package factory

import (
	"errors"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/parsers"
	dataBlock "github.com/ElrondNetwork/elrond-go/data/block"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/epochStart"
	"github.com/ElrondNetwork/elrond-go/epochStart/bootstrap/disabled"
	metachainEpochStart "github.com/ElrondNetwork/elrond-go/epochStart/metachain"
	"github.com/ElrondNetwork/elrond-go/genesis"
	processDisabled "github.com/ElrondNetwork/elrond-go/genesis/process/disabled"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/process/block"
	"github.com/ElrondNetwork/elrond-go/process/block/postprocess"
	"github.com/ElrondNetwork/elrond-go/process/block/preprocess"
	"github.com/ElrondNetwork/elrond-go/process/coordinator"
	"github.com/ElrondNetwork/elrond-go/process/factory"
	"github.com/ElrondNetwork/elrond-go/process/factory/metachain"
	"github.com/ElrondNetwork/elrond-go/process/factory/shard"
	"github.com/ElrondNetwork/elrond-go/process/rewardTransaction"
	"github.com/ElrondNetwork/elrond-go/process/scToProtocol"
	"github.com/ElrondNetwork/elrond-go/process/smartContract"
	"github.com/ElrondNetwork/elrond-go/process/smartContract/builtInFunctions"
	"github.com/ElrondNetwork/elrond-go/process/smartContract/hooks"
	"github.com/ElrondNetwork/elrond-go/process/throttle"
	"github.com/ElrondNetwork/elrond-go/process/transaction"
	"github.com/ElrondNetwork/elrond-go/process/txsimulator"
	"github.com/ElrondNetwork/elrond-go/vm"
)

func (pcf *processComponentsFactory) newBlockProcessor(
	requestHandler process.RequestHandler,
	forkDetector process.ForkDetector,
	epochStartTrigger epochStart.TriggerHandler,
	bootStorer process.BootStorer,
	validatorStatisticsProcessor process.ValidatorStatisticsProcessor,
	headerValidator process.HeaderConstructionValidator,
	blockTracker process.BlockTracker,
	pendingMiniBlocksHandler process.PendingMiniBlocksHandler,
	txSimulatorProcessorArgs *txsimulator.ArgsTxSimulator,
) (process.BlockProcessor, error) {
	if pcf.shardCoordinator.SelfId() < pcf.shardCoordinator.NumberOfShards() {
		return pcf.newShardBlockProcessor(
			requestHandler,
			forkDetector,
			epochStartTrigger,
			bootStorer,
			headerValidator,
			blockTracker,
			pcf.smartContractParser,
			txSimulatorProcessorArgs,
		)
	}
	if pcf.shardCoordinator.SelfId() == core.MetachainShardId {
		return pcf.newMetaBlockProcessor(
			requestHandler,
			forkDetector,
			validatorStatisticsProcessor,
			epochStartTrigger,
			bootStorer,
			headerValidator,
			blockTracker,
			pendingMiniBlocksHandler,
			txSimulatorProcessorArgs,
		)
	}

	return nil, errors.New("could not create block processor")
}

func (pcf *processComponentsFactory) newShardBlockProcessor(
	requestHandler process.RequestHandler,
	forkDetector process.ForkDetector,
	epochStartTrigger epochStart.TriggerHandler,
	bootStorer process.BootStorer,
	headerValidator process.HeaderConstructionValidator,
	blockTracker process.BlockTracker,
	smartContractParser genesis.InitialSmartContractParser,
	txSimulatorProcessorArgs *txsimulator.ArgsTxSimulator,
) (process.BlockProcessor, error) {
	argsParser := smartContract.NewArgumentParser()

	mapDNSAddresses, err := smartContractParser.GetDeployedSCAddresses(genesis.DNSType)
	if err != nil {
		return nil, err
	}

	argsBuiltIn := builtInFunctions.ArgsCreateBuiltInFunctionContainer{
		GasSchedule:     pcf.gasSchedule,
		MapDNSAddresses: mapDNSAddresses,
		Marshalizer:     pcf.coreData.InternalMarshalizer(),
		Accounts:        pcf.state.AccountsAdapter(),
	}

	builtInFuncFactory, err := builtInFunctions.NewBuiltInFunctionsFactory(argsBuiltIn)
	if err != nil {
		return nil, err
	}
	builtInFuncs, err := builtInFuncFactory.CreateBuiltInFunctionContainer()

	argsHook := hooks.ArgBlockChainHook{
		Accounts:           pcf.state.AccountsAdapter(),
		PubkeyConv:         pcf.coreData.AddressPubKeyConverter(),
		StorageService:     pcf.data.StorageService(),
		BlockChain:         pcf.data.Blockchain(),
		ShardCoordinator:   pcf.shardCoordinator,
		Marshalizer:        pcf.coreData.InternalMarshalizer(),
		Uint64Converter:    pcf.coreData.Uint64ByteSliceConverter(),
		BuiltInFunctions:   builtInFuncs,
		DataPool:           pcf.data.Datapool(),
		CompiledSCPool:     pcf.data.Datapool().SmartContracts(),
		WorkingDir:         pcf.workingDir,
		NilCompiledSCStore: false,
		ConfigSCStorage:    pcf.config.SmartContractsStorage,
	}

	argsNewVMFactory := shard.ArgVMContainerFactory{
		Config:                         pcf.config.VirtualMachine.Execution,
		BlockGasLimit:                  pcf.economicsData.MaxGasLimitPerBlock(pcf.shardCoordinator.SelfId()),
		GasSchedule:                    pcf.gasSchedule,
		ArgBlockChainHook:              argsHook,
		DeployEnableEpoch:              pcf.config.GeneralSettings.SCDeployEnableEpoch,
		AheadOfTimeGasUsageEnableEpoch: pcf.config.GeneralSettings.AheadOfTimeGasUsageEnableEpoch,
		ArwenV3EnableEpoch:             pcf.config.GeneralSettings.RepairCallbackEnableEpoch,
	}
	vmFactory, err := shard.NewVMContainerFactory(argsNewVMFactory)
	if err != nil {
		return nil, err
	}

	vmContainer, err := vmFactory.Create()
	if err != nil {
		return nil, err
	}

	err = builtInFunctions.SetPayableHandler(builtInFuncs, vmFactory.BlockChainHookImpl())
	if err != nil {
		return nil, err
	}

	interimProcFactory, err := shard.NewIntermediateProcessorsContainerFactory(
		pcf.shardCoordinator,
		pcf.coreData.InternalMarshalizer(),
		pcf.coreData.Hasher(),
		pcf.coreData.AddressPubKeyConverter(),
		pcf.data.StorageService(),
		pcf.data.Datapool(),
	)
	if err != nil {
		return nil, err
	}

	interimProcContainer, err := interimProcFactory.Create()
	if err != nil {
		return nil, err
	}

	scForwarder, err := interimProcContainer.Get(dataBlock.SmartContractResultBlock)
	if err != nil {
		return nil, err
	}

	receiptTxInterim, err := interimProcContainer.Get(dataBlock.ReceiptBlock)
	if err != nil {
		return nil, err
	}

	badTxInterim, err := interimProcContainer.Get(dataBlock.InvalidBlock)
	if err != nil {
		return nil, err
	}

	argsTxTypeHandler := coordinator.ArgNewTxTypeHandler{
		PubkeyConverter:  pcf.coreData.AddressPubKeyConverter(),
		ShardCoordinator: pcf.shardCoordinator,
		BuiltInFuncNames: builtInFuncs.Keys(),
		ArgumentParser:   parsers.NewCallArgsParser(),
	}
	txTypeHandler, err := coordinator.NewTxTypeHandler(argsTxTypeHandler)
	if err != nil {
		return nil, err
	}

	gasHandler, err := preprocess.NewGasComputation(
		pcf.coreData.EconomicsData(),
		txTypeHandler,
		pcf.coreData.EpochNotifier(),
		pcf.config.GeneralSettings.SCDeployEnableEpoch,
	)
	if err != nil {
		return nil, err
	}

	txFeeHandler, err := postprocess.NewFeeAccumulator()
	if err != nil {
		return nil, err
	}

	generalSettings := pcf.config.GeneralSettings

	argsNewScProcessor := smartContract.ArgsNewSmartContractProcessor{
		VmContainer:                    vmContainer,
		ArgsParser:                     argsParser,
		Hasher:                         pcf.coreData.Hasher(),
		Marshalizer:                    pcf.coreData.InternalMarshalizer(),
		AccountsDB:                     pcf.state.AccountsAdapter(),
		BlockChainHook:                 vmFactory.BlockChainHookImpl(),
		PubkeyConv:                     pcf.coreData.AddressPubKeyConverter(),
		ShardCoordinator:               pcf.shardCoordinator,
		ScrForwarder:                   scForwarder,
		TxFeeHandler:                   txFeeHandler,
		EconomicsFee:                   pcf.economicsData,
		GasHandler:                     gasHandler,
		GasSchedule:                    pcf.gasSchedule,
		BuiltInFunctions:               vmFactory.BlockChainHookImpl().GetBuiltInFunctions(),
		TxLogsProcessor:                pcf.txLogsProcessor,
		TxTypeHandler:                  txTypeHandler,
		DeployEnableEpoch:              generalSettings.SCDeployEnableEpoch,
		BuiltinEnableEpoch:             generalSettings.BuiltInFunctionsEnableEpoch,
		PenalizedTooMuchGasEnableEpoch: generalSettings.PenalizedTooMuchGasEnableEpoch,
		RepairCallbackEnableEpoch:      pcf.config.GeneralSettings.RepairCallbackEnableEpoch,
		BadTxForwarder:                 badTxInterim,
		EpochNotifier:                  pcf.epochNotifier,
		StakingV2EnableEpoch:           pcf.systemSCConfig.StakingSystemSCConfig.StakingV2Epoch,
	}
	scProcessor, err := smartContract.NewSmartContractProcessor(argsNewScProcessor)
	if err != nil {
		return nil, err
	}

	rewardsTxProcessor, err := rewardTransaction.NewRewardTxProcessor(
		pcf.state.AccountsAdapter(),
		pcf.coreData.AddressPubKeyConverter(),
		pcf.shardCoordinator,
	)
	if err != nil {
		return nil, err
	}

	argsNewTxProcessor := transaction.ArgsNewTxProcessor{
		Accounts:                       pcf.state.AccountsAdapter(),
		Hasher:                         pcf.coreData.Hasher(),
		PubkeyConv:                     pcf.coreData.AddressPubKeyConverter(),
		Marshalizer:                    pcf.coreData.InternalMarshalizer(),
		SignMarshalizer:                pcf.coreData.TxMarshalizer(),
		ShardCoordinator:               pcf.shardCoordinator,
		ScProcessor:                    scProcessor,
		TxFeeHandler:                   txFeeHandler,
		TxTypeHandler:                  txTypeHandler,
		EconomicsFee:                   pcf.economicsData,
		ReceiptForwarder:               receiptTxInterim,
		BadTxForwarder:                 badTxInterim,
		ArgsParser:                     argsParser,
		ScrForwarder:                   scForwarder,
		RelayedTxEnableEpoch:           generalSettings.RelayedTransactionsEnableEpoch,
		PenalizedTooMuchGasEnableEpoch: generalSettings.PenalizedTooMuchGasEnableEpoch,
		MetaProtectionEnableEpoch:      generalSettings.MetaProtectionEnableEpoch,
		EpochNotifier:                  pcf.epochNotifier,
	}
	transactionProcessor, err := transaction.NewTxProcessor(argsNewTxProcessor)
	if err != nil {
		return nil, errors.New("could not create transaction statisticsProcessor: " + err.Error())
	}

	err = pcf.createShardTxSimulatorProcessor(txSimulatorProcessorArgs, argsNewScProcessor, argsNewTxProcessor)
	if err != nil {
		return nil, err
	}

	blockSizeThrottler, err := throttle.NewBlockSizeThrottle(pcf.minSizeInBytes, pcf.maxSizeInBytes)
	if err != nil {
		return nil, err
	}

	blockSizeComputationHandler, err := preprocess.NewBlockSizeComputation(
		pcf.coreData.InternalMarshalizer(),
		blockSizeThrottler,
		pcf.maxSizeInBytes,
	)
	if err != nil {
		return nil, err
	}

	balanceComputationHandler, err := preprocess.NewBalanceComputation()
	if err != nil {
		return nil, err
	}

	preProcFactory, err := shard.NewPreProcessorsContainerFactory(
		pcf.shardCoordinator,
		pcf.data.StorageService(),
		pcf.coreData.InternalMarshalizer(),
		pcf.coreData.Hasher(),
		pcf.data.Datapool(),
		pcf.coreData.AddressPubKeyConverter(),
		pcf.state.AccountsAdapter(),
		requestHandler,
		transactionProcessor,
		scProcessor,
		scProcessor,
		rewardsTxProcessor,
		pcf.economicsData,
		gasHandler,
		blockTracker,
		blockSizeComputationHandler,
		balanceComputationHandler,
	)
	if err != nil {
		return nil, err
	}

	preProcContainer, err := preProcFactory.Create()
	if err != nil {
		return nil, err
	}

	argsTransactionCoordinator := coordinator.ArgTransactionCoordinator{
		Hasher:                            pcf.coreData.Hasher(),
		Marshalizer:                       pcf.coreData.InternalMarshalizer(),
		ShardCoordinator:                  pcf.shardCoordinator,
		Accounts:                          pcf.state.AccountsAdapter(),
		MiniBlockPool:                     pcf.data.Datapool().MiniBlocks(),
		RequestHandler:                    requestHandler,
		PreProcessors:                     preProcContainer,
		InterProcessors:                   interimProcContainer,
		GasHandler:                        gasHandler,
		FeeHandler:                        txFeeHandler,
		BlockSizeComputation:              blockSizeComputationHandler,
		BalanceComputation:                balanceComputationHandler,
		EconomicsFee:                      pcf.economicsData,
		TxTypeHandler:                     txTypeHandler,
		BlockGasAndFeesReCheckEnableEpoch: pcf.config.GeneralSettings.BlockGasAndFeesReCheckEnableEpoch,
	}
	txCoordinator, err := coordinator.NewTransactionCoordinator(argsTransactionCoordinator)
	if err != nil {
		return nil, err
	}

	accountsDb := make(map[state.AccountsDbIdentifier]state.AccountsAdapter)
	accountsDb[state.UserAccountsState] = pcf.state.AccountsAdapter()

	argumentsBaseProcessor := block.ArgBaseProcessor{
		CoreComponents:          pcf.coreData,
		DataComponents:          pcf.data,
		Version:                 pcf.version,
		AccountsDB:              accountsDb,
		ForkDetector:            forkDetector,
		ShardCoordinator:        pcf.shardCoordinator,
		NodesCoordinator:        pcf.nodesCoordinator,
		RequestHandler:          requestHandler,
		BlockChainHook:          vmFactory.BlockChainHookImpl(),
		TxCoordinator:           txCoordinator,
		RoundHandler:            pcf.roundHandler,
		EpochStartTrigger:       epochStartTrigger,
		HeaderValidator:         headerValidator,
		BootStorer:              bootStorer,
		BlockTracker:            blockTracker,
		FeeHandler:              txFeeHandler,
		StateCheckpointModulus:  pcf.stateCheckpointModulus,
		BlockSizeThrottler:      blockSizeThrottler,
		Indexer:                 pcf.indexer,
		TpsBenchmark:            pcf.tpsBenchmark,
		HistoryRepository:       pcf.historyRepo,
		EpochNotifier:           pcf.epochNotifier,
		HeaderIntegrityVerifier: pcf.headerIntegrityVerifier,
		AppStatusHandler:        pcf.coreData.StatusHandler(),
		VMContainersFactory:     vmFactory,
		VmContainer:             vmContainer,
	}
	arguments := block.ArgShardProcessor{
		ArgBaseProcessor: argumentsBaseProcessor,
	}

	blockProcessor, err := block.NewShardProcessor(arguments)
	if err != nil {
		return nil, errors.New("could not create block statisticsProcessor: " + err.Error())
	}

	return blockProcessor, nil
}

func (pcf *processComponentsFactory) newMetaBlockProcessor(
	requestHandler process.RequestHandler,
	forkDetector process.ForkDetector,
	validatorStatisticsProcessor process.ValidatorStatisticsProcessor,
	epochStartTrigger epochStart.TriggerHandler,
	bootStorer process.BootStorer,
	headerValidator process.HeaderConstructionValidator,
	blockTracker process.BlockTracker,
	pendingMiniBlocksHandler process.PendingMiniBlocksHandler,
	txSimulatorProcessorArgs *txsimulator.ArgsTxSimulator,
) (process.BlockProcessor, error) {

	argsBuiltIn := builtInFunctions.ArgsCreateBuiltInFunctionContainer{
		GasSchedule:     pcf.gasSchedule,
		MapDNSAddresses: make(map[string]struct{}), // no dns for meta
		Marshalizer:     pcf.coreData.InternalMarshalizer(),
		Accounts:        pcf.state.AccountsAdapter(),
	}
	builtInFuncFactory, err := builtInFunctions.NewBuiltInFunctionsFactory(argsBuiltIn)
	if err != nil {
		return nil, err
	}
	builtInFuncs, err := builtInFuncFactory.CreateBuiltInFunctionContainer()
	if err != nil {
		return nil, err
	}

	argsHook := hooks.ArgBlockChainHook{
		Accounts:           pcf.state.AccountsAdapter(),
		PubkeyConv:         pcf.coreData.AddressPubKeyConverter(),
		StorageService:     pcf.data.StorageService(),
		BlockChain:         pcf.data.Blockchain(),
		ShardCoordinator:   pcf.shardCoordinator,
		Marshalizer:        pcf.coreData.InternalMarshalizer(),
		Uint64Converter:    pcf.coreData.Uint64ByteSliceConverter(),
		BuiltInFunctions:   builtInFuncs,
		DataPool:           pcf.data.Datapool(),
		CompiledSCPool:     pcf.data.Datapool().SmartContracts(),
		ConfigSCStorage:    pcf.config.SmartContractsStorage,
		WorkingDir:         pcf.workingDir,
		NilCompiledSCStore: false,
	}

	argsNewVMContainer := metachain.ArgsNewVMContainerFactory{
		ArgBlockChainHook:   argsHook,
		Economics:           pcf.economicsData,
		MessageSignVerifier: pcf.crypto.MessageSignVerifier(),
		GasSchedule:         pcf.gasSchedule,
		NodesConfigProvider: pcf.coreData.GenesisNodesSetup(),
		Hasher:              pcf.coreData.Hasher(),
		Marshalizer:         pcf.coreData.InternalMarshalizer(),
		SystemSCConfig:      pcf.systemSCConfig,
		ValidatorAccountsDB: pcf.state.PeerAccounts(),
		ChanceComputer:      pcf.coreData.Rater(),
		EpochNotifier:       pcf.coreData.EpochNotifier(),
	}
	vmFactory, err := metachain.NewVMContainerFactory(argsNewVMContainer)
	if err != nil {
		return nil, err
	}

	argsParser := smartContract.NewArgumentParser()

	vmContainer, err := vmFactory.Create()
	if err != nil {
		return nil, err
	}

	interimProcFactory, err := metachain.NewIntermediateProcessorsContainerFactory(
		pcf.shardCoordinator,
		pcf.coreData.InternalMarshalizer(),
		pcf.coreData.Hasher(),
		pcf.coreData.AddressPubKeyConverter(),
		pcf.data.StorageService(),
		pcf.data.Datapool(),
	)
	if err != nil {
		return nil, err
	}

	interimProcContainer, err := interimProcFactory.Create()
	if err != nil {
		return nil, err
	}

	scForwarder, err := interimProcContainer.Get(dataBlock.SmartContractResultBlock)
	if err != nil {
		return nil, err
	}

	badTxForwarder, err := interimProcContainer.Get(dataBlock.InvalidBlock)
	if err != nil {
		return nil, err
	}

	argsTxTypeHandler := coordinator.ArgNewTxTypeHandler{
		PubkeyConverter:  pcf.coreData.AddressPubKeyConverter(),
		ShardCoordinator: pcf.shardCoordinator,
		BuiltInFuncNames: builtInFuncs.Keys(),
		ArgumentParser:   parsers.NewCallArgsParser(),
	}
	txTypeHandler, err := coordinator.NewTxTypeHandler(argsTxTypeHandler)
	if err != nil {
		return nil, err
	}

	gasHandler, err := preprocess.NewGasComputation(
		pcf.economicsData,
		txTypeHandler,
		pcf.epochNotifier,
		pcf.config.GeneralSettings.SCDeployEnableEpoch,
	)
	if err != nil {
		return nil, err
	}

	txFeeHandler, err := postprocess.NewFeeAccumulator()
	if err != nil {
		return nil, err
	}

	generalSettingsConfig := pcf.config.GeneralSettings
	argsNewScProcessor := smartContract.ArgsNewSmartContractProcessor{
		VmContainer:                    vmContainer,
		ArgsParser:                     argsParser,
		Hasher:                         pcf.coreData.Hasher(),
		Marshalizer:                    pcf.coreData.InternalMarshalizer(),
		AccountsDB:                     pcf.state.AccountsAdapter(),
		BlockChainHook:                 vmFactory.BlockChainHookImpl(),
		PubkeyConv:                     pcf.coreData.AddressPubKeyConverter(),
		ShardCoordinator:               pcf.shardCoordinator,
		ScrForwarder:                   scForwarder,
		TxFeeHandler:                   txFeeHandler,
		EconomicsFee:                   pcf.economicsData,
		TxTypeHandler:                  txTypeHandler,
		GasHandler:                     gasHandler,
		GasSchedule:                    pcf.gasSchedule,
		BuiltInFunctions:               vmFactory.BlockChainHookImpl().GetBuiltInFunctions(),
		TxLogsProcessor:                pcf.txLogsProcessor,
		DeployEnableEpoch:              generalSettingsConfig.SCDeployEnableEpoch,
		BuiltinEnableEpoch:             generalSettingsConfig.BuiltInFunctionsEnableEpoch,
		PenalizedTooMuchGasEnableEpoch: generalSettingsConfig.PenalizedTooMuchGasEnableEpoch,
		RepairCallbackEnableEpoch:      generalSettingsConfig.RepairCallbackEnableEpoch,
		BadTxForwarder:                 badTxForwarder,
		EpochNotifier:                  pcf.epochNotifier,
		StakingV2EnableEpoch:           pcf.systemSCConfig.StakingSystemSCConfig.StakingV2Epoch,
	}
	scProcessor, err := smartContract.NewSmartContractProcessor(argsNewScProcessor)
	if err != nil {
		return nil, err
	}

	argsNewMetaTxProcessor := transaction.ArgsNewMetaTxProcessor{
		Hasher:           pcf.coreData.Hasher(),
		Marshalizer:      pcf.coreData.InternalMarshalizer(),
		Accounts:         pcf.state.AccountsAdapter(),
		PubkeyConv:       pcf.coreData.AddressPubKeyConverter(),
		ShardCoordinator: pcf.shardCoordinator,
		ScProcessor:      scProcessor,
		TxTypeHandler:    txTypeHandler,
		EconomicsFee:     pcf.economicsData,
		ESDTEnableEpoch:  pcf.systemSCConfig.ESDTSystemSCConfig.EnabledEpoch,
		EpochNotifier:    pcf.epochNotifier,
	}

	transactionProcessor, err := transaction.NewMetaTxProcessor(argsNewMetaTxProcessor)
	if err != nil {
		return nil, errors.New("could not create transaction processor: " + err.Error())
	}

	err = pcf.createMetaTxSimulatorProcessor(txSimulatorProcessorArgs, argsNewScProcessor, txTypeHandler)
	if err != nil {
		return nil, err
	}

	blockSizeThrottler, err := throttle.NewBlockSizeThrottle(pcf.minSizeInBytes, pcf.maxSizeInBytes)
	if err != nil {
		return nil, err
	}

	blockSizeComputationHandler, err := preprocess.NewBlockSizeComputation(
		pcf.coreData.InternalMarshalizer(),
		blockSizeThrottler,
		pcf.maxSizeInBytes,
	)
	if err != nil {
		return nil, err
	}

	balanceComputationHandler, err := preprocess.NewBalanceComputation()
	if err != nil {
		return nil, err
	}

	preProcFactory, err := metachain.NewPreProcessorsContainerFactory(
		pcf.shardCoordinator,
		pcf.data.StorageService(),
		pcf.coreData.InternalMarshalizer(),
		pcf.coreData.Hasher(),
		pcf.data.Datapool(),
		pcf.state.AccountsAdapter(),
		requestHandler,
		transactionProcessor,
		scProcessor,
		pcf.economicsData,
		gasHandler,
		blockTracker,
		pcf.coreData.AddressPubKeyConverter(),
		blockSizeComputationHandler,
		balanceComputationHandler,
	)
	if err != nil {
		return nil, err
	}

	preProcContainer, err := preProcFactory.Create()
	if err != nil {
		return nil, err
	}

	argsTransactionCoordinator := coordinator.ArgTransactionCoordinator{
		Hasher:                            pcf.coreData.Hasher(),
		Marshalizer:                       pcf.coreData.InternalMarshalizer(),
		ShardCoordinator:                  pcf.shardCoordinator,
		Accounts:                          pcf.state.AccountsAdapter(),
		MiniBlockPool:                     pcf.data.Datapool().MiniBlocks(),
		RequestHandler:                    requestHandler,
		PreProcessors:                     preProcContainer,
		InterProcessors:                   interimProcContainer,
		GasHandler:                        gasHandler,
		FeeHandler:                        txFeeHandler,
		BlockSizeComputation:              blockSizeComputationHandler,
		BalanceComputation:                balanceComputationHandler,
		EconomicsFee:                      pcf.economicsData,
		TxTypeHandler:                     txTypeHandler,
		BlockGasAndFeesReCheckEnableEpoch: generalSettingsConfig.BlockGasAndFeesReCheckEnableEpoch,
	}
	txCoordinator, err := coordinator.NewTransactionCoordinator(argsTransactionCoordinator)
	if err != nil {
		return nil, err
	}

	argsStaking := scToProtocol.ArgStakingToPeer{
		PubkeyConv:       pcf.coreData.ValidatorPubKeyConverter(),
		Hasher:           pcf.coreData.Hasher(),
		Marshalizer:      pcf.coreData.InternalMarshalizer(),
		PeerState:        pcf.state.PeerAccounts(),
		BaseState:        pcf.state.AccountsAdapter(),
		ArgParser:        argsParser,
		CurrTxs:          pcf.data.Datapool().CurrentBlockTxs(),
		RatingsData:      pcf.ratingsData,
		EpochNotifier:    pcf.coreData.EpochNotifier(),
		StakeEnableEpoch: pcf.systemSCConfig.StakingSystemSCConfig.StakeEnableEpoch,
	}
	smartContractToProtocol, err := scToProtocol.NewStakingToPeer(argsStaking)
	if err != nil {
		return nil, err
	}

	genesisHdr := pcf.data.Blockchain().GetGenesisHeader()
	argsEpochStartData := metachainEpochStart.ArgsNewEpochStartData{
		Marshalizer:       pcf.coreData.InternalMarshalizer(),
		Hasher:            pcf.coreData.Hasher(),
		Store:             pcf.data.StorageService(),
		DataPool:          pcf.data.Datapool(),
		BlockTracker:      blockTracker,
		ShardCoordinator:  pcf.shardCoordinator,
		EpochStartTrigger: epochStartTrigger,
		RequestHandler:    requestHandler,
		GenesisEpoch:      genesisHdr.GetEpoch(),
	}
	epochStartDataCreator, err := metachainEpochStart.NewEpochStartData(argsEpochStartData)
	if err != nil {
		return nil, err
	}

	economicsDataProvider := metachainEpochStart.NewEpochEconomicsStatistics()
	argsEpochEconomics := metachainEpochStart.ArgsNewEpochEconomics{
		Marshalizer:           pcf.coreData.InternalMarshalizer(),
		Hasher:                pcf.coreData.Hasher(),
		Store:                 pcf.data.StorageService(),
		ShardCoordinator:      pcf.shardCoordinator,
		RewardsHandler:        pcf.economicsData,
		RoundTime:             pcf.roundHandler,
		GenesisNonce:          genesisHdr.GetNonce(),
		GenesisEpoch:          genesisHdr.GetEpoch(),
		GenesisTotalSupply:    pcf.economicsData.GenesisTotalSupply(),
		EconomicsDataNotified: economicsDataProvider,
		StakingV2EnableEpoch:  pcf.systemSCConfig.StakingSystemSCConfig.StakingV2Epoch,
	}
	epochEconomics, err := metachainEpochStart.NewEndOfEpochEconomicsDataCreator(argsEpochEconomics)
	if err != nil {
		return nil, err
	}

	systemVM, err := vmContainer.Get(factory.SystemVirtualMachine)
	if err != nil {
		return nil, err
	}

	// TODO: in case of changing the minimum node price, make sure to update the staking data provider
	stakingDataProvider, err := metachainEpochStart.NewStakingDataProvider(systemVM, pcf.systemSCConfig.StakingSystemSCConfig.GenesisNodePrice)
	if err != nil {
		return nil, err
	}

	rewardsStorage := pcf.data.StorageService().GetStorer(dataRetriever.RewardTransactionUnit)
	miniBlockStorage := pcf.data.StorageService().GetStorer(dataRetriever.MiniBlockUnit)
	argsEpochRewards := metachainEpochStart.RewardsCreatorProxyArgs{
		BaseRewardsCreatorArgs: metachainEpochStart.BaseRewardsCreatorArgs{
			ShardCoordinator:              pcf.shardCoordinator,
			PubkeyConverter:               pcf.coreData.AddressPubKeyConverter(),
			RewardsStorage:                rewardsStorage,
			MiniBlockStorage:              miniBlockStorage,
			Hasher:                        pcf.coreData.Hasher(),
			Marshalizer:                   pcf.coreData.InternalMarshalizer(),
			DataPool:                      pcf.data.Datapool(),
			ProtocolSustainabilityAddress: pcf.economicsData.ProtocolSustainabilityAddress(),
			NodesConfigProvider:           pcf.nodesCoordinator,
			UserAccountsDB:                pcf.state.AccountsAdapter(),
			RewardsFix1EpochEnable:        generalSettingsConfig.SwitchJailWaitingEnableEpoch,
			DelegationSystemSCEnableEpoch: pcf.systemSCConfig.StakingSystemSCConfig.StakingV2Epoch,
		},
		StakingDataProvider:   stakingDataProvider,
		TopUpRewardFactor:     pcf.economicsData.RewardsTopUpFactor(),
		TopUpGradientPoint:    pcf.economicsData.RewardsTopUpGradientPoint(),
		EconomicsDataProvider: economicsDataProvider,
		EpochEnableV2:         pcf.systemSCConfig.StakingSystemSCConfig.StakingV2Epoch,
	}
	epochRewards, err := metachainEpochStart.NewRewardsCreatorProxy(argsEpochRewards)
	if err != nil {
		return nil, err
	}

	argsEpochValidatorInfo := metachainEpochStart.ArgsNewValidatorInfoCreator{
		ShardCoordinator: pcf.shardCoordinator,
		MiniBlockStorage: miniBlockStorage,
		Hasher:           pcf.coreData.Hasher(),
		Marshalizer:      pcf.coreData.InternalMarshalizer(),
		DataPool:         pcf.data.Datapool(),
	}
	validatorInfoCreator, err := metachainEpochStart.NewValidatorInfoCreator(argsEpochValidatorInfo)
	if err != nil {
		return nil, err
	}

	accountsDb := make(map[state.AccountsDbIdentifier]state.AccountsAdapter)
	accountsDb[state.UserAccountsState] = pcf.state.AccountsAdapter()
	accountsDb[state.PeerAccountsState] = pcf.state.PeerAccounts()

	argumentsBaseProcessor := block.ArgBaseProcessor{
		CoreComponents:          pcf.coreData,
		DataComponents:          pcf.data,
		Version:                 pcf.version,
		AccountsDB:              accountsDb,
		ForkDetector:            forkDetector,
		ShardCoordinator:        pcf.shardCoordinator,
		NodesCoordinator:        pcf.nodesCoordinator,
		RequestHandler:          requestHandler,
		BlockChainHook:          vmFactory.BlockChainHookImpl(),
		TxCoordinator:           txCoordinator,
		EpochStartTrigger:       epochStartTrigger,
		RoundHandler:            pcf.roundHandler,
		HeaderValidator:         headerValidator,
		BootStorer:              bootStorer,
		BlockTracker:            blockTracker,
		FeeHandler:              txFeeHandler,
		StateCheckpointModulus:  pcf.stateCheckpointModulus,
		BlockSizeThrottler:      blockSizeThrottler,
		Indexer:                 pcf.indexer,
		TpsBenchmark:            pcf.tpsBenchmark,
		HistoryRepository:       pcf.historyRepo,
		EpochNotifier:           pcf.epochNotifier,
		HeaderIntegrityVerifier: pcf.headerIntegrityVerifier,
		AppStatusHandler:        pcf.coreData.StatusHandler(),
		VMContainersFactory:     vmFactory,
		VmContainer:             vmContainer,
	}

	argsEpochSystemSC := metachainEpochStart.ArgsNewEpochStartSystemSCProcessing{
		SystemVM:                               systemVM,
		UserAccountsDB:                         pcf.state.AccountsAdapter(),
		PeerAccountsDB:                         pcf.state.PeerAccounts(),
		Marshalizer:                            pcf.coreData.InternalMarshalizer(),
		StartRating:                            pcf.coreData.RatingsData().StartRating(),
		ValidatorInfoCreator:                   validatorStatisticsProcessor,
		EndOfEpochCallerAddress:                vm.EndOfEpochAddress,
		StakingSCAddress:                       vm.StakingSCAddress,
		ChanceComputer:                         pcf.coreData.Rater(),
		EpochNotifier:                          pcf.coreData.EpochNotifier(),
		SwitchJailWaitingEnableEpoch:           generalSettingsConfig.SwitchJailWaitingEnableEpoch,
		SwitchHysteresisForMinNodesEnableEpoch: generalSettingsConfig.SwitchHysteresisForMinNodesEnableEpoch,
		DelegationEnableEpoch:                  pcf.systemSCConfig.DelegationManagerSystemSCConfig.EnabledEpoch,
		StakingV2EnableEpoch:                   pcf.systemSCConfig.StakingSystemSCConfig.StakingV2Epoch,
		GenesisNodesConfig:                     pcf.coreData.GenesisNodesSetup(),
		MaxNodesEnableConfig:                   generalSettingsConfig.MaxNodesChangeEnableEpoch,
		StakingDataProvider:                    stakingDataProvider,
		NodesConfigProvider:                    pcf.nodesCoordinator,
		ShardCoordinator:                       pcf.shardCoordinator,
	}
	epochStartSystemSCProcessor, err := metachainEpochStart.NewSystemSCProcessor(argsEpochSystemSC)
	if err != nil {
		return nil, err
	}

	arguments := block.ArgMetaProcessor{
		ArgBaseProcessor:             argumentsBaseProcessor,
		SCToProtocol:                 smartContractToProtocol,
		PendingMiniBlocksHandler:     pendingMiniBlocksHandler,
		EpochStartDataCreator:        epochStartDataCreator,
		EpochEconomics:               epochEconomics,
		EpochRewardsCreator:          epochRewards,
		EpochValidatorInfoCreator:    validatorInfoCreator,
		ValidatorStatisticsProcessor: validatorStatisticsProcessor,
		EpochSystemSCProcessor:       epochStartSystemSCProcessor,
		RewardsV2EnableEpoch:         pcf.systemSCConfig.StakingSystemSCConfig.StakingV2Epoch,
	}

	metaProcessor, err := block.NewMetaProcessor(arguments)
	if err != nil {
		return nil, errors.New("could not create block processor: " + err.Error())
	}

	return metaProcessor, nil
}

func (pcf *processComponentsFactory) createShardTxSimulatorProcessor(
	txSimulatorProcessorArgs *txsimulator.ArgsTxSimulator,
	scProcArgs smartContract.ArgsNewSmartContractProcessor,
	txProcArgs transaction.ArgsNewTxProcessor,
) error {
	readOnlyAccountsDB, err := txsimulator.NewReadOnlyAccountsDB(pcf.state.AccountsAdapter())
	if err != nil {
		return err
	}

	interimProcFactory, err := shard.NewIntermediateProcessorsContainerFactory(
		pcf.shardCoordinator,
		pcf.coreData.InternalMarshalizer(),
		pcf.coreData.Hasher(),
		pcf.coreData.AddressPubKeyConverter(),
		disabled.NewChainStorer(),
		pcf.data.Datapool(),
	)
	if err != nil {
		return err
	}

	interimProcContainer, err := interimProcFactory.Create()
	if err != nil {
		return err
	}

	scForwarder, err := interimProcContainer.Get(dataBlock.SmartContractResultBlock)
	if err != nil {
		return err
	}
	scProcArgs.ScrForwarder = scForwarder

	receiptTxInterim, err := interimProcContainer.Get(dataBlock.ReceiptBlock)
	if err != nil {
		return err
	}
	txProcArgs.ReceiptForwarder = receiptTxInterim

	badTxInterim, err := interimProcContainer.Get(dataBlock.InvalidBlock)
	if err != nil {
		return err
	}
	scProcArgs.BadTxForwarder = badTxInterim
	txProcArgs.BadTxForwarder = badTxInterim

	scProcArgs.TxFeeHandler = &processDisabled.FeeHandler{}
	txProcArgs.TxFeeHandler = &processDisabled.FeeHandler{}

	scProcArgs.AccountsDB = readOnlyAccountsDB

	scProcessor, err := smartContract.NewSmartContractProcessor(scProcArgs)
	if err != nil {
		return err
	}
	txProcArgs.ScProcessor = scProcessor

	txProcArgs.Accounts = readOnlyAccountsDB

	txSimulatorProcessorArgs.TransactionProcessor, err = transaction.NewTxProcessor(txProcArgs)
	if err != nil {
		return err
	}

	txSimulatorProcessorArgs.IntermmediateProcContainer = interimProcContainer

	return nil
}

func (pcf *processComponentsFactory) createMetaTxSimulatorProcessor(
	txSimulatorProcessorArgs *txsimulator.ArgsTxSimulator,
	scProcArgs smartContract.ArgsNewSmartContractProcessor,
	txTypeHandler process.TxTypeHandler,
) error {
	interimProcFactory, err := shard.NewIntermediateProcessorsContainerFactory(
		pcf.shardCoordinator,
		pcf.coreData.InternalMarshalizer(),
		pcf.coreData.Hasher(),
		pcf.coreData.AddressPubKeyConverter(),
		disabled.NewChainStorer(),
		pcf.data.Datapool(),
	)
	if err != nil {
		return err
	}

	interimProcContainer, err := interimProcFactory.Create()
	if err != nil {
		return err
	}

	scForwarder, err := interimProcContainer.Get(dataBlock.SmartContractResultBlock)
	if err != nil {
		return err
	}
	scProcArgs.ScrForwarder = scForwarder

	badTxInterim, err := interimProcContainer.Get(dataBlock.InvalidBlock)
	if err != nil {
		return err
	}
	scProcArgs.BadTxForwarder = badTxInterim

	scProcArgs.TxFeeHandler = &processDisabled.FeeHandler{}

	scProcessor, err := smartContract.NewSmartContractProcessor(scProcArgs)
	if err != nil {
		return err
	}

	accountsWrapper, err := txsimulator.NewReadOnlyAccountsDB(pcf.state.AccountsAdapter())
	if err != nil {
		return err
	}

	argsNewMetaTx := transaction.ArgsNewMetaTxProcessor{
		Hasher:           pcf.coreData.Hasher(),
		Marshalizer:      pcf.coreData.InternalMarshalizer(),
		Accounts:         accountsWrapper,
		PubkeyConv:       pcf.coreData.AddressPubKeyConverter(),
		ShardCoordinator: pcf.shardCoordinator,
		ScProcessor:      scProcessor,
		TxTypeHandler:    txTypeHandler,
		EconomicsFee:     &processDisabled.FeeHandler{},
		ESDTEnableEpoch:  pcf.systemSCConfig.ESDTSystemSCConfig.EnabledEpoch,
		EpochNotifier:    pcf.epochNotifier,
	}

	txSimulatorProcessorArgs.TransactionProcessor, err = transaction.NewMetaTxProcessor(argsNewMetaTx)
	if err != nil {
		return err
	}

	txSimulatorProcessorArgs.IntermmediateProcContainer = interimProcContainer

	return nil
}
