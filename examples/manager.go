package main

/*
It's a simple example how to use "clonetube" manager.

1) We have to get complex structure ("cloneTestStruture" in our case).
We do it with "myTestFuncNewStruct". In real life it may be getting
data from database or remote server.
Out  "myTestFuncNewStruct" works each "newCloneTTL" milliseconds.

2) Before using the complex structure we have to clone it into
new structure (so we remove common references and
aren't afraid change copy of  structure).
Function for clone is "myTestFuncClone".

*/

import (
	"github.com/iostrovok/go-clonetube/clonetube"

	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"
)

const (
	// 10000 microsecons = 10 milliseconds
	readTTL             int           = 100
	newCloneTTL         time.Duration = 500 * time.Millisecond
	threadLen           int           = 5
	numberIters         int           = 1000
	debugDateTimeFormat string        = "2006-01-02 15:04:05.99999"
)

func main() {

	runtime.GOMAXPROCS(8)

	fmt.Println("START")

	m := clonetube.Manager(false)
	err := m.Start("point", newCloneTTL, threadLen, myTestFuncNewStruct, myTestFuncClone)
	if err != nil {
		log.Fatalln(err)
	}

	lastTimeRead := ""
	lastIDRead := ""

	countErrors := 0
	countOLD := 0
	for i := 0; i < numberIters; i++ {

		obj, err := m.Get("point", readTTL)
		if err != nil {
			countErrors++
			fmt.Printf("TestManager_Get: bad returned object")
			continue
		}

		if w, ok := obj.(cloneTestStruture); ok {
			newTime := w.Now.Format(debugDateTimeFormat)
			if lastTimeRead != newTime || lastIDRead != w.Text {
				fmt.Printf("Read new time. ID: %d. Time: %s [%d objects was readed with preview time.]\n", w.ID, newTime, countOLD)
				lastTimeRead = newTime
				lastIDRead = w.Text
				countOLD = 0
			} else {
				countOLD++
			}
		} else {
			countErrors++
			fmt.Printf("TestManager_Get: bad type returned object")
		}

		time.Sleep(time.Millisecond * 5)
	}

	fmt.Printf("\n\nFINISH %d cycles. Errors: %d\n", numberIters, countErrors)
}

/*------------------------------------------------*/

/*
	Test struct
*/
type cloneTestStruture struct {
	ID   int
	Now  time.Time
	Text string
	List []cloneTestStruture
}

func myTestFuncClone(in interface{}) (interface{}, error) {

	i, ok := in.(cloneTestStruture)
	if !ok {
		return nil, fmt.Errorf("Bad input params\n")
	}

	i.ID++

	return i.clone(), nil
}

func (cl cloneTestStruture) clone() cloneTestStruture {

	out := cloneTestStruture{
		ID:   cl.ID,
		Text: cl.Text,
		Now:  cl.Now,
		List: make([]cloneTestStruture, len(cl.List)),
	}

	for i := 0; i < len(cl.List); i++ {
		out.List[i] = cl.List[i].clone()
	}

	return out
}

//
func myTestFuncNewStruct(in interface{}) (interface{}, error) {
	fmt.Println("myTestFuncNewStruct. Start.")
	return genTestTree(2, 2), nil
}

func genTestTree(depth, width int) cloneTestStruture {

	out := &cloneTestStruture{
		Now:  time.Now(),
		Text: randSeq(30),
	}

	if depth <= 0 {
		return *out
	}

	depth--

	out.List = make([]cloneTestStruture, width)
	for i := 0; i < width; i++ {
		out.List[i] = genTestTree(depth, width)
	}

	return *out
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
