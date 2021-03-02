package heartbeat

import (
<<<<<<< Updated upstream
	"sync"
	"time"
)

=======
	"time"
)

var emptyTimestamp = time.Time{}

>>>>>>> Stashed changes
// heartbeatMessageInfo retain the message info received from another node (identified by a public key)
type heartbeatMessageInfo struct {
	maxDurationPeerUnresponsive time.Duration
	maxInactiveTime             Duration
	totalUpTime                 Duration
	totalDownTime               Duration

	getTimeHandler     func() time.Time
	timeStamp          time.Time
	isActive           bool
	receivedShardID    uint32
	computedShardID    uint32
	versionNumber      string
	nodeDisplayName    string
	isValidator        bool
	lastUptimeDowntime time.Time
<<<<<<< Updated upstream
	genesisTime        time.Time
	updateMutex        sync.Mutex
=======
>>>>>>> Stashed changes
}

// newHeartbeatMessageInfo returns a new instance of a heartbeatMessageInfo
func newHeartbeatMessageInfo(
	maxDurationPeerUnresponsive time.Duration,
	isValidator bool,
<<<<<<< Updated upstream
	genesisTime time.Time,
	timer Timer,
=======
>>>>>>> Stashed changes
) (*heartbeatMessageInfo, error) {

	if maxDurationPeerUnresponsive == 0 {
		return nil, ErrInvalidMaxDurationPeerUnresponsive
	}
<<<<<<< Updated upstream
	if timer == nil || timer.IsInterfaceNil() {
		return nil, ErrNilTimer
	}
=======
>>>>>>> Stashed changes

	hbmi := &heartbeatMessageInfo{
		maxDurationPeerUnresponsive: maxDurationPeerUnresponsive,
		maxInactiveTime:             Duration{0},
		isActive:                    false,
		receivedShardID:             uint32(0),
<<<<<<< Updated upstream
		timeStamp:                   genesisTime,
		lastUptimeDowntime:          timer.Now(),
=======
		timeStamp:                   emptyTimestamp,
		lastUptimeDowntime:          time.Now(),
>>>>>>> Stashed changes
		totalUpTime:                 Duration{0},
		totalDownTime:               Duration{0},
		versionNumber:               "",
		nodeDisplayName:             "",
		isValidator:                 isValidator,
<<<<<<< Updated upstream
		genesisTime:                 genesisTime,
		getTimeHandler:              timer.Now,
	}
=======
	}
	hbmi.getTimeHandler = hbmi.clockTime
>>>>>>> Stashed changes

	return hbmi, nil
}

<<<<<<< Updated upstream
func (hbmi *heartbeatMessageInfo) updateFields(crtTime time.Time) {
	validDuration := computeValidDuration(crtTime, hbmi)
	previousActive := hbmi.isActive && validDuration
	hbmi.isActive = true

	hbmi.updateTimes(crtTime, previousActive)
}

func (hbmi *heartbeatMessageInfo) computeActive(crtTime time.Time) {
	hbmi.updateMutex.Lock()
	validDuration := computeValidDuration(crtTime, hbmi)
	hbmi.isActive = hbmi.isActive && validDuration
	hbmi.updateTimes(crtTime, hbmi.isActive)
	hbmi.updateMutex.Unlock()
}

func (hbmi *heartbeatMessageInfo) updateTimes(crtTime time.Time, previousActive bool) {
	if crtTime.Sub(hbmi.genesisTime) < 0 {
		return
	}
	hbmi.updateMaxInactiveTimeDuration(crtTime)
	hbmi.updateUpAndDownTime(previousActive, crtTime)
}

func computeValidDuration(crtTime time.Time, hbmi *heartbeatMessageInfo) bool {
	crtDuration := crtTime.Sub(hbmi.timeStamp)
	crtDuration = maxDuration(0, crtDuration)
	validDuration := crtDuration <= hbmi.maxDurationPeerUnresponsive
	return validDuration
}

