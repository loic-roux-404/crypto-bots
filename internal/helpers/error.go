package helpers

import (
	"log"
)

// RecoverAndLog errors
func RecoverAndLog() {
	r := recover();
	if _, ok := r.(error); r != nil && ok {
		log.Printf("Error: %s", r)
	}
}
