package tool

func ClearChanAfterClosed(ch chan struct{}){
	for {
		if _, ok := <-ch; !ok{
			break
		}
	}
}
