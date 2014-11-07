// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"sort"
)

// hSweeper sweeps horizontally from top-to-bottom. It effectively just sorts
// the input points in this order but is representitive of important concept.
type hSweeper struct {
	points []Point
}

// sweep is short-handed for sort.Sort(s)
func (s *hSweeper) sweep()             { sort.Sort(s) }
func (s *hSweeper) Len() int           { return len(s.points) }
func (s *hSweeper) Swap(i, j int)      { s.points[i], s.points[j] = s.points[j], s.points[i] }
func (s *hSweeper) Less(i, j int) bool { return s.points[i].X > s.points[j].X }

// vSweeper sweeps vertically from left-to-right. It effectively just sorts the
// input points in this order but is representitive of important concept.
type vSweeper struct {
	points []Point
}

// sweep is short-handed for sort.Sort(s)
func (s *vSweeper) sweep()             { sort.Sort(s) }
func (s *vSweeper) Len() int           { return len(s.points) }
func (s *vSweeper) Swap(i, j int)      { s.points[i], s.points[j] = s.points[j], s.points[i] }
func (s *vSweeper) Less(i, j int) bool { return s.points[i].Y > s.points[j].Y }
