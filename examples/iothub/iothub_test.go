package iothub

import "testing"

func Test_Main(t *testing.T) {

	go ServerLoop()
	go ClientLoop("12345")
	go ClientLoop("67891")

	select {}
}
