package main

// #include "DisAsm.h"
import "C"

import "fmt"
//import "unsafe"

func main() {
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

	fmt.Printf("GADGET COUNT BETWEEN %x and %x: %d (%d%%)\n", start, end, gadgets, gadgets * 100 / int((uintptr(end) - uintptr(start))));
} // main()
