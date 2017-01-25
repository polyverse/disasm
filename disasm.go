package disasm

// #include "disasm.h"
import "C"

import "errors"
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

type Instruction struct {
	Address Ptr    `json: "address"`
	Octets  int    `json: "octets"`
	DisAsm  string `json: "disasm"`
}
type InstructionList []Instruction

type Gadget struct {
	Address      Ptr             `json: "address"`
	Octets       int             `json: "octets"`
	Instructions InstructionList `json: "instructions"`
}
type GadgetList []Gadget

func copy(s string) string {
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

func DecodeInstruction(info Info, pc Ptr) (instruction *Instruction, err error) {
        disAsmInfoPtr := info.info.info

        bytes := int(C.DisAsmDecodeInstruction(disAsmInfoPtr, pc))
        s := C.GoStringN(&disAsmInfoPtr.disAsmPrintBuffer.data[0], disAsmInfoPtr.disAsmPrintBuffer.index)
        s = strings.TrimSpace(s)

        return &Instruction{pc, bytes, s}, nil

} // DecodeInstruction()

func DecodeGadget(info Info, pc Ptr) (gadget *Gadget, err error) {
        disAsmInfoPtr := info.info.info
	g := Gadget{Address: pc, Octets: 0, Instructions: nil}

        for pc0 := pc; pc0 < Ptr(disAsmInfoPtr.end); {
                var b byte = *(*byte)(unsafe.Pointer(pc0))
                var good bool = b == 0xC3                                                 // ret
                var bad bool = ((b == 0xE9) || (b == 0xEA) || (b == 0xEB) || (b == 0xFF)) // jmps. ToDo: More work here

                instruction, _ := DecodeInstruction(info, pc0)
		g.Octets += instruction.Octets
                g.Instructions = append(g.Instructions, *instruction)

                pc0 = Ptr(uintptr(pc0) + uintptr(instruction.Octets))

                if good {
                        return &g, nil
                } else if bad {
                        return nil, errors.New("Encountered jmp instruction")
                }
        } // for

	return nil, errors.New("Nothing found")
} // DecodeGadget()

func InfoFree(i *iInfo) {
	C.DisAsmInfoFree(i.info)
	i.info = nil
} // InfoFree()
