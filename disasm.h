#define PACKAGE 1
#define PACKAGE_VERSION 1
#include "dis-asm.h"

typedef void *DisAsmPtr;

typedef struct DisAsmPrintBuffer {
	int	index;
	char	data[1024]; // Expects more than 100 bytes of free space
} DisAsmPrintBufferType, *DisAsmPrintBufferPtr;

typedef struct DisAsmInfo {
	disassemble_info info;
	DisAsmPrintBufferType disAsmPrintBuffer;
} DisAsmInfoType, *DisAsmInfoPtr;

extern DisAsmInfoPtr DisAsmInfoInit(DisAsmPtr start, DisAsmPtr end);
extern int DisAsmPrintInstruction(DisAsmInfoPtr disAsmInfoPtr, DisAsmPtr pc, int doPrint);
extern int DisAsmPrintGadget(DisAsmInfoPtr disAsmInfoPtr, DisAsmPtr pc, int doPrint);
extern void DisAsmInfoFree(DisAsmInfoPtr disAsmInfoPtr);
