package wallet

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-go/integrationTests"
	"github.com/stretchr/testify/assert"
)

func TestInterceptedTxFromFrontendGeneratedParamsAllParams(t *testing.T) {
	testInterceptedTxFromFrontendGeneratedParams(
		t,
		0,
		big.NewInt(10),
		"53669be65aac358a6add8e8a8b1251bb994dc1e4a0cc885956f5ecd53396f0d8",
		"a10e99839fe19bdb2ec8b22e0805da40053d4e5b6ace564949f26d49095e36e8",
		"e1e38ae48088baeca9da900cf054d71d7500171986a73cd04027d32fe3435241338979db530bd79e5148d8b0146204c9b2d985d201019a1728218841b8454a09",
		10,
		1000,
		"aa@bbbb@cccc",
	)
}

func TestInterceptedTxFromFrontendGeneratedParamsAllParams2(t *testing.T) {
	testInterceptedTxFromFrontendGeneratedParams(
		t,
		12,
		big.NewInt(2),
		"943643524936191d1c5627e044f7b5e4ca559c7d0ba1c2b85d1b2e6c299ebcd8",
		"943643524936191d1c5627e044f7b5e4ca559c7d0ba1c2b85d1b2e6c299ebcd8",
		"1ef83bae21227e93e9717f45a4ec34e3f5c6a110e31dfa438ac2b8c1f5459e5167fd8424d1dfa6de59756437fe599def6872217ddad5717fe61a41853606450c",
		1,
		10000,
		"aa@dd@cc",
	)
}

func TestInterceptedTxFromFrontendGeneratedParamsGasPriceGasLimitNoData(t *testing.T) {
	testInterceptedTxFromFrontendGeneratedParams(
		t,
		0,
		big.NewInt(10),
		"53669be65aac358a6add8e8a8b1251bb994dc1e4a0cc885956f5ecd53396f0d8",
		"6afb8018dcc5a53d22d4dcdda39ceaf25dafd1ea353a9bbe12073057f4e6d262",
		"1d96166ecd6cae86797046126b64028099fcd026a37a82c4bdd19700bd49828069a822fb5453e0b32f66ed895d4f162af35ea8aca862af498e2831c596250e03",
		10,
		1000,
		"",
	)
}

// testInterceptedTxFromFrontendGeneratedParams tests that a frontend generated tx will pass through an interceptor
// and ends up in the datapool, concluding the tx is correctly signed and follows our protocol
func testInterceptedTxFromFrontendGeneratedParams(
	t *testing.T,
	frontendNonce uint64,
	frontendValue *big.Int,
	frontendReceiverHex string,
	frontendSenderHex string,
	frontendSignature string,
	frontendGasPrice uint64,
	frontendGasLimit uint64,
	frontendData string,
) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	chDone := make(chan struct{})

	maxShards := uint32(1)
	nodeShardId := uint32(0)
	txSignPrivKeyShardId := uint32(0)
	initialNodeAddr := "nodeAddr"
	valMinting := big.NewInt(20000)

	node := integrationTests.NewTestProcessorNode(
		maxShards,
		nodeShardId,
		txSignPrivKeyShardId,
		initialNodeAddr,
	)

	txHexHash := ""

	err := node.SetAccountNonce(uint64(0))
	assert.Nil(t, err)

	node.ShardDataPool.Transactions().RegisterHandler(func(key []byte) {
		assert.Equal(t, txHexHash, hex.EncodeToString(key))

		dataRecovered, _ := node.ShardDataPool.Transactions().SearchFirstData(key)
		assert.NotNil(t, dataRecovered)

		txRecovered, ok := dataRecovered.(*transaction.Transaction)
		assert.True(t, ok)

		assert.Equal(t, frontendNonce, txRecovered.Nonce)
		assert.Equal(t, frontendValue, txRecovered.Value)

		sender, _ := hex.DecodeString(frontendSenderHex)
		assert.Equal(t, sender, txRecovered.SndAddr)

		receiver, _ := hex.DecodeString(frontendReceiverHex)
		assert.Equal(t, receiver, txRecovered.RcvAddr)

		sig, _ := hex.DecodeString(frontendSignature)
		assert.Equal(t, sig, txRecovered.Signature)
		assert.Equal(t, frontendData, txRecovered.Data)

		chDone <- struct{}{}
	})

	rcvAddrBytes, _ := hex.DecodeString(frontendReceiverHex)
	sndAddrBytes, _ := hex.DecodeString(frontendSenderHex)
	signatureBytes, _ := hex.DecodeString(frontendSignature)

	integrationTests.MintAddress(node.AccntState, sndAddrBytes, valMinting)

	txHexHash, err = node.SendTransaction(&transaction.Transaction{
		Nonce:     frontendNonce,
		Value:     frontendValue,
		RcvAddr:   rcvAddrBytes,
		SndAddr:   sndAddrBytes,
		GasPrice:  frontendGasPrice,
		GasLimit:  frontendGasLimit,
		Data:      frontendData,
		Signature: signatureBytes,
	})

	assert.Nil(t, err)

	select {
	case <-chDone:
	case <-time.After(time.Second * 2):
		assert.Fail(t, "timeout getting transaction")
	}
}
