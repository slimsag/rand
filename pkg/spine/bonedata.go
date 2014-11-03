package spine

type BoneData struct {
	Name                          string
	Parent                        *Bone
	Length                        float64
	X, Y                          float64
	Rotation                      float64
	ScaleX, ScaleY                float64
	InheritScale, InheritRotation bool
}

// NewBoneData returns a new *BoneData with the given name string and parent,
// using the default values for other struct members.
func NewBoneData(name string, parent *Bone) *BoneData {
	return &BoneData{
		Name:            name,
		Parent:          parent,
		ScaleX:          1.0,
		ScaleY:          1.0,
		InheritScale:    true,
		InheritRotation: true,
	}
}
