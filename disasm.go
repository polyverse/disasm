package disasm

// #include "disasm.h"
import "C"

import "runtime"

type Ptr uintptr

type iInfo struct {
	info *C.struct_DisAsmInfo
}

type Info struct {
	info *iInfo
}

type InfoPtr uintptr

func InfoInit(start Ptr, end Ptr) Info {
	cinfo := C.DisAsmInfoInit(start, end)
	iinfo := &iInfo{cinfo}
	runtime.SetFinalizer(iinfo, InfoFree)
	info := Info{iinfo}

	return info
} // InfoInit()

func PrintGadget(info Info, pc Ptr, doPrint bool) int {
	if doPrint {
		return int(C.DisAsmPrintGadget(info.info.info, pc, 1))
	} else {
		return int(C.DisAsmPrintGadget(info.info.info, pc, 0))
	}
} // PrintGadget()

func InfoFree(i *iInfo) {
	C.DisAsmInfoFree(i.info)
	i.info = nil
} // InfoFree()
