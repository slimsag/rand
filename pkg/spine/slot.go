package spine

import "time"

type Slot struct {
	Data           *SlotData
	Skeleton       *Skeleton
	Bone           *Bone
	R, G, B, A     float32
	attachmentTime time.Time
	Attachment     Attachment
}

// NewSlot returns a new *Slot with the given slot data, skeleton, and bone,
// using the default values for other struct members.
func NewSlot(slotData *SlotData, skeleton *Skeleton, bone *Bone) *Slot {
	s := &Slot{
		Data:     slotData,
		Skeleton: skeleton,
		Bone:     bone,
		R:        1,
		G:        1,
		B:        1,
		A:        1,
	}
	s.SetToSetupPose()
	return s
}

func (s *Slot) SetAttachment(attachment Attachment) {
	s.Attachment = attachment
	s.attachmentTime = s.skeleton.Time
}

func (s *Slot) SetAttachmentTime(t time.Time) {
	s.attachmentTime = s.Skeleton.Time - t
}

func (s *Slot) AttachmentTime() time.Time {
	return s.Skeleton.Time - s.attachmentTime
}

func (s *Slot) SetToSetupPose() {
	s.R = s.Data.R
	s.G = s.Data.G
	s.B = s.Data.B
	s.A = s.Data.A

	for i, slotData := range s.Skeleton.Data.Slots {
		if slotData == s.Data {
			s.SetAttachment(s.Skeleton.AttachmentBySlotIndex(i, data.attachmentName))
		}
	}
}
