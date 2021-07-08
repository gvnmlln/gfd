package main

import (
	"bufio"
	"dpb"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	pb "gopkg.in/cheggaaa/pb.v1"
)

const tmp = ".tmp"
const refreshRate = time.Millisecond * 100

// DownloadFile downloads a file from the provided URL to the provided filepath
func DownloadFile(filepath string, url string) error {
	// Create a temp file
	out, err := os.Create(filepath + tmp)
	if err != nil {
		return err
	}
	// Make a GET request to the provided URL
	response, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer response.Body.Close()

	// Get the file size and use it to initialise the counter
	fileSize, _ := strconv.Atoi(response.Header.Get("Content-Length"))
	counter := dpb.NewWriteCounter(fileSize)
	counter.Start()
	if _, err = io.Copy(out, io.TeeReader(response.Body, counter)); err != nil {
		out.Close()
		return err
	}
	counter.Finish()
	out.Close()
	if err = os.Rename(filepath+tmp, filepath); err != nil {
		return err
	}
	return nil
}

// ParseFileName accepts a path and returns the file name and file root (file name without extenstion)
// Accepts a directory path 'dir' to download the URLs to and file path to a text file 'list' containing a list of URLs
func ParseFileName(filePath string) (string, string) {
	fileNameAndExt := path.Base(filePath)
	fileNameWithOutExt := strings.TrimSuffix(fileNameAndExt, filepath.Ext(fileNameAndExt))
	return fileNameAndExt, fileNameWithOutExt
}

// ParseFlagsAndArgs parses the CLI flags and args to string values and returns them as string values and an array respectively
func ParseFlagsAndArgs() (*string, *string, []string) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// Flags
	// dir defaults to the current directory if not chosen as an option
	targetDirectory := flag.String("dir", workingDirectory, "the directory you want the file to be downloaded to")
	// list default value is empty
	urlTxt := flag.String("list", "", "a text file containing a list of urls")
	flag.Parse()

	// Args - each arg represents an individual URL to download
	args := flag.Args()

	return targetDirectory, urlTxt, args
}

// ParseURLsFromTextFile parses an array of strings representing URLs from a given text file
func ParseURLsFromTextFile(filepath string) ([]string, error) {
	var urls []string
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if _, err := url.ParseRequestURI(line); err != nil {
			fmt.Printf("\"%s\" is not a valid URL.\n", line)
			continue
		} else {
			urls = append(urls, line)
		}
	}
	if len(urls) == 0 {
		fmt.Println("No valid URLs were provided in the list - no files were downloaded.")
		return nil, nil
	}
	return urls, nil
}

// URLIsValid accepts a URL and returns true if it is valid and false if it is not
func URLIsValid(u string) bool {
	if _, err := url.ParseRequestURI(u); err != nil {
		return false
	}
	return true
}

// FileExists accepts a path and returns true if the path exists or points to a file, otherwise returns false
func FileExists(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}

// DownloadURLs accepts an array of strings representing URLs to download and a target directory to download them to
func DownloadURLs(urls []string, targetDirectory string) {

	for _, u := range urls {
		// Check if each of the provided URLs is valid
		if !URLIsValid(u) {
			fmt.Printf("Cannot download file from \"%s\" as it is not a valid URL.\n", u)
			continue
		}

		fileName, fileRootName := ParseFileName(u)
		downloadTargetFilePath := path.Join(targetDirectory, fileRootName, fileName)
		downloadDirectory := path.Dir(downloadTargetFilePath)

		// Check the file doesn't already exist at the target filepath
		fmt.Println("Downloading:", fileName)
		if !FileExists(downloadDirectory) {
			err := os.MkdirAll(downloadDirectory, 0755)
			if err != nil {
				log.Panicln(err)
			}
		}

		// Check if the file was downloaded successfully, throw an error if not
		err := DownloadFile(downloadTargetFilePath, u)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println("Finished downloading:", fileName)
	}
}

type WriteCounter struct {
	bytesRead   int
	progressBar *pb.ProgressBar
}

func NewWriteCounter(total int) *WriteCounter {
	bar := pb.New(total)
	bar.SetRefreshRate(refreshRate)
	bar.ShowTimeLeft = true
	bar.SetUnits(pb.U_BYTES)
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

func main() {
	// Get the target directory, the URL list text file path and any URLS passed as args
	directory, urlTxtPath, urls := ParseFlagsAndArgs()
	numberOfURLArgs := len(urls)
	urlList := *urlTxtPath
	targetDirectory := *directory

	// If no URLs are provided as args or via a list, close the program
	if numberOfURLArgs == 0 && urlList == "" {
		fmt.Println("No URL(s) provided. Closing file-downloader.")
		os.Exit(0)
	}

	// If a URL list text file is provided, add the URLs to the list to be downloaded
	if urlList != "" && FileExists(urlList) {
		urlsFromList, _ := ParseURLsFromTextFile(urlList)
		urls = append(urls, urlsFromList...)
	}

	// Download the URLs
	if len(urls) != 0 {
		DownloadURLs(urls, targetDirectory)
	}

}
