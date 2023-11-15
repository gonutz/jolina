package di8

import (
	"syscall"
	"unsafe"
)

type HWND uintptr
type HINSTANCE uintptr

type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]uint8
}

type DEVICEINSTANCE struct {
	Size         uint32
	GuidInstance GUID
	GuidProduct  GUID
	DevType      uint32
	InstanceName [MAX_PATH]uint16
	ProductName  [MAX_PATH]uint16
	GuidFFDriver GUID
	UsagePage    uint16
	Usage        uint16
}

type DATAFORMAT struct {
	Size     uint32
	ObjSize  uint32
	Flags    uint32
	DataSize uint32
	NumObjs  uint32
	Rgodf    *OBJECTDATAFORMAT
}

type OBJECTDATAFORMAT struct {
	Guid  *GUID
	Ofs   uint32
	Type  uint32
	Flags uint32
}

type PROPHEADER struct {
	Size       uint32
	HeaderSize uint32
	Obj        uint32
	How        uint32
}

func (p *PROPHEADER) PropHeader() *PROPHEADER { return p }

func NewPropCPoints(obj, how uint32, points []CPOINT) *PROPCPOINTS {
	var p PROPCPOINTS
	p.Obj = obj
	p.How = how
	p.Size = uint32(unsafe.Sizeof(p))
	p.HeaderSize = uint32(unsafe.Sizeof(p.PROPHEADER))
	if len(points) > MAXCPOINTSNUM {
		points = points[:MAXCPOINTSNUM]
	}
	p.CPointsNum = uint32(len(points))
	for i := range points {
		p.Points[i] = points[i]
	}
	return &p
}

type PROPCPOINTS struct {
	PROPHEADER
	CPointsNum uint32
	Points     [MAXCPOINTSNUM]CPOINT
}

func (p *PROPCPOINTS) PropHeader() *PROPHEADER { return &p.PROPHEADER }

type CPOINT struct {
	P   int32
	Log uint32
}

func NewPropDWord(obj, how, data uint32) *PROPDWORD {
	var p PROPDWORD
	p.Obj = obj
	p.How = how
	p.Size = uint32(unsafe.Sizeof(p))
	p.HeaderSize = uint32(unsafe.Sizeof(p.PROPHEADER))
	p.Data = data
	return &p
}

type PROPDWORD struct {
	PROPHEADER
	Data uint32
}

func (p *PROPDWORD) PropHeader() *PROPHEADER { return &p.PROPHEADER }

func NewPropRange(obj, how uint32, min, max int32) *PROPRANGE {
	var p PROPRANGE
	p.Obj = obj
	p.How = how
	p.Size = uint32(unsafe.Sizeof(p))
	p.HeaderSize = uint32(unsafe.Sizeof(p.PROPHEADER))
	p.Min = min
	p.Max = max
	return &p
}

type PROPRANGE struct {
	PROPHEADER
	Min int32
	Max int32
}

func (p *PROPRANGE) PropHeader() *PROPHEADER { return &p.PROPHEADER }

func NewPropCal(obj, how uint32, min, center, max int32) *PROPCAL {
	var p PROPCAL
	p.Obj = obj
	p.How = how
	p.Size = uint32(unsafe.Sizeof(p))
	p.HeaderSize = uint32(unsafe.Sizeof(p.PROPHEADER))
	p.Min = min
	p.Center = center
	p.Max = max
	return &p
}

type PROPCAL struct {
	PROPHEADER
	Min    int32
	Center int32
	Max    int32
}

func (p *PROPCAL) PropHeader() *PROPHEADER { return &p.PROPHEADER }

func NewPropCalPOV(obj, how uint32, min, max [5]int32) *PROPCALPOV {
	var p PROPCALPOV
	p.Obj = obj
	p.How = how
	p.Size = uint32(unsafe.Sizeof(p))
	p.HeaderSize = uint32(unsafe.Sizeof(p.PROPHEADER))
	p.Min = min
	p.Max = max
	return &p
}

