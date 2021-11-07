package helpers

import (
	"log"
)

// RecoverAndLog errors
func RecoverAndLog() {
	r := recover()
	if _, ok := r.(error); r != nil && ok {
		log.Printf("Error: %s", r.(error))
	}

	if _, ok := r.(string); r != nil && ok {
		log.Printf("Warn: %s", r.(string))
	}
}
