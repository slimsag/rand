package spine

import (
	"image/color"
	"io/ioutil"
	"testing"
)

func TestColorConversion(t *testing.T) {
	if (hexToRGBA("FF000000") != color.RGBA{255, 0, 0, 0}) {
		t.Fatal("Invalid (Red) Hex->RGBA conversion")
	}
	if (hexToRGBA("00FF0000") != color.RGBA{0, 255, 0, 0}) {
		t.Fatal("Invalid (Green) Hex->RGBA conversion")
	}
	if (hexToRGBA("0000FF00") != color.RGBA{0, 0, 255, 0}) {
		t.Fatal("Invalid (Blue) Hex->RGBA conversion")
	}
	if (hexToRGBA("000000FF") != color.RGBA{0, 0, 0, 255}) {
		t.Fatal("Invalid (Alpha) Hex->RGBA conversion")
	}
	if (hexToRGBA("") != color.RGBA{255, 255, 255, 255}) {
		t.Fatal("Invalid (empty string) Hex->RGBA conversion")
	}
}

func testSpineJSON(t *testing.T, filepath string) *Skeleton {
	// Read file data.
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
	}

	// Load a skeleton from a JSON file.
	skeleton, err := LoadJSON(data)
	if err != nil {
		t.Fatal(err)
	}

	// Check skeleton's bones.
	someParent := false
	for _, b := range skeleton.Bones {
		if b.Parent != nil {
			someParent = true
		}
	}
	if !someParent {
		t.Fatal("Skeleton bones are missing parent bones.")
	}

	// Check skeleton's slots.
	someBone := false
	for _, s := range skeleton.Slots {
		if s.Bone != nil {
			someBone = true
		}
	}
	if !someBone {
		t.Fatal("Skeleton's slots are missing their bones")
	}

	// Check skeleton's skins.
	d := skeleton.Skins["default"]
	for slotName, attachments := range d {
		t.Log(slotName)
		for name, attachment := range attachments {
			if attachment == nil {
				t.Log("Attachment name:", name)
				t.Fatal("Not loaded properly.")
			}
			t.Log("    ", name, attachment)
		}
	}
	return skeleton
}

func TestSpineJSONGoblin(t *testing.T) {
	skeleton := testSpineJSON(t, "goblins.json")
	_ = skeleton.Skins["default"]["eyes"]["eyes-closed"]
}

func TestSpineJSONSpineboy(t *testing.T) {
	skeleton := testSpineJSON(t, "spineboy.json")
	_ = skeleton.Skins["default"]["eyes"]["eyes-closed"]
}

func TestSpineJSONExample(t *testing.T) {
	skeleton := testSpineJSON(t, "example.json")
	_ = skeleton.Skins["default"]["eyes"]["eyes-closed"]
}
