package clonetube

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Preloaded params
const (
	DefaultTTL          int    = 10000
	DefaultThreadNumber int    = 100
	debugDateTimeFormat string = "2006-01-02 15:04:05.99999"
)

// Preloaded params
type managerStruct struct {
	debug bool
	sync.RWMutex
	threadNumber  int
	cloneTubeList map[string]*oneManagerStruct
}

type oneManagerStruct struct {
	cloneTube *Main
	stopCh    chan interface{}
}

// Manager - "constructor"
func Manager(debugs ...bool) *managerStruct {

	debug := false
	if len(debugs) > 0 {
		debug = debugs[0]
	}

	m := &managerStruct{
		debug:         debug,
		threadNumber:  DefaultThreadNumber,
		cloneTubeList: map[string]*oneManagerStruct{},
	}

	return m
}

// SetThreadNumber. It sets thread number for clontube
func (m *managerStruct) SetThreadNumber(n int) error {
	if n < 1 {
		return errors.New("ThreadNumber must be positive.")
	}
	m.threadNumber = n
	return nil
}

// GetThreadNumber - back of SetThreadNumber.
func (m *managerStruct) GetThreadNumber() int {
	return m.threadNumber
}

// internal function
func (m *managerStruct) printDebug(format string, params ...interface{}) {
	if !m.debug {
		return
	}
	fmt.Printf(format, params...)
}

// internal function
func (m *managerStruct) init(point string, count int, funcClone CloneFn,
	threadNumber int) (*oneManagerStruct, error) {

	// check inner params
	if threadNumber < 1 {
		threadNumber = DefaultThreadNumber
	}

	if point == "" {
		return nil, errors.New("Type of manager [point] must be defined")
	}

	if count < 1 {
		return nil, errors.New("Count of manager [point] must be lager than null")
	}

	// save mew manager structure
	m.Lock()
	defer m.Unlock()

	obj := &oneManagerStruct{
		cloneTube: New(count, funcClone, threadNumber),
		stopCh:    make(chan interface{}, 1),
	}

	m.cloneTubeList[point] = obj

	return obj, nil
}

// Start - adds and starts new clontube object with same init and clone functions
func (m *managerStruct) Start(point string, TTL time.Duration, count int,
	funcNew, funcClone CloneFn, params ...interface{}) error {

	// check inner params
	if funcNew == nil || funcClone == nil {
		return errors.New("function must be no nil")
	}

	if TTL == 0 {
		return errors.New("TTL must be more than zero")
	}

	threadNumber := m.threadNumber
	if len(params) > 0 {
		i, ok := params[0].(int)
		if ok {
			threadNumber = i
		}
		params = params[0:]
	}

	// init manager
	obj, err := m.init(point, count, funcClone, threadNumber)
	if err != nil {
		return err
	}

	// check init function and save first result for cloning
	err = m._initAutoClone(point, funcNew, params...)
	if err != nil {
		return err
	}

	go func(m *managerStruct, point string, stopCh chan interface{}, TTL time.Duration, funcNew CloneFn, params ...interface{}) {
		timer := time.NewTicker(TTL)

		for {
			select {
			case <-timer.C:
				m.printDebug("manager-clonetube. save new result for cloning: %s - %s\n", point, time.Now().Format(debugDateTimeFormat))
				// save new result for cloning
				m._initAutoClone(point, funcNew, params...)
			case _, stop := <-stopCh:
				m.printDebug("manager-clonetube: Start STOP: exit: %s - %s\n", point, time.Now().Format(debugDateTimeFormat))
				if !stop {
					m.printDebug("manager-clonetube: Start STOP: exit: %s - %s\n", point, time.Now().Format(debugDateTimeFormat))
					return
				}
			}
		}
	}(m, point, obj.stopCh, TTL, funcNew, params)

	return nil
}

// internal function
func (m *managerStruct) _initAutoClone(point string, f CloneFn, params ...interface{}) error {

	// check inner params
	if len(params) == 0 {
		params = append(params, nil)
	}

	in, err := f(params)

	if err != nil {
		return err
	}

	return m.put(point, in)
}

// Stop - stops and deletes one clontube object
func (m *managerStruct) Stop(point string) error {

	m.printDebug("manager-clonetube: Stop one: %s - %s\n", point, time.Now().Format(debugDateTimeFormat))

	m.Lock()
	obj, ok := m.cloneTubeList[point]
	if !ok {
		m.Unlock()
		return errors.New("Stop. Point " + point + " not found")
	}

	delete(m.cloneTubeList, point)
	m.Unlock()

	close(obj.stopCh)
	obj.cloneTube.Stop()

	return nil
}

// Stop - stops and deletes all clontube objects
func (m *managerStruct) StopAll() {

	m.printDebug("manager-clonetube: StopAll: all points - %s\n", time.Now().Format(debugDateTimeFormat))

	m.Lock()
	for _, obj := range m.cloneTubeList {
		close(obj.stopCh)
		obj.cloneTube.Stop()
	}

	m.cloneTubeList = map[string]*oneManagerStruct{}

	m.Unlock()
}

// Get - returns cloned object for one point
func (m *managerStruct) Get(point string, timeOutMicrosecond ...int) (interface{}, error) {

	m.RLock()
	obj, ok := m.cloneTubeList[point]
	m.RUnlock()

	if !ok {
		return nil, errors.New("Point " + point + " not found")
	}

	return obj.cloneTube.Get(timeOutMicrosecond...)
}

// internal function
func (m *managerStruct) put(point string, in interface{}) error {

	m.RLock()
	obj, ok := m.cloneTubeList[point]
	m.RUnlock()

	if !ok {
		return errors.New("Point " + point + " not found")
	}

	return obj.cloneTube.Put(in)
}
