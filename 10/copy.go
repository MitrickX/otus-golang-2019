package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const chunkSize = 512

type copyProgress struct {
	fileSize     int64
	chunkSize    int
	progressSize int64
}

func (cp *copyProgress) add(n int) {
	cp.progressSize += int64(n)
}

func (cp *copyProgress) progress() float32 {
	return float32(float64(cp.progressSize)/float64(cp.fileSize)) * 100.0
}

func (cp *copyProgress) step(n int) {
	cp.add(n)
	fmt.Print(cp.progressString(false))
}

func (cp *copyProgress) done() {
	cp.progressSize = cp.fileSize
	fmt.Print(cp.progressString(true))
}

func (cp *copyProgress) stop() {
	fmt.Print(cp.progressString(true))
}

func (cp *copyProgress) progressString(stop bool) string {
	progress := cp.progress()

	out := make([]byte, 0, 100)
	out = append(out, []byte(fmt.Sprintf("Progress: [%3.2f%%] ", progress))...)

	out = append(out, '[')
	progressInt := int(progress)
	for i := 1; i <= 100; i++ {
		if i <= progressInt {
			out = append(out, '#')
		} else {
			out = append(out, '.')
		}
	}
	out = append(out, ']')

	if stop {
		out = append(out, '\n')
	} else {
		out = append(out, '\r')
	}

	return string(out)
}

func newCopyProgress(fileSize int64, chunkSize int) *copyProgress {
	cp := &copyProgress{
		fileSize:  fileSize,
		chunkSize: chunkSize,
	}
	return cp
}

type copyOptions struct {
	srcFilePath  string
	dstFilePath  string
	withProgress bool
	chunkSize    int
}

func copyFile(options copyOptions) (resErr error) {

	// helper to close file
	closeFile := func(file *os.File) {
		err := file.Close()
		if err != nil {
			if resErr == nil {
				resErr = err
			}
		}
	}

	srcFile, err := os.Open(options.srcFilePath)
	if err != nil {
		resErr = err
		return
	}

	defer closeFile(srcFile)

	dstFile, err := os.Create(options.dstFilePath)
	if err != nil {
		resErr = err
		return
	}

	defer closeFile(dstFile)

	stat, err := srcFile.Stat()
	if err != nil {
		resErr = err
		return
	}
	fileSize := stat.Size()

	var progress *copyProgress
	if options.withProgress {
		progress = newCopyProgress(fileSize, chunkSize)
	}

	if options.chunkSize == 0 {
		options.chunkSize = chunkSize
	}

	buffer := make([]byte, options.chunkSize)

	for {
		n, err := srcFile.Read(buffer)
		if err == io.EOF {
			if progress != nil {
				progress.done()
			}
			break
		}
		if err != nil {
			if progress != nil {
				progress.stop()
			}
			resErr = fmt.Errorf("error while reading: %s", err)
			return
		}

		nw, err := dstFile.Write(buffer[0:n])

		if err != nil {
			if progress != nil {
				progress.stop()
			}
			resErr = fmt.Errorf("error while writing: %s", err)
			return
		}

		if n != nw {
			if progress != nil {
				progress.stop()
			}
			resErr = fmt.Errorf("error while writing: not all bytes be written (written = %d, read = %d)", nw, n)
			return
		}

		if progress != nil {
			progress.step(n)
		}
	}

	return

}

func parseArgs() copyOptions {
	options := copyOptions{}

	flagSetOut := &strings.Builder{}

	flagSet := flag.NewFlagSet("argParser", flag.ContinueOnError)
	flagSet.SetOutput(flagSetOut)

	srcFilePathPtr := flagSet.String("if", "", "path to input file")
	dstFilePathPtr := flagSet.String("of", "", "path to output file")

	chunkSizePtr := flagSet.Int("bs", 0, "chunk size of bytes to read and write at once")

	// helper to print defaults (extendend variant)
	printDefaults := func() {
		flagSetOut.Reset()
		flagSet.PrintDefaults()
		defaults := flagSetOut.String()
		defaults = strings.Replace(defaults, "-if string", "-if string (required)", 1)
		defaults = strings.Replace(defaults, "-of string", "-of string (required)", 1)
		defaults = strings.Replace(defaults, "-bs int", fmt.Sprintf("-bs int [optional] default is %d", chunkSize), 1)
		fmt.Fprintln(os.Stderr, defaults)
	}

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		printDefaults()
		os.Exit(2)
	}

	options.srcFilePath = *srcFilePathPtr
	options.dstFilePath = *dstFilePathPtr
	options.chunkSize = *chunkSizePtr

	if options.srcFilePath == "" {
		printDefaults()
		os.Exit(2)
	}

	if options.dstFilePath == "" {
		printDefaults()
		os.Exit(2)
	}

	options.withProgress = true

	return options
}

func main() {
	options := parseArgs()
	err := copyFile(options)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
