package iothub

import (
	_ "github.com/AnyISalIn/yrpc/shared"
	shared "github.com/AnyISalIn/yrpc/shared"
)

type Device struct {
	id    string
	state bool
}

const (
	DEVICE_OFF      = "Device.Off"
	DEVICE_ON       = "Device.On"
	DEVICE_STATE    = "Device.State"
	DEVICE_CATEGORY = "Device.Category"
	DEVICE_ID       = "Device.ID"
)

func (s *Device) Off(args *shared.Empty, reply *shared.Empty) error {
	s.state = false
	return nil
}

func (s *Device) On(args *shared.Empty, reply *shared.Empty) error {
	s.state = true
	return nil
}

func (s *Device) State(args *shared.Empty, reply *bool) error {
	reply = &s.state
	return nil
}
func (s *Device) Category(args *shared.Empty, reply *string) error {
	*reply = "switch"
	return nil
}

func (s *Device) ID(args *shared.Empty, reply *string) error {
	*reply = s.id
	return nil
}
