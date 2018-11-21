package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func CreatePackage(cfg *Config) error {

	if err := cfg.Validate(); err != nil {
		return err
	}

	if err := os.MkdirAll(cfg.Assetspath, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(cfg.Assetsfilepath)
	if err != nil {
		return err
	}

	err = writeAssetsFile(f, cfg)
	f.Close()
	if err != nil {
		return err
	}

	err = os.MkdirAll(cfg.Assetspath+"/data", os.ModePerm)
	if err != nil {
		return err
	}

	f, err = os.Create("pack.sh")
	if err != nil {
		return err
	}
	err = writePackFile(f, cfg)
	if err != nil {
		return err
	}
	f.Close()
	err = setExecute("pack.sh")
	if err != nil {
		return err
	}
	err = runPackFile(cfg)
	if err != nil {
		return err
	}
	return nil
}

func setExecute(path string) error {
	return exec.Command("chmod", "+x", path).Run()
}

func writePackFile(w io.Writer, cfg *Config) error {

	fmt.Fprintf(w, `#!/bin/bash

echo "packing %s"

if [[ -f assets/data/bindata.go  ]]; then
	rm assets/data/bindata.go
fi

go-bindata -pkg data %s  && \
mv bindata.go assets/data/. && \
cd assets/ && go install 

`, cfg.Folder, cfg.Folder)

	return nil
}
func runPackFile(cfg *Config) error {
	return exec.Command("bash", "-c", "./pack.sh").Run()
}

type Config struct {
	Folder         string
	folderpath     string
	Assetspath     string
	Assetsfilepath string
	PackagePrefix  string
	PackFilePath   string
}

func (c *Config) Validate() error {
	var err error

	c.PackagePrefix, err = GetPathPrefix(c)
	if err != nil {
		return errors.Wrap(err, "cfg.validate.pathprefix")
	}

	f, err := GetSanePath(c.Folder)
	if err != nil {
		return errors.Wrap(err, "cfg.validate:"+c.Folder)
	}
	c.folderpath = f

	if !Exists(c.folderpath) {
		return errors.New("folder does not exist")
	}

	c.Assetspath, err = GetSanePath("assets")
	if err != nil {
		return err
	}
	if Exists(c.Assetspath) {
		return errors.New("there is already a folder named assets here...")
	}

	if Exists("pack.sh") {
		return errors.New("there is already a file named pack.sh here...")
	}

	c.Assetsfilepath = c.Assetspath + "/assets.go"
	return nil
}

func GetSanePath(path string) (string, error) {
	f, err := filepath.Abs(path)
	if err != nil {
		return f, err
	}
	// windows. mingw. stuffs
	f = strings.Replace(f, "\\", "/", -1)
	return f, nil
}

func GetPathPrefix(cfg *Config) (string, error) {
	// already set
	if len(cfg.PackagePrefix) > 0 {
		return cfg.PackagePrefix, nil
	}

	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	path, err = GetSanePath(path)
	if err != nil {
		return "", err
	}
	gopath, err := GetSanePath(GetGoPath())
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(path, gopath) {
		return "", errors.Errorf("cwd not gopath")
	}

	s := path[len(gopath)+len("/src/"):]
	return s, err
}

func GetGoPath() string {
	return os.Getenv("GOPATH")
}

func writeAssetsFile(w io.Writer, cfg *Config) error {

	fmt.Fprintf(w, `package assets

import (
	"%s/assets/data"
	"net/http"
	"strconv"
	"strings"
)

var (
	ContentTypeMapping = map[string]string{
		".css":  "text/css",
		".js":   "application/javascript",
		".html": "text/html",
	}
)

func Get(key string) ([]byte, error) {
	return data.Asset(key)
}

func WriteAsHTTPResponse(key string, rw http.ResponseWriter) error {
	b, err := Get(key)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return err
	}
	found := false
	for k, v := range ContentTypeMapping {
		if strings.HasSuffix(key, k) {
			rw.Header().Set("Content-Type", v)
			found = true
			break
		}
	}
	if !found {
		rw.Header().Set("Content-Type", "application/octet-stream")
	}
	rw.Header().Set("Content-Length", strconv.Itoa(len(b)))
	_, err = rw.Write(b)
	return err
}

`, cfg.PackagePrefix)

	return nil
}

func logln() {
	log.Println(runtime.GOOS)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil && !os.IsNotExist(err)
}
