// Inferno utils/6c/6.out.h
// http://code.google.com/p/inferno-os/source/browse/utils/6c/6.out.h
//
//	Copyright © 1994-1999 Lucent Technologies Inc.  All rights reserved.
//	Portions Copyright © 1995-1997 C H Forsyth (forsyth@terzarima.net)
//	Portions Copyright © 1997-1999 Vita Nuova Limited
//	Portions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com)
//	Portions Copyright © 2004,2006 Bruce Ellis
//	Portions Copyright © 2005-2007 C H Forsyth (forsyth@terzarima.net)
//	Revisions Copyright © 2000-2007 Lucent Technologies Inc. and others
//	Portions Copyright © 2009 The Go Authors.  All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

#define	NSYM	50
#define	NSNAME	8
#include "../ld/textflag.h"

/*
 *	amd64
 */

enum	as
{
	AXXX,
	AAAA,
	AAAD,
	AAAM,
	AAAS,
	AADCB,
	AADCL,
	AADCW,
	AADDB,
	AADDL,
	AADDW,
	AADJSP,
	AANDB,
	AANDL,
	AANDW,
	AARPL,
	ABOUNDL,
	ABOUNDW,
	ABSFL,
	ABSFW,
	ABSRL,
	ABSRW,
	ABTL,
	ABTW,
	ABTCL,
	ABTCW,
	ABTRL,
	ABTRW,
	ABTSL,
	ABTSW,
	ABYTE,
	ACALL,
	ACLC,
	ACLD,
	ACLI,
	ACLTS,
	ACMC,
	ACMPB,
	ACMPL,
	ACMPW,
	ACMPSB,
	ACMPSL,
	ACMPSW,
	ADAA,
	ADAS,
	ADATA,
	ADECB,
	ADECL,
	ADECQ,
	ADECW,
	ADIVB,
	ADIVL,
	ADIVW,
	AENTER,
	AGLOBL,
	AGOK,
	AHISTORY,
	AHLT,
	AIDIVB,
	AIDIVL,
	AIDIVW,
	AIMULB,
	AIMULL,
	AIMULW,
	AINB,
	AINL,
	AINW,
	AINCB,
	AINCL,
	AINCQ,
	AINCW,
	AINSB,
	AINSL,
	AINSW,
	AINT,
	AINTO,
	AIRETL,
	AIRETW,
	AJCC,
	AJCS,
	AJCXZL,
	AJEQ,
	AJGE,
	AJGT,
	AJHI,
	AJLE,
	AJLS,
	AJLT,
	AJMI,
	AJMP,
	AJNE,
	AJOC,
	AJOS,
	AJPC,
	AJPL,
	AJPS,
	ALAHF,
	ALARL,
	ALARW,
	ALEAL,
	ALEAW,
	ALEAVEL,
	ALEAVEW,
	ALOCK,
	ALODSB,
	ALODSL,
	ALODSW,
	ALONG,
	ALOOP,
	ALOOPEQ,
	ALOOPNE,
	ALSLL,
	ALSLW,
	AMOVB,
	AMOVL,
	AMOVW,
	AMOVBLSX,
	AMOVBLZX,
	AMOVBQSX,
	AMOVBQZX,
	AMOVBWSX,
	AMOVBWZX,
	AMOVWLSX,
	AMOVWLZX,
	AMOVWQSX,
	AMOVWQZX,
	AMOVSB,
	AMOVSL,
	AMOVSW,
	AMULB,
	AMULL,
	AMULW,
	ANAME,
	ANEGB,
	ANEGL,
	ANEGW,
	ANOP,
	ANOTB,
	ANOTL,
	ANOTW,
	AORB,
	AORL,
	AORW,
	AOUTB,
	AOUTL,
	AOUTW,
	AOUTSB,
	AOUTSL,
	AOUTSW,
	APAUSE,
	APOPAL,
	APOPAW,
	APOPFL,
	APOPFW,
	APOPL,
	APOPW,
	APUSHAL,
	APUSHAW,
	APUSHFL,
	APUSHFW,
	APUSHL,
	APUSHW,
	ARCLB,
	ARCLL,
	ARCLW,
	ARCRB,
	ARCRL,
	ARCRW,
	AREP,
	AREPN,
	ARET,
	AROLB,
	AROLL,
	AROLW,
	ARORB,
	ARORL,
	ARORW,
	ASAHF,
	ASALB,
	ASALL,
	ASALW,
	ASARB,
	ASARL,
	ASARW,
	ASBBB,
	ASBBL,
	ASBBW,
	ASCASB,
	ASCASL,
	ASCASW,
	ASETCC,
	ASETCS,
	ASETEQ,
	ASETGE,
	ASETGT,
	ASETHI,
	ASETLE,
	ASETLS,
	ASETLT,
	ASETMI,
	ASETNE,
	ASETOC,
	ASETOS,
	ASETPC,
	ASETPL,
	ASETPS,
	ACDQ,
	ACWD,
	ASHLB,
	ASHLL,
	ASHLW,
	ASHRB,
	ASHRL,
	ASHRW,
	ASTC,
	ASTD,
	ASTI,
	ASTOSB,
	ASTOSL,
	ASTOSW,
	ASUBB,
	ASUBL,
	ASUBW,
	ASYSCALL,
	ATESTB,
	ATESTL,
	ATESTW,
	ATEXT,
	AVERR,
	AVERW,
	AWAIT,
	AWORD,
	AXCHGB,
	AXCHGL,
	AXCHGW,
	AXLAT,
	AXORB,
	AXORL,
	AXORW,

