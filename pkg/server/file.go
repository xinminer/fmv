package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"io"
	"log"
	"net"
	"os"
)

const toBytes = 1048576

type FileInfo struct {
	fileSize int64
	nameSize int64
	fileName string
}

type FileServer struct {
	conn         net.Conn
	dest         string
	info         FileInfo
	file         *os.File
	fileContents []byte
}

func NewFileServer(conn net.Conn, dest string, chunkSize int) *FileServer {
	return &FileServer{
		conn:         conn,
		dest:         dest,
		fileContents: make([]byte, chunkSize*toBytes),
		info: FileInfo{
			fileSize: 0,
			nameSize: 0,
			fileName: "",
		},
	}
}

func (fs *FileServer) Close() {
	_ = fs.file.Close()
	_ = fs.conn.Close()
}

func (fs *FileServer) fetchFileInfoSizes() error {
	err := binary.Read(fs.conn, binary.LittleEndian, &fs.info.fileSize)
	if err != nil {
		return err
	}
	err = binary.Read(fs.conn, binary.LittleEndian, &fs.info.nameSize)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileServer) fetchFileInfoName() error {
	fileName := new(bytes.Buffer)
	for {
		n, err := io.CopyN(fileName, fs.conn, fs.info.nameSize)
		if err != nil {
			return err
		}
		if n == fs.info.nameSize {
			break
		}
	}
	fs.info.fileName = fileName.String()
	return nil
}

func (fs *FileServer) initFileInfo() {
	if err := fs.fetchFileInfoSizes(); err != nil {
		return
	}
	if err := fs.fetchFileInfoName(); err != nil {

	}
}

func (fs *FileServer) createNewFile() {
	// Create destination file
	fName := fs.dest + string(os.PathSeparator) + fs.info.fileName

	var err error
	fs.file, err = os.Create(fName)
	if err != nil {
		log.Panicln("Could not create destination file in: ", fName)
	}
}

func (fs *FileServer) parseHeader() {
	// File format is:
	// 64 bits: File size = Z
	// 64 bits: Filename size = N
	// N  bits: Filename string
	// Z  bits: File contents

	fs.initFileInfo()

	log.Printf("Receiving file [%s]\n\t| Size: %.2fMB\n", fs.info.fileName, float64(fs.info.fileSize)/1_000_000)

	fs.createNewFile()
}

func (fs *FileServer) HandleFile() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("TERMINATED READING FILE")
		}
	}()
	defer fs.Close()

	fs.parseHeader()
	fs.parseBody()
}

func (fs *FileServer) parseBody() {
	// Get file contents
	var totalRead int64
	var totalWritten int64
	var totalBuffer int

	for totalWritten < fs.info.fileSize {
		for totalBuffer < cap(fs.fileContents) {
			n, err := fs.conn.Read(fs.fileContents[totalBuffer:])
			if err != nil {
				if err != io.EOF {
					log.Panicf("Could not read file contents for file [%s]\n", fs.info.fileName)
				}
			}
			totalBuffer += n
			totalRead += int64(n)
			if totalRead == fs.info.fileSize {
				break
			}
		}
		totalBufferLimit := totalBuffer
		totalBuffer = 0

		log.Printf("File [%s] \n\t| %.2f%%", fs.info.fileName, (float64(totalRead)/float64(fs.info.fileSize))*100.0)

		for totalBuffer < cap(fs.fileContents) {
			n, err := fs.file.Write(fs.fileContents[totalBuffer:totalBufferLimit])
			totalBuffer += n
			totalWritten += int64(n)
			if err != nil {
				log.Panicf("Could not write file contents for file [%s]\n", fs.info.fileName)
			}
			if totalWritten == fs.info.fileSize {
				break
			}
		}
		totalBuffer = 0
	}
	fs.conn.Write([]byte("y"))
	log.Printf("File [%s] DONE\n", fs.info.fileName)
	tmpName := fs.dest + string(os.PathSeparator) + fs.info.fileName
	finName := fs.dest + string(os.PathSeparator) + gstr.Replace(fs.info.fileName, ".fmv", "")
	_ = gfile.Move(tmpName, finName)
}
