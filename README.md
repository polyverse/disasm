# DEPRECATION NOTICE

Please note that this repository has been deprecated and is no longer actively maintained by Polyverse Corporation.  It may be removed in the future, but for now remains public for the benefit of any users.

Importantly, as the repository has not been maintained, it may contain unpatched security issues and other critical issues.  Use at your own risk.

While it is not maintained, we would graciously consider any pull requests in accordance with our Individual Contributor License Agreement.  https://github.com/polyverse/contributor-license-agreement

For any other issues, please feel free to contact info@polyverse.com

---

# polyverse/disasm

This package is a stand-alone (x86_64 only) disassembler and ROP gadget enumerator. It locates all valid gadgets at all alignments within a given block of memory.

A gadget is currently defined as a block of code at some address that terminates with a "ret" (0xC3) instruction and contains no intervening unconditional "jmp" instructions (0xE9 0xEA 0xEB 0xFF). (TBD: This needs to be tighted up)

## Test (While not a real go test module, it will disassemble a known buffer and dump the raw output.)
```
go test .
```

### The nuts and bolts of the disassembler were pulled from the GNU binutils package. Binutils is a heavily interdependent set of utilities that tends to pull in the kitchen sink, so I extracted the bare essentials for the sake of size and expediency. 
