package commonSubround

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ElrondNetwork/elrond-go/consensus/spos"
	"github.com/ElrondNetwork/elrond-go/core"
<<<<<<< Updated upstream
	"github.com/ElrondNetwork/elrond-go/core/indexer"
=======
>>>>>>> Stashed changes
	"github.com/ElrondNetwork/elrond-go/core/logger"
	"github.com/ElrondNetwork/elrond-go/statusHandler"
)

var log = logger.DefaultLogger()

// SubroundStartRound defines the data needed by the subround StartRound
type SubroundStartRound struct {
	*spos.Subround
	processingThresholdPercentage int
	getSubroundName               func(subroundId int) string
	executeStoredMessages         func()
<<<<<<< Updated upstream

	appStatusHandler core.AppStatusHandler
	indexer          indexer.Indexer
=======
	broadcastUnnotarisedBlocks    func()

	appStatusHandler core.AppStatusHandler
>>>>>>> Stashed changes
}

// NewSubroundStartRound creates a SubroundStartRound object
func NewSubroundStartRound(
	baseSubround *spos.Subround,
	extend func(subroundId int),
	processingThresholdPercentage int,
	getSubroundName func(subroundId int) string,
	executeStoredMessages func(),
<<<<<<< Updated upstream
) (*SubroundStartRound, error) {
	err := checkNewSubroundStartRoundParams(
		baseSubround,
=======
	broadcastUnnotarisedBlocks func(),
) (*SubroundStartRound, error) {
	err := checkNewSubroundStartRoundParams(
		baseSubround,
		broadcastUnnotarisedBlocks,
>>>>>>> Stashed changes
	)
	if err != nil {
		return nil, err
	}

	srStartRound := SubroundStartRound{
		baseSubround,
		processingThresholdPercentage,
		getSubroundName,
		executeStoredMessages,
<<<<<<< Updated upstream
		statusHandler.NewNilStatusHandler(),
		indexer.NewNilIndexer(),
=======
		broadcastUnnotarisedBlocks,
		statusHandler.NewNilStatusHandler(),
>>>>>>> Stashed changes
	}
	srStartRound.Job = srStartRound.doStartRoundJob
	srStartRound.Check = srStartRound.doStartRoundConsensusCheck
	srStartRound.Extend = extend

	return &srStartRound, nil
}

func checkNewSubroundStartRoundParams(
	baseSubround *spos.Subround,
<<<<<<< Updated upstream
=======
	broadcastUnnotarisedBlocks func(),
>>>>>>> Stashed changes
) error {
	if baseSubround == nil {
		return spos.ErrNilSubround
	}
	if baseSubround.ConsensusState == nil {
		return spos.ErrNilConsensusState
	}
<<<<<<< Updated upstream
=======
	if broadcastUnnotarisedBlocks == nil {
		return spos.ErrNilBroadcastUnnotarisedBlocks
	}
>>>>>>> Stashed changes

	err := spos.ValidateConsensusCore(baseSubround.ConsensusCoreHandler)

	return err
}

// SetAppStatusHandler method set appStatusHandler
func (sr *SubroundStartRound) SetAppStatusHandler(ash core.AppStatusHandler) error {
	if ash == nil || ash.IsInterfaceNil() {
		return spos.ErrNilAppStatusHandler
	}

	sr.appStatusHandler = ash
	return nil
}

<<<<<<< Updated upstream
// SetIndexer method set indexer
func (sr *SubroundStartRound) SetIndexer(indexer indexer.Indexer) {
	sr.indexer = indexer
}

=======
>>>>>>> Stashed changes
// doStartRoundJob method does the job of the subround StartRound
func (sr *SubroundStartRound) doStartRoundJob() bool {
	sr.ResetConsensusState()
	sr.RoundIndex = sr.Rounder().Index()
	sr.RoundTimeStamp = sr.Rounder().TimeStamp()
	return true
}

// doStartRoundConsensusCheck method checks if the consensus is achieved in the subround StartRound
func (sr *SubroundStartRound) doStartRoundConsensusCheck() bool {
	if sr.RoundCanceled {
		return false
	}

	if sr.Status(sr.Current()) == spos.SsFinished {
		return true
	}

	if sr.initCurrentRound() {
		return true
	}

	return false
}

func (sr *SubroundStartRound) initCurrentRound() bool {
	if sr.BootStrapper().ShouldSync() { // if node is not synchronized yet, it has to continue the bootstrapping mechanism
		return false
	}
	sr.appStatusHandler.SetStringValue(core.MetricConsensusRoundState, "")

	err := sr.generateNextConsensusGroup(sr.Rounder().Index())
	if err != nil {
		log.Error(err.Error())

		sr.RoundCanceled = true

		return false
	}

	leader, err := sr.GetLeader()
	if err != nil {
		log.Error(err.Error())

		sr.RoundCanceled = true

		return false
	}

	msg := ""
	if leader == sr.SelfPubKey() {
		sr.appStatusHandler.Increment(core.MetricCountLeader)
		sr.appStatusHandler.SetStringValue(core.MetricConsensusRoundState, "proposed")
		sr.appStatusHandler.SetStringValue(core.MetricConsensusState, "proposer")
		msg = " (my turn)"
	}

	log.Info(fmt.Sprintf("%sStep 0: preparing for this round with leader %s%s\n",
		sr.SyncTimer().FormattedCurrentTime(), core.GetTrimmedPk(hex.EncodeToString([]byte(leader))), msg))

	pubKeys := sr.ConsensusGroup()

<<<<<<< Updated upstream
	sr.indexRoundIfNeeded(pubKeys)

=======
>>>>>>> Stashed changes
	selfIndex, err := sr.SelfConsensusGroupIndex()
	if err != nil {
		log.Info(fmt.Sprintf("%scanceled round %d in subround %s, not in the consensus group\n",
			sr.SyncTimer().FormattedCurrentTime(), sr.Rounder().Index(), sr.getSubroundName(sr.Current())))

		sr.RoundCanceled = true

		sr.appStatusHandler.SetStringValue(core.MetricConsensusState, "not in consensus group")

		return false
	}

	sr.appStatusHandler.Increment(core.MetricCountConsensus)
	sr.appStatusHandler.SetStringValue(core.MetricConsensusState, "participant")

	err = sr.MultiSigner().Reset(pubKeys, uint16(selfIndex))
	if err != nil {
		log.Error(err.Error())

		sr.RoundCanceled = true

		return false
	}

	startTime := time.Time{}
	startTime = sr.RoundTimeStamp
	maxTime := sr.Rounder().TimeDuration() * time.Duration(sr.processingThresholdPercentage) / 100
	if sr.Rounder().RemainingTime(startTime, maxTime) < 0 {
		log.Info(fmt.Sprintf("%scanceled round %d in subround %s, time is out\n",
			sr.SyncTimer().FormattedCurrentTime(), sr.Rounder().Index(), sr.getSubroundName(sr.Current())))

		sr.RoundCanceled = true

		return false
	}

	sr.SetStatus(sr.Current(), spos.SsFinished)

<<<<<<< Updated upstream
=======
	if leader == sr.SelfPubKey() {
		//TODO: Should be analyzed if call of sr.broadcastUnnotarisedBlocks() is still necessary
	}

>>>>>>> Stashed changes
	// execute stored messages which were received in this new round but before this initialisation
	go sr.executeStoredMessages()

	return true
}

<<<<<<< Updated upstream
func (sr *SubroundStartRound) indexRoundIfNeeded(pubKeys []string) {
	if sr.indexer == nil || sr.indexer.IsNilIndexer() {
		return
	}

	shardId := sr.ShardCoordinator().SelfId()
	signersIndexes := sr.NodesCoordinator().GetValidatorsIndexes(pubKeys)
	round := sr.Rounder().Index()

	roundInfo := indexer.RoundInfo{
		Index:            uint64(round),
		SignersIndexes:   signersIndexes,
		BlockWasProposed: false,
		ShardId:          shardId,
		Timestamp:        time.Duration(sr.RoundTimeStamp.Unix()),
	}

	go sr.indexer.SaveRoundInfo(roundInfo)
}

=======
>>>>>>> Stashed changes
func (sr *SubroundStartRound) generateNextConsensusGroup(roundIndex int64) error {
	currentHeader := sr.Blockchain().GetCurrentBlockHeader()
	if currentHeader == nil {
		currentHeader = sr.Blockchain().GetGenesisHeader()
		if currentHeader == nil {
			return spos.ErrNilHeader
		}
	}

<<<<<<< Updated upstream
	randomSeed := currentHeader.GetRandSeed()

	log.Info(fmt.Sprintf("random source used to determine the next consensus group is: %s\n",
		core.ToB64(randomSeed)),
	)

	shardId := sr.ShardCoordinator().SelfId()

	nextConsensusGroup, _, err := sr.GetNextConsensusGroup(
		randomSeed,
		uint64(sr.RoundIndex),
		shardId,
		sr.NodesCoordinator(),
	)
=======
	randomSource := fmt.Sprintf("%d-%s", roundIndex, core.ToB64(currentHeader.GetRandSeed()))

	log.Info(fmt.Sprintf("random source used to determine the next consensus group is: %s\n", randomSource))

	nextConsensusGroup, err := sr.GetNextConsensusGroup(randomSource, sr.ValidatorGroupSelector())
>>>>>>> Stashed changes
	if err != nil {
		return err
	}

<<<<<<< Updated upstream
	log.Debug(fmt.Sprintf("consensus group for round %d is formed by next validators:\n",
		roundIndex))

	for i := 0; i < len(nextConsensusGroup); i++ {
		log.Debug(fmt.Sprintf("%s", core.GetTrimmedPk(hex.EncodeToString([]byte(nextConsensusGroup[i])))))
	}

	log.Debug(fmt.Sprintf("\n"))

	sr.SetConsensusGroup(nextConsensusGroup)

	sr.BlockProcessor().SetConsensusData(randomSeed, uint64(sr.RoundIndex), currentHeader.GetEpoch(), shardId)

=======
	log.Info(fmt.Sprintf("consensus group for round %d is formed by next validators:\n",
		roundIndex))

	for i := 0; i < len(nextConsensusGroup); i++ {
		log.Info(fmt.Sprintf("%s", core.GetTrimmedPk(hex.EncodeToString([]byte(nextConsensusGroup[i])))))
	}

	log.Info(fmt.Sprintf("\n"))

	sr.SetConsensusGroup(nextConsensusGroup)

>>>>>>> Stashed changes
	return nil
}
