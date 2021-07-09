![Go](https://github.com/gvnmlln/go-file-downloader/actions/workflows/go.yml/badge.svg)

# go-file-downloader
Go CLI program that downloads files at provided URLs to a provided directory.

# Setup

This project requires Go in order to build it from the source code. If you don't have Go installed on your machine, you can follow [these instructions](https://golang.org/doc/install) to get started.

Once you have Go installed on your machine:
1. Download the source code from this repository.
2. Run ```go build file_downloader``` in the directory in source code directory.
This will create a binary file that you can use in the CLI to download URLs.

To download one or more URLs simply run:
Run ```./file_downloader <URL_1> <URL_2>```

The download directory defaults to the directory in which the binary is found.
To specifiy a particular target directory use the ```dir``` flag, e.g. ```./file_downloader -dir=<Target Directory> <URL>```

To download a batch of URLs at once you can specify any number of URLs as CLI args or alternatively you can list the URLs in a text file and point to it use the ```list``` flag, e.g. ```./file_downloader -list=path/to/textfile.txt``` - the files will be downloaded to the binary's working directory or the provided ```dir```.