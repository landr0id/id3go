package id3go

import (
	"bytes"
	"errors"
	"os"
)

func byteString(b []byte) string {
	pos := bytes.IndexByte(b, 0)

	if pos == -1 {
		pos = len(b)
	}

	return string(b[0:pos])
}

func ReadId3V1Tag(filename string) (Id3V1Tag, error) {
	buff_ := make([]byte, tagSize)

	f, err := os.Open(filename)
	defer f.Close()

	if err != nil {
		return Id3V1Tag{}, err
	}

	// Read last 128 bytes of file to see ID3 tag
	f.Seek(-tagSize, 2)
	f.Read(buff_)

	// First 3 characters are static "TAG" 
	if byteString(buff_[0:tagStart]) != "TAG" {
		return Id3V1Tag{}, errors.New("No ID3 tag found")
	}

	buff := buff_[tagStart:]

	id3tag := Id3V1Tag{byteString(buff[0:titleEnd]), byteString(buff[titleEnd:artistEnd]), byteString(buff[artistEnd:albumEnd]), byteString(buff[albumEnd:yearEnd]), byteString(buff[yearEnd:commentEnd]), "", 0xFF, 0xFF}

	// Special case. If next-to-last comment byte is zero, then the last
	// comment byte is the track number
	if buff[commentEnd-2] == 0 {
		id3tag.Track = buff[commentEnd-1]
	}
	id3tag.Genre = buff[commentEnd]
	id3tag.GenreName = codeToName[id3tag.Genre]

	return id3tag, nil
}

func WriteId3V1Tag(filename string)
