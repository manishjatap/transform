package primitive

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

//Mode : mode
type Mode int

//Const : constants
const (
	ModeCombo Mode = iota
	ModeTriangle
)

var ioTempFile = func(dir string, prefix string) (*os.File, error) {
	return ioutil.TempFile(dir, prefix)
}

var ioCopyFile = func(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

var removeFile = func(file *os.File) error {
	return os.Remove(file.Name())
}

var createFile = func(filename string) (*os.File, error) {
	return os.Create(filename)
}

var getFileName = func(file *os.File) string {
	return file.Name()
}

var getCombinedOutput = func(cmd *exec.Cmd) ([]byte, error) {
	return cmd.CombinedOutput()
}

//Transform : It transform the image
func Transform(img io.Reader, ext string, shapes int) (io.Reader, error) {
	inFile, err := createTempFile("in_", ext)
	if err != nil {
		return nil, err
	}
	defer removeFile(inFile)
	outFile, err := createTempFile("out_", ext)
	if err != nil {
		return nil, err
	}
	defer removeFile(outFile)

	_, err = ioCopyFile(inFile, img)
	if err != nil {
		return nil, err
	}

	_, err = primitive(getFileName(inFile), getFileName(outFile), shapes, ModeCombo)
	if err != nil {
		return nil, err
	}

	b := bytes.NewBuffer(nil)
	_, err = ioCopyFile(b, outFile)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func primitive(inFile string, outFile string, shapes int, mode Mode) (string, error) {
	argStr := fmt.Sprintf(" -i %s -o %s -n %d -m %d", inFile, outFile, shapes, mode)
	cmd := exec.Command("primitive", strings.Fields(argStr)...)
	bytes, err := getCombinedOutput(cmd)
	return string(bytes), err
}

func createTempFile(prefix string, ext string) (*os.File, error) {
	file, err := ioTempFile("./", prefix)
	if err != nil {
		return nil, err
	}
	defer removeFile(file)
	return createFile(getFileName(file) + "." + ext)
}
