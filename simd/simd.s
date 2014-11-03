// avxSupport uses CPUID to check if the CPU supports AVX
// func avxSupport() bool
// Return: +0(FP)
TEXT ·avxSupport(SB),$0
	MOVL $0x1, AX
	CPUID

	// AND ECX, 0x18000000
	BYTE $0x81; BYTE $0xE1; BYTE $0x00; BYTE $0x00; BYTE $0x00; BYTE $0x18;

	// CMP ECX, 0x18000000
	BYTE $0x81; BYTE $0xF9; BYTE $0x00; BYTE $0x00; BYTE $0x00; BYTE $0x18;

	JNE notSupported

	MOVB $0x1, ret+0(FP)
	RET

notSupported:
	MOVB $0x0, ret+0(FP)
	RET

