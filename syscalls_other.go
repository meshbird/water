// +build !linux,!darwin

package water

import (
	"errors"
)

func newTAP(ifName string) (*Interface, error) {
	return nil, errors.New("water: tap interface not implemented on this platform")
}

func newTUN(ifName string) (*Interface, error) {
	return nil, errors.New("water: tap interface not implemented on this platform")
}

func setPersistent(fd uintptr, persistent bool) error {
	return errors.New("water: setPersistent not implemented on this platform")
}
