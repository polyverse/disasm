package disasm

// #include "disasm.h"
// #cgo CFLAGS: -std=c99
import "C"

import "errors"
import "runtime"
import "strings"
import "strconv"
import "unsafe"

type Ptr uintptr
type Len uint64
type Sig uint16

type iInfo struct {
	info *C.struct_DisAsmInfo
}

type Info struct {
	info   *iInfo
	start  Ptr
	end    Ptr
	length Len
	memory []byte
}

type Instruction struct {
	Address Ptr    `json:"address,string"`
	Octets  []byte `json:"octets"`
	DisAsm  string `json:"disasm"`
}

type Gadget struct {
	Address      Ptr            `json:"address,string"`
	Instructions []*Instruction `json:"instructions"`
}

func (p Ptr) String() string {
	str := strconv.FormatUint(uint64(p), 16)
	return "0x" + strings.Repeat("0", 12-len(str)) + str
}

func (i *Instruction) String() string {

	octetCount := len(i.Octets)
	b := i.Octets
	s := i.Address.String() + " "

	for o := 0; o < 8; o++ {
		if o < octetCount {
			if b[o] < 16 {
				s += "0"
			} // if
			s += strconv.FormatUint(uint64(b[o]), 16)
		} else {
			s += "  "
		} // else
	} // for

	return s + " " + i.DisAsm

	//return i.DisAsm
}

func (g *Gadget) String() string {
	instrStr := ""
	for _, instr := range g.Instructions {
		if instrStr != "" {
			instrStr = instrStr + ", "
		}
		instrStr = instrStr + instr.DisAsm
	}

	sAdr := strconv.FormatUint(uint64(g.Address), 16)
	return "0x" + strings.Repeat("0", 12-len(sAdr)) + sAdr + ": " + instrStr
}

func InfoInit(s Ptr, e Ptr) Info {
	l := Len(e - s + 1)

	cinfo := C.DisAsmInfoInit(C.DisAsmPtr(s), C.DisAsmPtr(e))
	iinfo := &iInfo{cinfo}
	runtime.SetFinalizer(iinfo, InfoFree)
	info := Info{info: iinfo, start: s, end: e, length: l, memory: C.GoBytes(unsafe.Pointer(s), C.int(l))}

	return info
} // InfoInit()

func InfoInitBytes(s Ptr, e Ptr, b []byte) Info {
	l := Len(e - s + 1)
	if l != Len(len(b)) {
		panic("Disallowed assertion")
	}

	cinfo := C.DisAsmInfoInitBytes(C.DisAsmPtr(s), C.DisAsmPtr(e), unsafe.Pointer(&b[0]))
	iinfo := &iInfo{cinfo}
	runtime.SetFinalizer(iinfo, InfoFree)
	info := Info{info: iinfo, start: s, end: e, length: l, memory: b}

	return info
} // InfoInitBytes()

func DecodeInstruction(info Info, pc Ptr) (instruction *Instruction, err error) {
	disAsmInfoPtr := info.info.info

	octetCount := int(C.DisAsmDecodeInstruction(disAsmInfoPtr, C.DisAsmPtr(pc)))
	if octetCount > 0 {
		s := C.GoStringN(&disAsmInfoPtr.disAsmPrintBuffer.data[0], disAsmInfoPtr.disAsmPrintBuffer.index)
		s = strings.TrimRight(s, " ")

		start := pc - info.start
		end := start + Ptr(octetCount)
		octets := info.memory[start:end]
		return &Instruction{Address: pc, Octets: octets, DisAsm: s}, nil
	} // if

	return nil, errors.New("Error with disassembly")
} // DecodeInstruction()

func DecodeGadget(info Info, pc Ptr, instructions int, numOctets int) (gadget *Gadget, err error) {
	g := &Gadget{
		Address:      pc,
		Instructions: []*Instruction{},
	}

	octetCount := 0

	for pc0 := pc; pc0 <= info.end; {
		var b = info.memory[pc0-info.start]
		var good bool = ((b == 0xC2) || (b == 0xC3) || (b == 0xCA) || (b == 0xCB))
		var bad bool = ((b == 0xE9) || (b == 0xEA) || (b == 0xEB) || (b == 0xFF)) // JMP, JMP, JMP, 0xFF

		instruction, err := DecodeInstruction(info, pc0)
		if err != nil {
			return nil, err
		}
		if strings.Contains(instruction.DisAsm, "(bad)") {
			return nil, errors.New("Encountered (bad) instruction")
		} // if

		octetCount += len(instruction.Octets)
		g.Instructions = append(g.Instructions, instruction)

		if (octetCount > numOctets) || (len(g.Instructions) > instructions) {
			return nil, errors.New("Limits exceeded")
		} // if

		pc0 = Ptr(uintptr(pc0) + uintptr(len(instruction.Octets)))

		if good {
			return g, nil
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
