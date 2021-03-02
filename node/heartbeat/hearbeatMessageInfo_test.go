<<<<<<< Updated upstream
package heartbeat_test
=======
package heartbeat
>>>>>>> Stashed changes

import (
	"testing"
	"time"

<<<<<<< Updated upstream
	"github.com/ElrondNetwork/elrond-go/node/heartbeat"
	"github.com/ElrondNetwork/elrond-go/node/mock"
	"github.com/stretchr/testify/assert"
)

//------- newHeartbeatMessageInfo

func TestNewHeartbeatMessageInfo_InvalidDurationShouldErr(t *testing.T) {
	t.Parallel()

	hbmi, err := heartbeat.NewHeartbeatMessageInfo(
		0,
		false,
		time.Time{},
		&mock.MockTimer{},
	)

	assert.Nil(t, hbmi)
	assert.Equal(t, heartbeat.ErrInvalidMaxDurationPeerUnresponsive, err)
}

func TestNewHeartbeatMessageInfo_NilGetTimeHandlerShouldErr(t *testing.T) {
	t.Parallel()

	hbmi, err := heartbeat.NewHeartbeatMessageInfo(
		1,
		false,
		time.Time{},
		nil,
	)

	assert.Nil(t, hbmi)
	assert.Equal(t, heartbeat.ErrNilTimer, err)
=======
	"github.com/stretchr/testify/assert"
)

//------ newHeartbeatMessageInfo
func TestNewHeartbeatMessageInfo_InvalidDurationShouldErr(t *testing.T) {
	t.Parallel()

	hbmi, err := newHeartbeatMessageInfo(0, false)

	assert.Nil(t, hbmi)
	assert.Equal(t, ErrInvalidMaxDurationPeerUnresponsive, err)
>>>>>>> Stashed changes
}

func TestNewHeartbeatMessageInfo_OkValsShouldWork(t *testing.T) {
	t.Parallel()

<<<<<<< Updated upstream
	hbmi, err := heartbeat.NewHeartbeatMessageInfo(
		1,
		false,
		time.Time{},
		&mock.MockTimer{},
	)
=======
	hbmi, err := newHeartbeatMessageInfo(1, false)
>>>>>>> Stashed changes

	assert.NotNil(t, hbmi)
	assert.Nil(t, err)
}

<<<<<<< Updated upstream
//------- HeartbeatReceived

