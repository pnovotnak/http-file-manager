package fman

import (
	"net/http"
	"os"
	"path"
	"syscall"
	"strings"
	"fmt"
	"io"
)

var uploadForm = []byte(`<!DOCTYPE html>
<html lang="en-US">
<body>
	<form method="POST" enctype="multipart/form-data">
		<input type="file" name="file" id="file">
		<input type="submit" value="Upload File" name="submit">
	</form>
</body>
</html>
`)

func receiveFile(w http.ResponseWriter, r *http.Request, root string, filePerm os.FileMode) {
	var diskFile *os.File
	var err error
	memFile, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer memFile.Close()

	diskFile, err = os.OpenFile(path.Join(root, "upload", path.Base(header.Filename)), os.O_CREATE|os.O_EXCL|os.O_WRONLY|os.O_TRUNC, filePerm)
	if err != nil {
		if err.(*os.PathError).Err == syscall.EEXIST {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer diskFile.Close()

	name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", name[0])
	_, err = io.Copy(diskFile, memFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("OK"))
}

func Router(root string, filePerm os.FileMode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			receiveFile(w, r, root, filePerm)
		} else if r.Method == "GET" {
			w.Write(uploadForm)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 - Only GET and POST are allowed."))
		}
	}
}
