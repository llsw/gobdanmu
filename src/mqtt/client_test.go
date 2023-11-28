package client

import (
	"sync"
	"testing"
)

func TestRun(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	Run()
	wg.Wait()
}
