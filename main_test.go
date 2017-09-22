package main

import (
	"os"
	"testing"
)

// Test methods

func TestMain(m *testing.M) {

	main()

	code := m.Run()

	os.Exit(code)
}