type PROPCALPOV struct {
	PROPHEADER
	Min [5]int32
	Max [5]int32
}

func (p *PROPCALPOV) PropHeader() *PROPHEADER { return &p.PROPHEADER }

func NewPropGuidAndPath(obj, how uint32, guid GUID, path string) *PROPGUIDANDPATH {
	var p PROPGUIDANDPATH
	p.Obj = obj
	p.How = how
	p.Size = uint32(unsafe.Sizeof(p))
	p.HeaderSize = uint32(unsafe.Sizeof(p.PROPHEADER))
	p.GuidClass = guid
	str, err := syscall.UTF16FromString(path)
	if err == nil {
		if len(str) > MAX_PATH {
			str = str[:MAX_PATH]
			str[MAX_PATH-1] = 0
		}
		for i := range str {
			p.Path[i] = str[i]
		}
	}
	return &p
}

type PROPGUIDANDPATH struct {
	PROPHEADER
	GuidClass GUID
	Path      [MAX_PATH]uint16
}

func (p *PROPGUIDANDPATH) PropHeader() *PROPHEADER { return &p.PROPHEADER }

func NewPropString(obj, how uint32, s string) *PROPSTRING {
	var p PROPSTRING
	p.Obj = obj
	p.How = how
	p.Size = uint32(unsafe.Sizeof(p))
	p.HeaderSize = uint32(unsafe.Sizeof(p.PROPHEADER))
	str, err := syscall.UTF16FromString(s)
	if err == nil {
		if len(str) > MAX_PATH {
			str = str[:MAX_PATH]
			str[MAX_PATH-1] = 0
		}
		for i := range str {
			p.String[i] = str[i]
		}
	}
	return &p
}

type PROPSTRING struct {
	PROPHEADER
	String [MAX_PATH]uint16
}

func (p *PROPSTRING) PropHeader() *PROPHEADER { return &p.PROPHEADER }

func NewPropPointer(obj, how uint32, pointer uintptr) *PROPPOINTER {
	var p PROPPOINTER
	p.Obj = obj
	p.How = how
	p.Size = uint32(unsafe.Sizeof(p))
	p.HeaderSize = uint32(unsafe.Sizeof(p.PROPHEADER))
	p.Data = pointer
	return &p
}

type PROPPOINTER struct {
	PROPHEADER
	Data uintptr
}

func (p *PROPPOINTER) PropHeader() *PROPHEADER { return &p.PROPHEADER }

type DEVICEOBJECTDATA struct {
	Ofs       uint32
	Data      uint32
	TimeStamp uint32
	Sequence  uint32
	AppData   uintptr
}

type MOUSESTATE struct {
	X       int32
	Y       int32
	Z       int32
	Buttons [4]byte
}

func (s *MOUSESTATE) Ptr() uintptr { return uintptr(unsafe.Pointer(s)) }
func (s *MOUSESTATE) Size() int    { return int(unsafe.Sizeof(*s)) }

type MOUSESTATE2 struct {
	X       int32
	Y       int32
	Z       int32
	Buttons [8]byte
}

func (s *MOUSESTATE2) Ptr() uintptr { return uintptr(unsafe.Pointer(s)) }
func (s *MOUSESTATE2) Size() int    { return int(unsafe.Sizeof(*s)) }

type KEYBOARDSTATE [256]byte

func (s *KEYBOARDSTATE) Ptr() uintptr { return uintptr(unsafe.Pointer(&s[0])) }
func (s *KEYBOARDSTATE) Size() int    { return len(*s) }

type JOYSTATE struct {
	X       int32
	Y       int32
	Z       int32
	Rx      int32
	Ry      int32
	Rz      int32
	Slider  [2]int32
	POV     [4]uint32
	Buttons [32]byte
}

func (s *JOYSTATE) Ptr() uintptr { return uintptr(unsafe.Pointer(s)) }
func (s *JOYSTATE) Size() int    { return int(unsafe.Sizeof(*s)) }

