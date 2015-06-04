package clonetube

import (
	"sync"
)

type threadsChannels struct {
	sync.RWMutex
	chs0 []chan *item
	chs1 []chan *item

	pi        int
	numCannel int
	lenCannel int
}

func newThreadsChannels(numCannel, lenCannel int) *threadsChannels {
	tc := &threadsChannels{
		pi:        0,
		numCannel: numCannel,
		lenCannel: lenCannel,
	}

	tc.generate()

	tc.pi++
	tc.generate()

	return tc
}

func (tc *threadsChannels) generate() {
	tc.Lock()
	defer tc.Unlock()

	chs := make([]chan *item, tc.numCannel)
	for i := 0; i < tc.numCannel; i++ {
		chs[i] = make(chan *item, tc.lenCannel)
	}
	if tc.pi%2 == 1 {
		tc.chs0 = chs
	}
	if tc.pi%2 == 0 {
		tc.chs1 = chs
	}
}

func (tc *threadsChannels) nextSet() []chan *item {
	tc.Lock()
	defer tc.Unlock()

	tc.pi++

	go tc.generate()

	if tc.pi%2 == 0 {
		return tc.chs0
	}
	if tc.pi%2 == 1 {
		return tc.chs1
	}

	return []chan *item{}
}
