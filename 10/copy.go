package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const bufferSize = 4096

type copyProgress struct {
	fileSize     int64
	bufferSize   int
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

func newCopyProgress(fileSize int64, bufferSize int) *copyProgress {
	cp := &copyProgress{
		fileSize:   fileSize,
		bufferSize: bufferSize,
	}
	return cp
}

func copyFile(srcFilePath string, dstFilePath string, withProgress bool) (resErr error) {

	// helper to close file
	closeFile := func(file *os.File) {
		err := file.Close()
		if err != nil {
			if resErr == nil {
				resErr = err
			}
		}
	}

	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		resErr = err
		return
	}

	defer closeFile(srcFile)

	dstFile, err := os.Create(dstFilePath)
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
	if withProgress {
		progress = newCopyProgress(fileSize, bufferSize)
	}

	buffer := make([]byte, bufferSize)

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

type args struct {
	srcFilePath *string
	dstFilePath *string
}

func parseArgs() args {
	a := args{}

	flagSetOut := &strings.Builder{}

	flagSet := flag.NewFlagSet("argParser", flag.ContinueOnError)
	flagSet.SetOutput(flagSetOut)

	a.srcFilePath = flagSet.String("if", "", "path to input file")
	a.dstFilePath = flagSet.String("of", "", "path to output file")

	// helper to print defaults (extendend variant)
	printDefaults := func() {
		flagSet.PrintDefaults()
		defaults := flagSetOut.String()
		defaults = strings.Replace(defaults, "-if string", "-if string (required)", 1)
		defaults = strings.Replace(defaults, "-of string", "-if string (required)", 1)
		fmt.Fprintln(os.Stderr, defaults)
	}

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		printDefaults()
		os.Exit(2)
	}

	if *a.srcFilePath == "" {
		printDefaults()
		os.Exit(2)
	}

	if *a.dstFilePath == "" {
		printDefaults()
		os.Exit(2)
	}

	return a
}

func main() {
	args := parseArgs()
	err := copyFile(*args.srcFilePath, *args.dstFilePath, true)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
