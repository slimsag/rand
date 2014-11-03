package spine

import (
	"fmt"
	"image/color"
	"strconv"
)

// hexToRGBA parses a RRGGBBAA string and returns a color.RGBA. If the input
// string is invalid color.RGBA{255, 255, 255, 255} will be returned.
func hexToRGBA(c string) color.RGBA {
	// Spine uses RRGGBBAA color strings with no prefixed # hash sign.
	// If an invalid length value then simply return
	if len(c) != 8 {
		return color.RGBA{255, 255, 255, 255}
	}

	// Parse color string.
	rgba, err := strconv.ParseUint(c, 16, 48)
	if err != nil {
		return color.RGBA{255, 255, 255, 255}
	}
	return color.RGBA{
		R: uint8(rgba >> 24),
		G: uint8(rgba >> 16),
		B: uint8(rgba >> 8),
		A: uint8(rgba),
	}
}

// Bone describes a single bone in it's setup pose.
type Bone struct {
	// Name of the bone, this is unique for the skeleton.
	Name string

	// The parent of the bone, nil if there is none.
	Parent *Bone

	// The length of the bone.
	Length float64

	// The X and Y position of the bone relative to the parent.
	X, Y float64

	// The X and Y scaling values, one by default.
	ScaleX, ScaleY float64

	// The rotation in degrees of the bone relative to the parent .
	Rotation float64
}

// String returns a string representation of this bone.
func (b *Bone) String() string {
	return fmt.Sprintf("Bone(%q, Length=%v, X=%v, Y=%v, ScaleX=%v, ScaleY=%v, Rotation=%v)", b.Name, b.Length, b.X, b.Y, b.ScaleX, b.ScaleY, b.Rotation)
}

// Slot represents a single slot where attachments can be assigned.
type Slot struct {
	// Name of the slot, this is unique for the skeleton.
	Name string

	// The bone of the slot.
	Bone *Bone

	// The color of the slot.
	Color color.RGBA

	// The attachment of the slot.
	Attachment string
}

// String returns a string representation of this slot.
func (s *Slot) String() string {
	return fmt.Sprintf("Slot(%q, Color=%v, Attachment=%q)", s.Name, s.Color, s.Attachment)
}

// Region represents a single image attachment.
type Region struct {
	// The name of the image region, either a key into a texture atlas or a
	// filepath.
	Name string

	// The position of the image relative to the slot's bone.
	X, Y float64

	// The scale of the image.
	ScaleX, ScaleY float64

	// Rotation of the image in degrees relative to the slot's bone.
	Rotation float64

	// The width and height of the image.
	Width, Height int
}

// String returns a string representation of this region.
func (r *Region) String() string {
	return fmt.Sprintf("Region(%q, Pos=[%v,%v], Scale=[%v,%v], Rotation=%v Size=[%v,%v])", r.Name, r.X, r.Y, r.ScaleX, r.ScaleY, r.Rotation, r.Width, r.Height)
}

// RegionSequence represents a single image sequence attachment.
type RegionSequence struct {
	// The name of the image region sequence, either a key into a texture atlas
	// or a filepath.
	Name string

	// The position of the image relative to the slot's bone.
	X, Y float64

	// The scale of the image.
	ScaleX, ScaleY float64

	// Rotation of the image in degrees relative to the slot's bone.
	Rotation float64

	// The width and height of the image.
	Width, Height int

	// The frame rate at which to display the sequence.
	FPS float64

	// The display mode of the image region sequence.
	Mode Mode
}

// String returns a string representation of this region.
func (r *RegionSequence) String() string {
	return fmt.Sprintf("RegionSequence(%q, Pos=[%v,%v], Scale=[%v,%v], Rotation=%v, Size=[%v,%v], FPS=%v, Mode=%v)", r.Name, r.X, r.Y, r.ScaleX, r.ScaleY, r.Rotation, r.Width, r.Height, r.FPS, r.Mode)
}

// Skin describes attachments that can be assigned to each slot, it is a map of
// attachments by slot name to a map of attachments by attachment name. Each
// attachment is a interface value which will always be one of the following
// types (a type assertion should be used):
//  *Region
//  *RegionSequence
//  *BoundingBox
type Skin map[string]map[string]interface{}

// Mode represents a single region sequence display mode.
type Mode int

const (
	Forward Mode = iota
	Backward
	ForwardLoop
	BackwardLoop
	PingPong
	Random
)

// String returns a string representation of this region sequence display mode.
func (m Mode) String() string {
	switch m {
	case Forward:
		return "Forward"
	case Backward:
		return "Backward"
	case ForwardLoop:
		return "ForwardLoop"
	case BackwardLoop:
		return "BackwardLoop"
	case PingPong:
		return "PingPong"
	case Random:
		return "Random"
	}
	return fmt.Sprintf("Mode(%d)", m)
}

// Mesh represents a single mesh attachment.
type Mesh struct {
	Name                            string
	Vertices, Triangles, UVs, Edges []float64
	Hull                            float64
	Width, Height                   int
}

// String returns a string representation of this mesh.
func (m *Mesh) String() string {
	return fmt.Sprintf("Mesh(%q, %v Vertices, %v Triangles, %v UVs, %v Edges, Hull=%v, Size=[%v,%v])", m.Name, len(m.Vertices), len(m.Triangles), len(m.UVs), len(m.Edges), m.Hull, m.Width, m.Height)
}

// BoundingBox represents a single bounding box attachment.
type BoundingBox struct {
	Name     string
	Vertices []float64
}

// Event represents a single event.
type Event struct {
	Name   string
	Int    int
	Float  float64
	String string
}

// Skeleton represents a single Spine skeleton.
type Skeleton struct {
	// The bones for the setup pose of the skeleton.
	Bones []*Bone

	// The slots of the skeleton, describes the draw order of the available
	// slots where attachments can be assigned.
	Slots []*Slot

	// A map of skins by name.
	Skins map[string]Skin

	// A map of events by name that can be triggered during animations.
	Events map[string]Event
}
