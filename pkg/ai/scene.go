package ai

/*
#include "assimp/scene.h"
*/
import "C"

import (
	"unsafe"
)

type SceneFlags int

const (
	// Specifies that the scene data structure that was imported is not complete.
	// This flag bypasses some internal validations and allows the import
	// of animation skeletons, material libraries or camera animation paths
	// using Assimp. Most applications won't support such data.
	SceneFlagsIncomplete SceneFlags = C.AI_SCENE_FLAGS_INCOMPLETE

	// This flag is set by the validation postprocess-step (aiPostProcess_ValidateDS)
	// if the validation is successful. In a validated scene you can be sure that
	// any cross references in the data structure (e.g. vertex indices) are valid.
	SceneFlagsValidated SceneFlags = C.AI_SCENE_FLAGS_VALIDATED

	// This flag is set by the validation postprocess-step (aiPostProcess_ValidateDS)
	// if the validation is successful but some issues have been found.
	// This can for example mean that a texture that does not exist is referenced
	// by a material or that the bone weights for a vertex don't sum to 1.0 ... .
	// In most cases you should still be able to use the import. This flag could
	// be useful for applications which don't capture Assimp's log output.
	SceneFlagsValidationWarning SceneFlags = C.AI_SCENE_FLAGS_VALIDATION_WARNING

	// This flag is currently only set by the aiProcess_JoinIdenticalVertices step.
	// It indicates that the vertices of the output meshes aren't in the internal
	// verbose format anymore. In the verbose format all vertices are unique,
	// no vertex is ever referenced by more than one face.
	SceneFlagsNonVerboseFormat SceneFlags = C.AI_SCENE_FLAGS_NON_VERBOSE_FORMAT

	// Denotes pure height-map terrain data. Pure terrains usually consist of quads,
	// sometimes triangles, in a regular grid. The x,y coordinates of all vertex
	// positions refer to the x,y coordinates on the terrain height map, the z-axis
	// stores the elevation at a specific point.
	//
	// TER (Terragen) and HMP (3D Game Studio) are height map formats.
	// @note Assimp is probably not the best choice for loading *huge* terrains -
	// fully triangulated data takes extremely much free store and should be avoided
	// as long as possible (typically you'll do the triangulation when you actually
	// need to render it).
	SceneFlagsTerrain SceneFlags = C.AI_SCENE_FLAGS_TERRAIN
)

// A node in the imported hierarchy.
//
// Each node has name, a parent node (except for the root node),
// a transformation relative to its parent and possibly several child nodes.
// Simple file formats don't support hierarchical structures - for these formats
// the imported scene does consist of only a single root node without children.
type Node struct {
	c *C.struct_aiNode
}

func aiNode(c *C.struct_aiNode) *Node {
	return &Node{
		c: c,
	}
}

// Any combination of the AI_SCENE_FLAGS_XXX flags. By default
// this value is 0, no flags are set. Most applications will
// want to reject all scenes with the AI_SCENE_FLAGS_INCOMPLETE
// bit set.
func (s *Scene) Flags() SceneFlags {
	return SceneFlags(s.access().mFlags)
}

// The root node of the hierarchy.
//
// There will always be at least the root node if the import
// was successful (and no special flags have been set).
// Presence of further nodes depends on the format and content
// of the imported file.
func (s *Scene) RootNode() *Node {
	return aiNode(s.c.mRootNode)
}

// Use the indices given in the *Node structure to accecss this slice. If the
// SceneFlagsIncomplete flag is not set there will always be at least ONE
// material.
func (s *Scene) Meshes() (ms []*Mesh) {
	for i := 0; i < s.c.mNumMeshes; i++ {
		meshes := uintptr(unsafe.Pointer(s.c.nMeshes))
		ms = append(ms, aiMesh(unsafe.Pointer(meshes+uintptr(i))))
	}
	return ms
}

// Use the index given in each *Mesh structure to access this slice. If the
// SceneFlagsIncomplete flag is not set there will always be at least ONE
// material.
func (s *Scene) Materials() (ms []*Material) {
	for i := 0; i < s.c.mNumMaterials; i++ {
		materials := uintptr(unsafe.Pointer(s.c.nMaterials))
		ms = append(ms, aiMaterial(unsafe.Pointer(materials+uintptr(i))))
	}
	return ms
}

// All animations imported from the given file are listed here.
func (s *Scene) Animations() (as []*Animation) {
	for i := 0; i < s.c.mNumAnimation; i++ {
		animations := uintptr(unsafe.Pointer(s.c.nAnimations))
		as = append(as, aiAnimation(unsafe.Pointer(animations+uintptr(i))))
	}
	return as
}

// Not many file formats embed their textures into the file. An example is
// Quake's MDL format (which is also used by some GameStudio versions).
func (s *Scene) Textures() (ts []*Texture) {
	for i := 0; i < s.c.mNumTextures; i++ {
		textures := uintptr(unsafe.Pointer(s.c.nTexture))
		ts = append(ts, aiTexture(unsafe.Pointer(textures+uintptr(i))))
	}
	return ts
}

// Since light sources are fully optional, in most cases the length of the
// slice will be zero.
func (s *Scene) Lights() (ls []*Light) {
	for i := 0; i < s.c.mNumLights; i++ {
		lights := uintptr(unsafe.Pointer(s.c.nLight))
		ls = append(ls, aiLight(unsafe.Pointer(lights+uintptr(i))))
	}
	return ls
}

// Since cameras are fully optional, in most cases the length of the
// slice will be zero. If there are any cameras the first one is the default
// camera view into the scene.
func (s *Scene) Cameras() (cs []*Camera) {
	for i := 0; i < s.c.mNumCameras; i++ {
		cameras := uintptr(unsafe.Pointer(s.c.nCameras))
		cs = append(cs, aiCamera(unsafe.Pointer(cameras+uintptr(i))))
	}
	return ls
}
