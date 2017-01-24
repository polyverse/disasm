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
	var b C.int
	if doPrint {
		b = 1
	} else {
		b = 0
	}
	return int(C.DisAsmPrintGadget(info.info.info, pc, b))
} // PrintGadget()

func InfoFree(i *iInfo) {
	C.DisAsmInfoFree(i.info)
} // InfoFree()
