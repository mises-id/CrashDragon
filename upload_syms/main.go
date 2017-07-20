package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

var file string
var repo string
var host string
var prod string
var ver string

func init() {
	flag.StringVar(&file, "file", "", "file to upload symbols of")
	flag.StringVar(&repo, "repo", "", "the git repo the source file is part of")
	flag.StringVar(&host, "host", "", "host (with protocol) to upload symbols to")
	flag.StringVar(&prod, "prod", "", "the product for which this symfile is")
	flag.StringVar(&ver, "ver", "", "the version of the product")
}

func main() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stderr)
	flag.Parse()
	defer os.RemoveAll(file + ".sym")
	defer os.RemoveAll(file + ".dSYM")
	dsymutil := exec.Command("dsymutil", file)
	dsymutil.Stderr = os.Stderr
	err := dsymutil.Run()
	if err != nil {
		log.Fatal(err)
	}

	dumpsyms := exec.Command("dump_syms", "-r", "-g", file+".dSYM", file)
	var in bytes.Buffer
	dumpsyms.Stdout = &in
	dumpsyms.Stderr = os.Stderr
	err = dumpsyms.Run()
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(in.Bytes()), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "FILE") {
			parts := strings.Split(line, " ")
			num := parts[1]
			opath := path.Join(path.Dir(parts[2]), path.Base(parts[2]))
			if !strings.HasPrefix(opath, repo) {
				continue
			}
			gitparts := strings.Split(repo, "/")
			origparts := strings.Split(opath, "/")
			npath := strings.Join(origparts[len(gitparts):], "/")
			lines[i] = fmt.Sprintf("FILE %s %s", num, npath)
		}
	}
	var out bytes.Buffer
	out.Write([]byte(strings.Join(lines, "\n")))
	body := upload(host, file, out)
	log.Printf("%s", string(body))
}

func upload(url, filename string, filedata bytes.Buffer) []byte {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("symfile", filename)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = io.Copy(fw, &filedata); err != nil {
		log.Fatal(err)
	}
	pw, err := w.CreateFormField("prod")
	if err != nil {
		log.Fatal(err)
	}
	pw.Write([]byte(prod))
	vw, err := w.CreateFormField("ver")
	if err != nil {
		log.Fatal(err)
	}
	vw.Write([]byte(ver))
	w.Close()

	req, err := http.NewRequest("POST", url+"/symfiles", &b)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var body bytes.Buffer
	io.Copy(&body, res.Body)
	return body.Bytes()
}
