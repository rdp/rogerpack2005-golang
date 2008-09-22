// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "runtime.h"

typedef	struct	Sigt	Sigt;
typedef	struct	Sigi	Sigi;
typedef	struct	Map	Map;

struct	Sigt
{
	byte*	name;
	uint32	hash;
	void	(*fun)(void);
};

struct	Sigi
{
	byte*	name;
	uint32	hash;
	uint32	offset;
};

struct	Map
{
	Sigi*	sigi;
	Sigt*	sigt;
	Map*	link;
	int32	bad;
	int32	unused;
	void	(*fun[])(void);
};

static	Map*	hash[1009];
static	int32	debug	= 0;

static Map*
hashmap(Sigi *si, Sigt *ss)
{
	int32 ns, ni;
	uint32 ihash, h;
	byte *sname, *iname;
	Map *m;

	h = ((uint32)(uint64)si + (uint32)(uint64)ss) % nelem(hash);
	for(m=hash[h]; m!=nil; m=m->link) {
		if(m->sigi == si && m->sigt == ss) {
			if(m->bad) {
				throw("bad hashmap");
				m = nil;
			}
			// prints("old hashmap\n");
			return m;
		}
	}

	ni = si[0].offset;	// first word has size
	m = mal(sizeof(*m) + ni*sizeof(m->fun[0]));
	m->sigi = si;
	m->sigt = ss;

	ni = 1;			// skip first word
	ns = 0;

loop1:
	// pick up next name from
	// interface signature
	iname = si[ni].name;
	if(iname == nil) {
		m->link = hash[h];
		hash[h] = m;
		// prints("new hashmap\n");
		return m;
	}
	ihash = si[ni].hash;

loop2:
	// pick up and comapre next name
	// from structure signature
	sname = ss[ns].name;
	if(sname == nil) {
		prints((int8*)iname);
		prints(": ");
		throw("hashmap: failed to find method");
		m->bad = 1;
		m->link = hash[h];
		hash[h] = m;
		return nil;
	}

	if(ihash != ss[ns].hash ||
	   strcmp(sname, iname) != 0) {
		ns++;
		goto loop2;
	}

	m->fun[si[ni].offset] = ss[ns].fun;
	ni++;
	goto loop1;
}

static void
printsigi(Sigi *si)
{
	sys·printpointer(si);
}

static void
printsigt(Sigt *st)
{
	sys·printpointer(st);
}

static void
printiface(Map *im, void *it)
{
	prints("(");
	sys·printpointer(im);
	prints(",");
	sys·printpointer(it);
	prints(")");
}

// ifaceT2I(sigi *byte, sigt *byte, elem any) (ret interface{});
void
sys·ifaceT2I(Sigi *si, Sigt *st, void *elem, Map *retim, void *retit)
{

	if(debug) {
		prints("T2I sigi=");
		printsigi(si);
		prints(" sigt=");
		printsigt(st);
		prints(" elem=");
		sys·printpointer(elem);
		prints("\n");
	}

	retim = hashmap(si, st);
	retit = elem;

	if(debug) {
		prints("T2I ret=");
		printiface(retim, retit);
		prints("\n");
	}

	FLUSH(&retim);
}

// ifaceI2T(sigt *byte, iface interface{}) (ret any);
void
sys·ifaceI2T(Sigt *st, Map *im, void *it, void *ret)
{

	if(debug) {
		prints("I2T sigt=");
		printsigt(st);
		prints(" iface=");
		printiface(im, it);
		prints("\n");
	}

	if(im == nil)
		throw("ifaceI2T: nil map");

	if(im->sigt != st)
		throw("ifaceI2T: wrong type");

	ret = it;
	if(debug) {
		prints("I2T ret=");
		sys·printpointer(ret);
		prints("\n");
	}

	FLUSH(&ret);
}

// ifaceI2I(sigi *byte, iface interface{}) (ret interface{});
void
sys·ifaceI2I(Sigi *si, Map *im, void *it, Map *retim, void *retit)
{

	if(debug) {
		prints("I2I sigi=");
		sys·printpointer(si);
		prints(" iface=");
		printiface(im, it);
		prints("\n");
	}

	if(im == nil) {
		throw("ifaceI2I: nil map");
		return;
	}

	retit = it;
	retim = im;
	if(im->sigi != si)
		retim = hashmap(si, im->sigt);

	if(debug) {
		prints("I2I ret=");
		printiface(retim, retit);
		prints("\n");
	}

	FLUSH(&retim);
}
