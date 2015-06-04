package clonetube

/*
	Test structure for clonetube_test.go
*/

import "fmt"

type CloneTestStruture struct {
	ID   int
	List []CloneTestStruture
}

func MyTestFuncClone(in interface{}) (interface{}, error) {

	i, ok := in.(CloneTestStruture)
	if !ok {
		return nil, fmt.Errorf("Bad input params\n")
	}

	i.ID++

	return i, nil
}

func genTestTree(depth, width int) CloneTestStruture {

	out := &CloneTestStruture{}

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
