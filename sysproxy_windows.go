package proxy

import (
	"runtime"
	"syscall"
	"unsafe"
)

// #include <stdlib.h>
import "C"

type internetOption int

var (
	internettOptionRefresh              internetOption = 37
	internettOptionPerConnectionOption  internetOption = 75
	internettOptionSettingsChanged      internetOption = 39
	internettOptionProxySettingsChanged internetOption = 95
)

type internetPerConnOption int

var (
	internettPerConnOptionFlags         internetPerConnOption = 1
	internettPerConnOptionProxyServer   internetPerConnOption = 2
	internettPerConnOptionProxyBypass   internetPerConnOption = 3
	internettPerConnOptionAutoConfigUrl internetPerConnOption = 4
)

type internettPerConnFlags int

var (
	internettPerConnFlagsDirect       internettPerConnFlags = 1
	internettPerConnFlagsProxy        internettPerConnFlags = 2
	internettPerConnFlagsAutoProxyUrl internettPerConnFlags = 4
	internettPerConnFlagsAutoDetect   internettPerConnFlags = 8
)

type internetPerConnOptionList struct {
	dwSize        uint32
	pszConnection uintptr
	dwOptionCount uint32
	dwOptionError uint32
	pOptions      uintptr
}

type internetConnOption struct {
	dwOption uint32
	dwValue  uintptr // DWORD | LPSTR
}

func reset() error {
	op := internetConnOption{
		dwOption: uint32(internettPerConnOptionFlags),
		dwValue:  uintptr(internettPerConnFlagsDirect),
	}

	// here must pin memory of op, prevent it from being moved by GC.
	var pin runtime.Pinner
	pin.Pin(&op)
	defer pin.Unpin()

	opl := internetPerConnOptionList{
		pszConnection: 0,
		dwOptionCount: 1,
		dwOptionError: 0,
		pOptions:      uintptr(unsafe.Pointer(&op)),
	}
	opl.dwSize = uint32(unsafe.Sizeof(&opl))

	return setSystemProxy(&opl)
}

func setGlobal(proxy string, bypass string) error {
	op0 := internetConnOption{
		dwOption: uint32(internettPerConnOptionFlags),
		dwValue:  uintptr(internettPerConnFlagsProxy | internettPerConnFlagsDirect),
	}

	cProxy := unsafe.Pointer(C.CString(proxy))
	defer C.free(cProxy)
	cBypass := unsafe.Pointer(C.CString(bypass))
	defer C.free(cBypass)

	op1 := internetConnOption{
		dwOption: uint32(internettPerConnOptionProxyServer),
		dwValue:  uintptr(cProxy), // cstring is allocated in C memory, it is safe
	}

	op2 := internetConnOption{
		dwOption: uint32(internettPerConnOptionProxyBypass),
		dwValue:  uintptr(cBypass),
	}

	ops := [3]internetConnOption{op0, op1, op2}

	var pin runtime.Pinner
	pin.Pin(&ops)
	defer pin.Unpin()

	opl := internetPerConnOptionList{
		pszConnection: 0,
		dwOptionCount: 3,
		dwOptionError: 0,
		pOptions:      uintptr(unsafe.Pointer(&ops)),
	}
	opl.dwSize = uint32(unsafe.Sizeof(&opl))

	return setSystemProxy(&opl)
}

func setPac(proxyUrl string) error {
	op0 := internetConnOption{
		dwOption: uint32(internettPerConnOptionFlags),
		dwValue:  uintptr(internettPerConnFlagsAutoProxyUrl | internettPerConnFlagsDirect),
	}

	cProxyUrl := unsafe.Pointer(C.CString(proxyUrl))
	defer C.free(cProxyUrl)

	op1 := internetConnOption{
		dwOption: uint32(internettPerConnOptionAutoConfigUrl),
		dwValue:  uintptr(cProxyUrl),
	}

	ops := [3]internetConnOption{op0, op1}

	var pin runtime.Pinner
	pin.Pin(&ops)
	defer pin.Unpin()

	opl := internetPerConnOptionList{
		pszConnection: 0,
		dwOptionCount: 2,
		dwOptionError: 0,
		pOptions:      uintptr(unsafe.Pointer(&ops)),
	}
	opl.dwSize = uint32(unsafe.Sizeof(&opl))

	return setSystemProxy(&opl)
}

func setSystemProxy(opl *internetPerConnOptionList) error {
	lib := syscall.MustLoadDLL("wininet.dll")
	defer lib.Release()
	fn := lib.MustFindProc("InternetSetOptionA")
	// converting pointer to uintptr is safe when syscall. referenced object is retained and not moved until the syscall completes.
	// see https://pkg.go.dev/unsafe#Pointer
	_, _, err := syscall.SyscallN(fn.Addr(), 0, uintptr(internettOptionPerConnectionOption), uintptr(unsafe.Pointer(opl)), unsafe.Sizeof(*opl))
	if err != 0 {
		return err
	}
	_, _, err = syscall.SyscallN(fn.Addr(), 0, uintptr(internettOptionProxySettingsChanged), 0, 0)
	if err != 0 {
		return err
	}
	_, _, err = syscall.SyscallN(fn.Addr(), 0, uintptr(internettOptionRefresh), 0, 0)
	if err != 0 {
		return err
	}
	return nil
}
