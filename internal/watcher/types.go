package watcher

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/loic-roux-404/crypto-bots/internal/model/kecacc"
)

// AccSubCallback function passed in subscriber callback
type AccSubCallback func(tx *kecacc.Tx)

// ScSubCallback processed
type ScSubCallback func(log types.Log)

// WatcherSc default interface
type WatcherSc interface {
	RunEventLoop(callback ScSubCallback)
}
// WatcherAcc default
type WatcherAcc interface {
	RunEventLoop(callback AccSubCallback)
}

// LogAdapter Use case for generics
type LogAdapter interface {
	Log() interface{}
}
