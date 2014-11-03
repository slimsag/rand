package spine

type SkeletonData struct {
	Bones       []*Bone
	Slots       []*Slot
	Skins       []*Skin
	Events      []*Event
	Animations  []*Animation
	DefaultSKin *Skin
}

// FindBone finds and returns the bone with the given name, or returns nil.
func (s *SkeletonData) FindBone(name string) *Bone {
	for _, bone := range s.Bones {
		if bone.Name == name {
			return bone
		}
	}
	return nil
}

// FindBoneIndex finds and returns the index for the bone with the given name
// or returns -1.
func (s *SkeletonData) FindBoneIndex(name string) int {
	for i, bone := range s.Bones {
		if bone.Name == name {
			return i
		}
	}
	return -1
}

// FindSlot finds and returns the slot with the given name, or returns nil.
func (s *SkeletonData) FindSlot(name string) *Slot {
	for _, slot := range s.Slot {
		if slot.Name == name {
			return slot
		}
	}
	return nil
}

// FindSlotIndex finds and returns the index for the slot with the given name
// or returns -1.
func (s *SkeletonData) FindSlotIndex(name string) int {
	for i, slot := range s.Slots {
		if slot.Name == name {
			return i
		}
	}
	return -1
}

// FindSkin finds and returns the skin with the given name, or returns nil.
func (s *SkeletonData) FindSkin(name string) *Skin {
	for _, skin := range s.Skins {
		if skin.Name == name {
			return skin
		}
	}
	return nil
}

// FindEvent finds and returns the event with the given name, or returns nil.
func (s *SkeletonData) FindEvent(name string) *Event {
	for _, event := range s.Events {
		if event.Name == name {
			return event
		}
	}
	return nil
}

// FindAnimation finds and returns the animation with the given name, or
// returns nil.
func (s *SkeletonData) FindAnimation(name string) *Animation {
	for _, animation := range s.Animations {
		if animation.Name == name {
			return animation
		}
	}
	return nil
}
