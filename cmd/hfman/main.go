package main

import (
	"net/http"
	"os"
	"path"
	"fmt"
	"flag"
	"github.com/pnovotnak/fman/internal/app/hfman"
	"github.com/pnovotnak/fman/internal/pkg/httplib"
)

func main()  {
	// Root directory to serve
	var dir string
	var err error

	var bind = *flag.String("bind", ":3000", "Bind specification")
	var parents = *flag.Bool("parents", true, "Create directory parents as required")
	var filePerm = os.FileMode(*flag.Uint("filemode", 0660, "Created file permissions"))
	var dirPerm = os.FileMode(*flag.Uint("dirmode", 0770, "Created folder permissions"))
	flag.Parse()

	if len(os.Args) != 2 {
		fmt.Println("Must pass directory to serve as only argument")
		os.Exit(1)
	}

	dir = os.Args[1]
	if !parents {
		if _, err = os.Stat(dir); os.IsNotExist(err) {
			fmt.Println("Directory to serve doesn't exist and create parents was not passed")
			os.Exit(1)
		}
	}
	err = hfman.CreateDirs(dir, dirPerm)

	http.HandleFunc("/", httplib.Redirector("/download/", 301))
	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(path.Join(dir, "download")))))
	http.HandleFunc("/upload/", hfman.Router(dir, filePerm))

	fmt.Println(os.Args[0], "starting, serving from:", dir, "bound to:", bind)
	if err = http.ListenAndServe(bind, httplib.RequestLoggerHandler(http.DefaultServeMux)); err != nil {
		panic(err)
	}
}