type JOYSTATE2 struct {
	X       int32
	Y       int32
	Z       int32
	Rx      int32
	Ry      int32
	Rz      int32
	Slider  [2]int32
	POV     [4]uint32
	Buttons [128]byte
	VX      int32
	VY      int32
	VZ      int32
	VRx     int32
	VRy     int32
	VRz     int32
	VSlider [2]int32
	AX      int32
	AY      int32
	AZ      int32
	ARx     int32
	ARy     int32
	ARz     int32
	ASlider [2]int32
	FX      int32
	FY      int32
	FZ      int32
	FRx     int32
	FRy     int32
	FRz     int32
	FSlider [2]int32
}

func (s *JOYSTATE2) Ptr() uintptr { return uintptr(unsafe.Pointer(s)) }
func (s *JOYSTATE2) Size() int    { return int(unsafe.Sizeof(*s)) }

///*
//#include "dinput_wrapper.h"
//
//LPCTSTR getActionName(DIACTION* action) {
//	return action->lptszActionName;
//}
//
//UINT getResIdString(DIACTION* action) {
//	return action->uResIdString;
//}
//*/
//import "C"
//import "unsafe"
//
//type Action struct {
//	AppData      uintptr
//	Semantic     uint32
//	Flags        uint32
//	ResIdString  uint
//	GuidInstance GUID
//	ObjID        uint32
//	How          uint32
//}
//
//func (s *Action) fromC(c *C.DIACTION) {
//	s.AppData = uintptr(c.uAppData)
//	s.Semantic = uint32(c.dwSemantic)
//	s.Flags = uint32(c.dwFlags)
//	// NOTE using c.uResIdString or c.union_uResIdString gives error (undefined)
//	s.ResIdString = uint(C.getResIdString(c))
//	s.GuidInstance.fromC(&c.guidInstance)
//	s.ObjID = uint32(c.dwObjID)
//	s.How = uint32(c.dwHow)
//}
//

//
//func (s *GUID) toC() C.GUID {
//	var c C.GUID
//	c.Data1 = (C.ulong)(s.Data1)
//	c.Data2 = (C.ushort)(s.Data2)
//	c.Data3 = (C.ushort)(s.Data3)
//	for i := range s.Data4 {
//		c.Data4[i] = (C.uchar)(s.Data4[i])
//	}
//	return c
//}
//
//func (s *GUID) fromC(c *C.GUID) {
//	s.Data1 = (uint32)(c.Data1)
//	s.Data2 = (uint16)(c.Data2)
//	s.Data3 = (uint16)(c.Data3)
//	for i := range s.Data4 {
//		s.Data4[i] = (uint8)(c.Data4[i])
//	}
//}
//

