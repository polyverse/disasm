#include <assert.h>
#include <memory.h>
#include <stdarg.h>
#include "DisAsm.h"

static void DisAsmPrintAddress(bfd_vma addr, struct disassemble_info *info)
{
        info->fprintf_func(info->stream, "0x%lx", addr);
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

void DisAsmCommencement(void)
{
} // DisAsmCommencement()

void DisAsmInfoInit(DisAsmInfoType *disAsmInfoPtr, DisAsmPtr start, DisAsmPtr end)
{
	memset(disAsmInfoPtr, 0, sizeof(*disAsmInfoPtr));

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
        disAsmInfoPtr->info.buffer_length             = end-start;
        disAsmInfoPtr->info.buffer                    = start;
} // DisAsmInfoInit()

int DisAsmPrintGadget(DisAsmInfoType *disAsmInfoPtr, DisAsmPtr pc, int doPrint)
{
        DisAsmPtr end = disAsmInfoPtr->info.buffer + disAsmInfoPtr->info.buffer_length;
	int instructions = 0;

        for (DisAsmPtr pc0 = pc; pc0 < end;)
        {
                unsigned char b = *((unsigned char *) pc0);
                int good = b == 0xC3; // ret
                int bad  = ((b == 0xE9) || (b == 0xEA) || (b == 0xEB) || (b == 0xFF)); // jmps. ToDo: More work here

                disAsmInfoPtr->disAsmPrintBuffer.index = 0;

                if (doPrint)
			printf("0x%p  ", pc0);
                int count = (int) print_insn_i386((unsigned long) pc0, &disAsmInfoPtr->info);
                assert(count != 0);

		if (doPrint)
		{
                	for (int i = 0; i < 8; i++)
                		if (i < count)
                        		printf("%02x", *((unsigned char *) ((unsigned long) pc0 ) + i));
                		else
                        		printf("  ");
               		printf("  %s\n", disAsmInfoPtr->disAsmPrintBuffer.data);
		} // if

                pc0 += count;

		instructions++;

                if (good ^ bad)
                        return good ? instructions : 0;
        } // for

        return 0;
} // DisAsmPrintGadget()

void DisAsmFin(void)
{
} // DisAsmFin()
