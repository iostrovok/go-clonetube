package clonetube

import (
	"sync"
)

type threadsChannels struct {
	sync.RWMutex
	chs0 []chan Item
	chs1 []chan Item
	pi   int
	num  int
}

func newThreadsChannels(num int) *threadsChannels {
	tc := &threadsChannels{
		pi:  0,
		num: num,
	}

	tc.generate()

	tc.pi++
	tc.generate()

	return tc
}

func (tc *threadsChannels) generate() {
	tc.Lock()
	defer tc.Unlock()

	chs := make([]chan Item, tc.num)
	for i := 0; i < tc.num; i++ {
		chs[i] = make(chan Item)
	}
	if tc.pi%2 == 1 {
		tc.chs0 = chs
	}
	if tc.pi%2 == 0 {
		tc.chs1 = chs
	}
}

func (tc *threadsChannels) nextSet() []chan Item {
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

	return []chan Item{}
}

/*
	newThreadi is simple func.
	Clone interface{} and result send into channel.
	If we have any incoming message we get out from function.
*/
func newThread(id int64, body interface{}, f CloneFn, ChIn, ChOut chan Item) {
	// infinity loop
	i := 0
	for {

		i++

		b, err := f(body)
		if err != nil {
			return
		}

		// If we have good clone copy
		it := Item{
			ID:   id,
			Body: b,
			I:    i,
		}

		select {
		case _, ok := <-ChIn:
			// go out by any messages
			if ok || !ok {
				return
			}
		case ChOut <- it:
			// Nothiing
		}
	}
}
