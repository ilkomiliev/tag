// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tag

import (
	"errors"
	"io"
	"strconv"
)

// ReadDSFTags reads DSF metadata from the io.ReadSeeker, returning the resulting
// metadata in a Metadata implementation, or non-nil error if there was a problem.
func ReadDSFTags(r io.ReadSeeker) (Metadata, error) {
	dsd, err := readString(r, 4)
	if err != nil {
		return nil, err
	}
	if dsd != "DSD " {
		return nil, errors.New("expected 'DSD '")
	}

	n2, err := readBytes(r, 8)
	if err != nil {
		return nil, err
	}
	chunkSize := lsb(n2)

	n3, err := readBytes(r, 8)
	if err != nil {
		return nil, err
	}
	fileSize := lsb(n3)

	n4, err := readBytes(r, 8)
	if err != nil {
		return nil, err
	}
	id3Pointer := lsb(n4)

	_, err = r.Seek(int64(id3Pointer), 0)
	if err != nil {
		return nil, err
	}

	m := new(metadataDSF)
	m.fileSize = fileSize
	m.chunkSize = chunkSize

	err = m.readDSFMetadata(r)
	if err != nil {
		return nil, err
	}

	return m, nil
}

type metadataDSF struct {
	id3       Metadata
	fileSize  int
	chunkSize int
}

func (m *metadataDSF) readDSFMetadata(r io.ReadSeeker) error {
	id3, err := ReadID3v2Tags(r)
	if err != nil {
		return err
	}
	m.id3 = id3
	return nil
}

func (m metadataDSF) Format() Format {
	return m.id3.Format()
}

func (m metadataDSF) FileType() FileType {
	return DSF
}

func (m metadataDSF) Title() string {
	return m.id3.Title()
}

func (m metadataDSF) Album() string {
	return m.id3.Album()
}

func (m metadataDSF) Artist() string {
	return m.id3.Artist()
}

func (m metadataDSF) AlbumArtist() string {
	return m.id3.AlbumArtist()
}

func (m metadataDSF) Composer() string {
	return m.id3.Composer()
}

func (m metadataDSF) Year() int {
	return m.id3.Year()
}

func (m metadataDSF) Genre() string {
	return m.id3.Genre()
}

func (m metadataDSF) Track() (int, int) {
	return m.id3.Track()
}

func (m metadataDSF) Disc() (int, int) {
	return m.id3.Disc()
}

func (m metadataDSF) Picture() *Picture {
	return m.id3.Picture()
}

func (m metadataDSF) Lyrics() string {
	return m.id3.Lyrics()
}

func (m metadataDSF) Raw() map[string]interface{} {
	return m.id3.Raw()
}

func hex(i int) string {
	i64 := int64(i)
	var s string
	s = strconv.FormatInt(i64, 16)
	if len(s) == 1 { // otherwise we are loosing the zeroes
		s = "0" + s
	}
	return s
}

func hex2int(hexStr string) int {
	result, _ := strconv.ParseInt(hexStr, 16, 64)
	return int(result)
}

func lsb(n3 []byte) int {
	var s string
	for _, b := range n3 {
		s = hex(int(b)) + s
	}
	res := hex2int(s)
	return res
}