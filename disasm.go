package disasm

// #include "disasm.h"
import "C"

type DisAsmPtr uintptr
type DisAsmInfoType C.struct_DisAsmInfo;

func DisAsmInfoInit(disAsmInfoPtr *DisAsmInfoType, start DisAsmPtr, end DisAsmPtr) {
	C.DisAsmInfoInit(disAsmInfoPtr, start, end);
	return;
} // DisAsmInfoInit()

func DisAsmPrintGadget(disAsmInfoPtr *DisAsmInfoType, pc DisAsmPtr, doPrint bool) int {
	var b C.int; if doPrint { b = 1; } else { b = 0; } 
	return int(C.DisAsmPrintGadget(disAsmInfoPtr, pc, b));
} // DisAsmPrintGadget()