func TestHeartbeatMessageInfo_HeartbeatReceivedShouldUpdate(t *testing.T) {
	t.Parallel()

	mockTimer := &mock.MockTimer{}
	genesisTime := mockTimer.Now()

	hbmi, _ := heartbeat.NewHeartbeatMessageInfo(
		10*time.Second,
		false,
		genesisTime,
		mockTimer,
	)

	assert.Equal(t, genesisTime, hbmi.GetTimeStamp())

	mockTimer.IncrementSeconds(1)

	expectedTime := time.Unix(1, 0)
	hbmi.HeartbeatReceived(uint32(0), uint32(0), "v0.1", "undefined")
	assert.Equal(t, expectedTime, hbmi.GetTimeStamp())
	assert.Equal(t, uint32(0), hbmi.GetReceiverShardId())

	mockTimer.IncrementSeconds(1)
	expectedTime = time.Unix(2, 0)
	hbmi.HeartbeatReceived(uint32(0), uint32(1), "v0.1", "undefined")
	assert.Equal(t, expectedTime, hbmi.GetTimeStamp())
	assert.Equal(t, uint32(1), hbmi.GetReceiverShardId())
=======
func TestHeartbeatMessageInfo_HeartbeatReceivedShouldUpdate(t *testing.T) {
	t.Parallel()

	hbmi, _ := newHeartbeatMessageInfo(time.Duration(10), false)
	incrementalTime := int64(0)
	hbmi.getTimeHandler = func() time.Time {
		if incrementalTime < 2 {
			incrementalTime++
		}
		return time.Unix(0, incrementalTime)
	}

	assert.Equal(t, emptyTimestamp, hbmi.timeStamp)

	hbmi.HeartbeatReceived(uint32(0), uint32(0), "v0.1", "undefined")
	assert.NotEqual(t, emptyTimestamp, hbmi.timeStamp)
	assert.Equal(t, uint32(0), hbmi.receivedShardID)

	hbmi.HeartbeatReceived(uint32(0), uint32(1), "v0.1", "undefined")
	assert.NotEqual(t, emptyTimestamp, hbmi.timeStamp)
	assert.Equal(t, uint32(1), hbmi.receivedShardID)
>>>>>>> Stashed changes
}

func TestHeartbeatMessageInfo_HeartbeatUpdateFieldsShouldWork(t *testing.T) {
	t.Parallel()

<<<<<<< Updated upstream
	mockTimer := &mock.MockTimer{}
	genesisTime := mockTimer.Now()
	hbmi, _ := heartbeat.NewHeartbeatMessageInfo(
		100*time.Second,
		false,
		genesisTime,
		mockTimer,
	)

	assert.Equal(t, genesisTime, hbmi.GetTimeStamp())

	mockTimer.IncrementSeconds(1)

	expectedTime := time.Unix(1, 0)
	expectedUptime := time.Duration(0)
	expectedDownTime := time.Duration(1 * time.Second)
	hbmi.HeartbeatReceived(uint32(0), uint32(3), "v0.1", "undefined")
	assert.Equal(t, expectedTime, hbmi.GetTimeStamp())
	assert.Equal(t, true, hbmi.GetIsActive())
	assert.Equal(t, expectedUptime, hbmi.GetTotalUpTime().Duration)
	assert.Equal(t, expectedDownTime, hbmi.GetTotalDownTime().Duration)
}

func TestHeartbeatMessageInfo_HeartbeatShouldUpdateUpDownTime(t *testing.T) {
	t.Parallel()

	mockTimer := &mock.MockTimer{}
	genesisTime := mockTimer.Now()
	hbmi, _ := heartbeat.NewHeartbeatMessageInfo(
		100*time.Second,
		false,
		genesisTime,
		mockTimer,
	)

	assert.Equal(t, genesisTime, hbmi.GetTimeStamp())

	// send heartbeat twice in order to calculate the duration between thm
	mockTimer.IncrementSeconds(1)
	hbmi.HeartbeatReceived(uint32(0), uint32(1), "v0.1", "undefined")
	mockTimer.IncrementSeconds(1)
	hbmi.HeartbeatReceived(uint32(0), uint32(2), "v0.1", "undefined")

	expectedDownDuration := time.Duration(1 * time.Second)
	expectedUpDuration := time.Duration(1 * time.Second)
	assert.Equal(t, expectedUpDuration, hbmi.GetTotalUpTime().Duration)
	assert.Equal(t, expectedDownDuration, hbmi.GetTotalDownTime().Duration)
	expectedTime := time.Unix(2, 0)
	assert.Equal(t, expectedTime, hbmi.GetTimeStamp())
}

func TestHeartbeatMessageInfo_HeartbeatLongerDurationThanMaxShouldUpdateDownTime(t *testing.T) {
	t.Parallel()

	mockTimer := &mock.MockTimer{}
	genesisTime := mockTimer.Now()
	hbmi, _ := heartbeat.NewHeartbeatMessageInfo(
		500*time.Millisecond,
		false,
		genesisTime,
		mockTimer,
	)

	assert.Equal(t, genesisTime, hbmi.GetTimeStamp())

	// send heartbeat twice in order to calculate the duration between thm
	mockTimer.IncrementSeconds(1)
	hbmi.HeartbeatReceived(uint32(0), uint32(1), "v0.1", "undefined")
	mockTimer.IncrementSeconds(1)
	hbmi.HeartbeatReceived(uint32(0), uint32(2), "v0.1", "undefined")

	expectedDownDuration := time.Duration(2 * time.Second)
	expectedUpDuration := time.Duration(0)
	assert.Equal(t, expectedDownDuration, hbmi.GetTotalDownTime().Duration)
	assert.Equal(t, expectedUpDuration, hbmi.GetTotalUpTime().Duration)
	expectedTime := time.Unix(2, 0)
	assert.Equal(t, expectedTime, hbmi.GetTimeStamp())
}

func TestHeartbeatMessageInfo_HeartbeatBeforeGenesisShouldNotUpdateUpDownTime(t *testing.T) {
	t.Parallel()

	mockTimer := &mock.MockTimer{}
	genesisTime := time.Unix(5, 0)
	hbmi, _ := heartbeat.NewHeartbeatMessageInfo(
		100*time.Second,
		false,
		genesisTime,
		mockTimer,
	)

	assert.Equal(t, genesisTime, hbmi.GetTimeStamp())

	// send heartbeat twice in order to calculate the duration between thm
	mockTimer.IncrementSeconds(1)
	hbmi.HeartbeatReceived(uint32(0), uint32(1), "v0.1", "undefined")
	mockTimer.IncrementSeconds(1)
	hbmi.HeartbeatReceived(uint32(0), uint32(2), "v0.1", "undefined")

	expectedDuration := time.Duration(0)
	assert.Equal(t, expectedDuration, hbmi.GetTotalDownTime().Duration)
	assert.Equal(t, expectedDuration, hbmi.GetTotalUpTime().Duration)
	expectedTime := time.Unix(2, 0)
	assert.Equal(t, expectedTime, hbmi.GetTimeStamp())
}

func TestHeartbeatMessageInfo_HeartbeatEqualGenesisShouldHaveUpDownTimeZero(t *testing.T) {
	t.Parallel()

	mockTimer := &mock.MockTimer{}
	genesisTime := time.Unix(1, 0)
	hbmi, _ := heartbeat.NewHeartbeatMessageInfo(
		100*time.Second,
		false,
		genesisTime,
		mockTimer,
	)

	assert.Equal(t, genesisTime, hbmi.GetTimeStamp())
	mockTimer.IncrementSeconds(1)
	hbmi.HeartbeatReceived(uint32(0), uint32(1), "v0.1", "undefined")

	expectedDuration := time.Duration(0)
	assert.Equal(t, expectedDuration, hbmi.GetTotalUpTime().Duration)
	assert.Equal(t, expectedDuration, hbmi.GetTotalDownTime().Duration)
	expectedTime := time.Unix(1, 0)
	assert.Equal(t, expectedTime, hbmi.GetTimeStamp())
=======
	hbmi, _ := newHeartbeatMessageInfo(time.Duration(1), false)
	incrementalTime := int64(0)
	hbmi.getTimeHandler = func() time.Time {
		tReturned := time.Unix(0, incrementalTime)
		incrementalTime += 10

		return tReturned
	}

	assert.Equal(t, emptyTimestamp, hbmi.timeStamp)

	hbmi.HeartbeatReceived(uint32(0), uint32(3), "v0.1", "undefined")
	assert.NotEqual(t, emptyTimestamp, hbmi.timeStamp)
}

func TestHeartbeatMessageInfo_HeartbeatShouldUpdateUpTime(t *testing.T) {
	t.Parallel()

	hbmi, _ := newHeartbeatMessageInfo(time.Duration(10), false)
	incrementalTime := int64(0)
	hbmi.getTimeHandler = func() time.Time {
		tReturned := time.Unix(0, incrementalTime)
		incrementalTime += 1

		return tReturned
	}

	assert.Equal(t, emptyTimestamp, hbmi.timeStamp)

	// send heartbeat twice in order to calculate the duration between thm
	hbmi.HeartbeatReceived(uint32(0), uint32(1), "v0.1", "undefined")
	hbmi.HeartbeatReceived(uint32(0), uint32(2), "v0.1", "undefined")

	assert.True(t, hbmi.totalUpTime.Duration > time.Duration(0))
	assert.NotEqual(t, emptyTimestamp, hbmi.timeStamp)
>>>>>>> Stashed changes
}
