//go:build linux
// +build linux

package flite

// #cgo CFLAGS: -I /usr/include/flite/
// #cgo LDFLAGS: -lflite -lflite_cmu_us_kal
// #include "flite.h"
// cst_voice* register_cmu_us_kal(const char *voxdir);
import "C"

import (
	"fmt"
	"unsafe"
)

// #cgo CFLAGS: -I /usr/include/flite/
//	- inform cgo to include all files in flite directory where "flite.h" is present
// #cgo LDFLAGS: -lflite -lflite_cmu_us_kal
// 	- inform cgo to include "flite" and "flite_cmu_us_kal" libraries
// 	  flite_cmu_us_kal - contains the defition for `voice`
//    How did we come to know?
// 		ls /usr/lib/x86_64-linux-gnu/libflite*
//		- gave the location of `register_cmu_us_kal`
// 		- /usr/lib/x86_64-linux-gnu/libflite_cmu_us_kal
// #include "flite.h"
// cst_voice* register_cmu_us_kal(const char *voxdir);
// - is the declaration for the function.
//   How did we get this?
// 	 	Google : https://stackoverflow.com/questions/8510195/register-cmu-us-kal-function-in-flite

var voice *C.cst_voice

func init() {
	C.flite_init()
	voice = C.register_cmu_us_kal(nil)
}

// TextToSpeech takes a `text` to speak and stores it in `path`
func TextToSpeech(text, path string) error {
	if voice == nil {
		return fmt.Errorf("could not find default voide")
	}

	// find / | grep flite
	// - we get flite.h in /usr/include/flite/
	// cat /usr/include/flite/flite.h
	// float flite_text_to_speech(const char *filename,
	// 		cst_voice *voice,
	// 		const char *outtype);
	ctext := C.CString(text)
	cout := C.CString(path)
	C.flite_text_to_speech(ctext, voice, cout)
	C.free(unsafe.Pointer(ctext))
	C.free(unsafe.Pointer(cout))
	return nil
}
