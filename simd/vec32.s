// func avxVec32Add(a, b Vec32) Vec32
// a: +0(SP)
// b: +32(SP)
// Return: +64(FP)
TEXT ·avxVec32Add(SB),$0-128
	BYTE $0xc5; BYTE $0xfd; BYTE $0x10; BYTE $0x44; BYTE $0x24; BYTE $0x08; // vmovupd 0x8(%rsp),%ymm0
	BYTE $0xc5; BYTE $0xfd; BYTE $0x58; BYTE $0x4c; BYTE $0x24; BYTE $0x28; // vaddpd 0x28(%rsp),%ymm0,%ymm1
	BYTE $0xc5; BYTE $0xfd; BYTE $0x11; BYTE $0x4c; BYTE $0x24; BYTE $0x48; // vmovupd %ymm1,0x48(%rsp)
	RET


// func sse2Vec32Add(a, b Vec32) Vec32
// a: +0(FP)
// b: +32(FP)
// Return: +64(FP)
TEXT ·sse2Vec32Add(SB),$0-128
	MOVUPD    a+0(FP), X0
	MOVUPD    b+32(FP), X2
	ADDPD     X2, X0
	MOVUPD    X0, ret+64(FP)
	MOVUPD    a+16(FP), X0
	MOVUPD    b+48(FP), X2
	ADDPD     X2, X0
	MOVUPD    X0, ret+80(FP)
	RET

// func avxVec32Sub(a, b Vec32) Vec32
// a: +0(SP)
// b: +32(SP)
// Return: +64(FP)
TEXT ·avxVec32Sub(SB),$0-128
	BYTE $0xc5; BYTE $0xfd; BYTE $0x10; BYTE $0x44; BYTE $0x24; BYTE $0x08; // vmovupd 0x8(%rsp),%ymm0
	BYTE $0xc5; BYTE $0xfd; BYTE $0x5c; BYTE $0x4c; BYTE $0x24; BYTE $0x28; // vsubpd 0x28(%rsp),%ymm0,%ymm1
	BYTE $0xc5; BYTE $0xfd; BYTE $0x11; BYTE $0x4c; BYTE $0x24; BYTE $0x48; // vmovupd %ymm1,0x48(%rsp)
	RET

// func sse2Vec32Sub(a, b Vec32) Vec32
// a: +0(FP)
// b: +32(FP)
// Return: +64(FP)
TEXT ·sse2Vec32Sub(SB),$0-128
	MOVUPD    a+0(FP), X0
	MOVUPD    b+32(FP), X2
	SUBPD     X2, X0
	MOVUPD    X0, ret+64(FP)
	MOVUPD    a+16(FP), X0
	MOVUPD    b+48(FP), X2
	SUBPD     X2, X0
	MOVUPD    X0, ret+80(FP)
	RET

// func avxVec32Mul(a, b Vec32) Vec32
// a: +0(SP)
// b: +32(SP)
// Return: +64(FP)
TEXT ·avxVec32Mul(SB),$0-128
	BYTE $0xc5; BYTE $0xfd; BYTE $0x10; BYTE $0x44; BYTE $0x24; BYTE $0x08; // vmovupd 0x8(%rsp),%ymm0
	BYTE $0xc5; BYTE $0xfd; BYTE $0x59; BYTE $0x4c; BYTE $0x24; BYTE $0x28; // vmulpd 0x28(%rsp),%ymm0,%ymm1
	BYTE $0xc5; BYTE $0xfd; BYTE $0x11; BYTE $0x4c; BYTE $0x24; BYTE $0x48; // vmovupd %ymm1,0x48(%rsp)
	RET

// func sse2Vec32Mul(a, b Vec32) Vec32
// a: +0(FP)
// b: +32(FP)
// Return: +64(FP)
TEXT ·sse2Vec32Mul(SB),$0-128
	MOVUPD    a+0(FP), X0
	MOVUPD    b+32(FP), X2
	MULPD     X2, X0
	MOVUPD    X0, ret+64(FP)
	MOVUPD    a+16(FP), X0
	MOVUPD    b+48(FP), X2
	MULPD     X2, X0
	MOVUPD    X0, ret+80(FP)
	RET

// func avxVec32Div(a, b Vec32) Vec32
// a: +0(SP)
// b: +32(SP)
// Return: +64(FP)
TEXT ·avxVec32Div(SB),$0-128
	BYTE $0xc5; BYTE $0xfd; BYTE $0x10; BYTE $0x44; BYTE $0x24; BYTE $0x08; // vmovupd 0x8(%rsp),%ymm0
	BYTE $0xc5; BYTE $0xfd; BYTE $0x5e; BYTE $0x4c; BYTE $0x24; BYTE $0x28; // vdivpd 0x28(%rsp),%ymm0,%ymm1
	BYTE $0xc5; BYTE $0xfd; BYTE $0x11; BYTE $0x4c; BYTE $0x24; BYTE $0x48; // vmovupd %ymm1,0x48(%rsp)
	RET