	AFMOVB,
	AFMOVBP,
	AFMOVD,
	AFMOVDP,
	AFMOVF,
	AFMOVFP,
	AFMOVL,
	AFMOVLP,
	AFMOVV,
	AFMOVVP,
	AFMOVW,
	AFMOVWP,
	AFMOVX,
	AFMOVXP,

	AFCOMB,
	AFCOMBP,
	AFCOMD,
	AFCOMDP,
	AFCOMDPP,
	AFCOMF,
	AFCOMFP,
	AFCOML,
	AFCOMLP,
	AFCOMW,
	AFCOMWP,
	AFUCOM,
	AFUCOMP,
	AFUCOMPP,

	AFADDDP,
	AFADDW,
	AFADDL,
	AFADDF,
	AFADDD,

	AFMULDP,
	AFMULW,
	AFMULL,
	AFMULF,
	AFMULD,

	AFSUBDP,
	AFSUBW,
	AFSUBL,
	AFSUBF,
	AFSUBD,

	AFSUBRDP,
	AFSUBRW,
	AFSUBRL,
	AFSUBRF,
	AFSUBRD,

	AFDIVDP,
	AFDIVW,
	AFDIVL,
	AFDIVF,
	AFDIVD,

	AFDIVRDP,
	AFDIVRW,
	AFDIVRL,
	AFDIVRF,
	AFDIVRD,

	AFXCHD,
	AFFREE,

	AFLDCW,
	AFLDENV,
	AFRSTOR,
	AFSAVE,
	AFSTCW,
	AFSTENV,
	AFSTSW,

	AF2XM1,
	AFABS,
	AFCHS,
	AFCLEX,
	AFCOS,
	AFDECSTP,
	AFINCSTP,
	AFINIT,
	AFLD1,
	AFLDL2E,
	AFLDL2T,
	AFLDLG2,
	AFLDLN2,
	AFLDPI,
	AFLDZ,
	AFNOP,
	AFPATAN,
	AFPREM,
	AFPREM1,
	AFPTAN,
	AFRNDINT,
	AFSCALE,
	AFSIN,
	AFSINCOS,
	AFSQRT,
	AFTST,
	AFXAM,
	AFXTRACT,
	AFYL2X,
	AFYL2XP1,

	AEND,

