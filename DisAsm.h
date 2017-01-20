#define PACKAGE 1
#define PACKAGE_VERSION 1
#include "dis-asm.h"

typedef void *DisAsmPtr;

static unsigned char *memory = 0;

typedef struct DisAsmPrintBuffer {
	int	index;
	char	data[1024];
} DisAsmPrintBufferType, *DisAsmPrintBufferPtr;

typedef struct DisAsmInfo {
	disassemble_info info;
	DisAsmPrintBufferType disAsmPrintBuffer;
} DisAsmInfoType, *DisAsmInfoPtr;

extern void DisAsmInfoInit(DisAsmInfoPtr disAsmInfoPtr, DisAsmPtr start, DisAsmPtr end);
extern int DisAsmPrintGadget(DisAsmInfoPtr disAsmInfoPtr, DisAsmPtr pc, int doPrint);
