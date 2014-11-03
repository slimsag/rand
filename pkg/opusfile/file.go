package opusfile

/*
#include <opusfile.h>
*/
import "C"

import (
	"errors"
	"io"
	"unsafe"
)

/*
read() must be implemented

seek and tell may be NULL, or may always return -1 to indicate a source is unseekable.
but if seek() is implemented and succeeds on a particular source, then tell() must also.

close() may be NULL, bit if it is not, it will be called when the file is destroyed.
   -> op_free()

error is
   <dt>#OP_EREAD</dt>
   <dd>An underlying read, seek, or tell operation
    failed when it should have succeeded, or we failed
    to find data in the stream we had seen before.</dd>
   <dt>#OP_EFAULT</dt>
   <dd>There was a memory allocation failure, or an
    internal library error.</dd>
   <dt>#OP_EIMPL</dt>
   <dd>The stream used a feature that is not
    implemented, such as an unsupported channel
    family.</dd>
   <dt>#OP_EINVAL</dt>
   <dd><code><a href="#op_seek_func">seek()</a></code>
    was implemented and succeeded on this source, but
    <code><a href="#op_tell_func">tell()</a></code>
    did not, or the starting position indicator was
    not equal to \a _initial_bytes.</dd>
   <dt>#OP_ENOTFORMAT</dt>
   <dd>The stream contained a link that did not have
    any logical Opus streams in it.</dd>
   <dt>#OP_EBADHEADER</dt>
   <dd>A required header packet was not properly
    formatted, contained illegal values, or was missing
    altogether.</dd>
   <dt>#OP_EVERSION</dt>
   <dd>An ID header contained an unrecognized version
    number.</dd>
   <dt>#OP_EBADLINK</dt>
   <dd>We failed to find data we had seen before after
    seeking.</dd>
   <dt>#OP_EBADTIMESTAMP</dt>
   <dd>The first or last timestamp in a link failed
    basic validity checks.</dd>
 </dl>


OggOpusFile *op_open_callbacks(void *_source,
 const OpusFileCallbacks *_cb,const unsigned char *_initial_data,
 size_t _initial_bytes,int *_error);
*/

var opusFileCallbacks *C.OpusFileCallbacks

func init() {
	opusFileCallbacks = new(C.OpusFileCallbacks)
	//azul3d_opus_file_read,
	//azul3d_opus_file_seek,
	//azul3d_opus_file_tell,
	//azul3d_opus_file_close,
}

type FileInterface interface {
	io.Reader
	io.Seeker
	io.Closer
}

type File struct {
	fp   *FileInterface
	cptr *C.OggOpusFile
}

func Open(f FileInterface) (*File, error) {
	fp := &f
	var err C.int
	cptr := C.op_open_callbacks(
		unsafe.Pointer(fp),
		opusFileCallbacks,
		nil,
		0,
		&err,
	)
	if cptr == nil {
		f.Close()
		return nil, errors.New("FIXME")
	}
	return &File{
		fp:   fp,
		cptr: cptr,
	}, nil
}
