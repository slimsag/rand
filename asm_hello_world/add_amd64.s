TEXT ·add(SB),7,$0-12
	MOVL	a+0(FP),AX
	MOVL	AX,BX
	ADDL	b+4(FP),BX
	MOVL	BX,result+8(FP)
	RET

