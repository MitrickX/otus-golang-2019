package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const statusCodeInvalidArguments = 1 // print usage
const statusCodeFail = 111           // envdir on fail exit with 111 status code

//envdir d child
//d is a single argument.  child consists of one or more arguments.
//envdir sets various environment variables as specified by files in the directory named  d.
//It then runs child.
//If d contains a file named s whose first line is t, envdir removes an environment variable
//named s if one exists, and then adds an environment variable named s with  value  t.   The
//name  s  must  not  contain =. Spaces and tabs at the end of t are removed. Nulls in t are
//changed to newlines in the environment variable.
//If the file s is completely empty (0 bytes long), envdir removes an  environment  variable
//named s if one exists, without adding a new variable.
//envdir  exits  111  if  it has trouble reading d, if it runs out of memory for environment
//variables, or if it cannot run child.  Otherwise its exit code is  the  same  as  that  of
//child.

type envVal struct {
	val    string
	remove bool
}

type envSet map[string]*envVal

func newEnvSet() envSet {
	return make(envSet)
}

func parseDir(dir string) (envSet, error) {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	res := newEnvSet()

	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}

		key := fileInfo.Name()

		if fileInfo.Size() == 0 {
			res[key] = &envVal{
				remove: true,
			}
			continue
		}

		filepath := filepath.Join(dir, fileInfo.Name())

		file, err := os.Open(filepath)
		if err != nil {
			return nil, err
		}

		reader := bufio.NewReader(file)
		line, _, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}

		val := string(line)
		val = strings.TrimRight(val, " \t")
		val = strings.Replace(val, "\x00", "\n", -1)

		res[key] = &envVal{
			val: val,
		}
	}

	return res, nil
}

func setEnv(newEnv envSet) {
	for key, val := range newEnv {
		if !val.remove {
			os.Setenv(key, val.val)
		} else {
			os.Unsetenv(key)
		}
	}
}

func getExitCode(runErr error) int {
	if runErr == nil {
		return 0
	}

	exitErr, ok := runErr.(*exec.ExitError)
	if !ok {
		return 0
	}

	return exitErr.ExitCode()
}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: envdir d child\n")
		os.Exit(statusCodeInvalidArguments)
	}

	dir := args[0]
	child := args[1]

	var childArgs []string
	if len(args) > 2 {
		childArgs = args[2:]
	}

	env, err := parseDir(dir)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error happend while parse dir %s: %s\n", dir, err)
		os.Exit(statusCodeFail)
	}

	setEnv(env)

	cmd := exec.Command(child, childArgs...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		if len(childArgs) > 0 {
			fmt.Fprintf(os.Stderr, "error happend while parse run child %s with arguments %s: %s\n", child, childArgs, err)
		} else {
			fmt.Fprintf(os.Stderr, "error happend while parse run child %s: %s\n", child, err)
		}
		os.Exit(statusCodeFail)
	}

	runErr := cmd.Wait()
	exitCode := getExitCode(runErr)

	os.Exit(exitCode)

}
