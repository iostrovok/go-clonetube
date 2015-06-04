package clonetube

/*
	newThreadi is simple func.
	Clone interface{} and result send into channel.
	If we have any incoming message we get out from function.
*/
func newThread(id int64, body interface{}, f CloneFn, ChIn, ChOut chan *item) {
	// infinity loop
	i := 0
	for {

		i++

		b, err := f(body)
		if err != nil {
			return
		}

		// If we have good clone copy
		it := &item{
			id:   id,
			body: b,
			i:    i,
		}

		select {
		case _, ok := <-ChIn:
			// go out by any messages
			if ok || !ok {
				return
			}
		case ChOut <- it:
			// Nothiing
		}
	}
}
