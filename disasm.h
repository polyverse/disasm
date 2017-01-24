#define PACKAGE 1
#define PACKAGE_VERSION 1
#include "dis-asm.h"

typedef void *DisAsmPtr;

typedef struct DisAsmPrintBuffer {
	int	index;
	char	data[1024];
} DisAsmPrintBufferType, *DisAsmPrintBufferPtr;

typedef struct DisAsmInfo {
	disassemble_info info;
	DisAsmPrintBufferType disAsmPrintBuffer;
} DisAsmInfoType, *DisAsmInfoPtr;

extern DisAsmInfoPtr DisAsmInfoInit(DisAsmPtr start, DisAsmPtr end);
extern int DisAsmPrintGadget(DisAsmInfoType *disAsmInfoPtr, DisAsmPtr pc, int doPrint);
extern void DisAsmInfoFree(DisAsmInfoPtr disAsmInfoPtr);
