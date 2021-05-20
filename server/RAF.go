package main

import (
	"fmt"
	"os"
)

type RandomAccessFile struct {
	file_name       string
	file_pointer    int64
	file_descriptor *os.File
	file_info       os.FileInfo
}

type MatchRange struct {
	left  int64
	right int64
}

func (self *RandomAccessFile) close() error {
	return self.file_descriptor.Close()
}

func (self *RandomAccessFile) clear() {
	self.seek(0)
	self.file_descriptor.Truncate(0)
	self.commit()
}

func (self *RandomAccessFile) getRange(target string) *MatchRange {
	var btarget []byte = []byte(target)
	var file_content []byte = self.readFrom(0, 0)
	for h := range file_content {
		for k, b2 := range btarget {
			if (h+k) >= len(file_content) || b2 != file_content[h+k] {
				break
			}
			if k == len(btarget)-1 {
				return &MatchRange{left: int64(h), right: int64(h + k)}
			}
		}
	}
	return &MatchRange{left: -1, right: -1}
}

func (self *RandomAccessFile) seek(file_pointer int64) (err error) {
	if file_pointer <= self.file_info.Size() {
		self.file_pointer = file_pointer
	} else {
		err = fmt.Errorf("file pointer to big for file with length %d", self.file_info.Size())
	}
	return err
}

func (self *RandomAccessFile) commit() {
	self.file_info, _ = self.file_descriptor.Stat()
}

func (self *RandomAccessFile) readChar() byte {
	var content []byte = make([]byte, 1)
	self.file_descriptor.ReadAt(content, self.file_pointer)
	return content[0]
}

func (self *RandomAccessFile) readFrom(start int64, end int64) []byte {
	var file_content []byte

	if end == 0 {
		end = self.file_info.Size()
	}

	if end <= start {
		return make([]byte, 0)
	}

	file_content = make([]byte, end-start)
	self.file_descriptor.ReadAt(file_content, start)
	return file_content
}

func (self *RandomAccessFile) Size() int64 {
	return self.file_info.Size()
}

func (self *RandomAccessFile) truncate(new_size int64, inplace bool) {
	self.file_descriptor.Truncate(new_size)
	if inplace {
		self.file_descriptor.Sync()
	}
	self.commit()
}

func (self *RandomAccessFile) write(content string) (err error) {
	_, err = self.file_descriptor.WriteAt([]byte(content), self.file_pointer)
	self.commit()
	return
}

func createRAF(filename string) *RandomAccessFile {
	var new_raf *RandomAccessFile = new(RandomAccessFile)
	new_raf.file_name = filename
	new_raf.file_pointer = 0

	if !pathExists(filename) {
		f, _ := os.OpenFile(filename, os.O_CREATE, 0777)
		f.Close()
	}

	descriptor, err := os.OpenFile(filename, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	new_raf.file_descriptor = descriptor
	new_raf.file_info, _ = descriptor.Stat()
	return new_raf
}
