#include <assert.h>
#include <stdio.h>
#include "DisAsm.h"

extern void fin(void);

int main()
{
	DisAsmPtr start = &main;
	DisAsmPtr end   = &fin;

	DisAsmInfoType disAsmInfo;
	DisAsmInfoInit(&disAsmInfo, (DisAsmPtr) start, end);

	int gadgets = 0;
	for (DisAsmPtr pc = start; pc < end; pc++)
	{
		int instructions = DisAsmPrintGadget(&disAsmInfo, pc, 0);

		if (instructions > 0)
		{
			printf("GADGET AT: %p (Length: %d)\n", pc, instructions);
			assert(DisAsmPrintGadget(&disAsmInfo, pc, 1));
			printf("\n");
			gadgets++;
		} // if
	} // for 

	printf("GADGET COUNT: %d (%ld%%)\n", gadgets, gadgets * 100 / (end - start));
} // main()

void fin()
{
} // fin()