	ADYNT_,
	AINIT_,

	ASIGNAME,

	/* extra 32-bit operations */
	ACMPXCHGB,
	ACMPXCHGL,
	ACMPXCHGW,
	ACMPXCHG8B,
	ACPUID,
	AINVD,
	AINVLPG,
	ALFENCE,
	AMFENCE,
	AMOVNTIL,
	ARDMSR,
	ARDPMC,
	ARDTSC,
	ARSM,
	ASFENCE,
	ASYSRET,
	AWBINVD,
	AWRMSR,
	AXADDB,
	AXADDL,
	AXADDW,

	/* conditional move */
	ACMOVLCC,
	ACMOVLCS,
	ACMOVLEQ,
	ACMOVLGE,
	ACMOVLGT,
	ACMOVLHI,
	ACMOVLLE,
	ACMOVLLS,
	ACMOVLLT,
	ACMOVLMI,
	ACMOVLNE,
	ACMOVLOC,
	ACMOVLOS,
	ACMOVLPC,
	ACMOVLPL,
	ACMOVLPS,
	ACMOVQCC,
	ACMOVQCS,
	ACMOVQEQ,
	ACMOVQGE,
	ACMOVQGT,
	ACMOVQHI,
	ACMOVQLE,
	ACMOVQLS,
	ACMOVQLT,
	ACMOVQMI,
	ACMOVQNE,
	ACMOVQOC,
	ACMOVQOS,
	ACMOVQPC,
	ACMOVQPL,
	ACMOVQPS,
	ACMOVWCC,
	ACMOVWCS,
	ACMOVWEQ,
	ACMOVWGE,
	ACMOVWGT,
	ACMOVWHI,
	ACMOVWLE,
	ACMOVWLS,
	ACMOVWLT,
	ACMOVWMI,
	ACMOVWNE,
	ACMOVWOC,
	ACMOVWOS,
	ACMOVWPC,
	ACMOVWPL,
	ACMOVWPS,

	/* 64-bit */
	AADCQ,
	AADDQ,
	AANDQ,
	ABSFQ,
	ABSRQ,
	ABTCQ,
	ABTQ,
	ABTRQ,
	ABTSQ,
	ACMPQ,
	ACMPSQ,
	ACMPXCHGQ,
	ACQO,
	ADIVQ,
	AIDIVQ,
	AIMULQ,
	AIRETQ,
	AJCXZQ,
	ALEAQ,
	ALEAVEQ,
	ALODSQ,
	AMOVQ,
	AMOVLQSX,
	AMOVLQZX,
	AMOVNTIQ,
	AMOVSQ,
	AMULQ,
	ANEGQ,
	ANOTQ,
	AORQ,
	APOPFQ,
	APOPQ,
	APUSHFQ,
	APUSHQ,
	ARCLQ,
	ARCRQ,
	AROLQ,
	ARORQ,
	AQUAD,
	ASALQ,
	ASARQ,
	ASBBQ,
	ASCASQ,
	ASHLQ,
	ASHRQ,
	ASTOSQ,
	ASUBQ,
	ATESTQ,
	AXADDQ,
	AXCHGQ,
	AXORQ,

