// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgfx

import (
	"sort"
)

// hSweeper sweeps horizontally from top-to-bottom. It effectively just sorts
// the input points in this order but is representitive of important concept.
type hSweeper Polygon

// sweep is short-handed for sort.Sort(s)
func (s hSweeper) sweep()             { sort.Sort(s) }
func (s hSweeper) Len() int           { return len(s) }
func (s hSweeper) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s hSweeper) Less(i, j int) bool { return s[i].X > s[j].X }

// vSweeper sweeps vertically from left-to-right. It effectively just sorts the
// input points in this order but is representitive of important concept.
type vSweeper Polygon

// sweep is short-handed for sort.Sort(s)
func (s vSweeper) sweep()             { sort.Sort(s) }
func (s vSweeper) Len() int           { return len(s) }
func (s vSweeper) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s vSweeper) Less(i, j int) bool { return s[i].Y > s[j].Y }
