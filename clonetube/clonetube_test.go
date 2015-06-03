package clonetube

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestCloneTube(t *testing.T) {
	TestingT(t)
}

const (
	// 10000 microsecons = 10 milliseconds
	TTL   int = 10000
	ChLen int = 10
)

type CloneTubeTestsSuite struct{}

var _ = Suite(&CloneTubeTestsSuite{})

func (s *CloneTubeTestsSuite) TestCloneTubeNew(c *C) {

	//c.Skip("Not now")
	ct := New(ChLen, MyTestFuncClone)

	c.Assert(ct, NotNil)
}

func (s *CloneTubeTestsSuite) TestCloneTubeStartStop(c *C) {

	//c.Skip("Not now")
	ct := New(ChLen, MyTestFuncClone).Start()
	c.Assert(ct, NotNil)

	ct.Stop()
}

func (s *CloneTubeTestsSuite) TestCloneTubePut(c *C) {

	//c.Skip("Not now")
	ct := New(ChLen, MyTestFuncClone).Start()
	c.Assert(ct, NotNil)

	cl := genTestTree(5, 5)

	err := ct.Put(cl)
	c.Assert(err, IsNil)

	ct.Stop()
}

func (s *CloneTubeTestsSuite) TestCloneTubeGet_100(c *C) {

	//c.Skip("Not now")
	ct := New(ChLen, MyTestFuncClone).Start()
	c.Assert(ct, NotNil)

	cl := genTestTree(5, 5)

	err := ct.Put(cl)
	c.Assert(err, IsNil)

	for i := 0; i < 100; i++ {

		w, err := ct.Get(TTL)
		c.Assert(err, IsNil)
		c.Assert(w, NotNil)

		inter, ok := w.(CloneTestStruture)
		c.Assert(ok, Equals, true)
		c.Assert(inter.ID, Equals, 1)
	}

	ct.Stop()
}

func (s *CloneTubeTestsSuite) TestCloneTube_get_100_Put(c *C) {

	//c.Skip("Not now")
	ct := New(ChLen, MyTestFuncClone).Start()
	c.Assert(ct, NotNil)

	for j := 1; j < 10; j++ {
		cl := genTestTree(5, 5)

		ct.Put(cl)
		for i := 0; i < 100; i++ {
			w, err := ct._get(TTL)
			c.Assert(err, IsNil)

			if i > 2*ChLen {
				// Skip channel length * 2 so we can take old value.
				c.Assert(w.ID, Equals, int64(j))
			}
		}
	}

	ct.Stop()
}
