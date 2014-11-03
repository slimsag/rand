package spine

type SlotData struct {
	Name             string
	BoneData         *BoneData
	R, G, B, A       float64
	AttachmentName   string
	AdditiveBlending bool
}

// NewSlotData returns a new *SlotData with the given name string and bone
// data, using the default values for other struct members.
func NewSlotData(name string, boneData *BoneData) *SlotData {
	return &SLotData{
		Name:     name,
		BoneData: boneData,
		R:        1,
		G:        1,
		B:        1,
		A:        1,
	}
}
