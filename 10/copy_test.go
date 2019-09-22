package main

import (
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

// seed rand
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// test that copy works
func TestCopy(t *testing.T) {
	srcFilePath := generateFile(2048)
	srcFileSum := checkSum(srcFilePath)

	dstFilePath := generateFile(0)
	dstFileSum := checkSum(dstFilePath)

	if srcFileSum == dstFileSum {
		t.Fatal("generated files check sums are equal")
	}

	options := copyOptions{
		srcFilePath:  srcFilePath,
		dstFilePath:  dstFilePath,
		withProgress: false,
	}
	err := copyFile(options)
	if err != nil {
		t.Errorf("while copy file error happend %s\n", err)
	}

	dstFileSum = checkSum(dstFilePath)
	if srcFileSum != dstFileSum {
		t.Errorf("after copy files check sums are not equal, copy is failed\n")
	}

}

// test chunk size
// idea that small chunk size lead to slow execution time and same chunk sizes lead to same execution times (give or take)
func TestChunkSize(t *testing.T) {
	srcFilePath := generateFile(2048)
	srcFileSum := checkSum(srcFilePath)

	dstFilePath := generateFile(0)
	dstFileSum := checkSum(dstFilePath)

	if srcFileSum == dstFileSum {
		t.Fatal("generated files check sums are equal")
	}

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
		chunkSize:    2048,
	})
	if err != nil {
		t.Errorf("while copy file error happend %s\n", err)
	}

	elapsed2 := time.Since(start2)

	ratio := float64(elapsed1) / float64(elapsed2)
	if ratio < 2 {
		t.Errorf("copy with chunSize = 1 must be slower that copy with chunSize = 2048. Looks like the same and chunk size doesn't affect on copy process, ratio is %.2f", ratio)
	}
}

// generate test file or this size
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
func checkSum(filepath string) string {
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
