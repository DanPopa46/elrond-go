package mockVM

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/integrationTests/vm"
	"github.com/ElrondNetwork/elrond-go/process/factory"
	"github.com/ElrondNetwork/elrond-go/process/smartContract"
	"github.com/stretchr/testify/assert"
)

func TestVmGetShouldReturnValue(t *testing.T) {
	accnts, destinationAddressBytes, expectedValueForVar := deploySmartContract(t)

	mockVM := vm.CreateOneSCExecutorMockVM(accnts)
	scgd, _ := smartContract.NewSCDataGetter(mockVM)

	functionName := "Get"
	returnedVals, err := scgd.Get(destinationAddressBytes, functionName)

	assert.Nil(t, err)
	assert.Equal(t, expectedValueForVar.Bytes(), returnedVals)
}

func deploySmartContract(t *testing.T) (state.AccountsAdapter, []byte, *big.Int) {
	vmOpGas := uint64(0)
	senderAddressBytes := []byte("12345678901234567890123456789012")
	senderNonce := uint64(11)
	senderBalance := big.NewInt(100000000)
	round := uint64(444)
	gasPrice := uint64(1)
	gasLimit := vmOpGas
	transferOnCalls := big.NewInt(0)

	initialValueForInternalVariable := uint64(45)
	scCode := fmt.Sprintf("aaaa@%s@%X", hex.EncodeToString(factory.InternalTestingVM), initialValueForInternalVariable)

	tx := vm.CreateTx(
		t,
		senderAddressBytes,
		vm.CreateEmptyAddress().Bytes(),
		senderNonce,
		transferOnCalls,
		gasPrice,
		gasLimit,
		scCode,
	)

	txProc, accnts := vm.CreatePreparedTxProcessorAndAccountsWithMockedVM(t, vmOpGas, senderNonce, senderAddressBytes, senderBalance)

	err := txProc.ProcessTransaction(tx, round)
	assert.Nil(t, err)

	_, err = accnts.Commit()
	assert.Nil(t, err)

<<<<<<< Updated upstream
	destinationAddressBytes, _ := hex.DecodeString("0000000000000000ffff1a2983b179a480a60c4308da48f13b4480dbb4d33132")
=======
	destinationAddressBytes, _ := hex.DecodeString("000000000000000000001a2983b179a480a60c4308da48f13b4480dbb4d33132")
>>>>>>> Stashed changes
	return accnts, destinationAddressBytes, big.NewInt(0).SetUint64(initialValueForInternalVariable)
}
