package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// default size of generated file for tests
const testFileSize = 8192 // 8 Kb

// seed rand
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// test that copy works
func TestDefault(t *testing.T) {
	srcFilePath := generateFile(testFileSize)
	srcFileSum := fileCheckSum(srcFilePath)

	dstFilePath := generateFile(0)

	options := copyOptions{
		srcFilePath:  srcFilePath,
		dstFilePath:  dstFilePath,
		withProgress: false,
	}
	err := copyFile(options)
	if err != nil {
		t.Errorf("while copy file error happend %s\n", err)
	}

	dstFileSum := fileCheckSum(dstFilePath)
	if srcFileSum != dstFileSum {
		t.Errorf("after copy files check sums are not equal, copy is failed\n")
	}

}

// test chunk size
// idea that small chunk size lead to slow execution time and same chunk sizes lead to same execution times (give or take)
func TestChunkSize(t *testing.T) {
	srcFilePath := generateFile(testFileSize)

	dstFilePath := generateFile(0)

	var err error

	start1 := time.Now()

	err = copyFile(copyOptions{
		srcFilePath:  srcFilePath,
		dstFilePath:  dstFilePath,
		withProgress: false,
		chunkSize:    1,
	})
	if err != nil {
		t.Errorf("while copy file error happend %s\n", err)
	}

	elapsed1 := time.Since(start1)

	start2 := time.Now()

	err = copyFile(copyOptions{
		srcFilePath:  srcFilePath,
		dstFilePath:  dstFilePath,
		withProgress: false,
		chunkSize:    testFileSize,
	})
	if err != nil {
		t.Errorf("while copy file error happend %s\n", err)
	}

	elapsed2 := time.Since(start2)

	ratio := float64(elapsed1) / float64(elapsed2)
	if ratio < 2 {
		t.Errorf(`copy with chunSize = 1 must be slower that copy with chunSize = %d. 
Looks like the same and chunk size doesn't affect on copy process, ratio is %.2f`, testFileSize, ratio)
	}
}

func TestSkip(t *testing.T) {
	srcFilePath := generateFile(testFileSize)

	dstFilePath := generateFile(0)

	options := copyOptions{
		srcFilePath:  srcFilePath,
		dstFilePath:  dstFilePath,
		withProgress: false,
		skip:         1,
	}

	err := copyFile(options)
	if err != nil {
		t.Errorf("while copy file error happend %s\n", err)
	}

	srcData := readFile(srcFilePath)
	dstData := readFile(dstFilePath)

	expectedData := srcData[options.skip*defaultChunkSize:]

	dstDataSum := checkSum(dstData)
	expectedDataSum := checkSum(expectedData)

	if dstDataSum != expectedDataSum {
		t.Errorf("copy with skip = %d is failed", options.skip)
	}
}

func TestSeek(t *testing.T) {
	srcFilePath := generateFile(testFileSize)

	dstFilePath := generateFile(0)

	options := copyOptions{
		srcFilePath:  srcFilePath,
		dstFilePath:  dstFilePath,
		withProgress: false,
		seek:         1,
	}

	err := copyFile(options)
	if err != nil {
		t.Errorf("while copy file error happend %s\n", err)
	}

	srcData := readFile(srcFilePath)
	dstData := readFile(dstFilePath)

	// expected data has defaultChunkSize of zeros in the beginning and than all src data
	expectedData := make([]byte, defaultChunkSize, len(srcData)+defaultChunkSize)
	expectedData = append(expectedData, srcData...)

	dstDataSum := checkSum(dstData)
	expectedDataSum := checkSum(expectedData)

	if dstDataSum != expectedDataSum {
		t.Errorf("copy with seek = %d is failed", options.seek)
	}
}

func TestCount(t *testing.T) {
	srcFilePath := generateFile(testFileSize)

	dstFilePath := generateFile(0)

	options := copyOptions{
		srcFilePath:  srcFilePath,
		dstFilePath:  dstFilePath,
		withProgress: false,
		count:        2,
	}

	err := copyFile(options)
	if err != nil {
		t.Errorf("while copy file error happend %s\n", err)
	}

	srcData := readFile(srcFilePath)
	dstData := readFile(dstFilePath)

	expectedData := srcData[0 : options.count*defaultChunkSize]

	dstDataSum := checkSum(dstData)
	expectedDataSum := checkSum(expectedData)

	if dstDataSum != expectedDataSum {
		t.Errorf("copy with count = %d is failed", options.count)
	}
}

func TestSkipSeekCount(t *testing.T) {
	srcFilePath := generateFile(testFileSize)

	dstFilePath := generateFile(0)

	options := copyOptions{
		srcFilePath:  srcFilePath,
		dstFilePath:  dstFilePath,
		withProgress: false,
		skip:         2,
		seek:         3,
		count:        2,
	}

	err := copyFile(options)
	if err != nil {
		t.Errorf("while copy file error happend %s\n", err)
	}

	srcData := readFile(srcFilePath)
	dstData := readFile(dstFilePath)

	// skip, seek, count converted to bytes
	seek := options.seek * defaultChunkSize
	skip := options.skip * defaultChunkSize
	count := options.count * defaultChunkSize

	// build expected data
	// expected data has seek of zeros in the beginning and piece of srcData
	// peace of srcData is count bytes started from skip bytes
	expectedData := make([]byte, seek, seek+count)
	piece := srcData[skip : skip+count]
	expectedData = append(expectedData, piece...)

	dstDataSum := checkSum(dstData)
	expectedDataSum := checkSum(expectedData)

	if dstDataSum != expectedDataSum {
		t.Errorf("copy with skip = %d, seek = %d and count = %d is failed", options.skip, options.seek, options.count)
	}
}

// generate test file of this size
func generateFile(size int64) string {
	file, err := ioutil.TempFile(os.TempDir(), "copy_unit_test")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	for i := int64(0); i < size; i++ {
		d := rand.Intn(10)
		b := '0' + byte(d)
		_, err := file.Write([]byte{b})
		if err != nil {
			log.Fatal(err)
		}
	}

	path := filepath.Join(os.TempDir(), stat.Name())
	return path
}

// check sum of file
func fileCheckSum(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

// check some of bytes data
func checkSum(data []byte) string {
	h := md5.New()
	r := bytes.NewReader(data)
	if _, err := io.Copy(h, r); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// read file into bytes data
func readFile(filepath string) []byte {
	srcData, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	return srcData
}
