package data

import (
	"errors"
)

// ErrNilHeadersDataPool signals that a nil header pool has been provided
var ErrNilHeadersDataPool = errors.New("nil headers data pool")

// ErrNilHeadersNoncesDataPool signals that a nil header - nonce cache
var ErrNilHeadersNoncesDataPool = errors.New("nil headers nonces cache")

// ErrNilCacher signals that a nil cache has been provided
var ErrNilCacher = errors.New("nil cacher")

// ErrInvalidHeaderType signals an invalid header pointer was provided
var ErrInvalidHeaderType = errors.New("invalid header type")

// ErrInvalidBodyType signals an invalid header pointer was provided
var ErrInvalidBodyType = errors.New("invalid body type")

// ErrNilBlockBody signals that block body is nil
var ErrNilBlockBody = errors.New("nil block body")

// ErrMiniBlockEmpty signals that mini block is empty
var ErrMiniBlockEmpty = errors.New("mini block is empty")

// ErrWrongTypeAssertion signals that wrong type was provided
var ErrWrongTypeAssertion = errors.New("wrong type assertion")

// ErrNilElrondAddress signals that nil elrond address was provided
var ErrNilElrondAddress = errors.New("nil elrond address")

// ErrNilBurnAddress signals that nil burn address was provided
var ErrNilBurnAddress = errors.New("nil burn address")

// ErrNilAddressConverter signals that nil address converter was provided
var ErrNilAddressConverter = errors.New("nil address converter")

// ErrNilShardCoordinator signals that nil shard coordinator was provided
var ErrNilShardCoordinator = errors.New("nil shard coordinator")

// ErrNilNodesCoordinator signals that nil shard coordinator was provided
var ErrNilNodesCoordinator = errors.New("nil nodes coordinator")
