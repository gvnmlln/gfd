package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	ProgressBar "gopkg.in/cheggaaa/pb.v1"
)

const REFRESH_RATE = time.Millisecond * 100
const TMP_FILE_EXT = ".tmp"

type WriteCounter struct {
	bytesRead   int
	progressBar *ProgressBar.ProgressBar
}

func NewWriteCounter(total int) *WriteCounter {
	bar := ProgressBar.New(total)
	bar.SetRefreshRate(REFRESH_RATE)
	bar.ShowTimeLeft = true
	bar.SetUnits(ProgressBar.U_BYTES)
	return &WriteCounter{
		progressBar: bar,
	}
}

func (writeCounter *WriteCounter) Write(p []byte) (int, error) {
	writeCounter.bytesRead += len(p)
	writeCounter.progressBar.Set(writeCounter.bytesRead)
	return writeCounter.bytesRead, nil
}

func (writeCounter *WriteCounter) Start() {
	writeCounter.progressBar.Start()
}

func (writeCounter *WriteCounter) Finish() {
	writeCounter.progressBar.Finish()
}

// Downloads a file at the provided URL to the provided directory
func DownloadFile(filepath string, url string) error {
	out, err := os.Create(filepath + TMP_FILE_EXT)
	if err != nil {
		return err
	}
	response, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer response.Body.Close()
	fileSize, _ := strconv.Atoi(response.Header.Get("Content-Length"))
	counter := NewWriteCounter(fileSize)
	counter.Start()
	if _, err = io.Copy(out, io.TeeReader(response.Body, counter)); err != nil {
		out.Close()
		return err
	}
	counter.Finish()
	out.Close()
	if err = os.Rename(filepath+TMP_FILE_EXT, filepath); err != nil {
		return err
	}
	return nil
}

// Returns 'filename.extension' and 'filename' from a given url
func ParseFileName(url string) (string, string) {
	fileNameAndExt := path.Base(url)
	fileNameWithOutExt := strings.Split(fileNameAndExt, ".")[0]
	return fileNameAndExt, fileNameWithOutExt
}

func ParseFlagsAndArgs() (*string, []string) {
	workingDirectory, errWD := os.Getwd()
	if errWD != nil {
		log.Println(errWD)
	}
	targetDirectory := flag.String("dir", workingDirectory, "the directory you want the file to be downloaded to")
	flag.Parse()
	args := flag.Args()
	return targetDirectory, args
}

func main() {
	targetDirectory, urls := ParseFlagsAndArgs()
	numberOfURLArgs := len(urls)
	if numberOfURLArgs == 0 {
		fmt.Println("No URL(s) provided. Closing file-downloader.")
		os.Exit(0)
	}

	for _, urlPath := range urls {
		if _, err := url.ParseRequestURI(urlPath); err != nil {
			fmt.Printf("***Cannot download file from \"%s\" as it is not a valid URL.\n", urlPath)
			continue
		}
		fileName, directoryName := ParseFileName(urlPath)
		downloadFilePath := path.Join(*targetDirectory, directoryName, fileName)
		downloadDirectory := path.Dir(downloadFilePath)
		fmt.Println("Downloading:", fileName)
		if _, err := os.Stat(downloadDirectory); os.IsNotExist(err) {
			os.Mkdir(downloadDirectory, 0777)
		}

		err := DownloadFile(downloadFilePath, urlPath)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		fmt.Println("Finished downloading:", fileName)
	}
}
