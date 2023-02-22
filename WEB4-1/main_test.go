package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadTest(t *testing.T) {
	assert := assert.New(t)

	//file read
	path := "C:/Users/조재모/Downloads/포트폴리오.pptx"
	file, _ := os.Open(path)
	defer file.Close()

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	// formfile 생성
	multi, err := writer.CreateFormFile("upload_file", filepath.Base(path))

	assert.NoError(err)
	io.Copy(multi, file)

	writer.Close()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/uploads", buf)
	req.Header.Set("content-type", writer.FormDataContentType())

	uploadsHandler(res, req)

	assert.Equal(http.StatusOK, res.Code)

	uploadFilePath := "./uploads/" + filepath.Base(path)
	// file 의 info를 가져다준다
	finfo, err := os.Stat(uploadFilePath)
	assert.NoError(err)

	os.Open(uploadFilePath)

	uploadFile, _ := os.Open(uploadFilePath)
	originFile, _ := os.Open(path)

	defer uploadFile.Close()
	defer originFile.Close()
	uploadData := []byte{}
	originData := []byte{}

	uploadFile.Read(uploadData)
	originFile.Read(originData)

	fmt.Print(finfo)

	assert.Equal(originData, uploadData)

}
