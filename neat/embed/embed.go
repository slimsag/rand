// Copyright 2012 Lightpoke. All rights reserved.
// This source code is subject to the terms and
// conditions defined in the "License.txt" file.

// embed is a small tool to embed a single file within a Go binary.
package main

import (
	"azul3d.org/embed/embed"
	"flag"
	"io"
	"log"
	"os"
)

var (
	input, output, replace string
	defaultOutputFile      = "<input file>"
)

func init() {
	log.SetFlags(0)

	flag.StringVar(&input, "i", "", "Input binary file to operate on")
	flag.StringVar(&output, "o", defaultOutputFile, "Output file to write to")
	flag.StringVar(&replace, "replace", "", "File to replace the embedded file with")
}

func main() {
	flag.Parse()

	if len(input) == 0 {
		log.Println("Must specify input binary file\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if output == defaultOutputFile {
		output = input
	}

	_, err := os.Stat(input)
	if err != nil {
		log.Fatal(err)
	}

	if len(replace) > 0 {
		// Open the file we are to replace the embedded file with
		replaceFile, err := os.Open(replace)
		if err != nil {
			log.Fatal(err)
		}
		defer replaceFile.Close()

		// Create or use already existing output file
		outputFile, err := os.OpenFile(output, os.O_RDWR, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer outputFile.Close()

		// Open up the input binary file
		inputFile, err := os.Open(input)
		if err != nil {
			log.Fatal(err)
		}
		defer inputFile.Close()

		// See if the input file already has an embedded file
		_, err = embed.Offset(inputFile)
		if err != nil && err != embed.InvalidFileErr {
			// An I/O error
			log.Fatal(err)
		}

		if err == nil {
			// The file contains an embedded file already, so copy everything except the old
			// embedded file into the output file.
			preamble, err := embed.Preamble(inputFile)
			if err != nil {
				log.Fatal(err)
			}

			bytesWrote, err := io.Copy(outputFile, preamble)
			if err != nil {
				log.Fatal(err)
			}

			err = outputFile.Truncate(bytesWrote)
			if err != nil {
				log.Fatal(err)
			}
			_, err = outputFile.Seek(0, 2)
			if err != nil {
				log.Fatal(err)
			}

			// Now update the embedded file inside the output file
			bytesCopied, err := io.Copy(outputFile, replaceFile)
			if err != nil {
				log.Fatal(err)
			}

			// And update the embedded file footer
			err = embed.WriteFooter(outputFile, bytesCopied)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Replaced existing embedded file with %v bytes\n", bytesWrote)
		} else {
			// The file has no embedded file yet, so create an new one.

			// Copy the preamble binary file.
			bytesCopied, err := io.Copy(outputFile, inputFile)
			if err != nil {
				log.Fatal(err)
			}

			// Now create the embedded file inside the output file
			bytesCopied, err = io.Copy(outputFile, replaceFile)
			if err != nil {
				log.Fatal(err)
			}

			// Now write the embedded file footer
			err = embed.WriteFooter(outputFile, bytesCopied)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Created new embedded file with %v bytes\n", bytesCopied)
		}

	} else {
		// Extracting the embedded file

		// Open up the input binary file
		inputFile, err := os.Open(input)
		if err != nil {
			log.Fatal(err)
		}
		defer inputFile.Close()

		// Ensure that the input file actually has an embedded file
		offset, err := embed.Offset(inputFile)
		if err != nil {
			log.Fatal(err)
		}

		// Seek to where the embedded file begins
		_, err = inputFile.Seek(offset, 0)
		if err != nil {
			log.Fatal(err)
		}

		// Create or use already existing output file
		outputFile, err := os.Create(output)
		if err != nil {
			log.Fatal(err)
		}
		defer outputFile.Close()

		// Copy the embedded file into the output file
		bytesCopied, err := io.Copy(outputFile, inputFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Extracted %v bytes to %s\n", bytesCopied, output)
	}
}
