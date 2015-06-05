package clonetube

import (
	"errors"
	"sync"
	"time"
)

type CloneFn func(interface{}) (interface{}, error)

type isStop struct {
	// empty
}

type Main struct {
	sync.RWMutex
	currentID  int64
	threads    []chan *item
	lenChannel int
	numThreads int
	cFunc      CloneFn
	chOut      chan *item
	chIn       chan interface{}
	tCh        *threadsChannels
	wg         *sync.WaitGroup
	body       interface{}
}

type item struct {
	i    int
	id   int64
	body interface{}
}

func New(l int, f CloneFn, threadNumber ...int) *Main {

	chIn := make(chan interface{}, 10)
	chOut := make(chan *item, l)

	cloneTube := &Main{
		currentID:  0,
		lenChannel: l,
		numThreads: 1,
		cFunc:      f,
		chOut:      chOut,
		chIn:       chIn,
		body:       nil,
		wg:         &sync.WaitGroup{},
	}

	if len(threadNumber) > 0 && threadNumber[0] > 0 {
		cloneTube.numThreads = threadNumber[0]
	}

	cloneTube.tCh = newThreadsChannels(cloneTube.numThreads, cloneTube.lenChannel)
	cloneTube.threads = make([]chan *item, 0)

	go cloneTube._start()

	return cloneTube
}

func (cloneTube *Main) _start() {
	for {
		select {
		case in, ok := <-cloneTube.chIn:
			if !cloneTube._readIn(in, ok) {
				return
			}
		}
	}
}

func (cloneTube *Main) _readIn(in interface{}, ok bool) bool {

	if !ok {
		cloneTube.closeChanals()
		return false
	}

	if _, hasStop := in.(isStop); hasStop {
		cloneTube.closeChanals()
		return false
	}

	// Clean old data
	for _, ch := range cloneTube.threads {
		ch <- &item{}
	}

	go cloneTube.clean(cloneTube.currentID)

	cloneTube.body = in
	cloneTube.currentID++
	threads := cloneTube.tCh.nextSet()
	cloneTube.wg.Add(cloneTube.numThreads)
	for i := 0; i < cloneTube.numThreads; i++ {
		go func(id int64, body interface{}, f CloneFn, ChIn, ChOut chan *item) {
			newThread(id, body, f, ChIn, ChOut)
			cloneTube.wg.Done()
		}(cloneTube.currentID, cloneTube.body, cloneTube.cFunc, threads[i], cloneTube.chOut)
	}

	cloneTube.threads = threads

	return true
}

func (cloneTube *Main) Get(timeOutMicrosecond ...int) (interface{}, error) {
	out, err := cloneTube._get(timeOutMicrosecond...)
	if err != nil {
		return nil, err
	}
	return out.body, nil
}

func (cloneTube *Main) _get(timeOutMicrosecond ...int) (*item, error) {
	var (
		err          error
		out          *item
		ok           bool
		microseconds int
	)

	if len(timeOutMicrosecond) > 0 {
		microseconds = timeOutMicrosecond[0]
	}

	if microseconds > 0 {
		select {
		case out, ok = <-cloneTube.chOut:
			if !ok {
				err = errors.New("cloneTube.chOut is closed. You have to run cloneTube.New(....)")
			}
		case <-time.After(time.Microsecond * time.Duration(microseconds)):
			err = errors.New("timeout cloneTube.Get")
		}
	} else {
		select {
		case out, ok = <-cloneTube.chOut:
			if !ok {
				err = errors.New("cloneTube.chOut is closed. You have to run cloneTube.New(....)")
			}
		}
	}
	return out, err

}

func (cloneTube *Main) Put(in interface{}) error {

	if _, err := cloneTube.cFunc(in); err != nil {
		return err
	}

	cloneTube.chIn <- in

	return nil
}

/*
	Skip old values from output channel
*/
func (cloneTube *Main) clean(id int64) {

	var lastID int64 = 0
	for lastID < id {
		select {
		case it, ok := <-cloneTube.chOut:
			if !ok {
				return
			}
			lastID = it.id
		}
	}
}

/*
	Close chanel when stop work
*/
func (cloneTube *Main) closeChanals() {
	cloneTube.Lock()
	defer cloneTube.Unlock()

	for _, ch := range cloneTube.threads {
		ch <- &item{}
	}

	cloneTube.wg.Wait()

	if cloneTube.chIn != nil {
		close(cloneTube.chIn)
	}

	if cloneTube.chOut != nil {
		close(cloneTube.chOut)
	}
}

/*
	Stop our gorroutens*
*/
func (cloneTube *Main) Stop() {
	cloneTube.chIn <- isStop{}
}
