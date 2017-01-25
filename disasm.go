package disasm

// #include "disasm.h"
import "C"

import "reflect"
import "runtime"
import "strings"
import "unsafe"

type Ptr uintptr
type Len uint64

type iInfo struct {
	info *C.struct_DisAsmInfo
}

type Info struct {
	info *iInfo
}

func Copy(s string) string { 
        var b []byte 
        h := (*reflect.SliceHeader)(unsafe.Pointer(&b)) 
        h.Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data 
        h.Len = len(s) 
        h.Cap = len(s) 
        return string(b) 
} 
func InfoInit(start Ptr, length Len) Info {
	cinfo := C.DisAsmInfoInit(start, C.DisAsmLen(length))
	iinfo := &iInfo{cinfo}
	runtime.SetFinalizer(iinfo, InfoFree)
	info := Info{iinfo}

	return info
} // InfoInit()

func DecodeInstruction(info Info, pc Ptr) (int, string) {
	disAsmInfoPtr := info.info.info;

	bytes := int(C.DisAsmDecodeInstruction (disAsmInfoPtr, pc))
	s := C.GoStringN(&disAsmInfoPtr.disAsmPrintBuffer.data[0], disAsmInfoPtr.disAsmPrintBuffer.index)
	s = strings.TrimSpace(s)

	return bytes, s 

} // DisAsmPrintInstruction()

func DecodeGadget(info Info, pc Ptr) (int, []string) {
	var disAsmInfoPtr C.DisAsmInfoPtr = info.info.info
	var instructions []string
	var bytesTotal int = 0

        for pc0 := pc; pc0 < Ptr(disAsmInfoPtr.end);
        {
		var b byte = *(*byte)(unsafe.Pointer(pc0))
                var good bool = b == 0xC3; // ret
                var bad bool = ((b == 0xE9) || (b == 0xEA) || (b == 0xEB) || (b == 0xFF)); // jmps. ToDo: More work here

                bytes, s := DecodeInstruction(info, pc0)
		bytesTotal += bytes
		instructions = append(instructions, s)

                pc0 = Ptr(uintptr(pc0) + uintptr(bytes))

		if good {
			return bytesTotal, instructions
		} else if bad {
			return 0, nil
		}
        } // for

        return 0, nil
} // DisAsmPrintGadget()

func PrintGadget(info Info, pc Ptr, doPrint bool) int {
        disAsmInfoPtr := info.info.info;
        if doPrint {
                return int(C.DisAsmPrintGadget(disAsmInfoPtr, pc, 1))
        } else {
                return int(C.DisAsmPrintGadget(disAsmInfoPtr, pc, 0))
        }
} // PrintGadget()

func InfoFree(i *iInfo) {
        C.DisAsmInfoFree(i.info)
        i.info = nil
} // InfoFree()
