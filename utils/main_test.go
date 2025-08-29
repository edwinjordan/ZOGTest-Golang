package utils

import (
	"io"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Suppress logs in tests
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}
