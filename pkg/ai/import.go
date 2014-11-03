package ai

/*
#include "assimp/cimport.h"
*/
import "C"

import (
	"errors"
	"runtime"
	"strings"
	"unsafe"
)

var Extensions map[string]bool

func init() {
	s := new(C.struct_aiString)
	C.aiGetExtensionList(s)
	exts := goString(s)
	for _, ext := range strings.Split(exts, ";") {
		Extensions[ext] = true
	}
}

type MemoryInfo struct {
	// Storage allocated for texture data
	Textures uint

	// Storage allocated for material data
	Materials uint

	// Storage allocated for mesh data
	Meshes uint

	// Storage allocated for node data
	Nodes uint

	// Storage allocated for animation data
	Animations uint

	// Storage allocated for camera data
	Cameras uint

	// Storage allocated for light data
	Lights uint

	// Total storage allocated for the full import.
	Total uint
}

type Scene struct {
	c     *C.struct_aiScene
	fs    *fileIOWrap
	props *propertyStore
}

// Get the approximated storage required by an imported asset
func (s *Scene) MemoryRequirements() MemoryInfo {
	mi := new(C.struct_aiMemoryInfo)
	C.aiGetMemoryRequirements(
		s.c,
		mi,
	)
	return MemoryInfo{
		uint(mi.textures),
		uint(mi.materials),
		uint(mi.meshes),
		uint(mi.nodes),
		uint(mi.animations),
		uint(mi.cameras),
		uint(mi.lights),
		uint(mi.total),
	}
}

// Apply post-processing to an already-imported scene.
//
// This is strictly equivalent to calling Import() with the
// same flags. However, you can use this separate function to inspect the imported
// scene first to fine-tune your post-processing setup.
//
// The flags parameter is a bitwise combination of the #aiPostProcessSteps flags.
//
// Returns a pointer to the post-processed data. Post processing is done in-place,
// meaning this is still the same #aiScene which you passed for pScene. However,
// _if_ post-processing failed, the scene could now be NULL. That's quite a rare
// case, post processing steps are not really designed to 'fail'. To be exact,
// the #aiProcess_ValidateDS flag is currently the only post processing step
// which can actually cause the scene to be reset to NULL.
func (s *Scene) ApplyPostProcessing(flags PostFlags) (ok bool) {
	ptr := C.aiApplyPostProcessing(
		s.c,
		C.uint(flags),
	)
	if ptr == nil {
		return false
	}
	return true
}

// Import reads the file using the FileIO interface and returns it's content.
//
// The flags parameter is a bitwise combination of pre-defined PostFlags which
// determine if and how post processing of the data should occur.
//
// The properties map defines various importer properties to utilize.
//
// If successfull the data is returned in the Scene structure which should be
// considered read-only (as the assimp library owns the memory), and a nil
// error is returned.
//
// If failure occurs a nil *Scene and a human-readable error is returned.
func Import(filepath string, flags PostFlags, fs FileIO, properties map[Prop]interface{}) (*Scene, error) {
	if len(filepath) == 0 {
		return nil, errors.New("empty file path")
	}
	s := &Scene{
		fs: aiFileIO(fs),
	}
	if len(properties) > 0 {
		s.props = createPropertyStore()
		for k, v := range properties {
			s.props.Set(string(k), v)
		}
		s.c = C.aiImportFileExWithProperties(
			C.CString(filepath),
			C.uint(flags),
			(*[0]byte)(unsafe.Pointer(s.fs.c)),
			s.props.c,
		)
	} else {
		s.c = C.aiImportFileEx(
			C.CString(filepath),
			C.uint(flags),
			(*[0]byte)(unsafe.Pointer(s.fs.c)),
		)
	}
	if s.c == nil {
		err := errors.New(C.GoString(C.aiGetErrorString()))
		return nil, err
	}

	runtime.SetFinalizer(s, func(f *Scene) {
		C.aiReleaseImport(f.c)
	})
	return s, nil
}
