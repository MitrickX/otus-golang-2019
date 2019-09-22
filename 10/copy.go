package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// default chunk size of bytes to read and write at once
const defaultChunkSize = 512

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
	skip         int    // skip n chunks in source file
	seek         int    // skip n chunks in destination file
	count        int    // copy count chunks only
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

	// open src file
	srcFile, err := os.Open(options.srcFilePath)
	if err != nil {
		resErr = err
		return
	}

	defer closeFile(srcFile)

	// open dst file
	dstFile, err := os.Create(options.dstFilePath)
	if err != nil {
		resErr = err
		return
	}

	defer closeFile(dstFile)

	// define size of src file
	stat, err := srcFile.Stat()
	if err != nil {
		resErr = err
		return
	}
	fileSize := stat.Size()

	// chunk size for copy
	chunkSize := options.chunkSize
	if options.chunkSize == 0 {
		chunkSize = defaultChunkSize
	}

	// skip in bytes
	var skip int64
	if options.skip > 0 {
		skip = int64(options.skip) * int64(chunkSize)
	}

	// seek in bytes
	var seek int64
	if options.seek > 0 {
		seek = int64(options.seek) * int64(chunkSize)
	}

	// progress bar, if no need progress bar there is nil
	var progress *copyProgress
	if options.withProgress {
		// total size of all work
		var totalSize int64
		if options.count > 0 {
			totalSize = int64(options.count) * int64(chunkSize)
		} else {
			totalSize = fileSize - skip
		}
		progress = newCopyProgress(totalSize, chunkSize)
	}

	// init buffer (slice) of chunk size
	buffer := make([]byte, chunkSize)

	// skip in src file
	_, err = srcFile.Seek(skip, 0)
	if err != nil {
		resErr = err
		return
	}

	// seek in dst file
	_, err = dstFile.Seek(seek, 0)
	if err != nil {
		resErr = err
		return
	}

	// count of chunk copied so far
	count := 0

	// main copy loop (read chunk from src and write it to dst)
	for {

		if options.count > 0 && count >= options.count {
			break
		}
		count++

		// read chunk from src
		n, err := srcFile.Read(buffer)

		// eof - copy is done
		if err == io.EOF {
			if progress != nil {
				progress.done()
			}
			break
		}

		// some problem while reading, progress stop, return error
		if err != nil {
			if progress != nil {
				progress.stop()
			}
			resErr = fmt.Errorf("error while reading: %s", err)
			return
		}

		// write chunk to src
		nw, err := dstFile.Write(buffer[0:n])

		// some problem while writing, progress stop, return error
		if err != nil {
			if progress != nil {
				progress.stop()
			}
			resErr = fmt.Errorf("error while writing: %s", err)
			return
		}

		// numbers ofwritten and read bytes are not equals, stop progress, return error
		if n != nw {
			if progress != nil {
				progress.stop()
			}
			resErr = fmt.Errorf("error while writing: not all bytes be written (written = %d, read = %d)", nw, n)
			return
		}

		// print progress
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

	skipPtr := flagSet.Int("skip", 0, "skip n chunks in input file")
	seekPtr := flagSet.Int("seek", 0, "skip n chunks in output file")
	countPtr := flagSet.Int("count", 0, "copy count chunks only")

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
		defaults = strings.Replace(defaults, "-bs int", fmt.Sprintf("-bs int [optional] default is %d", defaultChunkSize), 1)
		defaults = strings.Replace(defaults, "-skip int", "-skip int [optional] default is no skip", 1)
		defaults = strings.Replace(defaults, "-seek int", "-seek int [optional] default is no seek", 1)
		defaults = strings.Replace(defaults, "-count int", "-count int [optional] default is copy whole input file", 1)

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
	options.skip = *skipPtr
	options.seek = *seekPtr
	options.count = *countPtr

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
