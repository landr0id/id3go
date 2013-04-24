package id3go

import (
	"bytes"
	"errors"
	"fmt"
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
	if err != nil {
		return Id3V1Tag{}, err
	}
	defer f.Close()

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

func WriteId3V1Tag(filename string, tag Id3V1Tag) error {
	// Make sure the tag is all in order
	if len(tag.Title) > 30 {
		return errors.New("Invalid tag value: Title (length > 30)")
	} else if len(tag.Artist) > 30 {
		return errors.New("Invalid tag value: Artist (length > 30)")
	} else if len(tag.Comment) > 30 || (tag.Track != 0 && len(tag.Comment) > 28) {
		return errors.New("Invalid tag value: Comment (length > 30, or 28 (if track is set to non-zero value))")
	} else if len(tag.Year) > 4 {
		return errors.New("Invalid tag value: Year (length > 4)")
	}
	// This buffer will hold the data that's being written
	buff := make([]byte, tagSize)

	buffer := bytes.NewBuffer(buff)
	buffer.WriteString("TAG")
	buffer.WriteString(tag.Title)
	// write the difference of null bytes
	fmt.Println(len(tag.Title))
	buffer.Write(make([]byte, 30-len(tag.Title)))
	buffer.WriteString(tag.Artist)
	buffer.Write(make([]byte, 30-len(tag.Artist)))
	buffer.WriteString(tag.Album)
	buffer.Write(make([]byte, 30-len(tag.Album)))
	buffer.WriteString(tag.Year)
	buffer.Write(make([]byte, 4-len(tag.Year)))
	buffer.WriteString(tag.Comment)
	buffer.Write(make([]byte, 28-len(tag.Comment)))
	if tag.Track != 0 {
		buffer.WriteByte(byte(0))
		buffer.WriteByte(tag.Track)
	} else {
		buffer.WriteByte(byte(1))
		buffer.WriteByte(byte(0))
	}
	buffer.WriteByte(tag.Genre)

	f, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Seek(-tagSize, 2)
	buffer.WriteTo(f)
	return nil
}
