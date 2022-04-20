package cmd

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
)

var updateGolden = flag.Bool("update", false, "update golden files")

const (
	ProjectName = "test"
)

func runTestCmd(t *testing.T, cmd string) {
	buf := new(bytes.Buffer)
	args, err := shellwords.Parse(cmd)
	if err != nil {
		t.Error(err)
		return
	}

	rootCmd := newRootCmd(args)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)

	if err := rootCmd.Execute(); err != nil {
		t.Error(err)
		return
	}
}

// makeTestSingleWorkspace will change current workspace
func makeTestSingleWorkspace(t *testing.T) string {
	p := path.Join(os.TempDir(), "kratos-scaffold/test-single")
	err := os.MkdirAll(p, 0o700)
	if err != nil {
		if os.IsExist(err) {
			_ = os.RemoveAll(p)
			_ = os.MkdirAll(p, 0o700)
		} else {
			t.Error(err)
			return ""
		}

	}
	_ = os.Chdir(p)
	cmd := "new " + ProjectName
	runTestCmd(t, cmd)
	return path.Join(p, ProjectName)
}

// AssertGoldenBytes asserts that the give actual content matches the contents of the given filename
func AssertGoldenBytes(t *testing.T, actual []byte, filename string) {
	t.Helper()

	if err := compare(actual, realRath(filename)); err != nil {
		t.Fatalf("%v", err)
	}
}

// AssertGoldenFile asserts that the content of the actual file matches the contents of the expected file
func AssertGoldenFile(t *testing.T, actualFileName string, expectedFilename string) {
	t.Helper()

	actual, err := ioutil.ReadFile(actualFileName)
	if err != nil {
		t.Fatalf("%v", err)
	}
	AssertGoldenBytes(t, actual, expectedFilename)
}

func realRath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join("testdata", filename)
}

func compare(actual []byte, filename string) error {
	actual = normalize(actual)
	if err := update(filename, actual); err != nil {
		return err
	}

	expected, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrapf(err, "unable to read testdata %s", filename)
	}
	expected = normalize(expected)
	if !bytes.Equal(expected, actual) {
		return errors.Errorf("does not match file %s\n\nWANT:\n'%s'\n\nGOT:\n'%s'\n", filename, expected, actual)
	}
	return nil
}

func update(filename string, in []byte) error {
	if !*updateGolden {
		return nil
	}
	return ioutil.WriteFile(filename, normalize(in), 0666)
}

func normalize(in []byte) []byte {
	return bytes.Replace(in, []byte("\r\n"), []byte("\n"), -1)
}
