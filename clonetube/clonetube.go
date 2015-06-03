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
	Len       int
	Func      CloneFn
	ChOut     chan Item
	ChIn      chan interface{}
	Body      interface{}
}

type Item struct {
	ID   int64
	Body interface{}
}

func New(l int, f CloneFn) *Main {

	chIn := make(chan interface{}, 10)
	chOut := make(chan Item, l)

	cloneTube := &Main{
		CurrentID: 0,
		Len:       l,
		Func:      f,
		ChOut:     chOut,
		ChIn:      chIn,
		Body:      nil,
	}

	return cloneTube
}

func (cloneTube *Main) Start() *Main {

	// Ready for work
	go cloneTube._start()
	return cloneTube
}

func (cloneTube *Main) _readIn(in interface{}, ok bool) bool {

	// Clean old data
	go cloneTube.clean(cloneTube.CurrentID)

	cloneTube.Body = in
	cloneTube.CurrentID++

	if !ok {
		cloneTube.closeChanals()
		return false
	}

	if _, hasStop := in.(isStop); hasStop {
		cloneTube.closeChanals()
		return false
	}

	return true
}

func (cloneTube *Main) _start() {
	//
	for {

		b, err := cloneTube.Func(cloneTube.Body)
		if err != nil {
			select {
			case in, ok := <-cloneTube.ChIn:
				if !cloneTube._readIn(in, ok) {
					return
				}
			}
			continue
		}

		// If we have good clone copy
		it := Item{
			ID:   cloneTube.CurrentID,
			Body: b,
		}

		select {
		case in, ok := <-cloneTube.ChIn:
			if !cloneTube._readIn(in, ok) {
				return
			}
		case cloneTube.ChOut <- it:
			// Nothiing
		}
	}
}

func (cloneTube *Main) Get(timeOutMicrosecond int) (interface{}, error) {
	//
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

func (cloneTube *Main) clean(id int64) {

	var lastID int64 = 0
	for lastID <= id {
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

	if cloneTube.ChOut != nil {
		close(cloneTube.ChOut)
	}

	if cloneTube.ChIn != nil {
		close(cloneTube.ChIn)
	}
}

/*
	Close chanel when stop work
*/
func (cloneTube *Main) Stop() {
	cloneTube.ChIn <- isStop{}
}
