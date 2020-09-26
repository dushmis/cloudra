package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/user"
	"strings"
)

var (
	file     string
	lifeTime string
	bucket string
)

func main() {
	flag.StringVar(&file, "f", "/dev/null", "file")
	flag.StringVar(&lifeTime, "life", "200", "life in minutes")
	flag.StringVar(&bucket, "bucket", "", "file bucket")
	flag.Parse()

	key := GetLocalUserKey()

	if key == "" {
		fmt.Println("Invalid Key")
		os.Exit(5)
	}

	f, err := os.Open(file)

	if err != nil {
		os.Exit(2)
	}

	fileInfo, err := f.Stat()

	if err != nil {
		os.Exit(3)
	}

	if fileInfo.IsDir() {
		fmt.Printf("Please provide file path only, directory names are not accepted.\n")
		os.Exit(1)
	}

	if !fileInfo.Mode().IsRegular() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(4)
	}

	data := make([]byte, fileInfo.Size())
	f.Read(data)
	defer f.Close()

	uploadLocation, err := GetUploadURL()
	if err == nil {
		Upload(fmt.Sprintf("%s", uploadLocation), file, lifeTime, key,bucket)
	}
}

func GetUploadURL() (string, error) {
	url := "http://cloudra.in/uploadurl"
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return fmt.Sprintf("%s", body), err
}

func Upload(url, file, lifeTime, key,bucket string) (err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := os.Open(file)
	if err != nil {
		return
	}
	fw, err := w.CreateFormFile("image", file)
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}

	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("LifeTime", lifeTime)
	req.Header.Set("Key", key)
	req.Header.Set("Bucket", bucket)
	

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}

	if key == "NEWKEY" {
		newKey:=res.Header.Get("Key")
		SaveLocalUserKey(newKey)
	}

	fmt.Printf("%s\n", string(contents))
	return
}

func SaveLocalUserKey(key string) {
	f, err := os.OpenFile(GetKeyFileName(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err == nil {
		f.Write([]byte(key))
	} else {
		fmt.Printf("%s\n", err.Error())
	}
	defer f.Close()
}

func GetKeyFileName() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fileName := ".stored"
	fileName = fmt.Sprintf("%s/%s", usr.HomeDir, fileName)
	return fileName
}

func GetLocalUserKey() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	fileName := ".stored"
	fileName = fmt.Sprintf("%s/%s", usr.HomeDir, fileName)

	if _, err := os.Stat(fileName); err == nil {
		dat, _ := ioutil.ReadFile(fileName)
		key := string(dat)
		return strings.Replace(key, "\n", "", -1)
	} else {
		return "NEWKEY"
	}

}