// Will update the total time a node was up and down
func (hbmi *heartbeatMessageInfo) updateUpAndDownTime(previousActive bool, crtTime time.Time) {
	if hbmi.lastUptimeDowntime.Sub(hbmi.genesisTime) < 0 {
		hbmi.lastUptimeDowntime = hbmi.genesisTime
	}

	lastDuration := crtTime.Sub(hbmi.lastUptimeDowntime)
	lastDuration = maxDuration(0, lastDuration)

	if previousActive && hbmi.isActive {
=======
func (hbmi *heartbeatMessageInfo) clockTime() time.Time {
	return time.Now()
}

func (hbmi *heartbeatMessageInfo) updateFields() {
	crtDuration := hbmi.getTimeHandler().Sub(hbmi.timeStamp)
	crtDuration = maxDuration(0, crtDuration)

	hbmi.isActive = crtDuration < hbmi.maxDurationPeerUnresponsive
	hbmi.updateUpAndDownTime()
	hbmi.updateMaxInactiveTimeDuration()
}

// Wil update the total time a node was up and down
func (hbmi *heartbeatMessageInfo) updateUpAndDownTime() {
	lastDuration := hbmi.clockTime().Sub(hbmi.lastUptimeDowntime)
	lastDuration = maxDuration(0, lastDuration)

	if hbmi.isActive {
>>>>>>> Stashed changes
		hbmi.totalUpTime.Duration += lastDuration
	} else {
		hbmi.totalDownTime.Duration += lastDuration
	}

<<<<<<< Updated upstream
	hbmi.lastUptimeDowntime = crtTime
}

// HeartbeatReceived processes a new message arrived from a peer
func (hbmi *heartbeatMessageInfo) HeartbeatReceived(
	computedShardID uint32,
	receivedshardID uint32,
	version string,
	nodeDisplayName string,
) {
	crtTime := hbmi.getTimeHandler()
	hbmi.updateFields(crtTime)
	hbmi.computedShardID = computedShardID
	hbmi.receivedShardID = receivedshardID
=======
	hbmi.lastUptimeDowntime = time.Now()
}

// HeartbeatReceived processes a new message arrived from a peer
func (hbmi *heartbeatMessageInfo) HeartbeatReceived(computedShardID, receivedshardID uint32, version string,
	nodeDisplayName string) {
	crtTime := hbmi.getTimeHandler()
	hbmi.updateFields()
	hbmi.computedShardID = computedShardID
	hbmi.receivedShardID = receivedshardID
	hbmi.updateMaxInactiveTimeDuration()
>>>>>>> Stashed changes
	hbmi.timeStamp = crtTime
	hbmi.versionNumber = version
	hbmi.nodeDisplayName = nodeDisplayName
}

<<<<<<< Updated upstream
func (hbmi *heartbeatMessageInfo) updateMaxInactiveTimeDuration(currentTime time.Time) {
	crtDuration := currentTime.Sub(hbmi.timeStamp)
	crtDuration = maxDuration(0, crtDuration)

	greaterDurationThanMax := hbmi.maxInactiveTime.Duration < crtDuration
	currentTimeAfterGenesis := hbmi.genesisTime.Sub(currentTime) < 0

	if greaterDurationThanMax && currentTimeAfterGenesis {
=======
func (hbmi *heartbeatMessageInfo) updateMaxInactiveTimeDuration() {
	crtDuration := hbmi.getTimeHandler().Sub(hbmi.timeStamp)
	crtDuration = maxDuration(0, crtDuration)

	if hbmi.maxInactiveTime.Duration < crtDuration && hbmi.timeStamp != emptyTimestamp {
>>>>>>> Stashed changes
		hbmi.maxInactiveTime.Duration = crtDuration
	}
}

func maxDuration(first, second time.Duration) time.Duration {
	if first > second {
		return first
	}

	return second
}
