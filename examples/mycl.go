package main

/*
	It's a simple example how you can use "clonetube"
*/

import (
	"github.com/iostrovok/go-clonetube/clonetube"

	"fmt"
	"math/rand"
	"runtime"
)

const (
	// 10000 microsecons = 10 milliseconds
	TTL         int = 10000
	ChLen       int = 10
	ThreadLen   int = 5
	NumberIters int = 1000
)

func main() {

	runtime.GOMAXPROCS(8)

	fmt.Println("START\n")

	countErrors := 0
	obj := genTestTree(5, 5)
	cl := clonetube.New(ChLen, MyTestFuncClone, ThreadLen)
	cl.Put(obj)

	for i := 0; i < NumberIters; i++ {
		_, err := cl.Get(TTL)
		if err != nil {
			countErrors++
			fmt.Printf("Reading error--> %d: %s\n", i, err)
		}
	}

	fmt.Printf("FINISH %d cycles. Errors: %d\n", NumberIters, countErrors)

}

/*------------------------------------------------*/

type CloneTestStruture struct {
	ID   int
	Text string
	List []CloneTestStruture
}

func MyTestFuncClone(in interface{}) (interface{}, error) {

	i, ok := in.(CloneTestStruture)
	if !ok {
		return nil, fmt.Errorf("Bad input params\n")
	}

	i.ID++

	return i.Clone(), nil
}

func (cl CloneTestStruture) Clone() CloneTestStruture {

	out := CloneTestStruture{
		ID:   cl.ID,
		Text: cl.Text,
		List: make([]CloneTestStruture, len(cl.List)),
	}

	for i := 0; i < len(cl.List); i++ {
		out.List[i] = cl.List[i].Clone()
	}

	return out
}

func genTestTree(depth, width int) CloneTestStruture {

	out := CloneTestStruture{
		Text: randSeq(30),
	}

	if depth <= 0 {
		return out
	}

	depth--

	out.List = make([]CloneTestStruture, width)
	for i := 0; i < width; i++ {
		out.List[i] = genTestTree(depth, width)
	}

	return out
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
