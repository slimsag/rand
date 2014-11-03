package spine

import "math"

var BoneYDown bool

type Bone struct {
	Data   *BoneData
	Parent *Bone

	X, Y           float64
	Rotation       float64
	ScaleX, ScaleY float64

	M00, M01, WorldX         float64
	M10, M11, WorldY         float64
	WorldRotation            float64
	WorldScaleX, WorldScaleY float64
}

// NewBone returns a new *Bone with the given bone data and parent, using the
// default values for other struct members.
func NewBone(boneData *BoneData, parent *Bone) *Bone {
	b := &Bone{
		ScaleX:      1,
		ScaleY:      1,
		WorldScaleX: 1,
		WorldScaleY: 1,
		Data:        boneData,
		Parent:      parent,
	}
	b.SetToSetupPose()
	return b
}

func (b *Bone) SetToSetupPose() {
	b.X = b.Data.X
	b.Y = b.Data.Y
	b.Rotation = b.Data.Rotation
	b.ScaleX = b.Data.ScaleX
	b.ScaleY = b.Data.ScaleY
}

func (b *Bone) UpdateWorldTransform(flipX, flipY bool) {
	if b.Parent != nil {
		b.WorldX = b.X*b.Parent.M00 + b.Y*b.Parent.M01 + b.Parent.WorldX
		b.WorldY = b.X*b.Parent.M10 + b.Y*b.Parent.M11 + b.Parent.WorldY
		if b.Data.InheritScale {
			b.WorldScaleX = b.Parent.WorldScaleX * b.ScaleX
			b.WorldScaleY = b.Parent.WorldScaleY * b.ScaleY
		} else {
			b.WorldScaleX = b.ScaleX
			b.WorldScaleY = b.ScaleY
		}
		if b.Data.InheritRotation {
			b.WorldRotation = b.Parent.WorldRotation + b.Rotation
		} else {
			b.WorldRotation = b.Rotation
		}
	} else {
		if flipX {
			b.WorldX = -b.X
		} else {
			b.WorldX = b.X
		}
		if flipY != BoneYDown {
			b.WorldY = -b.Y
		} else {
			b.WorldY = b.Y
		}
		b.WorldScaleX = b.ScaleX
		b.WorldScaleY = b.ScaleY
		b.WorldRotation = b.Rotation
	}
	radians := b.WorldRotation * math.Pi / 180
	cos := math.Cos(radians)
	sin := math.Sin(radians)
	b.M00 = cos * b.WorldScaleX
	b.M10 = sin * b.WorldScaleX
	b.M01 = -sin * b.WorldScaleY
	b.M11 = cos * b.WorldScaleY
	if flipX {
		b.M00 = -b.M00
		b.M01 = -b.M01
	}
	if flipY != BoneYDown {
		b.M10 = -b.M10
		b.M11 = -b.M11
	}
}
