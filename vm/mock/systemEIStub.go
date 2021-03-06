package mock

import (
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"math/big"
)

type SystemEIStub struct {
	TransferCalled       func(destination []byte, sender []byte, value *big.Int, input []byte) error
	GetBalanceCalled     func(addr []byte) *big.Int
	SetStorageCalled     func(key []byte, value []byte)
	GetStorageCalled     func(key []byte) []byte
	SelfDestructCalled   func(beneficiary []byte)
	CreateVMOutputCalled func() *vmcommon.VMOutput
	CleanCacheCalled     func()
}

func (s *SystemEIStub) SetSCAddress(addr []byte) {
}

func (s *SystemEIStub) Transfer(destination []byte, sender []byte, value *big.Int, input []byte) error {
	if s.TransferCalled != nil {
		return s.TransferCalled(destination, sender, value, input)
	}
	return nil
}

func (s *SystemEIStub) GetBalance(addr []byte) *big.Int {
	if s.GetBalanceCalled != nil {
		return s.GetBalanceCalled(addr)
	}
	return big.NewInt(0)
}

func (s *SystemEIStub) SetStorage(key []byte, value []byte) {
	if s.SetStorageCalled != nil {
		s.SetStorageCalled(key, value)
	}
}

func (s *SystemEIStub) GetStorage(key []byte) []byte {
	if s.GetStorageCalled != nil {
		return s.GetStorageCalled(key)
	}
	return nil
}

func (s *SystemEIStub) SelfDestruct(beneficiary []byte) {
	if s.SelfDestructCalled != nil {
		s.SelfDestructCalled(beneficiary)
	}
	return
}

func (s *SystemEIStub) CreateVMOutput() *vmcommon.VMOutput {
	if s.CreateVMOutputCalled != nil {
		return s.CreateVMOutputCalled()
	}

	return &vmcommon.VMOutput{}
}

func (s *SystemEIStub) CleanCache() {
	if s.CleanCacheCalled != nil {
		s.CleanCacheCalled()
	}
	return
}

func (s *SystemEIStub) IsInterfaceNil() bool {
	if s == nil {
		return true
	}
	return false
}
