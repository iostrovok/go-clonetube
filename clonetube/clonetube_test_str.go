package clonetube

/*
	Test structure for clonetube_test.go
*/

import (
	"fmt"
	"time"
)

type CloneTestStruture struct {
	ID   int
	now  time.Time
	List []CloneTestStruture
}

func MyTestFuncClone(in interface{}) (interface{}, error) {

	i, ok := in.(CloneTestStruture)
	if !ok {
		return nil, fmt.Errorf("MyTestFuncClone: Bad input params: %T\n", in)
	}

	i.ID++

	return i, nil
}

func MyTestFuncNewStruct(in interface{}) (interface{}, error) {
	return genTestTree(1, 1), nil
}

func genTestTree(depth, width int) CloneTestStruture {

	out := &CloneTestStruture{
		now: time.Now(),
	}

	if depth <= 0 {
		return *out
	}

	depth--

	out.List = make([]CloneTestStruture, width)
	for i := 0; i < width; i++ {
		out.List[i] = genTestTree(depth, width)
	}

	return *out
}
