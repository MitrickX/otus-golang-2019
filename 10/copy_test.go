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

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func TestCopy(t *testing.T) {
	srcFilePath := generateFile(2048)
	srcFileSum := checkSum(srcFilePath)

	dstFilePath := generateFile(0)
	dstFileSum := checkSum(dstFilePath)

	if srcFileSum == dstFileSum {
		t.Fatal("generated files check sums are equal")
	}

	err := copyFile(srcFilePath, dstFilePath, false)
	if err != nil {
		t.Errorf("while copy file error happend %s\n", err)
	}

	dstFileSum = checkSum(dstFilePath)
	if srcFileSum != dstFileSum {
		t.Errorf("after copy files check sums are not equal, copy is failed\n")
	}

}

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
