// +build tools

// Tool - Finds the BPM of a wav file.
package main

import (
	"azul3d.org/v1/audio"
	"azul3d.org/v1/audio/bpm"
	_ "azul3d.org/v1/audio/wav"
	"log"
	"os"
)

func test(fileName string) {
	log.Println(fileName)

	file, err := os.Open("src/azul3d.org/v1/assets/audio/" + fileName)
	if err != nil {
		log.Fatal(err)
	}

	// Create an decoder for the audio source
	decoder, format, err := audio.NewDecoder(file)
	if err != nil {
		log.Fatal(err)
	}

	// Grab the decoder's configuration
	config := decoder.Config()
	log.Println("Decoding an", format, "file.")
	log.Println(config)

	// Create an buffer that can hold 1 second of audio samples
	bufSize := 2 * config.SampleRate * config.Channels
	buf := make(audio.F64Samples, bufSize)

	// Fill the buffer with as many audio samples as we can
	read, err := decoder.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Read", read, "audio samples.")
	log.Println("")

	combSize := 128
	combDelay := 32
	bpm.Chunk(buf, combSize, combDelay)

	// readBuf := buf.Slice(0, read)
	// for i := 0; i < readBuf.Len(); i++ {
	//     sample := readBuf.At(i)
	// }
}

func main() {
	test("tune_stereo_44100hz_int16.wav")
	test("skrillex.wav")
}
