package di8

import (
	"syscall"
	"unsafe"
)

var (
	dll                = syscall.NewLazyDLL("dinput8.dll")
	directInput8Create = dll.NewProc("DirectInput8Create")
)

func Create(windowInstance HINSTANCE) (*DirectInput, error) {
	var obj *DirectInput
	ret, _, _ := directInput8Create.Call(
		uintptr(windowInstance),
		VERSION,
		uintptr(unsafe.Pointer(&IID_IDirectInput8W)),
		uintptr(unsafe.Pointer(&obj)),
		0,
	)
	return obj, toErr(ret)
}

type DirectInput struct {
	vtbl *directInputVtbl
}

type directInputVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	CreateDevice    uintptr
	EnumDevices     uintptr
	GetDeviceStatus uintptr
	RunControlPanel uintptr
	Initialize      uintptr
}

// AddRef increments the reference count for an interface on an object. This
// method should be called for every new copy of a pointer to an interface on an
// object.
func (obj *DirectInput) AddRef() uint32 {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.AddRef,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	return uint32(ret)
}

// Release has to be called when finished using the object to free its
// associated resources.
func (obj *DirectInput) Release() uint32 {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	return uint32(ret)
}

func (obj *DirectInput) CreateDevice(guid GUID) (device *Device, err error) {
	ret, _, _ := syscall.Syscall6(
		obj.vtbl.CreateDevice,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&guid)),
		uintptr(unsafe.Pointer(&device)),
		0,
		0,
		0,
	)
	err = toErr(ret)
	return
}

// devType: DEVCLASS_* or DEVTYPE_*
// flags: EDFL_*
func (obj *DirectInput) EnumDevices(
	devType uint32,
	callback func(*DEVICEINSTANCE, uintptr) uintptr,
	ref uintptr,
	flags uint32) (
	err error,
) {
	ret, _, _ := syscall.Syscall6(
		obj.vtbl.EnumDevices,
		5,
		uintptr(unsafe.Pointer(obj)),
		uintptr(devType),
		syscall.NewCallback(callback),
		ref,
		uintptr(flags),
		0,
	)
	err = toErr(ret)
	return
}

//func (obj DirectInput) FindDevice(
//	guid GUID,
//	name string,
//) (
//	guidDevice GUID,
//	err error,
//) {
//	cGuid := guid.toC()
//	cName := C.CString(name)
//	defer C.free(unsafe.Pointer(cName))
//	var cGuidDevice C.GUID
//	err = toError(C.IDirectInput8FindDevice(
//		obj.handle,
//		&cGuid,
//		(*C.CHAR)(cName),
//		&cGuidDevice,
//	))
//	guidDevice.fromC(&cGuidDevice)
//	return
//}
//
//func (obj DirectInput) GetDeviceStatus(guid GUID) (err error) {
//	cGuid := guid.toC()
//	err = toError(C.IDirectInput8GetDeviceStatus(obj.handle, &cGuid))
//	return
//}
//
//func (obj DirectInput) RunControlPanel(
//	ownerWindow unsafe.Pointer,
//	flags uint32,
//) (
//	err error,
//) {
//	err = toError(C.IDirectInput8RunControlPanel(
//		obj.handle,
//		C.HWND(ownerWindow),
//		C.DWORD(flags),
//	))
//	return
//}
