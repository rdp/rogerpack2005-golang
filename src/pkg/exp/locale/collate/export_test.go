// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collate

// Export for testing.

import (
	"exp/norm"
	"fmt"
)

type Weights struct {
	Primary, Secondary, Tertiary, Quaternary int
}

func W(ce ...int) Weights {
	w := Weights{ce[0], defaultSecondary, defaultTertiary, 0}
	if len(ce) > 1 {
		w.Secondary = ce[1]
	}
	if len(ce) > 2 {
		w.Tertiary = ce[2]
	}
	if len(ce) > 3 {
		w.Quaternary = ce[3]
	}
	return w
}
func (w Weights) String() string {
	return fmt.Sprintf("[%d.%d.%d.%d]", w.Primary, w.Secondary, w.Tertiary, w.Quaternary)
}

type Table struct {
	t *table
	w []weights
}

func GetTable(c *Collator) *Table {
	return &Table{c.t, nil}
}

func convertToWeights(ws []weights) []Weights {
	out := make([]Weights, len(ws))
	for i, w := range ws {
		out[i] = Weights{int(w.primary), int(w.secondary), int(w.tertiary), int(w.quaternary)}
	}
	return out
}

func convertFromWeights(ws []Weights) []weights {
	out := make([]weights, len(ws))
	for i, w := range ws {
		out[i] = weights{uint32(w.Primary), uint16(w.Secondary), uint8(w.Tertiary), uint32(w.Quaternary)}
	}
	return out
}

func (t *Table) AppendNext(s []byte) ([]Weights, int) {
	w, n := t.t.appendNext(nil, s)
	return convertToWeights(w), n
}

func SetTop(c *Collator, top int) {
	c.variableTop = uint32(top)
}

func InitCollator(c *Collator) {
	c.Strength = Quaternary
	c.f = norm.NFD
	c.t.maxContractLen = 30
}

func GetColElems(c *Collator, buf *Buffer, str []byte) []Weights {
	buf.ResetKeys()
	InitCollator(c)
	c.getColElems(buf, str)
	return convertToWeights(buf.ce)
}

func ProcessWeights(h AlternateHandling, top int, w []Weights) []Weights {
	in := convertFromWeights(w)
	processWeights(h, uint32(top), in)
	return convertToWeights(in)
}

func KeyFromElems(c *Collator, buf *Buffer, w []Weights) []byte {
	k := len(buf.key)
	c.keyFromElems(buf, convertFromWeights(w))
	return buf.key[k:]
}