// func sse2Vec32Div(a, b Vec32) Vec32
// a: +0(FP)
// b: +32(FP)
// Return: +64(FP)
TEXT ·sse2Vec32Div(SB),$0-128
	MOVUPD    a+0(FP), X0
	MOVUPD    b+32(FP), X2
	DIVPD     X2, X0
	MOVUPD    X0, ret+64(FP)
	MOVUPD    a+16(FP), X0
	MOVUPD    b+48(FP), X2
	DIVPD     X2, X0
	MOVUPD    X0, ret+80(FP)
	RET

// func avxVec32Eq(a, b Vec32) bool
// a: +0(FP)
// b: +32(FP)
// Return: +64(FP)
TEXT ·avxVec32Eq(SB),$0-128
	BYTE $0xc5; BYTE $0xfd; BYTE $0x10; BYTE $0x44; BYTE $0x24; BYTE $0x08; // vmovupd 0x8(%rsp),%ymm0
	BYTE $0xc5; BYTE $0xfd; BYTE $0xc2; BYTE $0x4c; BYTE $0x24; BYTE $0x28; BYTE $0x00; // vcmpeqpd 0x28(%rsp),%ymm0,%ymm1

	JNE avxNotEqual

	MOVB $0x1, ret+64(FP)
	RET

avxNotEqual:
	MOVB $0x0, ret+64(FP)
	RET

// func sse2Vec32Eq(a, b Vec32) bool
// a: +0(FP)
// b: +32(FP)
// Return: +64(FP)
TEXT ·sse2Vec32Eq(SB),$0-128
	MOVUPD    a+0(FP), X0
	MOVUPD    b+32(FP), X2
	BYTE $0x66; BYTE $0x0f; BYTE $0xc2; BYTE $0xd0; BYTE $0x00; // cmpeqpd %xmm0,%xmm2
	JNE sse2NotEqual

	MOVUPD    a+16(FP), X0
	MOVUPD    b+48(FP), X2
	BYTE $0x66; BYTE $0x0f; BYTE $0xc2; BYTE $0xd0; BYTE $0x00; // cmpeqpd %xmm0,%xmm2
	JNE sse2NotEqual

	MOVB $0x1, ret+64(FP)
	RET

sse2NotEqual:
	MOVB $0x0, ret+64(FP)
	RET

// func avxVec32Min(a, b Vec32) Vec32
// a: +0(SP)
// b: +32(SP)
// Return: +64(FP)
TEXT ·avxVec32Min(SB),$0-128
	BYTE $0xc5; BYTE $0xfd; BYTE $0x10; BYTE $0x44; BYTE $0x24; BYTE $0x08; // vmovupd 0x8(%rsp),%ymm0
	BYTE $0xc5; BYTE $0xfd; BYTE $0x5d; BYTE $0x4c; BYTE $0x24; BYTE $0x28; // vminpd 0x28(%rsp),%ymm0,%ymm1
	BYTE $0xc5; BYTE $0xfd; BYTE $0x11; BYTE $0x4c; BYTE $0x24; BYTE $0x48; // vmovupd %ymm1,0x48(%rsp)
	RET

// func sse2Vec32Min(a, b Vec32) Vec32
// a: +0(FP)
// b: +32(FP)
// Return: +64(FP)
TEXT ·sse2Vec32Min(SB),$0-128
	MOVUPD    a+0(FP), X0
	MOVUPD    b+32(FP), X2
	MINPD     X2, X0
	MOVUPD    X0, ret+64(FP)
	MOVUPD    a+16(FP), X0
	MOVUPD    b+48(FP), X2
	MINPD     X2, X0
	MOVUPD    X0, ret+80(FP)
	RET

// func avxVec32Max(a, b Vec32) Vec32
// a: +0(SP)
// b: +32(SP)
// Return: +64(FP)
TEXT ·avxVec32Max(SB),$0-128
	BYTE $0xc5; BYTE $0xfd; BYTE $0x10; BYTE $0x44; BYTE $0x24; BYTE $0x08; // vmovupd 0x8(%rsp),%ymm0
	BYTE $0xc5; BYTE $0xfd; BYTE $0x5f; BYTE $0x4c; BYTE $0x24; BYTE $0x28; // vmaxpd 0x28(%rsp),%ymm0,%ymm1
	BYTE $0xc5; BYTE $0xfd; BYTE $0x11; BYTE $0x4c; BYTE $0x24; BYTE $0x48; // vmovupd %ymm1,0x48(%rsp)
	RET

// func sse2Vec32Max(a, b Vec32) Vec32
// a: +0(FP)
// b: +32(FP)
// Return: +64(FP)
TEXT ·sse2Vec32Max(SB),$0-128
	MOVUPD    a+0(FP), X0
	MOVUPD    b+32(FP), X2
	MAXPD     X2, X0
	MOVUPD    X0, ret+64(FP)
	MOVUPD    a+16(FP), X0
	MOVUPD    b+48(FP), X2
	MAXPD     X2, X0
	MOVUPD    X0, ret+80(FP)
	RET

