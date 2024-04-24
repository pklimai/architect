package closer

var globalCloser = New()

func Add(f func() error) {
	globalCloser.Add(f)
}

func Wait() {
	globalCloser.Wait()
}

func CloseAll() {
	globalCloser.CloseAll()
}
