package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// default chunk size of bytes to read and write at once
const chunkSize = 512

// copy process struct to represent progress bar in terminal
type copyProgress struct {
	fileSize     int64
	chunkSize    int
	progressSize int64
}

// copy process constructor
func newCopyProgress(fileSize int64, chunkSize int) *copyProgress {
	cp := &copyProgress{
		fileSize:  fileSize,
		chunkSize: chunkSize,
	}
	return cp
}

// add n bytes copied on step (generally it less or equal chunkSize)
func (cp *copyProgress) add(n int) {
	cp.progressSize += int64(n)
}

// total progress in percentages
func (cp *copyProgress) progress() float32 {
	return float32(float64(cp.progressSize)/float64(cp.fileSize)) * 100.0
}

// print progress of one step on coping (n bytes has been copied)
func (cp *copyProgress) step(n int) {
	cp.add(n)
	fmt.Print(cp.progressString(false))
}

// print 100% done progress bar
func (cp *copyProgress) done() {
	cp.progressSize = cp.fileSize
	fmt.Print(cp.progressString(true))
}

// print some stoped progress bar (if was failed somewhere in copy)
func (cp *copyProgress) stop() {
	fmt.Print(cp.progressString(true))
}

// helper that prints progress string
// - stop:
//     means it will not priting anymore (print out \n symbol at the end)
//     otherwise print \r and it will printing yet in this line of terminal
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

// options of copy function
type copyOptions struct {
	srcFilePath  string // path to source file
	dstFilePath  string // path to destination file
	withProgress bool   // need print progress bar?
	chunkSize    int    // chunk size of coping
}

// copy file
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

// parse cli arguments and return copyOptions struct for calling copyFile with it
func parseArgs() copyOptions {
	// options struct that will be filled after flag parsing
	options := copyOptions{}

	// output for flag set, need to print extended variant of defaults (aka usage or help)
	flagSetOut := &strings.Builder{}

	flagSet := flag.NewFlagSet("argParser", flag.ContinueOnError)
	flagSet.SetOutput(flagSetOut)

	srcFilePathPtr := flagSet.String("if", "", "path to input file")
	dstFilePathPtr := flagSet.String("of", "", "path to output file")

	chunkSizePtr := flagSet.Int("bs", 0, "chunk size of bytes to read and write at once")

	// helper to print defaults (extendend variant with extra info like required and etc)
	printDefaults := func() {

		// there is call inside of flagSet.Parse on error, so prevent double printing we should clear output
		flagSetOut.Reset()

		// print defaults in our output
		flagSet.PrintDefaults()

		// extend out print defaults
		defaults := flagSetOut.String()
		defaults = strings.Replace(defaults, "-if string", "-if string (required)", 1)
		defaults = strings.Replace(defaults, "-of string", "-of string (required)", 1)
		defaults = strings.Replace(defaults, "-bs int", fmt.Sprintf("-bs int [optional] default is %d", chunkSize), 1)

		// print our defaults into stderr
		fmt.Fprintln(os.Stderr, defaults)
	}

	// parse arguments and print defaults on error (with this error)
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		printDefaults()
		os.Exit(2)
	}

	// fill options after parse
	options.srcFilePath = *srcFilePathPtr
	options.dstFilePath = *dstFilePathPtr
	options.chunkSize = *chunkSizePtr

	// required
	if options.srcFilePath == "" {
		printDefaults()
		os.Exit(2)
	}

	// required
	if options.dstFilePath == "" {
		printDefaults()
		os.Exit(2)
	}

	// when call from cli need print progress
	options.withProgress = true

	return options
}

// parse arguments to copy options and call our copy with progress bar
func main() {
	options := parseArgs()
	err := copyFile(options)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
