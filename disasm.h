#define PACKAGE 1
#define PACKAGE_VERSION 1
#include "dis-asm.h"

typedef void *DisAsmPtr;
typedef unsigned long DisAsmLen;

typedef struct DisAsmPrintBuffer {
	int	index;
	char	data[1024];
} DisAsmPrintBufferType, *DisAsmPrintBufferPtr;

typedef struct DisAsmInfo {
	disassemble_info info;
	DisAsmPrintBufferType disAsmPrintBuffer;
} DisAsmInfoType, *DisAsmInfoPtr;

extern DisAsmInfoPtr DisAsmInfoInit(DisAsmPtr start, DisAsmLen length);
extern int DisAsmPrintInstruction(DisAsmInfoPtr disAsmInfoPtr, DisAsmPtr pc, int doPrint);
extern int DisAsmPrintGadget(DisAsmInfoPtr disAsmInfoPtr, DisAsmPtr pc, int doPrint);
extern void DisAsmInfoFree(DisAsmInfoPtr disAsmInfoPtr);
