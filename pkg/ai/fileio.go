// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ai

/*
#include "assimp/cfileio.h"

size_t azul_ai_file_write(struct aiFile*, char*, size_t, size_t);
size_t azul_ai_file_read(struct aiFile*, char*, size_t, size_t);
size_t azul_ai_file_tell(struct aiFile*);
size_t azul_ai_file_size(struct aiFile*);
void azul_ai_file_flush(struct aiFile*);
aiReturn azul_ai_file_seek(struct aiFile*, size_t, aiOrigin);

struct aiFile* azul_ai_fileio_open(struct aiFileIO* fio, char* path, char* mode);
void azul_ai_fileio_close(struct aiFileIO* fio, C_STRUCT aiFile* path);
*/
import "C"

import (
	"io"
	"log"
	"os"
	"reflect"
	"unsafe"
)

type File interface {
	io.Writer
	io.Reader
	io.Seeker
	io.Closer

	// e.g. for an actual Go file one could simply:
	//  fi, err := file.Stat()
	//  if err != nil {
	//      return 0, err
	//  }
	//  return fi.Size(), nil
	Size() (int64, error)

	Sync() error
}

//export azul_ai_file_write
func azul_ai_file_write(f *C.struct_aiFile, ptr *C.char, size, count C.size_t) C.size_t {
	goFile := *((*File)(unsafe.Pointer(f.UserData)))
	var buf []byte
	h := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	h.Data = uintptr(unsafe.Pointer(ptr))
	h.Len = int(size * count)
	h.Cap = int(size * count)
	bytesWrote, err := goFile.Write(buf)
	if err != nil {
		log.Println("assimp:", err)
	}
	return C.size_t(bytesWrote)
}

//export azul_ai_file_read
func azul_ai_file_read(f *C.struct_aiFile, ptr *C.char, size, count C.size_t) C.size_t {
	goFile := *((*File)(unsafe.Pointer(f.UserData)))
	var buf []byte
	h := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	h.Data = uintptr(unsafe.Pointer(ptr))
	h.Len = int(size * count)
	h.Cap = int(size * count)
	bytesRead, err := goFile.Read(buf)
	if err != nil {
		log.Println("assimp:", err)
	}
	return C.size_t(bytesRead)
}

//export azul_ai_file_tell
func azul_ai_file_tell(f *C.struct_aiFile) C.size_t {
	goFile := *((*File)(unsafe.Pointer(f.UserData)))
	offset, err := goFile.Seek(0, os.SEEK_CUR)
	if err != nil {
		log.Println("assimp:", err)
	}
	return C.size_t(offset)
}

//export azul_ai_file_size
func azul_ai_file_size(f *C.struct_aiFile) C.size_t {
	goFile := *((*File)(unsafe.Pointer(f.UserData)))
	size, err := goFile.Size()
	if err != nil {
		log.Println("assimp:", err)
	}
	return C.size_t(size)
}

//export azul_ai_file_flush
func azul_ai_file_flush(f *C.struct_aiFile) {
	goFile := *((*File)(unsafe.Pointer(f.UserData)))
	err := goFile.Sync()
	if err != nil {
		log.Println("assimp:", err)
	}
}

//export azul_ai_file_seek
func azul_ai_file_seek(f *C.struct_aiFile, offset C.size_t, corigin C.aiOrigin) C.aiReturn {
	var origin int
	switch corigin {
	case C.aiOrigin_SET:
		origin = os.SEEK_SET
	case C.aiOrigin_CUR:
		origin = os.SEEK_CUR
	case C.aiOrigin_END:
		origin = os.SEEK_END
	}

	goFile := *((*File)(unsafe.Pointer(f.UserData)))
	_, err := goFile.Seek(int64(offset), origin)
	if err != nil {
		log.Println("assimp:", err)
		return C.aiReturn_FAILURE
	}
	return C.aiReturn_SUCCESS
}

type FileIO interface {
	// Open should open the given file path using the given mode string (e.g.
	// "rb" for read-binary) and return it.
	Open(path, mode string) File
}

//export azul_ai_fileio_open
func azul_ai_fileio_open(fio *C.struct_aiFileIO, path, mode *C.char) *C.struct_aiFile {
	f := (*fileIOWrap)(unsafe.Pointer(fio.UserData))
	goFile := f.Open(C.GoString(path), C.GoString(mode))
	f.files[goFile] = true
	c := new(C.struct_aiFile)
	c.ReadProc = (C.aiFileReadProc)(C.azul_ai_file_read)
	c.WriteProc = (C.aiFileWriteProc)(C.azul_ai_file_write)
	c.TellProc = (C.aiFileTellProc)(C.azul_ai_file_tell)
	c.FileSizeProc = (C.aiFileTellProc)(C.azul_ai_file_size)
	c.FlushProc = (C.aiFileFlushProc)(C.azul_ai_file_flush)
	c.SeekProc = (C.aiFileSeek)(C.azul_ai_file_seek)
	c.UserData = (C.aiUserData)(unsafe.Pointer(&goFile))
	return c
}

//export azul_ai_fileio_close
func azul_ai_fileio_close(fio *C.struct_aiFileIO, f *C.struct_aiFile) {
	fw := (*fileIOWrap)(unsafe.Pointer(fio.UserData))
	goFile := *((*File)(unsafe.Pointer(f.UserData)))
	goFile.Close()
	delete(fw.files, goFile)
	return
}

type fileIOWrap struct {
	FileIO
	c     *C.struct_aiFileIO
	files map[File]bool
}

// aiFileIO wraps a Go FileIO interface into a C one and returns it. You must
// hold reference to the returned Go type for as long as the C types remain
// needed.
func aiFileIO(g FileIO) *fileIOWrap {
	f := &fileIOWrap{
		FileIO: g,
		files:  make(map[File]bool),
	}
	f.c = new(C.struct_aiFileIO)
	f.c.OpenProc = (C.aiFileOpenProc)(C.azul_ai_fileio_open)
	f.c.CloseProc = (C.aiFileCloseProc)(C.azul_ai_fileio_close)
	f.c.UserData = (C.aiUserData)(unsafe.Pointer(f))
	return f
}
