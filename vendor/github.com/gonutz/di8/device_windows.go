package di8

import (
	"syscall"
	"unsafe"
)

type Device struct {
	vtbl *deviceVtbl
}

type deviceVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetCapabilities      uintptr
	EnumObjects          uintptr
	GetProperty          uintptr
	SetProperty          uintptr
	Acquire              uintptr
	Unacquire            uintptr
	GetDeviceState       uintptr
	GetDeviceData        uintptr
	SetDataFormat        uintptr
	SetEventNotification uintptr
	SetCooperativeLevel  uintptr
	GetObjectInfo        uintptr
	GetDeviceInfo        uintptr
	RunControlPanel      uintptr
	Initialize           uintptr
}

// AddRef increments the reference count for an interface on an object. This
// method should be called for every new copy of a pointer to an interface on an
// object.
func (obj *Device) AddRef() uint32 {
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
func (obj *Device) Release() uint32 {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	return uint32(ret)
}

func (obj *Device) Acquire() Error {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Acquire,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	return toErr(ret)
}

func (obj *Device) Unacquire() Error {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Unacquire,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	return toErr(ret)
}

func (obj *Device) SetCooperativeLevel(window HWND, flags uint32) Error {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.SetCooperativeLevel,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(window),
		uintptr(flags),
	)
	return toErr(ret)
}

func (obj *Device) SetDataFormat(format *DATAFORMAT) Error {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.SetDataFormat,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(format)),
		0,
	)
	return toErr(ret)
}

var deviceobjectdata DEVICEOBJECTDATA

func (obj *Device) GetDeviceData(data []DEVICEOBJECTDATA, flags uint32) (int, Error) {
	var dataPtr uintptr
	if len(data) > 0 {
		dataPtr = uintptr(unsafe.Pointer(&data[0]))
	}
	count := uint32(len(data))
	ret, _, _ := syscall.Syscall6(
		obj.vtbl.GetDeviceData,
		5,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Sizeof(deviceobjectdata)),
		dataPtr,
		uintptr(unsafe.Pointer(&count)),
		uintptr(flags),
		0,
	)
	return int(count), toErr(ret)
}

type DeviceState interface {
	Ptr() uintptr
	Size() int
}

func (obj *Device) GetDeviceState(state DeviceState) Error {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.GetDeviceState,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(state.Size()),
		state.Ptr(),
	)
	return toErr(ret)
}

type Property interface {
	PropHeader() *PROPHEADER
}

// SetProperty sets one of the PROP_* properties for the device. Predefined
// property types are: PROPCPOINTS, PROPDWORD, PROPRANGE, PROPCAL, PROPCALPOV,
// PROPGUIDANDPATH, PROPSTRING and PROPPOINTER. Create them with the NewProp*
// functions.
func (obj *Device) SetProperty(guid *GUID, prop Property) Error {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.SetProperty,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(guid)),
		uintptr(unsafe.Pointer(prop.PropHeader())),
	)
	return toErr(ret)
}

