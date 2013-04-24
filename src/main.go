package main

import (
	"./id3go"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func VisitFile(path string, f *os.FileInfo) {
	printTag(path)
}

func VisitDir(path string, f os.FileInfo, er error) error {
	fmt.Println(path)
	return nil
}

func printTag(filename string) {
	fmt.Println(filename)

	res, err := id3go.ReadId3V1Tag(filename)
	if err != nil {
		log.Print(err)
	}

	fmt.Println("Title:", res.Title)
	fmt.Println("Artist:", res.Artist)
	fmt.Println("Album:", res.Album)
	fmt.Println("Comment:", res.Comment)
	fmt.Println()
}

func main() {
	flag.Parse()

	for _, filename := range flag.Args() {
		finfo, err := os.Stat(filename)

		if err != nil {
			log.Print(err)
			continue
		}

		// file
		if !finfo.IsDir() {
			res, _ := id3go.ReadId3V1Tag(filename) //printTag(filename)
			res.Title = "HELLO THIS IS TEST"
			err := id3go.WriteId3V1Tag(filename, res)
			if err != nil {
				fmt.Println(err)
			}

		} else { // Folder
			errChan := make(chan error, 64)
			filepath.Walk(filename, VisitDir)
			select {
			case err := <-errChan:
				log.Print(err)
			default:
			}
		}
	}
}
