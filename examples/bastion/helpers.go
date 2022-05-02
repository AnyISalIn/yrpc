package bastion

import (
	"io"
	"sync"
)

func Bridge(a, b io.ReadWriteCloser) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, _ = io.Copy(a, b)
		a.Close()
		b.Close()
	}()
	go func() {
		defer wg.Done()
		_, _ = io.Copy(b, a)
		a.Close()
		b.Close()
	}()
	wg.Wait()
}