///*
//#include "dinput_wrapper.h"
//
//HRESULT IDirectInputDevice8Acquire(IDirectInputDevice8* obj) {
//	return obj->lpVtbl->Acquire(obj);
//}
//
//HRESULT IDirectInputDevice8BuildActionMap(
//		IDirectInputDevice8* obj,
//		LPDIACTIONFORMAT lpdiaf,
//		LPCTSTR lpszUserName,
//		DWORD dwFlags) {
//	return obj->lpVtbl->BuildActionMap(obj, lpdiaf, lpszUserName, dwFlags);
//}
//
//BOOL enumEffectsCallbackGo(LPCDIEFFECTINFO, void*);
//
//HRESULT IDirectInputDevice8EnumEffects(
//		IDirectInputDevice8* obj,
//		DWORD dwEffType) {
//	return obj->lpVtbl->EnumEffects(
//		obj,
//		(LPDIENUMEFFECTSCALLBACK)enumEffectsCallbackGo,
//		0,
//		dwEffType);
//}
//
//BOOL enumObjectsCallbackGo(LPCDIDEVICEOBJECTINSTANCE, void*);
//
//HRESULT IDirectInputDevice8EnumObjects(IDirectInputDevice8* obj, DWORD flags) {
//	return obj->lpVtbl->EnumObjects(
//		obj,
//		(LPDIENUMDEVICEOBJECTSCALLBACK)enumObjectsCallbackGo,
//		0,
//		flags);
//}
//
//HRESULT IDirectInputDevice8GetDeviceState(
//		IDirectInputDevice8* obj,
//		DWORD cbData,
//		void* lpvData) {
//	return obj->lpVtbl->GetDeviceState(obj, cbData, lpvData);
//}
//
//HRESULT IDirectInputDevice8SetCooperativeLevel(
//		IDirectInputDevice8* obj,
//		HWND hwnd,
//		DWORD dwFlags) {
//	return obj->lpVtbl->SetCooperativeLevel(obj, hwnd, dwFlags);
//}
//
//HRESULT IDirectInputDevice8SetDataFormat(
//		IDirectInputDevice8* obj,
//		DIDATAFORMAT* lpdf) {
//	return obj->lpVtbl->SetDataFormat(obj, lpdf);
//}
//
//HRESULT IDirectInputDevice8GetDeviceData(
//		IDirectInputDevice8* obj,
//		DWORD cbObjectData,
//		LPDIDEVICEOBJECTDATA rgdod,
//		LPDWORD pdwInOut,
//		DWORD dwFlags) {
//	return obj->lpVtbl->GetDeviceData(obj, cbObjectData, rgdod, pdwInOut, dwFlags);
//}
//
//HRESULT IDirectInputDevice8GetProperty(
//		IDirectInputDevice8* obj,
//		REFGUID rguidProp,
//		LPDIPROPHEADER pdiph){
//	return obj->lpVtbl->GetProperty(obj, rguidProp, pdiph);
//}
//
//HRESULT IDirectInputDevice8SetProperty(
//		IDirectInputDevice8* obj,
//		REFGUID rguidProp,
//		LPCDIPROPHEADER pdiph){
//	return obj->lpVtbl->SetProperty(obj, rguidProp, pdiph);
//}
//
//HRESULT IDirectInputDevice8GetPredefinedProperty(
//		IDirectInputDevice8* obj,
//		void* rguidProp,
//		LPDIPROPHEADER pdiph){
//	return obj->lpVtbl->GetProperty(obj, (GUID*)rguidProp, pdiph);
//}
//
//HRESULT IDirectInputDevice8SetPredefinedProperty(
//		IDirectInputDevice8* obj,
//		void* rguidProp,
//		LPCDIPROPHEADER pdiph){
//	return obj->lpVtbl->SetProperty(obj, (GUID*)rguidProp, pdiph);
//}
//
//HRESULT IDirectInputDevice8Unacquire(IDirectInputDevice8* obj) {
//	return obj->lpVtbl->Unacquire(obj);
//}
//
//void IDirectInputDevice8Release(IDirectInputDevice8* obj) {
//	obj->lpVtbl->Release(obj);
//}
//*/
//import "C"
//import (
//	"errors"
//	"strconv"
//	"syscall"
//	"unsafe"
//)
//
//var (
//	GUID_SysKeyboard GUID
//	GUID_SysMouse    GUID
//
//	GUID_ConstantForce GUID
//	GUID_RampForce     GUID
//	GUID_Square        GUID
//	GUID_Sine          GUID
//	GUID_Triangle      GUID
//	GUID_SawtoothUp    GUID
//	GUID_SawtoothDown  GUID
//	GUID_Spring        GUID
//	GUID_Damper        GUID
//	GUID_Inertia       GUID
//	GUID_Friction      GUID
//	GUID_CustomForce   GUID
//)
//
//func init() {
//	GUID_SysKeyboard.fromC(&C.GUID_SysKeyboard)
//	GUID_SysMouse.fromC(&C.GUID_SysMouse)
//
//	GUID_ConstantForce.fromC(&C.GUID_ConstantForce)
//	GUID_RampForce.fromC(&C.GUID_RampForce)
//	GUID_Square.fromC(&C.GUID_Square)
//	GUID_Sine.fromC(&C.GUID_Sine)
//	GUID_Triangle.fromC(&C.GUID_Triangle)
//	GUID_SawtoothUp.fromC(&C.GUID_SawtoothUp)
//	GUID_SawtoothDown.fromC(&C.GUID_SawtoothDown)
//	GUID_Spring.fromC(&C.GUID_Spring)
//	GUID_Damper.fromC(&C.GUID_Damper)
//	GUID_Inertia.fromC(&C.GUID_Inertia)
//	GUID_Friction.fromC(&C.GUID_Friction)
//	GUID_CustomForce.fromC(&C.GUID_CustomForce)
//}
//
//type Device struct {
//	handle *C.IDirectInputDevice8
//}
//
//func (obj Device) Release() {
//	C.IDirectInputDevice8Release(obj.handle)
//}
//
//func (obj Device) Acquire() (err error) {
//	err = toError(C.IDirectInputDevice8Acquire(obj.handle))
//	return
//}
//
////func (obj Device) BuildActionMap(userName string, flags uint32) (actionFormat ACTIONFORMAT, err error) {
////	var cActionFormat C.DIACTIONFORMAT
////	// TODO does cActionFormat.dwSize need to be set?
////	if len(userName) == 0 {
////		err = toError(C.IDirectInputDevice8BuildActionMap(obj.handle, &cActionFormat, nil, C.DWORD(flags)))
////	} else {
////		cUserName := C.CString(userName)
////		defer C.free(unsafe.Pointer(cUserName))
////		err = toError(C.IDirectInputDevice8BuildActionMap(obj.handle, &cActionFormat, (*C.CHAR)(cUserName), C.DWORD(flags)))
////	}
////	actionFormat.fromC(&cActionFormat)
////	return
////}
//
////func (obj Device) CreateEffect() (err error) {
////	err = toError(C.IDirectInputDevice8CreateEffect(obj.handle))
////  return
////}
//
////func (obj Device) EnumCreatedEffectObjects() (err error) {
////	err = toError(C.IDirectInputDevice8EnumCreatedEffectObjects(obj.handle))
////  return
////}
//
////func (obj Device) EnumEffects() (err error) {
////	err = toError(C.IDirectInputDevice8EnumEffects(obj.handle))
////  return
////}
//
////func (obj Device) EnumEffectsInFile() (err error) {
////	err = toError(C.IDirectInputDevice8EnumEffectsInFile(obj.handle))
////  return
////}
//
//func (obj Device) EnumObjects(
//	callback EnumObjectsCallback,
//	flags uint32,
//) (
//	err error,
//) {
//	currentEnumObjectsCallback = callback
//	err = toError(C.IDirectInputDevice8EnumObjects(obj.handle, C.DWORD(flags)))
//	return
//}
//
//func (obj Device) EnumEffects(
//	callback EnumEffectsCallback,
//	effectType uint32,
//) (
//	err error,
//) {
//	currentEnumEffectsCallback = callback
//	err = toError(C.IDirectInputDevice8EnumEffects(
//		obj.handle,
//		C.DWORD(effectType),
//	))
//	return
//}
//
////func (obj Device) Escape() (err error) {
////	err = toError(C.IDirectInputDevice8Escape(obj.handle))
////  return
////}
//
////func (obj Device) GetCapabilities() (err error) {
////	err = toError(C.IDirectInputDevice8GetCapabilities(obj.handle))
////  return
////}
//
//func (obj Device) GetDeviceData(bufferSize int) (data []DEVICEOBJECTDATA, err error) {
//	cData := make([]C.DIDEVICEOBJECTDATA, bufferSize)
//	objectCount := C.DWORD(bufferSize)
//
//	err = toError(C.IDirectInputDevice8GetDeviceData(
//		obj.handle,
//		C.sizeof_DIDEVICEOBJECTDATA,
//		&cData[0],
//		&objectCount,
//		0,
//	))
//
//	if err != nil {
//		return
//	}
//
//	data = make([]DEVICEOBJECTDATA, objectCount)
//	for i := range data {
//		data[i].fromC(&cData[i])
//	}
//
//	return
//}
//
////func (obj Device) GetDeviceInfo() (err error) {
////	err = toError(C.IDirectInputDevice8GetDeviceInfo(obj.handle))
////  return
////}
//
////func (obj Device) GetDeviceState() (err error) {
////	err = toError(C.IDirectInputDevice8GetDeviceState(obj.handle))
////	return
////}
//
//func (obj Device) GetKeyboardState(state *KeyboardState) (err error) {
//	err = toError(C.IDirectInputDevice8GetDeviceState(
//		obj.handle,
//		256,
//		unsafe.Pointer(&state[0]),
//	))
//	return
//}
//
//type KeyboardState [256]byte
//
//func (k *KeyboardState) IsDown(key int) bool {
//	if key < 0 || key >= 256 {
//		return false
//	}
//	return (*k)[key]&0x80 != 0
//}
//
//func (obj Device) GetMouseState(state *MouseState) (err error) {
//	err = toError(C.IDirectInputDevice8GetDeviceState(
//		obj.handle,
//		C.sizeof_DIMOUSESTATE,
//		unsafe.Pointer(&state.state),
//	))
//	return
//}
//
//type MouseState struct {
//	state C.DIMOUSESTATE
//}
//
//func (m *MouseState) X() int {
//	return int(m.state.lX)
//}
//
//func (m *MouseState) Y() int {
//	return int(m.state.lY)
//}
//
//func (m *MouseState) Z() int {
//	return int(m.state.lZ)
//}
//
//func (m *MouseState) IsDown(mouseButton int) bool {
//	if mouseButton < 0 || mouseButton > 3 {
//		return false
//	}
//	return m.state.rgbButtons[mouseButton]&0x80 != 0
//}
//
//const (
//	MouseButtonLeft  = 0
//	MouseButtonRight = 1
//	MouseButtonMiddl = 2
//)
//
//func (obj Device) GetMouseState2(state *MouseState2) (err error) {
//	err = toError(C.IDirectInputDevice8GetDeviceState(
//		obj.handle,
//		C.sizeof_DIMOUSESTATE2,
//		unsafe.Pointer(&state.state),
//	))
//	return
//}
//
//type MouseState2 struct {
//	state C.DIMOUSESTATE2
//}
//
//func (m *MouseState2) X() int {
//	return int(m.state.lX)
//}
//
//func (m *MouseState2) Y() int {
//	return int(m.state.lY)
//}
//
//func (m *MouseState2) Z() int {
//	return int(m.state.lZ)
//}
//
//func (m *MouseState2) IsDown(mouseButton int) bool {
//	if mouseButton < 0 || mouseButton >= 8 {
//		return false
//	}
//	return m.state.rgbButtons[mouseButton]&0x80 != 0
//}
//
//func (obj Device) GetJoyState(state *JoyState) (err error) {
//	err = toError(C.IDirectInputDevice8GetDeviceState(
//		obj.handle,
//		C.sizeof_DIJOYSTATE,
//		unsafe.Pointer(&state.state),
//	))
//	return
//}
//
//type JoyState struct {
//	state C.DIJOYSTATE
//}
//
//func (j *JoyState) X() int {
//	return int(j.state.lX)
//}
//
//func (j *JoyState) Y() int {
//	return int(j.state.lY)
//}
//
//func (j *JoyState) Z() int {
//	return int(j.state.lZ)
//}
//
//func (j *JoyState) Rx() int {
//	return int(j.state.lRx)
//}
//
//func (j *JoyState) Ry() int {
//	return int(j.state.lRy)
//}
//
//func (j *JoyState) Rz() int {
//	return int(j.state.lRz)
//}
//
//func (j *JoyState) Slider(index int) int {
//	if index < 0 || index >= 2 {
//		return 0
//	}
//	return int(j.state.rglSlider[index])
//}
//
//// POV returns the direction in which the position control is currently pushed,
//// in hundredths of a degree clockwise from north (away from the user). The
//// center position is reported as -1 (or POVCentered).
//func (j *JoyState) POV(pov int) int {
//	if pov < 0 || pov >= 4 {
//		return POVCentered
//	}
//	// centered is either reported as -1 or 65535, this is the recommended way
//	// to check for centered (see MSDN).
//	if j.state.rgdwPOV[pov]&0xFFFF != 0 {
//		return POVCentered
//	}
//	return int(j.state.rgdwPOV[pov])
//}
//
//const (
//	POVCentered = -1
//	POVNorth    = 0
//	POVEast     = 9000
//	POVSouth    = 18000
//	POVWest     = 27000
//)
//
//func (j *JoyState) IsDown(button int) bool {
//	if button < 0 || button >= 32 {
//		return false
//	}
//	return j.state.rgbButtons[button]&0x80 != 0
//}
//
//func (obj Device) GetJoyState2(state *JoyState2) (err error) {
//	err = toError(C.IDirectInputDevice8GetDeviceState(
//		obj.handle,
//		C.sizeof_DIJOYSTATE2,
//		unsafe.Pointer(&state.state),
//	))
//	return
//}
//
//type JoyState2 struct {
//	state C.DIJOYSTATE2
//}
//
//func (j *JoyState2) X() int {
//	return int(j.state.lX)
//}
//
//func (j *JoyState2) Y() int {
//	return int(j.state.lY)
//}
//
//func (j *JoyState2) Z() int {
//	return int(j.state.lZ)
//}
//
//func (j *JoyState2) Rx() int {
//	return int(j.state.lRx)
//}
//
//func (j *JoyState2) Ry() int {
//	return int(j.state.lRy)
//}
//
//func (j *JoyState2) Rz() int {
//	return int(j.state.lRz)
//}
//
//func (j *JoyState2) Vx() int {
//	return int(j.state.lVX)
//}
//
//func (j *JoyState2) Vy() int {
//	return int(j.state.lVY)
//}
//
//func (j *JoyState2) Vz() int {
//	return int(j.state.lVZ)
//}
//
//func (j *JoyState2) VRx() int {
//	return int(j.state.lVRx)
//}
//
//func (j *JoyState2) VRy() int {
//	return int(j.state.lVRy)
//}
//
//func (j *JoyState2) VRz() int {
//	return int(j.state.lVRz)
//}
//
//func (j *JoyState2) Ax() int {
//	return int(j.state.lAX)
//}
//
//func (j *JoyState2) Ay() int {
//	return int(j.state.lAY)
//}
//
//func (j *JoyState2) Az() int {
//	return int(j.state.lAZ)
//}
//
//func (j *JoyState2) ARx() int {
//	return int(j.state.lARx)
//}
//
//func (j *JoyState2) ARy() int {
//	return int(j.state.lARy)
//}
//
//func (j *JoyState2) ARz() int {
//	return int(j.state.lARz)
//}
//
//func (j *JoyState2) Fx() int {
//	return int(j.state.lFX)
//}
//
//func (j *JoyState2) Fy() int {
//	return int(j.state.lFY)
//}
//
//func (j *JoyState2) Fz() int {
//	return int(j.state.lFZ)
//}
//
//func (j *JoyState2) FRx() int {
//	return int(j.state.lFRx)
//}
//
//func (j *JoyState2) FRy() int {
//	return int(j.state.lFRy)
//}
//
//func (j *JoyState2) FRz() int {
//	return int(j.state.lFRz)
//}
//
//func (j *JoyState2) Slider(index int) int {
//	if index < 0 || index >= 2 {
//		return 0
//	}
//	return int(j.state.rglSlider[index])
//}
//
//func (j *JoyState2) VSlider(index int) int {
//	if index < 0 || index >= 2 {
//		return 0
//	}
//	return int(j.state.rglVSlider[index])
//}
//
//func (j *JoyState2) ASlider(index int) int {
//	if index < 0 || index >= 2 {
//		return 0
//	}
//	return int(j.state.rglASlider[index])
//}
//
//func (j *JoyState2) FSlider(index int) int {
//	if index < 0 || index >= 2 {
//		return 0
//	}
//	return int(j.state.rglFSlider[index])
//}
//
//// POV returns the direction in which the position control is currently pushed,
//// in hundredths of a degree clockwise from north (away from the user). The
//// center position is reported as -1 (or POVCentered).
//func (j *JoyState2) POV(pov int) int {
//	if pov < 0 || pov >= 4 {
//		return POVCentered
//	}
//	// centered is either reported as -1 or 65535, this is the recommended way
//	// to check for centered (see MSDN).
//	if j.state.rgdwPOV[pov]&0xFFFF != 0 {
//		return POVCentered
//	}
//	return int(j.state.rgdwPOV[pov])
//}
//
//func (j *JoyState2) IsDown(button int) bool {
//	if button < 0 || button >= 128 {
//		return false
//	}
//	return j.state.rgbButtons[button]&0x80 != 0
//}
//
////func (obj Device) GetEffectInfo() (err error) {
////	err = toError(C.IDirectInputDevice8GetEffectInfo(obj.handle))
////  return
////}
//
////func (obj Device) GetForceFeedbackState() (err error) {
////	err = toError(C.IDirectInputDevice8GetForceFeedbackState(obj.handle))
////  return
////}
//
////func (obj Device) GetImageInfo() (err error) {
////	err = toError(C.IDirectInputDevice8GetImageInfo(obj.handle))
////  return
////}
//
////func (obj Device) GetObjectInfo() (err error) {
////	err = toError(C.IDirectInputDevice8GetObjectInfo(obj.handle))
////  return
////}
//
////func (obj Device) GetProperty() (err error) {
////	err = toError(C.IDirectInputDevice8GetProperty(obj.handle))
////  return
////}
//
////func (obj Device) Initialize() (err error) {
////	err = toError(C.IDirectInputDevice8Initialize(obj.handle))
////  return
////}
//
////func (obj Device) Poll() (err error) {
////	err = toError(C.IDirectInputDevice8Poll(obj.handle))
////  return
////}
//
////func (obj Device) RunControlPanel() (err error) {
////	err = toError(C.IDirectInputDevice8RunControlPanel(obj.handle))
////  return
////}
//
////func (obj Device) SendDeviceData() (err error) {
////	err = toError(C.IDirectInputDevice8SendDeviceData(obj.handle))
////  return
////}
//
////func (obj Device) SendForceFeedbackCommand() (err error) {
////	err = toError(C.IDirectInputDevice8SendForceFeedbackCommand(obj.handle))
////  return
////}
//
////func (obj Device) SetActionMap() (err error) {
////	err = toError(C.IDirectInputDevice8SetActionMap(obj.handle))
////  return
////}
//
//func (obj Device) SetCooperativeLevel(
//	windowHandle unsafe.Pointer,
//	flags uint32,
//) (
//	err error,
//) {
//	err = toError(C.IDirectInputDevice8SetCooperativeLevel(
//		obj.handle,
//		C.HWND(windowHandle),
//		C.DWORD(flags),
//	))
//	return
//}
//
//func (obj Device) SetPredefinedDataFormat(
//	format PredefinedDataFormat,
//) (
//	err error,
//) {
//	var cFormat *C.DIDATAFORMAT
//	switch format {
//	case DataFormatKeyboard:
//		cFormat = &C.c_dfDIKeyboard
//	case DataFormatMouse:
//		cFormat = &C.c_dfDIMouse
//	case DataFormatMouse2:
//		cFormat = &C.c_dfDIMouse2
//	case DataFormatJoystick:
//		cFormat = &C.c_dfDIJoystick
//	case DataFormatJoystick2:
//		cFormat = &C.c_dfDIJoystick2
//	default:
//		return toError(ERR_INVALIDPARAM)
//	}
//	err = toError(C.IDirectInputDevice8SetDataFormat(obj.handle, cFormat))
//	return
//}
//
//type PredefinedDataFormat int
//
//const (
//	DataFormatKeyboard PredefinedDataFormat = iota + 1
//	DataFormatMouse
//	DataFormatMouse2
//	DataFormatJoystick
//	DataFormatJoystick2
//)
//
////func (obj Device) SetDataFormat(format DATAFORMAT) (err error) {
////	cFormat := format.toC()
////	err = toError(C.IDirectInputDevice8SetDataFormat(obj.handle, &cFormat))
////	return
////}
//
////func (obj Device) SetEventNotification() (err error) {
////	err = toError(C.IDirectInputDevice8SetEventNotification(obj.handle))
////  return
////}
//
////func (obj Device) SetProperty(prop GUID) (err error) {
////	err = toError(C.IDirectInputDevice8SetProperty(obj.handle))
////	return
////}
//
//// NOTE GetProperty might return S_FALSE even on success, see MSDN
//
//func (device Device) GetPredefinedDwordProperty(
//	prop, obj, how int,
//) (
//	value uint32,
//	err error,
//) {
//	var cProp C.DIPROPDWORD
//	cProp.diph.dwSize = C.sizeof_DIPROPDWORD
//	cProp.diph.dwHeaderSize = C.sizeof_DIPROPHEADER
//	cProp.diph.dwObj = C.DWORD(obj)
//	cProp.diph.dwHow = C.DWORD(how)
//	err = toGetPropError(C.IDirectInputDevice8GetPredefinedProperty(
//		device.handle,
//		unsafe.Pointer(uintptr(prop)),
//		&cProp.diph,
//	))
//	value = uint32(cProp.dwData)
//	return
//}
//
//func (device Device) GetPredefinedPointerProperty(
//	prop, obj, how int,
//) (
//	value uintptr,
//	err error,
//) {
//	var cProp C.DIPROPPOINTER
//	cProp.diph.dwSize = C.sizeof_DIPROPPOINTER
//	cProp.diph.dwHeaderSize = C.sizeof_DIPROPHEADER
//	cProp.diph.dwObj = C.DWORD(obj)
//	cProp.diph.dwHow = C.DWORD(how)
//	err = toGetPropError(C.IDirectInputDevice8GetPredefinedProperty(
//		device.handle,
//		unsafe.Pointer(uintptr(prop)),
//		&cProp.diph,
//	))
//	value = uintptr(cProp.uData)
//	return
//}
//
//func (device Device) GetPredefinedRangeProperty(
//	prop, obj, how int,
//) (
//	min, max int,
//	err error,
//) {
//	var cProp C.DIPROPRANGE
//	cProp.diph.dwSize = C.sizeof_DIPROPRANGE
//	cProp.diph.dwHeaderSize = C.sizeof_DIPROPHEADER
//	cProp.diph.dwObj = C.DWORD(obj)
//	cProp.diph.dwHow = C.DWORD(how)
//	err = toGetPropError(C.IDirectInputDevice8GetPredefinedProperty(
//		device.handle,
//		unsafe.Pointer(uintptr(prop)),
//		&cProp.diph,
//	))
//	min, max = int(cProp.lMin), int(cProp.lMax)
//	return
//}
//
//func (device Device) GetPredefinedStringProperty(
//	prop, obj, how int,
//) (
//	value string,
//	err error,
//) {
//	var cProp C.DIPROPSTRING
//	cProp.diph.dwSize = C.sizeof_DIPROPSTRING
//	cProp.diph.dwHeaderSize = C.sizeof_DIPROPHEADER
//	cProp.diph.dwObj = C.DWORD(obj)
//	cProp.diph.dwHow = C.DWORD(how)
//	err = toGetPropError(C.IDirectInputDevice8GetPredefinedProperty(
//		device.handle,
//		unsafe.Pointer(uintptr(prop)),
//		&cProp.diph,
//	))
//	var buf [maxPath]uint16
//	length := 0
//	for ; length < maxPath; length++ {
//		buf[length] = uint16(cProp.wsz[length])
//		if cProp.wsz[length] == 0 {
//			break
//		}
//	}
//	value = syscall.UTF16ToString(buf[:length])
//	return
//}
//
//func (device Device) GetPredefinedGuidAndPathProperty(
//	prop, obj, how int,
//) (
//	guid GUID,
//	path string,
//	err error,
//) {
//	var cProp C.DIPROPGUIDANDPATH
//	cProp.diph.dwSize = C.sizeof_DIPROPGUIDANDPATH
//	cProp.diph.dwHeaderSize = C.sizeof_DIPROPHEADER
//	cProp.diph.dwObj = C.DWORD(obj)
//	cProp.diph.dwHow = C.DWORD(how)
//	err = toGetPropError(C.IDirectInputDevice8GetPredefinedProperty(
//		device.handle,
//		unsafe.Pointer(uintptr(prop)),
//		&cProp.diph,
//	))
//	guid.fromC(&cProp.guidClass)
//	var buf [maxPath]uint16
//	length := 0
//	for ; length < maxPath; length++ {
//		buf[length] = uint16(cProp.wszPath[length])
//		if cProp.wszPath[length] == 0 {
//			break
//		}
//	}
//	path = syscall.UTF16ToString(buf[:length])
//	return
//}
//
//func toGetPropError(hr C.HRESULT) error {
//	if hr == C.S_FALSE {
//		return nil
//	}
//	return toError(hr)
//}
//
//func (obj Device) SetPredefinedProperty(
//	prop int,
//	value propHeader) (
//	err error,
//) {
//	if prop < 1 || prop > 26 {
//		return errors.New(strconv.Itoa(prop) + " is not a predefined property")
//	}
//	if value == nil {
//		return toError(ERR_INVALIDPARAM)
//	}
//	err = toError(C.IDirectInputDevice8SetPredefinedProperty(
//		obj.handle,
//		unsafe.Pointer(uintptr(prop)),
//		value.headerAddress(),
//	))
//	return
//}
//
//type propHeader interface {
//	headerAddress() *C.DIPROPHEADER
//}
//
//// NOTE DIPROPCPOINTS is deprecated, so there is no equivalent for it.
//
//func NewPropPointer(obj, how int, data uintptr) propHeader {
//	var p C.DIPROPPOINTER
//	p.diph.dwSize = C.sizeof_DIPROPPOINTER
//	p.diph.dwHeaderSize = C.sizeof_DIPROPHEADER
//	p.diph.dwObj = C.DWORD(obj)
//	p.diph.dwHow = C.DWORD(how)
//	p.uData = C.UINT_PTR(data)
//	return &propPointer{p}
//}
//
//type propPointer struct {
//	data C.DIPROPPOINTER
//}
//
//func (p *propPointer) headerAddress() *C.DIPROPHEADER {
//	return &p.data.diph
//}
//
//func NewPropDword(obj, how int, data uint32) propHeader {
//	var p C.DIPROPDWORD
//	p.diph.dwSize = C.sizeof_DIPROPDWORD
//	p.diph.dwHeaderSize = C.sizeof_DIPROPHEADER
//	p.diph.dwObj = C.DWORD(obj)
//	p.diph.dwHow = C.DWORD(how)
//	p.dwData = C.DWORD(data)
//	return &propDword{p}
//}
//
//type propDword struct {
//	data C.DIPROPDWORD
//}
//
//func (p *propDword) headerAddress() *C.DIPROPHEADER {
//	return &p.data.diph
//}
//
//func NewPropString(obj, how int, data string) propHeader {
//	var p C.DIPROPSTRING
//	p.diph.dwSize = C.sizeof_DIPROPSTRING
//	p.diph.dwHeaderSize = C.sizeof_DIPROPHEADER
//	p.diph.dwObj = C.DWORD(obj)
//	p.diph.dwHow = C.DWORD(how)
//	str := syscall.StringToUTF16(data)
//	if len(str) > maxPath-1 {
//		str = str[:maxPath-1]
//	}
//	for i, r := range str {
//		p.wsz[i] = C.WCHAR(r)
//	}
//	p.wsz[len(str)] = 0
//	return &propString{p}
//}
//
//type propString struct {
//	data C.DIPROPSTRING
//}
//
//func (p *propString) headerAddress() *C.DIPROPHEADER {
//	return &p.data.diph
//}
//
//func NewPropRange(obj, how, min, max int) propHeader {
//	var p C.DIPROPRANGE
//	p.diph.dwSize = C.sizeof_DIPROPRANGE
//	p.diph.dwHeaderSize = C.sizeof_DIPROPHEADER
//	p.diph.dwObj = C.DWORD(obj)
//	p.diph.dwHow = C.DWORD(how)
//	p.lMin = C.LONG(min)
//	p.lMax = C.LONG(max)
//	return &propRange{p}
//}
//
//type propRange struct {
//	data C.DIPROPRANGE
//}
//
//func (p *propRange) headerAddress() *C.DIPROPHEADER {
//	return &p.data.diph
//}
//
//const (
//	PROP_BUFFERSIZE         = 1
//	PROP_AXISMODE           = 2
//	PROP_GRANULARITY        = 3
//	PROP_RANGE              = 4
//	PROP_DEADZONE           = 5
//	PROP_SATURATION         = 6
//	PROP_FFGAIN             = 7
//	PROP_FFLOAD             = 8
//	PROP_AUTOCENTER         = 9
//	PROP_CALIBRATIONMODE    = 10
//	PROP_CALIBRATION        = 11
//	PROP_GUIDANDPATH        = 12
//	PROP_INSTANCENAME       = 13
//	PROP_PRODUCTNAME        = 14
//	PROP_JOYSTICKID         = 15
//	PROP_GETPORTDISPLAYNAME = 16
//	PROP_PHYSICALRANGE      = 18
//	PROP_LOGICALRANGE       = 19
//	PROP_KEYNAME            = 20
//	PROP_CPOINTS            = 21
//	PROP_APPDATA            = 22
//	PROP_SCANCODE           = 23
//	PROP_VIDPID             = 24
//	PROP_USERNAME           = 25
//	PROP_TYPENAME           = 26
//)
//
////func (obj Device) SetPredefinedProperty(prop Property) (err error) {
////	if prop < 0 || prop >= len(predefinedProperties) {
////		// TODO describe the error here!
////		return toError(ERR_GENERIC)
////	}
////	err = toError(C.IDirectInputDevice8SetProperty(obj.handle,
////		predefinedProperties[prop]))
////	return
////}
//
////type Property int
//
////const (
////	PROP_APPDATA Property = iota
////	PROP_AUTOCENTER
////	PROP_AXISMODE
////	PROP_BUFFERSIZE
////	PROP_CALIBRATION
////	PROP_CALIBRATIONMODE
////	PROP_CPOINTS
////	PROP_DEADZONE
////	PROP_FFGAIN
////	PROP_INSTANCENAME
////	PROP_PRODUCTNAME
////	PROP_RANGE
////	PROP_SATURATION
////)
//
////var predefinedProperties = []C.REFGUID{
////	C.DIPROP_APPDATA,
////	C.DIPROP_AUTOCENTER,
////	C.DIPROP_AXISMODE,
////	C.DIPROP_BUFFERSIZE,
////	C.DIPROP_CALIBRATION,
////	C.DIPROP_CALIBRATIONMODE,
////	C.DIPROP_CPOINTS,
////	C.DIPROP_DEADZONE,
////	C.DIPROP_FFGAIN,
////	C.DIPROP_INSTANCENAME,
////	C.DIPROP_PRODUCTNAME,
////	C.DIPROP_RANGE,
////	C.DIPROP_SATURATION,
////}
//
//func (obj Device) Unacquire() (err error) {
//	err = toError(C.IDirectInputDevice8Unacquire(obj.handle))
//	return
//}
//
////func (obj Device) WriteEffectToFile() (err error) {
////	err = toError(C.IDirectInputDevice8WriteEffectToFile(obj.handle))
////  return
////}
//
