package mock

import (
	"encoding/binary"
	"errors"
)

type nonceHashConverterMock struct {
}

<<<<<<< Updated upstream
=======
// NewNonceHashConverterMock -
>>>>>>> Stashed changes
func NewNonceHashConverterMock() *nonceHashConverterMock {
	return &nonceHashConverterMock{}
}

<<<<<<< Updated upstream
=======
// ToByteSlice -
>>>>>>> Stashed changes
func (*nonceHashConverterMock) ToByteSlice(value uint64) []byte {
	buff := make([]byte, 8)

	binary.BigEndian.PutUint64(buff, value)

	return buff
}

<<<<<<< Updated upstream
=======
// ToUint64 -
>>>>>>> Stashed changes
func (*nonceHashConverterMock) ToUint64(buff []byte) (uint64, error) {
	if buff == nil {
		return 0, errors.New("failure, nil slice")
	}

	if len(buff) != 8 {
		return 0, errors.New("failure, len not 8")
	}

	return binary.BigEndian.Uint64(buff), nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (nhcm *nonceHashConverterMock) IsInterfaceNil() bool {
<<<<<<< Updated upstream
	if nhcm == nil {
		return true
	}
	return false
=======
	return nhcm == nil
>>>>>>> Stashed changes
}
