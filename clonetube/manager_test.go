package clonetube

import (
	. "gopkg.in/check.v1"
	"testing"
	"time"

	"fmt"
)

func TestManager(t *testing.T) {
	TestingT(t)
}

const (
	mTTL       time.Duration = 500 * time.Millisecond
	mChLen     int           = 10
	mThreadLen int           = 5
)

type ManagerTestsSuite struct{}

var _ = Suite(&ManagerTestsSuite{})

func _add_list_points(c *C, m *managerStruct, list []string) {
	for _, point := range list {
		err := m.Start(point, mTTL, 100, MyTestFuncNewStruct, MyTestFuncClone)
		c.Assert(err, IsNil)
	}
}

func _test_list_points(c *C, m *managerStruct, point string, count int) {
	for i := 0; i < 10; i++ {
		obj, err := m.Get(point)
		fmt.Printf("\n=======\nobj: %+v\n", obj)

		c.Assert(err, IsNil)
		if w, ok := obj.(CloneTestStruture); ok {

			fmt.Printf("ID: %d - %s\n", w.ID, w.now.Format(debugDateTimeFormat))
		} else {
			c.Errorf("TestManager_Get: bad returned object")
		}

		//time.Sleep(300 * time.Millisecond)
	}
}

func (s *ManagerTestsSuite) TestManager_Manager(c *C) {
	//c.Skip("Not now")
	m := Manager()
	c.Assert(m, NotNil)
}

func (s *ManagerTestsSuite) TestManager_Manager_Debug(c *C) {
	//c.Skip("Not now")
	m := Manager(false)
	c.Assert(m, NotNil)
}

func (s *ManagerTestsSuite) TestManager_init(c *C) {

	//c.Skip("Not now")
	m := Manager(false)
	obj, err := m.init("point", 5, MyTestFuncClone, mThreadLen)

	c.Assert(obj, NotNil)
	c.Assert(err, IsNil)
}

func (s *ManagerTestsSuite) TestManager_Start(c *C) {

	//c.Skip("Not now")
	m := Manager()

	err := m.Start("point", mTTL, 100, MyTestFuncNewStruct, MyTestFuncClone)

	c.Assert(err, IsNil)
	c.Assert(len(m.cloneTubeList), Equals, 1)
}

func (s *ManagerTestsSuite) TestManager_Start_err0(c *C) {
	//c.Skip("Not now")
	m := Manager(false)
	err := m.Start("", mTTL, 100, MyTestFuncNewStruct, MyTestFuncClone)
	c.Assert(err, NotNil)
}

func (s *ManagerTestsSuite) TestManager_Start_err1(c *C) {
	//c.Skip("Not now")
	m := Manager(false)
	err := m.Start("point", mTTL, 100, nil, MyTestFuncClone)
	c.Assert(err, NotNil)
}

func (s *ManagerTestsSuite) TestManager_Start_err2(c *C) {
	//c.Skip("Not now")
	m := Manager()
	err := m.Start("point", mTTL, 100, MyTestFuncNewStruct, nil)
	c.Assert(err, NotNil)
}

func (s *ManagerTestsSuite) TestManager_Start_err3(c *C) {
	//c.Skip("Not now")
	m := Manager()
	err := m.Start("point", mTTL, 0, MyTestFuncNewStruct, MyTestFuncClone)
	c.Assert(err, NotNil)
}

func (s *ManagerTestsSuite) TestManager_Start_err4(c *C) {
	//c.Skip("Not now")
	m := Manager()
	err := m.Start("point", 0, 100, MyTestFuncNewStruct, MyTestFuncClone)
	c.Assert(err, NotNil)
}

func (s *ManagerTestsSuite) TestManager_Get(c *C) {
	//c.Skip("Not now")
	m := Manager(true)
	err := m.Start("point", mTTL, 100, MyTestFuncNewStruct, MyTestFuncClone)
	c.Assert(err, IsNil)

	for i := 0; i < 10; i++ {
		obj, err := m.Get("point")
		fmt.Printf("\n=======\nobj: %+v\n", obj)

		c.Assert(err, IsNil)
		if w, ok := obj.(CloneTestStruture); ok {

			fmt.Printf("ID: %d - %s\n", w.ID, w.now.Format(debugDateTimeFormat))
		} else {
			c.Errorf("TestManager_Get: bad returned object")
		}

		time.Sleep(300 * time.Millisecond)
	}
}

func (s *ManagerTestsSuite) TestManager_Stop(c *C) {
	//c.Skip("Not now")
	m := Manager(true)
	_add_list_points(c, m, []string{"point", "point-second"})
	c.Assert(len(m.cloneTubeList), Equals, 2)

	_test_list_points(c, m, "point", 10)

	err := m.Stop("point")
	c.Assert(err, IsNil)

	c.Assert(len(m.cloneTubeList), Equals, 1)

	obj, err := m.Get("point")
	c.Assert(obj, IsNil)
	c.Assert(err, NotNil)

	err = m.Stop("point")
	c.Assert(err, NotNil)
}

func (s *ManagerTestsSuite) TestManager_StopAll(c *C) {
	//c.Skip("Not now")
	m := Manager(true)
	_add_list_points(c, m, []string{"point", "point-2", "point-3"})
	c.Assert(len(m.cloneTubeList), Equals, 3)

	_test_list_points(c, m, "point", 10)

	m.StopAll()

	c.Assert(len(m.cloneTubeList), Equals, 0)

	obj, err := m.Get("point")
	c.Assert(obj, IsNil)
	c.Assert(err, NotNil)

	err = m.Stop("point")
	c.Assert(err, NotNil)

	time.Sleep(5 * time.Second)
}
