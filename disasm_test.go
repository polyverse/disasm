package disasm

import "testing"
import "fmt"

func TestDisAsm(t *testing.T) {
	// Two dummy functions that "bracket" disasm.c (for demonstration purposes)
        start := Ptr(0x40029bb); // Clearly, these two values only worked once.
        end   := Ptr(0x4002a09); // Fix these, go get a real test.

	var info InfoType
	InfoInit(&info, start, end);

	gadgets := 0;
	for pc := start; pc < end; pc = pc + 1 {
		instructions := PrintGadget(&info, pc, false);

		if (instructions > 0) {
			fmt.Printf("GADGET AT: 0x%x (Length: %d)\n", pc, instructions);
			PrintGadget(&info, pc, true);
			fmt.Printf("\n");
			gadgets++;
		} // if
	} // for 

	fmt.Printf("GADGET COUNT BETWEEN 0x%x and 0x%x: %d (%d%%)\n", start, end, gadgets, gadgets * 100 / int((uintptr(end) - uintptr(start))));

	// Fix me. start and end need to be set up with a dummy buffer that has predictable content
	if gadgets == 0 {
		t.Error("Failing, because start and end are very likely invalid addresses.");
	} // if
} // TestDisAsm()
