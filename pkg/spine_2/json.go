package spine

import (
	"encoding/json"
	"fmt"
)

type jsonBone struct {
	Name           string
	Parent         string
	Length         float64
	X, Y           float64
	ScaleX, ScaleY float64
	Rotation       float64
}

type jsonSlot struct {
	Name       string
	Bone       string
	Color      string
	Attachment string
}

type jsonAttachment struct {
	Name                            string
	Type                            string
	X, Y                            float64
	ScaleX, ScaleY                  float64
	Rotation                        float64
	Width, Height                   int
	FPS                             float64
	Mode                            string
	Vertices, Triangles, UVs, Edges []float64
	Hull                            float64
}

type jsonEvent struct {
	Int    int
	Float  float64
	String string
}

type jsonFile struct {
	Bones  []jsonBone
	Slots  []jsonSlot
	Skins  map[string]map[string]map[string]jsonAttachment
	Events map[string]jsonEvent
}

func convertMode(mode string) Mode {
	switch mode {
	case "backward":
		return Backward
	case "forwardLoop":
		return ForwardLoop
	case "backwardLoop":
		return BackwardLoop
	case "pingPong":
		return PingPong
	case "random":
		return Random
	default:
		return Forward
	}
}

func LoadJSON(data []byte) (*Skeleton, error) {
	j := new(jsonFile)
	err := json.Unmarshal(data, j)
	if err != nil {
		return nil, err
	}

	s := &Skeleton{
		Bones: make([]*Bone, len(j.Bones)),
		Slots: make([]*Slot, len(j.Slots)),
		Skins: make(map[string]Skin, len(j.Skins)),
	}

	// Create a map of bones by name.
	bonesByName := make(map[string]*Bone, len(j.Bones))

	// Copy over bones.
	for i, jb := range j.Bones {
		b := &Bone{
			jb.Name,
			nil,
			//jb.Parent,
			jb.Length,
			jb.X,
			jb.Y,
			jb.ScaleX,
			jb.ScaleY,
			jb.Rotation,
		}
		bonesByName[b.Name] = b

		// Default scale values are one.
		if b.ScaleX == 0 {
			b.ScaleX = 1
		}
		if b.ScaleY == 0 {
			b.ScaleY = 1
		}

		s.Bones[i] = b
	}

	// Map parents to bones.
	for i, b := range s.Bones {
		parentName := j.Bones[i].Parent
		if len(parentName) > 0 {
			// Search for parent bone.
			parent, ok := bonesByName[parentName]
			if ok {
				b.Parent = parent
			}
		}
	}

	// Copy over slots.
	for i, jslot := range j.Slots {
		bone, _ := bonesByName[jslot.Bone]
		s.Slots[i] = &Slot{
			jslot.Name,
			bone,
			hexToRGBA(jslot.Color),
			jslot.Attachment,
		}
	}

	// Copy over skins.
	for skinName, js := range j.Skins {
		skin := make(Skin, len(js))
		for slotName, jsAttachments := range js {
			attachments := make(map[string]interface{}, len(jsAttachments))
			for jsAttachmentName, ja := range jsAttachments {
				// Assume name of slot if attachment name is ommited, according to docs.
				if len(ja.Name) == 0 {
					ja.Name = jsAttachmentName
				}

				// Scale is one if ommited.
				if ja.ScaleX == 0 {
					ja.ScaleX = 1
				}
				if ja.ScaleY == 0 {
					ja.ScaleY = 1
				}

				// Create attachment based on type string.
				var attachment interface{}
				switch ja.Type {
				case "", "region":
					attachment = &Region{
						ja.Name,
						ja.X,
						ja.Y,
						ja.ScaleX,
						ja.ScaleY,
						ja.Rotation,
						ja.Width,
						ja.Height,
					}
				case "regionsequence":
					attachment = &RegionSequence{
						ja.Name,
						ja.X,
						ja.Y,
						ja.ScaleX,
						ja.ScaleY,
						ja.Rotation,
						ja.Width,
						ja.Height,
						ja.FPS,
						convertMode(ja.Mode),
					}
				case "boundingbox":
					attachment = &BoundingBox{
						ja.Name,
						ja.Vertices,
					}
				case "mesh":
					attachment = &Mesh{
						ja.Name,
						ja.Vertices,
						ja.Triangles,
						ja.UVs,
						ja.Edges,
						ja.Hull,
						ja.Width,
						ja.Height,
					}
				}

				// Insert attachment.
				attachments[jsAttachmentName] = attachment
			}
			skin[slotName] = attachments
		}
		s.Skins[skinName] = skin
	}

	// Copy over events.
	s.Events = make(map[string]Event, len(j.Events))
	for name, ev := range j.Events {
		s.Events[name] = Event{
			name,
			ev.Int,
			ev.Float,
			ev.String,
		}
	}

	fmt.Println(j.Events)
	_ = fmt.Println
	return s, nil
}
