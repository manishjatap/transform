package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/manishjagtap/transform/primitive"
)

var openFile = func(path string) (*os.File, error) {
	return os.Open(path)
}

var transformImage = func(img io.Reader, ext string, shapes int) (io.Reader, error) {
	return primitive.Transform(img, ext, shapes)
}

var ioCopyFile = func(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

func home(w http.ResponseWriter, r *http.Request) {
	html := `<html>
				<body>
					<form action="/upload" method="post" enctype="multipart/formdata">
						<input type="text" name="image">
						<button type="submit">Upload Image</button>
					</form>
				</body>
			 </html>`
	fmt.Fprint(w, html)
}

func upload(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("image")

	inputF, err := openFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer inputF.Close()

	ext := filepath.Ext(path)[1:]

	out, err := transformImage(inputF, ext, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "image/png")
	ioCopyFile(w, out)
}

func getHandler() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", home).Methods(http.MethodGet)
	router.HandleFunc("/upload", upload).Methods(http.MethodPost)

	return router
}

func main() {
	log.Fatal(http.ListenAndServe(":3000", getHandler()))
}
