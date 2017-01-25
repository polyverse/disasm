#include <assert.h>
#include <memory.h>
#include <stdarg.h>
#include <stdlib.h>
#include "disasm.h"

static void DisAsmPrintAddress(bfd_vma addr, struct disassemble_info *info)
{
        info->fprintf_func(info->stream, "%p", addr);
} // DisAsmPrintAddress()

static int DisAsmPrintf(void *b, const char *fmt, ...)
{
	DisAsmPrintBufferPtr pbPtr = (DisAsmPrintBufferPtr) b;

	va_list arglist;
	va_start(arglist, fmt);
	int result = vsprintf(pbPtr->data + pbPtr->index, fmt, arglist);
	assert(pbPtr->index + result < sizeof(pbPtr->data));
	pbPtr->index += result;
	va_end(arglist);

	return result;
}

DisAsmInfoPtr DisAsmInfoInit(DisAsmPtr start, DisAsmLen length)
{
	DisAsmInfoPtr disAsmInfoPtr = calloc(1, sizeof(*disAsmInfoPtr));

        disAsmInfoPtr->info.flavour                   = bfd_target_unknown_flavour;
        disAsmInfoPtr->info.arch                      = bfd_arch_i386;
        disAsmInfoPtr->info.mach                      = bfd_mach_x86_64_intel_syntax;
        disAsmInfoPtr->info.endian                    = BFD_ENDIAN_LITTLE;
        disAsmInfoPtr->info.octets_per_byte           = 1;
        disAsmInfoPtr->info.fprintf_func              = DisAsmPrintf;
        disAsmInfoPtr->info.stream                    = &disAsmInfoPtr->disAsmPrintBuffer;
        disAsmInfoPtr->info.read_memory_func          = buffer_read_memory;
        disAsmInfoPtr->info.memory_error_func         = perror_memory;
        disAsmInfoPtr->info.print_address_func        = DisAsmPrintAddress;
        disAsmInfoPtr->info.symbol_at_address_func    = generic_symbol_at_address;
        disAsmInfoPtr->info.symbol_is_valid           = generic_symbol_is_valid;
        disAsmInfoPtr->info.display_endian            = BFD_ENDIAN_LITTLE;
        disAsmInfoPtr->info.buffer_vma                = (unsigned long) start;
        disAsmInfoPtr->info.buffer_length             = length;
        disAsmInfoPtr->info.buffer                    = start;
	
	disAsmInfoPtr->start = start;
	disAsmInfoPtr->end = start + length;

	return disAsmInfoPtr;
} // DisAsmInfoInit()

int DisAsmDecodeInstruction(DisAsmInfoType *disAsmInfoPtr, DisAsmPtr pc)
{
	disAsmInfoPtr->disAsmPrintBuffer.index = 0;

        //DisAsmPrintf(disAsmInfoPtr->info.stream, "%p ", pc);

	int count = (int) print_insn_i386((unsigned long) pc, &disAsmInfoPtr->info);
	assert(count != 0);

        //DisAsmPrintf(disAsmInfoPtr->info.stream, "\n");

        return count;
} // DisAsmDecodeInstruction()

int DisAsmPrintGadget(DisAsmInfoType *disAsmInfoPtr, DisAsmPtr pc, int doPrint)
{
        DisAsmPtr end = disAsmInfoPtr->info.buffer + disAsmInfoPtr->info.buffer_length;
	int instructions = 0;

        for (DisAsmPtr pc0 = pc; pc0 < end;)
        {
                unsigned char b = *((unsigned char *) pc0);
                int good = b == 0xC3; // ret
                int bad  = ((b == 0xE9) || (b == 0xEA) || (b == 0xEB) || (b == 0xFF)); // jmps. ToDo: More work here

		int count = DisAsmDecodeInstruction(disAsmInfoPtr, pc0);

		if (doPrint)
			printf("%s\n", disAsmInfoPtr->disAsmPrintBuffer.data);

                pc0 += count;

		instructions++;

                if (good ^ bad)
                        return good ? instructions : 0;
        } // for

        return 0;
} // DisAsmPrintGadget()

void DisAsmInfoFree(DisAsmInfoPtr disAsmInfoPtr)
{
	free(disAsmInfoPtr);
} // DisAsmInfoFree()
