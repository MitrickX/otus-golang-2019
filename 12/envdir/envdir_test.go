package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseDir(t *testing.T) {
	dirname := createDir("envdir_test_parse_dir")

	createFile(filepath.Join(dirname, "abc"), []byte("123"))
	createFile(filepath.Join(dirname, "def"), []byte("test\t  \t\nxyz\nok"))
	createFile(filepath.Join(dirname, "xyz"), []byte{'h', 0, 'i', '\n'})
	createFile(filepath.Join(dirname, "123"), nil)

	dirEnv, err := parseDir(dirname)
	if err != nil {
		t.Errorf("parseDir shouldn't return error in this case: %s", err)
		return
	}

	if len(dirEnv) == 0 {
		t.Error("parseDir should return not empty map")
		return
	}

	expected := envSet{
		"abc": envVal{
			val: "123",
		},
		"def": envVal{
			val: "test",
		},
		"xyz": envVal{
			val: "h\ni",
		},
		"123": envVal{
			remove: true,
		},
	}

	if !reflect.DeepEqual(expected, dirEnv) {
		t.Errorf("expected %v no equals to result %v\n", expected, dirEnv)
	}
}

func TestGetEnv(t *testing.T) {
	env := getEnv()
	if len(env) <= 0 {
		t.Errorf("Env is almost never is empty\n")
	}
}

func TestSetEnv(t *testing.T) {
	newEnv := envSet{
		"TEST_METHOD": envVal{
			val: "TestSetEnv",
		},
		"TEST_FILE": envVal{
			val: "envdir_test.go",
		},
	}

	setEnv(newEnv, false)

	currentEnv := getEnv()

	for key, val := range newEnv {
		curVal, ok := currentEnv[key]
		if !ok {
			t.Errorf("Key %s must be in current environment after setEnv", key)
		}
		if curVal.val != val.val {
			t.Errorf("Value on key %s must be equal %s", key, curVal.val)
		}
	}

}

func createDir(dirnamePrefix string) string {
	dirname, err := ioutil.TempDir(os.TempDir(), dirnamePrefix)
	if err != nil {
		log.Fatalf("couldn't create dir by prefix %s: %s\n", dirnamePrefix, err)
	}
	return dirname
}

func createFile(filename string, content []byte) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("couldn't create file %s: %s\n", filename, err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		log.Fatalf("couldn't write into new created file %s: %s\n", filename, err)
	}
}
