// +build darwin

package water

/*
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/inercia/kernctl"
)

const utunControlName = "com.apple.net.utun_control"
const utunOptIfName = 2

var errorNotDeviceFound = errors.New("could not find valid tun/tap device")

// Create a new TAP interface whose name is ifName.
// If ifName is empty, a default name (tap0, tap1, ... ) will be assigned.
// ifName should not exceed 16 bytes.
func newTAP(ifName string) (ifce *Interface, err error) {
	name, file, err := createInterface(ifName)
	if err != nil {
		return nil, err
	}
	ifce = &Interface{isTAP: true, file: file, name: name}
	return
}

// Create a new TUN interface whose name is ifName.
// If ifName is empty, a default name (tap0, tap1, ... ) will be assigned.
// ifName should not exceed 16 bytes.
func newTUN(ifName string) (ifce *Interface, err error) {
	name, file, err := createInterface(ifName)
	if err != nil {
		return nil, err
	}
	ifce = &Interface{isTAP: false, file: file, name: name}
	return
}

func createInterface(ifName string) (createdIFName string, file *os.File, err error) {
	file = nil
	err = errorNotDeviceFound

	var readBufLen C.int = 20
	var readBuf = C.CString("                    ")
	defer C.free(unsafe.Pointer(readBuf))

	for utunnum := 0; utunnum < 255; utunnum++ {
		conn := kernctl.NewConnByName(utunControlName)
		conn.UnitId = uint32(utunnum + 1)
		conn.Connect()

		_, _, gserr := syscall.Syscall6(syscall.SYS_GETSOCKOPT,
			uintptr(conn.Fd),
			uintptr(kernctl.SYSPROTO_CONTROL), uintptr(utunOptIfName),
			uintptr(unsafe.Pointer(readBuf)), uintptr(unsafe.Pointer(&readBufLen)), 0)
		if gserr != 0 {
			continue
		} else {
			createdIFName := C.GoStringN(readBuf, C.int(readBufLen))

			fmt.Printf("Try num: %d\n", utunnum)
			fmt.Printf("Fd: %d\n", conn.Fd)
			fmt.Printf("Dev name: %s [%d]\n", createdIFName, readBufLen)

			file = os.NewFile(uintptr(conn.Fd), createdIFName)
			err = nil
			break
		}
	}

	return createdIFName, file, err
}

func setPersistent(fd uintptr, persistent bool) error {
	return errors.New("setPersistent not defined on OS X")
}
