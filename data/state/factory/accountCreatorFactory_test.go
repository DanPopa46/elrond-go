package factory_test

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go/data/mock"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/data/state/factory"
<<<<<<< Updated upstream
	"github.com/stretchr/testify/assert"
)

func TestNewAccountFactoryCreator_NormalAccount(t *testing.T) {
	t.Parallel()

	accF, err := factory.NewAccountFactoryCreator(factory.UserAccount)
	assert.Nil(t, err)

	accWrp, err := accF.CreateAccount(mock.NewAddressMock(), &mock.AccountTrackerStub{})
	_, ok := accWrp.(*state.Account)
	assert.Equal(t, true, ok)

	assert.Nil(t, err)
	assert.NotNil(t, accF)
}

func TestNewAccountFactoryCreator_MetaAccount(t *testing.T) {
	t.Parallel()

	accF, err := factory.NewAccountFactoryCreator(factory.ShardStatistics)
	assert.Nil(t, err)

	accWrp, err := accF.CreateAccount(mock.NewAddressMock(), &mock.AccountTrackerStub{})
	_, ok := accWrp.(*state.MetaAccount)
=======
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/stretchr/testify/assert"
)

func TestNewAccountFactoryCreator_NilShardCoordinator(t *testing.T) {
	t.Parallel()

	accF, err := factory.NewAccountFactoryCreator(nil)

	assert.Equal(t, err, state.ErrNilShardCoordinator)
	assert.Nil(t, accF)
}

func TestNewAccountFactoryCreator_NormalAccount(t *testing.T) {
	t.Parallel()

	shardC := &mock.ShardCoordinatorMock{
		SelfID:     0,
		NrOfShards: 1,
	}
	accF, err := factory.NewAccountFactoryCreator(shardC)
	assert.Nil(t, err)

	accWrp, err := accF.CreateAccount(mock.NewAddressMock(), &mock.AccountTrackerStub{})
	_, ok := accWrp.(*state.Account)
>>>>>>> Stashed changes
	assert.Equal(t, true, ok)

	assert.Nil(t, err)
	assert.NotNil(t, accF)
}

<<<<<<< Updated upstream
func TestNewAccountFactoryCreator_PeerAccount(t *testing.T) {
	t.Parallel()

	accF, err := factory.NewAccountFactoryCreator(factory.ValidatorAccount)
	assert.Nil(t, err)

	accWrp, err := accF.CreateAccount(mock.NewAddressMock(), &mock.AccountTrackerStub{})
	_, ok := accWrp.(*state.PeerAccount)
=======
func TestNewAccountFactoryCreator_MetaAccount(t *testing.T) {
	t.Parallel()

	shardC := &mock.ShardCoordinatorMock{
		SelfID:     sharding.MetachainShardId,
		NrOfShards: 1,
	}
	accF, err := factory.NewAccountFactoryCreator(shardC)
	assert.Nil(t, err)

	accWrp, err := accF.CreateAccount(mock.NewAddressMock(), &mock.AccountTrackerStub{})
	_, ok := accWrp.(*state.MetaAccount)
>>>>>>> Stashed changes
	assert.Equal(t, true, ok)

	assert.Nil(t, err)
	assert.NotNil(t, accF)
}

<<<<<<< Updated upstream
func TestNewAccountFactoryCreator_UnknownType(t *testing.T) {
	t.Parallel()

	accF, err := factory.NewAccountFactoryCreator(10)
	assert.Nil(t, accF)
	assert.Equal(t, state.ErrUnknownAccountType, err)
=======
func TestNewAccountFactoryCreator_BadShardID(t *testing.T) {
	t.Parallel()

	shardC := &mock.ShardCoordinatorMock{
		SelfID:     10,
		NrOfShards: 5,
	}
	accF, err := factory.NewAccountFactoryCreator(shardC)
	assert.Nil(t, accF)
	assert.Equal(t, state.ErrUnknownShardId, err)
>>>>>>> Stashed changes
}
