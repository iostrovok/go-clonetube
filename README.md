## Simple module which prepares clones in background and returns by demand ##

### Installing ###
```bash
go get github.com/iostrovok/go-clonetube/clonetube
```
### How use example ###

```go
package main

import (
	"fmt"
	"github.com/iostrovok/go-clonetube/clonetube"
	"runtime"
)

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

	return i.Clone(), nil
}

func main() {

	runtime.GOMAXPROCS(8)

	obj := generateCloneTestStruture(...)
	cl := clonetube.New(10, MyTestFuncClone, 10)
	cl.Put(obj)

	for i := 0; i < 1000; i++ {
		_, err := cl.Get(10000)
		if err != nil {
			fmt.Printf("--> %d: %s\n", i, err)
		}
		/*
		    Do something here with your CloneTestStruture copy
		*/
	}
}

```

### Using ###
#### Import ####
```go

import "github.com/iostrovok/go-clonetube/clonetube"

```
### Interface ###

####  type CloneFn ####

```go

type CloneFn func(interface{}) (interface{}, error)

```
Type of function which is using for cloning

#### func New(l int, f CloneFn, threadNumber ...int) *Main ####
##### Params #####
- l - number of cached params in the storage
- f - cloning function (type CloneFn)
- threadNumber - number of clone threads which execute clone function

##### Returning #####
Ref to new object

##### Action #####
Creates object

```go

cloneTube := clonetube.New(100)

```

#### func (cloneTube *Main) Put(in interface{}) error ####
##### Params #####
- interface{}

##### Returning #####
Error after trying of first cloning object.

##### Action #####
Set or replace object which is cloning

Notece!
We have short time term between "Put(new_object)" and "Get() as new_object":
some of first requests of Get() can return old object.

It depends on clone function speed and number of cached params in the storage.

```go

s := map[string]bool{
    "a": true,
    "b": falsw,
}

if err := cloneTube.Put(s); err != nil {
    // do error
}

```

#### func (cloneTube *Main) Get(timeOutMicrosecond ...int) (interface{}, error) ####
##### Params #####
- timeout in microseconds / optional

##### Returning #####
- interface{} -  clone of object
- error

##### Action #####
Retruns a clone object.

If program often uses the Get(...) method and there are no prepared clones:

if  timeout doesn't set:

	we are waiting while clone function makes objects.
	
else:

	we get a clone or "timeout" error.

```go
/*
s := map[string]int{
    "a": 1,
    "b": 2,
}
*/

copyS, err := cloneTube.Get(10000);
if err != nil {
    // do error
}
copy, good := copyS(map[string]int)
if !good {
    // do error
}



```
