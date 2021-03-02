package factory_test

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go/data/mock"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/data/state/factory"
<<<<<<< Updated upstream
=======
	"github.com/ElrondNetwork/elrond-go/sharding"
>>>>>>> Stashed changes
	"github.com/stretchr/testify/assert"
)

func TestMetaAccountCreator_CreateAccountNilAddress(t *testing.T) {
	t.Parallel()

<<<<<<< Updated upstream
	accF, err := factory.NewAccountFactoryCreator(factory.ShardStatistics)
=======
	shardC := &mock.ShardCoordinatorMock{
		SelfID:     sharding.MetachainShardId,
		NrOfShards: 1,
	}
	accF, err := factory.NewAccountFactoryCreator(shardC)
>>>>>>> Stashed changes
	assert.Nil(t, err)

	_, ok := accF.(*factory.MetaAccountCreator)
	assert.Equal(t, true, ok)

	acc, err := accF.CreateAccount(nil, &mock.AccountTrackerStub{})

	assert.Nil(t, acc)
	assert.Equal(t, err, state.ErrNilAddressContainer)
}

func TestMetaAccountCreator_CreateAccountNilAccountTraccer(t *testing.T) {
	t.Parallel()

<<<<<<< Updated upstream
	accF, err := factory.NewAccountFactoryCreator(factory.ShardStatistics)
=======
	shardC := &mock.ShardCoordinatorMock{
		SelfID:     sharding.MetachainShardId,
		NrOfShards: 1,
	}
	accF, err := factory.NewAccountFactoryCreator(shardC)
>>>>>>> Stashed changes
	assert.Nil(t, err)

	_, ok := accF.(*factory.MetaAccountCreator)
	assert.Equal(t, true, ok)

	acc, err := accF.CreateAccount(&mock.AddressMock{}, nil)

	assert.Nil(t, acc)
	assert.Equal(t, err, state.ErrNilAccountTracker)
}

func TestMetaAccountCreator_CreateAccountOk(t *testing.T) {
	t.Parallel()

<<<<<<< Updated upstream
	accF, err := factory.NewAccountFactoryCreator(factory.ShardStatistics)
=======
	shardC := &mock.ShardCoordinatorMock{
		SelfID:     sharding.MetachainShardId,
		NrOfShards: 1,
	}
	accF, err := factory.NewAccountFactoryCreator(shardC)
>>>>>>> Stashed changes
	assert.Nil(t, err)

	_, ok := accF.(*factory.MetaAccountCreator)
	assert.Equal(t, true, ok)

	acc, err := accF.CreateAccount(&mock.AddressMock{}, &mock.AccountTrackerStub{})

	assert.NotNil(t, acc)
	assert.Nil(t, err)
}