//
//func (s *DeviceInstance) fromC(c *C.DIDEVICEINSTANCE) {
//	s.GuidInstance.fromC(&c.guidInstance)
//	s.GuidProduct.fromC(&c.guidProduct)
//	s.DevType = uint32(c.dwDevType)
//	s.InstanceName = maxPathStringToGoString(&c.tszInstanceName)
//	s.ProductName = maxPathStringToGoString(&c.tszProductName)
//	s.GuidFFDriver.fromC(&c.guidFFDriver)
//	s.UsagePage = uint16(c.wUsagePage)
//	s.Usage = uint16(c.wUsage)
//}
//
//type DeviceObjectInstance struct {
//	GuidType          GUID
//	Ofs               uint32
//	Type              uint32
//	Flags             uint32
//	Name              string
//	FFMaxForce        uint32
//	FFForceResolution uint32
//	CollectionNumber  uint16
//	DesignatorIndex   uint16
//	UsagePage         uint16
//	Usage             uint16
//	Dimension         uint32
//	Exponent          uint16
//	ReportId          uint16
//}
//
//func (s *DeviceObjectInstance) fromC(c *C.DIDEVICEOBJECTINSTANCE) {
//	s.GuidType.fromC(&c.guidType)
//	s.Ofs = uint32(c.dwOfs)
//	s.Type = uint32(c.dwType)
//	s.Flags = uint32(c.dwFlags)
//	s.Name = maxPathStringToGoString(&c.tszName)
//	s.FFMaxForce = uint32(c.dwFFMaxForce)
//	s.FFForceResolution = uint32(c.dwFFForceResolution)
//	s.CollectionNumber = uint16(c.wCollectionNumber)
//	s.DesignatorIndex = uint16(c.wDesignatorIndex)
//	s.UsagePage = uint16(c.wUsagePage)
//	s.Usage = uint16(c.wUsage)
//	s.Dimension = uint32(c.dwDimension)
//	s.Exponent = uint16(c.wExponent)
//	s.ReportId = uint16(c.wReportId)
//}
//
//type EffectInfo struct {
//	Guid          GUID
//	EffType       uint32
//	StaticParams  uint32
//	DynamicParams uint32
//	Name          string
//	//TCHAR tszName[MAX_PATH];
//}
//
//func (s *EffectInfo) fromC(c *C.DIEFFECTINFO) {
//	s.Guid.fromC(&c.guid)
//	s.EffType = uint32(c.dwEffType)
//	s.StaticParams = uint32(c.dwStaticParams)
//	s.DynamicParams = uint32(c.dwDynamicParams)
//	s.Name = maxPathStringToGoString(&c.tszName)
//
//}
//
//func maxPathStringToGoString(str *[C.MAX_PATH]C.CHAR) string {
//	buffer := make([]byte, C.MAX_PATH)
//	length := 0
//	for ; length < C.MAX_PATH; length++ {
//		if str[length] == 0 {
//			break
//		}
//		buffer[length] = byte(str[length])
//	}
//	return string(buffer[:length])
//}
//
//type FILETIME struct {
//	LowDateTime  uint32
//	HighDateTime uint32
//}
//
//func (s *FILETIME) toC() C.FILETIME {
//	var c C.FILETIME
//	c.dwLowDateTime = C.DWORD(s.LowDateTime)
//	c.dwHighDateTime = C.DWORD(s.HighDateTime)
//	return c
//}
//
//func (s *FILETIME) fromC(c *C.FILETIME) {
//	s.LowDateTime = uint32(c.dwLowDateTime)
//	s.HighDateTime = uint32(c.dwHighDateTime)
//}
//
//type CPOINT struct {
//	P   int32
//	Log uint32
//}
//
//type ACTION struct {
//	AppData      uintptr
//	Semantic     uint32
//	Flags        uint32
//	GuidInstance GUID
//	ObjID        uint32
//	How          uint32
//}
//
//type ACTIONFORMAT struct {
//	ActionSize    uint32
//	DataSize      uint32
//	NumActions    uint32
//	Action        *ACTION
//	GuidActionMap GUID
//	Genre         uint32
//	BufferSize    uint32
//	AxisMin       int32
//	AxisMax       int32
//	InstString    unsafe.Pointer
//	TimeStamp     FILETIME
//	CRC           uint32
//	ActionMap     [260]byte
//}
//
//type COLORSET struct {
//	TextFore         uint32
//	TextHighlight    uint32
//	CalloutLine      uint32
//	CalloutHighlight uint32
//	Border           uint32
//	ControlFill      uint32
//	HighlightFill    uint32
//	AreaFill         uint32
//}
//
//type CONDITION struct {
//	Offset              int32
//	PositiveCoefficient int32
//	NegativeCoefficient int32
//	PositiveSaturation  uint32
//	NegativeSaturation  uint32
//	DeadBand            int32
//}
//
//type CONFIGUREDEVICESPARAMS struct {
//	UsersSize    uint32
//	UserNames    string
//	FormatsSize  uint32
//	Formats      *ACTIONFORMAT
//	Hwnd         unsafe.Pointer
//	Dics         COLORSET
//	UnkDDSTarget unsafe.Pointer
//}
//
//type CONSTANTFORCE struct {
//	Magnitude int32
//}
//
//type CUSTOMFORCE struct {
//	Channels     uint32
//	SamplePeriod uint32
//	Samples      uint32
//	ForceData    *int32
//}
//