	/* media */
	AADDPD,
	AADDPS,
	AADDSD,
	AADDSS,
	AANDNPD,
	AANDNPS,
	AANDPD,
	AANDPS,
	ACMPPD,
	ACMPPS,
	ACMPSD,
	ACMPSS,
	ACOMISD,
	ACOMISS,
	ACVTPD2PL,
	ACVTPD2PS,
	ACVTPL2PD,
	ACVTPL2PS,
	ACVTPS2PD,
	ACVTPS2PL,
	ACVTSD2SL,
	ACVTSD2SQ,
	ACVTSD2SS,
	ACVTSL2SD,
	ACVTSL2SS,
	ACVTSQ2SD,
	ACVTSQ2SS,
	ACVTSS2SD,
	ACVTSS2SL,
	ACVTSS2SQ,
	ACVTTPD2PL,
	ACVTTPS2PL,
	ACVTTSD2SL,
	ACVTTSD2SQ,
	ACVTTSS2SL,
	ACVTTSS2SQ,
	ADIVPD,
	ADIVPS,
	ADIVSD,
	ADIVSS,
	AEMMS,
	AFXRSTOR,
	AFXRSTOR64,
	AFXSAVE,
	AFXSAVE64,
	ALDMXCSR,
	AMASKMOVOU,
	AMASKMOVQ,
	AMAXPD,
	AMAXPS,
	AMAXSD,
	AMAXSS,
	AMINPD,
	AMINPS,
	AMINSD,
	AMINSS,
	AMOVAPD,
	AMOVAPS,
	AMOVOU,
	AMOVHLPS,
	AMOVHPD,
	AMOVHPS,
	AMOVLHPS,
	AMOVLPD,
	AMOVLPS,
	AMOVMSKPD,
	AMOVMSKPS,
	AMOVNTO,
	AMOVNTPD,
	AMOVNTPS,
	AMOVNTQ,
	AMOVO,
	AMOVQOZX,
	AMOVSD,
	AMOVSS,
	AMOVUPD,
	AMOVUPS,
	AMULPD,
	AMULPS,
	AMULSD,
	AMULSS,
	AORPD,
	AORPS,
	APACKSSLW,
	APACKSSWB,
	APACKUSWB,
	APADDB,
	APADDL,
	APADDQ,
	APADDSB,
	APADDSW,
	APADDUSB,
	APADDUSW,
	APADDW,
	APANDB,
	APANDL,
	APANDSB,
	APANDSW,
	APANDUSB,
	APANDUSW,
	APANDW,
	APAND,
	APANDN,
	APAVGB,
	APAVGW,
	APCMPEQB,
	APCMPEQL,
	APCMPEQW,
	APCMPGTB,
	APCMPGTL,
	APCMPGTW,
	APEXTRW,
	APFACC,
	APFADD,
	APFCMPEQ,
	APFCMPGE,
	APFCMPGT,
	APFMAX,
	APFMIN,
	APFMUL,
	APFNACC,
	APFPNACC,
	APFRCP,
	APFRCPIT1,
	APFRCPI2T,
	APFRSQIT1,
	APFRSQRT,
	APFSUB,
	APFSUBR,
	APINSRW,
	APINSRD,
	APINSRQ,
	APMADDWL,
	APMAXSW,
	APMAXUB,
	APMINSW,
	APMINUB,
	APMOVMSKB,
	APMULHRW,
	APMULHUW,
	APMULHW,
	APMULLW,
	APMULULQ,
	APOR,
	APSADBW,
	APSHUFHW,
	APSHUFL,
	APSHUFLW,
	APSHUFW,
	APSHUFB,
	APSLLO,
	APSLLL,
	APSLLQ,
	APSLLW,
	APSRAL,
	APSRAW,
	APSRLO,
	APSRLL,
	APSRLQ,
	APSRLW,
	APSUBB,
	APSUBL,
	APSUBQ,
	APSUBSB,
	APSUBSW,
	APSUBUSB,
	APSUBUSW,
	APSUBW,
	APSWAPL,
	APUNPCKHBW,
	APUNPCKHLQ,
	APUNPCKHQDQ,
	APUNPCKHWL,
	APUNPCKLBW,
	APUNPCKLLQ,
	APUNPCKLQDQ,
	APUNPCKLWL,
	APXOR,
	ARCPPS,
	ARCPSS,
	ARSQRTPS,
	ARSQRTSS,
	ASHUFPD,
	ASHUFPS,
	ASQRTPD,
	ASQRTPS,
	ASQRTSD,
	ASQRTSS,
	ASTMXCSR,
	ASUBPD,
	ASUBPS,
	ASUBSD,
	ASUBSS,
	AUCOMISD,
	AUCOMISS,
	AUNPCKHPD,
	AUNPCKHPS,
	AUNPCKLPD,
	AUNPCKLPS,
	AXORPD,
	AXORPS,

