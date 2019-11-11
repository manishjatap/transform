package primitive

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIoTempFileActualImplementation(t *testing.T) {
	defer func() {
		err := recover()
		assert.Truef(t, err == nil, "Expected err message : %v", err)
	}()
	ioTempFile("a\xc5z", "a\xc5z")
}

func TestIoCopyFileActualImplementation(t *testing.T) {
	defer func() {
		err := recover()
		assert.Truef(t, err != nil, "Expected err message : %v", err)
	}()
	ioCopyFile(nil, nil)
}

func TestRemoveFileActualImplementation(t *testing.T) {
	defer func() {
		err := recover()
		assert.Truef(t, err != nil, "Expected err message : %v", err)
	}()
	removeFile(nil)
}

func TestGetFileNameActualImpl(t *testing.T) {
	defer func() {
		err := recover()
		assert.Truef(t, err != nil, "Expected err message : %v", err)
	}()
	getFileName(nil)
}

func TestCreateFileActualImpl(t *testing.T) {
	defer func() {
		err := recover()
		assert.Truef(t, err == nil, "Expected err message : %v", err)
	}()
	createFile("a\xc5z")
}

func TestTransformComCombinedOuputError(t *testing.T) {
	expectedError := "exit status 1"
	mockIoTempFile()
	mockRemoveFile()
	mockCreateFile()
	mockIoCopyFile()
	mockGetFileName()
	_, err := Transform(new(os.File), "fake", 0)

	assert.Equalf(t, expectedError, err.Error(), "Expected : %v", expectedError)
}

func TestTransformSuccess(t *testing.T) {
	mockIoTempFile()
	mockRemoveFile()
	mockCreateFile()
	mockIoCopyFile()
	mockGetFileName()
	mockGetCombinedOutput()
	_, err := Transform(new(os.File), "fake", 0)

	assert.NoError(t, err, "Expected : No Error")
}

func TestTransformIoTempFileEror(t *testing.T) {
	expectedError := "Error creating temp file"
	mockRemoveFile()
	mockCreateFile()
	mockIoCopyFile()
	mockGetFileName()
	errorIoTempFile(expectedError)
	_, err := Transform(new(os.File), "fake", 0)

	assert.Equalf(t, expectedError, err.Error(), "Expected : %v", expectedError)
}

func TestTransformIoCopyFileEror(t *testing.T) {
	expectedError := "Error coping file"
	mockIoTempFile()
	mockRemoveFile()
	mockCreateFile()
	errorIoCopyFile(expectedError)
	mockGetFileName()
	_, err := Transform(new(os.File), "fake", 0)

	assert.Equalf(t, err.Error(), expectedError, "Expected : %v", expectedError)
}

func TestTransformCreateFileEror(t *testing.T) {
	expectedError := "Error creating file"
	mockIoTempFile()
	mockRemoveFile()
	mockIoCopyFile()
	mockGetFileName()
	errorCreateFile(expectedError)
	_, err := Transform(new(os.File), "fake", 0)

	assert.Equalf(t, err.Error(), expectedError, "Expected : %v", expectedError)
}

func mockIoTempFile() {
	ioTempFile = func(dir string, prefix string) (*os.File, error) {
		return new(os.File), nil
	}
}

func mockIoCopyFile() {
	ioCopyFile = func(dst io.Writer, src io.Reader) (int64, error) {
		return 1, nil
	}
}

func mockRemoveFile() {
	removeFile = func(f *os.File) error {
		return nil
	}
}

func mockCreateFile() {
	createFile = func(filename string) (*os.File, error) {
		return new(os.File), nil
	}
}

func mockGetFileName() {
	getFileName = func(file *os.File) string {
		return "fake-file"
	}
}

func mockGetCombinedOutput() {
	getCombinedOutput = func(cmd *exec.Cmd) ([]byte, error) {
		var bytes []byte
		return bytes, nil
	}
}

func errorIoTempFile(errMsg string) {
	ioTempFile = func(dir string, prefix string) (*os.File, error) {
		return new(os.File), errors.New(errMsg)
	}
}

func errorIoCopyFile(errMsg string) {
	ioCopyFile = func(dst io.Writer, src io.Reader) (int64, error) {
		return 1, errors.New(errMsg)
	}
}

func errorCreateFile(errMsg string) {
	createFile = func(filename string) (*os.File, error) {
		return new(os.File), errors.New(errMsg)
	}
}