//
//type DEVCAPS struct {
//	Flags               uint32
//	DevType             uint32
//	Axes                uint32
//	Buttons             uint32
//	POVs                uint32
//	FFSamplePeriod      uint32
//	FFMinTimeResolution uint32
//	FirmwareRevision    uint32
//	HardwareRevision    uint32
//	FFDriverVersion     uint32
//}
//
//type DEVICEIMAGEINFO struct {
//	ImagePath    [260]byte
//	Flags        uint32
//	ViewID       uint32
//	Overlay      RECT
//	ObjID        uint32
//	ValidPtsSize uint32
//	CalloutLine  [5]POINT
//	CalloutRect  RECT
//	TextAlign    uint32
//}
//
//type DEVICEIMAGEINFOHEADER struct {
//	SizeImageInfo  uint32
//	ViewsSize      uint32
//	ButtonsSize    uint32
//	AxesSize       uint32
//	POVsSize       uint32
//	BufferSize     uint32
//	BufferUsed     uint32
//	ImageInfoArray DEVICEIMAGEINFO
//}
//
//type DEVICEINSTANCE struct {
//	GuidInstance GUID
//	GuidProduct  GUID
//	DevType      uint32
//	InstanceName [260]byte
//	ProductName  [260]byte
//	GuidFFDriver GUID
//	UsagePage    uint16
//	Usage        uint16
//}
//
//type DEVICEOBJECTDATA struct {
//	Ofs       uint32
//	Data      uint32
//	TimeStamp uint32
//	Sequence  uint32
//	AppData   uintptr
//}
//
//func (s *DEVICEOBJECTDATA) fromC(c *C.DIDEVICEOBJECTDATA) {
//	s.Ofs = uint32(c.dwOfs)
//	s.Data = uint32(c.dwData)
//	s.TimeStamp = uint32(c.dwTimeStamp)
//	s.Sequence = uint32(c.dwSequence)
//	s.AppData = uintptr(c.uAppData)
//}
//
//type DEVICEOBJECTINSTANCE struct {
//	GuidType          GUID
//	Ofs               uint32
//	Type              uint32
//	Flags             uint32
//	Name              [260]byte
//	FFMaxForce        uint32
//	FFForceResolution uint32
//	CollectionNumber  uint16
//	DesignatorIndex   uint16
//	UsagePage         uint16
//	Usage             uint16
//	Dimension         uint32
//	Exponent          uint16
//	ReportId          uint16
//}
//
//type EFFECT struct {
//	Flags                  uint32
//	Duration               uint32
//	SamplePeriod           uint32
//	Gain                   uint32
//	TriggerButton          uint32
//	TriggerRepeatInterval  uint32
//	AxesCount              uint32
//	Axes                   *uint32
//	Direction              *int32
//	Envelope               *ENVELOPE
//	TypeSpecificParamsSize uint32
//	TypeSpecificParams     unsafe.Pointer
//	StartDelay             uint32
//}
//
//type EFFECTINFO struct {
//	Guid          GUID
//	EffType       uint32
//	StaticParams  uint32
//	DynamicParams uint32
//	Name          [260]byte
//}
//
//type EFFESCAPE struct {
//	Command       uint32
//	InBuffer      unsafe.Pointer
//	InBufferSize  uint32
//	OutBuffer     unsafe.Pointer
//	OutBufferSize uint32
//}
//
//type ENVELOPE struct {
//	AttackLevel uint32
//	AttackTime  uint32
//	FadeLevel   uint32
//	FadeTime    uint32
//}
//
//type FILEEFFECT struct {
//	GuidEffect   GUID
//	Effect       *EFFECT
//	FriendlyName [260]byte
//}
//
//type PERIODIC struct {
//	Magnitude uint32
//	Offset    int32
//	Phase     uint32
//	Period    uint32
//}
//
//type RAMPFORCE struct {
//	Start int32
//	End   int32
//}
//
//type RECT struct {
//	Left   int32
//	Top    int32
//	Right  int32
//	Bottom int32
//}
//
//type POINT struct {
//	X int32
//	Y int32
//}
//
