package watcher

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/loic-roux-404/crypto-bots/internal/model/transaction"
)

// AccSubCallback function passed in subscriber callback
type AccSubCallback func(tx *transaction.Tx)

// ScSubCallback processed
type ScSubCallback func(log types.Log)

// Watcher default interface
type WatcherSc interface{
	RunEventLoop(callback ScSubCallback)
}

type WatcherAcc interface{
	RunEventLoop(callback AccSubCallback)
}

// Use case for generics
type LogAdapter interface{
	Log() interface{}
}
