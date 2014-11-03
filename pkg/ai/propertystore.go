package ai

/*
#include "assimp/cimport.h"

void azul_ai_log_stream(char*, char*);
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type propertyStore struct {
	c *C.struct_aiPropertyStore
}

// CreatePropertyStore creates an empty property store. Property stores are
// used to collect import settings.
//
// Returns a new property store. Property stores need to be manually destroyed
// using the ReleasePropertyStore() function.
//
// Reference to the property store must be held for as long as it must remain
// valid on the C side.
func createPropertyStore() *propertyStore {
	g := new(propertyStore)
	g.c = C.aiCreatePropertyStore()
	if g.c == nil {
		return nil
	}
	runtime.SetFinalizer(g, func(f *propertyStore) {
		C.aiReleasePropertyStore(f.c)
	})
	return g
}

// Set sets the named property to the given int, float32, or string value.
//
// Because this is a wrapper to the C-version of
// Assimp::Importer::SetPropertyInteger properties are always shared by all
// imports. It is not possible to specify them per import.
func (p *propertyStore) Set(name string, value interface{}) {
	switch t := value.(type) {
	case int:
		setImportPropertyInteger(p, name, t)
	case float32:
		setImportPropertyFloat(p, name, t)
	case string:
		setImportPropertyString(p, name, t)
	default:
		panic("Invalid property value type.")
	}
}

// SetImportPropertyInteger sets an integer property.
//
// This is the C-version of #Assimp::Importer::SetPropertyInteger(). In the C
// interface, properties are always shared by all imports. It is not possible
// to specify them per import.
//
// The szName parameter specifies the name of the configuration property to be
// set. All supported public properties are defined in the config.h header
// file.
//jhk
// The value paremeter specifies the new value for the property.
func setImportPropertyInteger(store *propertyStore, szName string, value int) {
	C.aiSetImportPropertyInteger(
		store.c,
		C.CString(szName),
		C.int(value),
	)
}

// SetImportPropertyFloat sets an float property.
//
// This is the C-version of #Assimp::Importer::SetPropertyFloat(). In the C
// interface, properties are always shared by all imports. It is not possible
// to specify them per import.
//
// The szName parameter specifies the name of the configuration property to be
// set. All supported public properties are defined in the config.h header
// file.
//
// The value paremeter specifies the new value for the property.
func setImportPropertyFloat(store *propertyStore, szName string, value float32) {
	C.aiSetImportPropertyFloat(
		(*C.struct_aiPropertyStore)(unsafe.Pointer(store)),
		C.CString(szName),
		C.float(value),
	)
}

// SetImportPropertyString sets an string property.
//
// This is the C-version of #Assimp::Importer::SetPropertyString(). In the C
// interface, properties are always shared by all imports. It is not possible
// to specify them per import.
//
// The szName parameter specifies the name of the configuration property to be
// set. All supported public properties are defined in the config.h header
// file.
//
// The value paremeter specifies the new value for the property.
func setImportPropertyString(store *propertyStore, szName string, value string) {
	C.aiSetImportPropertyString(
		(*C.struct_aiPropertyStore)(unsafe.Pointer(store)),
		C.CString(szName),
		aiString(value),
	)
}

/*
type aiLogStreamWrap struct {
	c *C.struct_aiLogStream
	l *log.Logger
}

//export azul_ai_log_stream
func azul_ai_log_stream(msg, user *C.char) {
	w := (*aiLogStreamWrap)(unsafe.Pointer(user))
	w.l.Print(C.GoString(msg))
}

func aiLogStream(l *log.Logger) *aiLogStreamWrap {
	w := &aiLogStreamWrap{
		l: l,
	}
	w.c = new(C.struct_aiLogStream)
	w.c.callback = (C.aiLogStreamCallback)(C.azul_ai_log_stream)
	w.c.user = (*C.char)(unsafe.Pointer(w))
	return w
}
*/
