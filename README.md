## Simple module which prepares clones in background and returns by demand ##

### Installing ###
```bash
go get github.com/iostrovok/go-clonetube/clonetube
```
### How use example (not add) ###


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

#### func New(l int, f CloneFn) *Main ####
##### Params #####
- number of cached params in the storage
- cloning function (type CloneFn)

##### Returning #####
Ref to new object

##### Action #####
Creates object

```go

cloneTube := clonetube.New(100)

```

#### func (cloneTube *Main) Start() *Main ####
##### Params #####
None
##### Returning #####
Ref to object

##### Action #####
Starts cloning process
```go

cloneTube := clonetube.New(100)
cloneTube.Start()

// OR
cloneTube := clonetube.New(100).Start()

```

#### func Put ####
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

#### func (cloneTube *Main) Get(timeOutMicrosecond int) (interface{}, error) ####
##### Params #####
- timeout in microseconds

##### Returning #####
- interface{} -  clone of object
- error

##### Action #####
Set or replace object which is cloning

Notece!
We have short time term between "Put(new_object)" and "Get() as new_object":
some of first requests of Get() can return old object.

It depends on clone function speed and number of cached params in the storage.

```go
/*
s := map[string]bool{
    "a": true,
    "b": falsw,
}
*/

copyS, err := cloneTube.Get(10000);
if err != nil {
    // do error
}
copy, good := copyS(map[string]bool)
if !good {
    // do error
}



```
