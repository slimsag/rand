package spine

import (
	"math"
)

type Attachment interface {
	AttachmentName() string
}

type RegionAttachment struct {
	Name           string
	X, Y           float64
	Rotation       float64
	ScaleX, ScaleY float64
	Width, Height  float64
	Offset         [8]float64
	UVs            [8]float64

	RendererObject                            interface{}
	RegionOffsetX, RegionOffsetY              float64
	RegionWidth, RegionHeight                 float64
	RegionOriginalWidth, RegionOriginalHeight float64
}

// Implements Attachment interface.
func (r *RegionAttachment) AttachmentName() string {
	return r.Name
}

// NewRegionAttachment returns a new *RegionAttachment with the given name,
// using the default values for other struct members.
func NewRegionAttachment(name string) *RegionAttachment {
	return &RegionAttachment{
		Name:   name,
		ScaleX: 1,
		ScaleY: 1,
	}
}

func (r *RegionAttachment) SetUVs(u, v, u2, v2, rotate bool) {
	if rotate {
		r.UVS[2] = u
		r.UVS[3] = v2
		r.UVS[4] = u
		r.UVS[5] = v
		r.UVS[6] = u2
		r.UVS[7] = v
		r.UVS[0] = u2
		r.UVS[1] = v2
	} else {
		r.UVS[0] = u
		r.UVS[1] = v2
		r.UVS[2] = u
		r.UVS[3] = v
		r.UVS[4] = u2
		r.UVS[5] = v
		r.UVS[6] = u2
		r.UVS[7] = v2
	}
}

func (r *RegionAttachment) UpdateOffset() {
	regionScaleX := r.Width / r.RegionOriginalWidth * r.ScaleX
	regionScaleY := r.Height / r.RegionOriginalHeight * r.ScaleY
	localX := -r.Width/2*r.ScaleX + r.RegionOffsetX*regionScaleX
	localY := -r.Height/2*r.ScaleY + r.RegionOffsetY*regionScaleY
	localX2 := localX + r.RegionWidth*regionScaleX
	localY2 := localY + r.RegionHeight*regionScaleY
	radians := r.Rotation * math.PI / 180
	cos := math.Cos(radians)
	sin := math.Sin(radians)
	localXCos := localX*cos + r.X
	localXSin := localX * sin
	localYCos := localY*cos + r.Y
	localYSin := localY * sin
	localX2Cos := localX2*cos + r.X
	localX2Sin := localX2 * sin
	localY2Cos := localY2*cos + r.Y
	localY2Sin := localY2 * sin
	r.Offset[0] = localXCos - localYSin
	r.Offset[1] = localYCos + localXSin
	r.Offset[2] = localXCos - localY2Sin
	r.Offset[3] = localY2Cos + localXSin
	r.Offset[4] = localX2Cos - localY2Sin
	r.Offset[5] = localY2Cos + localX2Sin
	r.Offset[6] = localX2Cos - localYSin
	r.Offset[7] = localYCos + localX2Sin
}

func (r *RegionAttachment) ComputeVerticecs(x, y float64, bone *Bone) (vertices [8]float64) {
	x += bone.WorldX
	y += bone.WorldY
	m00 := bone.M00
	m01 = bone.M01
	m10 = bone.M10
	m11 = bone.M11
	vertices[0] = r.Offset[0]*m00 + r.Offset[1]*m01 + x
	vertices[1] = r.Offset[0]*m10 + r.Offset[1]*m11 + y
	vertices[2] = r.Offset[2]*m00 + r.Offset[3]*m01 + x
	vertices[3] = r.Offset[2]*m10 + r.Offset[3]*m11 + y
	vertices[4] = r.Offset[4]*m00 + r.Offset[5]*m01 + x
	vertices[5] = r.Offset[4]*m10 + r.Offset[5]*m11 + y
	vertices[6] = r.Offset[6]*m00 + r.Offset[7]*m01 + x
	vertices[7] = r.Offset[6]*m10 + r.Offset[7]*m11 + y
}

type BoundingBoxAttachment struct {
	Name     string
	Vertices []float64
}

// Implements Attachment interface.
func (b *BoundingBoxAttachment) AttachmentName() string {
	return b.Name
}

// NewBoundingBoxAttachment returns a new *BoundingBoxAttachment with the given
// name, using the default values for other struct members.
func NewBoundingBoxAttachment(name string) *BoundingBoxAttachment {
	return &BoundingBoxAttachment{
		Name: name,
	}
}

func (b *BoundingBoxAttachment) ComputeWorldVertices(x, y float64, bone *Bone) (worldVertices []float64) {
	x += bone.WorldX
	y += bone.WorldY
	m00 := bone.M00
	m01 := bone.M01
	m10 := bone.M10
	m11 := bone.M11
	worldVertices = make([]float64, len(b.Vertices))
	for i := 0; i < len(b.Vertices); i += 2 {
		px := b.Vertices[i]
		py := b.Vertices[i+1]
		worldVertices[i] = px*m00 + py*m01 + x
		worldVertices[i+1] = px*m10 + py*m11 + y
	}
	return
}
