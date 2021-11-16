package sub

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/loic-roux-404/crypto-bots/internal/model/tx"
)

// AccSubCallback function passed in subscriber callback
type AccSubCallback func(tx *tx.Adapter)

// ScSubCallback processed
// TODO create adapter to log, raw rlp data ?
type ScSubCallback func(log types.Log)

// Sc default interface
type Sc interface {
	RunEventLoop(callback ScSubCallback)
}

// Acc default
type Acc interface {
	RunEventLoop(callback AccSubCallback)
}

// LogAdapter Use case for generics
type LogAdapter interface {
	Log() interface{}
}
