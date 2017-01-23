package disasm

// #include "disasm.h"
import "C"

type Ptr uintptr
type InfoType C.struct_DisAsmInfo

func InfoInit(infoPtr *InfoType, start Ptr, end Ptr) {
	C.DisAsmInfoInit(infoPtr, start, end)
	return
} // InfoInit()

func PrintGadget(infoPtr *InfoType, pc Ptr, doPrint bool) int {
	var b C.int
	if doPrint {
		b = 1
	} else {
		b = 0
	}
	return int(C.DisAsmPrintGadget(infoPtr, pc, b))
} // PrintGadget()