	APF2IW,
	APF2IL,
	API2FW,
	API2FL,
	ARETFW,
	ARETFL,
	ARETFQ,
	ASWAPGS,

	AMODE,
	ACRC32B,
	ACRC32Q,
	AIMUL3Q,
	
	APREFETCHT0,
	APREFETCHT1,
	APREFETCHT2,
	APREFETCHNTA,
	
	AMOVQL,
	ABSWAPL,
	ABSWAPQ,
	
	AUNDEF,

	AAESENC,
	AAESENCLAST,
	AAESDEC,
	AAESDECLAST,
	AAESIMC,
	AAESKEYGENASSIST,

	APSHUFD,
	APCLMULQDQ,
	
	AUSEFIELD,
	ATYPE,
	AFUNCDATA,
	APCDATA,
	ACHECKNIL,
	AVARDEF,
	AVARKILL,
	ADUFFCOPY,
	ADUFFZERO,
	
	ALAST
};

enum
{

	D_AL		= 0,
	D_CL,
	D_DL,
	D_BL,
	D_SPB,
	D_BPB,
	D_SIB,
	D_DIB,
	D_R8B,
	D_R9B,
	D_R10B,
	D_R11B,
	D_R12B,
	D_R13B,
	D_R14B,
	D_R15B,

	D_AX		= 16,
	D_CX,
	D_DX,
	D_BX,
	D_SP,
	D_BP,
	D_SI,
	D_DI,
	D_R8,
	D_R9,
	D_R10,
	D_R11,
	D_R12,
	D_R13,
	D_R14,
	D_R15,

	D_AH		= 32,
	D_CH,
	D_DH,
	D_BH,

	D_F0		= 36,

	D_M0		= 44,

	D_X0		= 52,
	D_X1,
	D_X2,
	D_X3,
	D_X4,
	D_X5,
	D_X6,
	D_X7,
	D_X8,
	D_X9,
	D_X10,
	D_X11,
	D_X12,
	D_X13,
	D_X14,
	D_X15,

	D_CS		= 68,
	D_SS,
	D_DS,
	D_ES,
	D_FS,
	D_GS,

	D_GDTR,		/* global descriptor table register */
	D_IDTR,		/* interrupt descriptor table register */
	D_LDTR,		/* local descriptor table register */
	D_MSW,		/* machine status word */
	D_TASK,		/* task register */

	D_CR		= 79,
	D_DR		= 95,
	D_TR		= 103,

	D_NONE		= 111,

	D_BRANCH	= 112,
	D_EXTERN	= 113,
	D_STATIC	= 114,
	D_AUTO		= 115,
	D_PARAM		= 116,
	D_CONST		= 117,
	D_FCONST	= 118,
	D_SCONST	= 119,
	D_ADDR		= 120,

	D_FILE,
	D_FILE1,

	D_INDIR,	/* additive */

	D_SIZE = D_INDIR + D_INDIR,	/* 6l internal */
	D_PCREL,
	D_TLS,

	T_TYPE		= 1<<0,
	T_INDEX		= 1<<1,
	T_OFFSET	= 1<<2,
	T_FCONST	= 1<<3,
	T_SYM		= 1<<4,
	T_SCONST	= 1<<5,
	T_64		= 1<<6,
	T_GOTYPE	= 1<<7,

	REGARG		= -1,
	REGRET		= D_AX,
	FREGRET		= D_X0,
	REGSP		= D_SP,
	REGTMP		= D_DI,
	REGEXT		= D_R15,	/* compiler allocates external registers R15 down */
	FREGMIN		= D_X0+5,	/* first register variable */
	FREGEXT		= D_X0+15	/* first external register */
};

/*
 * this is the ranlib header
 */
#define	SYMDEF	"__.GOSYMDEF"
