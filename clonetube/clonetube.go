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
	CurrentID int64
	Threads   []chan Item
	Len       int
	NumThr    int
	Func      CloneFn
	ChOut     chan Item
	ChIn      chan interface{}
	Body      interface{}

	tCh *threadsChannels
	WG  *sync.WaitGroup
}

type Item struct {
	I    int
	ID   int64
	Body interface{}
}

func New(l int, f CloneFn, threadNumber ...int) *Main {

	chIn := make(chan interface{}, 10)
	chOut := make(chan Item, l)

	cloneTube := &Main{
		CurrentID: 0,
		Len:       l,
		NumThr:    1,
		Func:      f,
		ChOut:     chOut,
		ChIn:      chIn,
		Body:      nil,
		WG:        &sync.WaitGroup{},
	}

	if len(threadNumber) > 0 && threadNumber[0] > 0 {
		cloneTube.NumThr = threadNumber[0]
	}

	cloneTube.tCh = newThreadsChannels(cloneTube.NumThr)
	cloneTube.Threads = make([]chan Item, 0)

	return cloneTube
}

func (cloneTube *Main) Start() *Main {
	// Ready for work
	go cloneTube._start()
	return cloneTube
}

func (cloneTube *Main) _readIn(in interface{}, ok bool) bool {

	//cloneTube.Lock()
	//defer cloneTube.Unlock()

	if !ok {
		cloneTube.closeChanals()
		return false
	}

	if _, hasStop := in.(isStop); hasStop {
		cloneTube.closeChanals()
		return false
	}

	// Clean old data
	for _, ch := range cloneTube.Threads {
		ch <- Item{}
	}

	go cloneTube.clean(cloneTube.CurrentID)

	cloneTube.Body = in
	cloneTube.CurrentID++
	threads := cloneTube.tCh.nextSet()
	cloneTube.WG.Add(cloneTube.NumThr)
	for i := 0; i < cloneTube.NumThr; i++ {
		go func(id int64, body interface{}, f CloneFn, ChIn, ChOut chan Item) {
			newThread(id, body, f, ChIn, ChOut)
			cloneTube.WG.Done()
		}(cloneTube.CurrentID, cloneTube.Body, cloneTube.Func, threads[i], cloneTube.ChOut)
	}

	cloneTube.Threads = threads

	return true
}

func (cloneTube *Main) _start() {
	//
	for {
		select {
		case in, ok := <-cloneTube.ChIn:
			if !cloneTube._readIn(in, ok) {
				return
			}
		}
	}
}

func (cloneTube *Main) Get(timeOutMicrosecond int) (interface{}, error) {
	out, err := cloneTube._get(timeOutMicrosecond)
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

func (cloneTube *Main) _get(timeOutMicrosecond int) (Item, error) {
	var (
		err error
		out Item
		ok  bool
	)

	select {
	case out, ok = <-cloneTube.ChOut:
		if !ok {
			err = errors.New("cloneTube.ChOut is closed. You have to run cloneTube.New(....)")
		}
	case <-time.After(time.Microsecond * time.Duration(timeOutMicrosecond)):
		err = errors.New("timeout cloneTube.Get")
	}

	return out, err
}

func (cloneTube *Main) Put(in interface{}) error {

	if _, err := cloneTube.Func(in); err != nil {
		return err
	}

	cloneTube.ChIn <- in

	return nil
}

/*
	Skip old values from output channel
*/
func (cloneTube *Main) clean(id int64) {

	var lastID int64 = 0
	for lastID < id {
		select {
		case it, ok := <-cloneTube.ChOut:
			if !ok {
				return
			}
			lastID = it.ID
		}
	}
}

/*
	Close chanel when stop work
*/
func (cloneTube *Main) closeChanals() {
	cloneTube.Lock()
	defer cloneTube.Unlock()

	for _, ch := range cloneTube.Threads {
		ch <- Item{}
	}

	cloneTube.WG.Wait()

	if cloneTube.ChIn != nil {
		close(cloneTube.ChIn)
	}

	if cloneTube.ChOut != nil {
		close(cloneTube.ChOut)
	}
}

/*
	Stop our gorroutens*
*/
func (cloneTube *Main) Stop() {
	cloneTube.ChIn <- isStop{}
}
