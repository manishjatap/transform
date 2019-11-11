package main

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenFileActualImplementation(t *testing.T) {
	defer func() {
		err := recover()
		assert.Truef(t, err == nil, "Expected no error : %v", err)
	}()
	openFile("a\xc5z")
}

func TestIoCopyFileActualImplementation(t *testing.T) {
	defer func() {
		err := recover()
		assert.Truef(t, err != nil, "Expected err message : %v", err)
	}()
	ioCopyFile(nil, nil)
}

func TestTransformImageActualImplementation(t *testing.T) {
	defer func() {
		err := recover()
		assert.Truef(t, err != nil, "Expected err message : %v", err)
	}()
	transformImage(nil, "", 0)
}

func TestUpload(t *testing.T) {

	mockOpenFile()
	mockTransformImage()
	mockIoCopyFile()

	reader := strings.NewReader("image=fake.fake") //post data
	req, _ := http.NewRequest("POST", "/upload", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(upload)
	handler.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code, "Expected : 200 ok")
}

func TestHome(t *testing.T) {

	req, _ := http.NewRequest("GET", "/home", nil)

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(home)
	handler.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code, "Expected : 200 ok")
}
func TestUploadOpenFileError(t *testing.T) {

	errorOpenFile("Error while opening file")
	mockTransformImage()
	mockIoCopyFile()

	reader := strings.NewReader("image=fake.fake") //post data
	req, _ := http.NewRequest("POST", "/upload", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(upload)
	handler.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "Expected : 400 Bad Request")
}

func TestUploadImageTransformError(t *testing.T) {

	mockOpenFile()
	errorTransformImage("Error while image transform")
	mockIoCopyFile()

	reader := strings.NewReader("image=fake.fake") //post data
	req, _ := http.NewRequest("POST", "/upload", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(upload)
	handler.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "Expected : 500 Internal Server Error")
}

func TestHandler(t *testing.T) {
	router := getHandler()
	assert.NotNil(t, router, "Expected : Router should not be nil")
}

func mockOpenFile() {
	openFile = func(path string) (*os.File, error) {
		return new(os.File), nil
	}
}

func mockTransformImage() {
	transformImage = func(img io.Reader, ext string, shapes int) (io.Reader, error) {
		return new(os.File), nil
	}
}

func mockIoCopyFile() {
	ioCopyFile = func(dst io.Writer, src io.Reader) (int64, error) {
		return 1, nil
	}
}

func errorOpenFile(errMsg string) {
	openFile = func(path string) (*os.File, error) {
		return new(os.File), errors.New(errMsg)
	}
}

func errorTransformImage(errMsg string) {
	transformImage = func(img io.Reader, ext string, shapes int) (io.Reader, error) {
		return new(os.File), errors.New(errMsg)
	}
}
