package main

// #include "disasm.h"
import "C"

import "fmt"

func main() {
	// Two dummy functions that "bracket" disasm.c (for demonstration purposes)
        start := uintptr(C.DisAsmCommencement);
        end   := uintptr(C.DisAsmFin);

	var DisAsmInfo C.struct_DisAsmInfo
	C.DisAsmInfoInit(&DisAsmInfo, start, end);

	gadgets := 0;
	for pc := start; pc < end; pc = pc + 1 {
		instructions := C.DisAsmPrintGadget(&DisAsmInfo, pc, 0);

		if (instructions > 0) {
			fmt.Printf("GADGET AT: %x (Length: %d)\n", pc, instructions);
			C.DisAsmPrintGadget(&DisAsmInfo, pc, 1);
			fmt.Printf("\n");
			gadgets++;
		} // if
	} // for 

	fmt.Printf("GADGET COUNT BETWEEN 0x%x and 0x%x: %d (%d%%)\n", start, end, gadgets, gadgets * 100 / int((uintptr(end) - uintptr(start))));
} // main()
