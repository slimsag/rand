package ai

/*
#include "assimp/cimport.h"
#include <stdlib.h>
*/
import "C"

// aiString allocates and returns a new aiString, it must later be free'd.
func aiString(g string) *C.struct_aiString {
	s := (*C.struct_aiString)(C.calloc(1, C.size_t(len(g))))
	s.length = C.size_t(len(g))
	for i, c := range g {
		s.data[i] = C.char(c)
	}
	return s
}

func goString(c *C.struct_aiString) string {
	buf := make([]byte, 0, c.length)
	for i, c := range c.data {
		if c == 0 {
			break
		}
		buf[i] = byte(c)
	}
	return string(buf)
}
