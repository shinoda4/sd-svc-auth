package logger

import (
	"log"
)

func Init() {
	// placeholder: use a structured logger in prod
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
