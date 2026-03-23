package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/ngyewch/fjage-go/services/shell"
)

type ProgressHandler func(current int64, total int64)

type Helper struct {
	client         *Client
	copyBufferSize int64
}

func NewHelper(client *Client, copyBufferSize int64) *Helper {
	return &Helper{
		client:         client,
		copyBufferSize: copyBufferSize,
	}
}

func (helper *Helper) GetFile(ctx context.Context, remotePath string, localPath string, progressHandler ProgressHandler) error {
	remoteDir := path.Dir(remotePath)
	if remoteDir == "" {
		remoteDir = "."
	}
	var dirEntry *shell.DirEntry
	{
		getFileRsp, err := helper.client.GetFile(ctx, remoteDir, 0, 0)
		if err != nil {
			return err
		}
		dirEntries, err := getFileRsp.DirEntries()
		if err != nil {
			return err
		}
		remoteBase := path.Base(remotePath)
		for _, entry := range dirEntries {
			if entry.Name == remoteBase {
				dirEntry = &entry
				break
			}
		}
	}
	if dirEntry == nil {
		return fmt.Errorf("remote file not found")
	}

	f, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	fileSize := dirEntry.Size
	if progressHandler != nil {
		progressHandler(0, fileSize)
	}

	copyBufferSize := helper.copyBufferSize
	if copyBufferSize <= 0 {
		copyBufferSize = fileSize
	}

	var offset int64 = 0
	for offset < fileSize {
		length := copyBufferSize
		if offset+length > dirEntry.Size {
			length = dirEntry.Size - offset
		}
		response, err := helper.client.GetFile(ctx, remotePath, offset, length)
		if err != nil {
			return err
		}

		if response.Directory {
			return fmt.Errorf("%s is a directory", remotePath)
		}

		_, err = f.Write(response.Contents)
		if err != nil {
			return err
		}

		offset += int64(len(response.Contents))
		if progressHandler != nil {
			progressHandler(offset, fileSize)
		}
	}

	return nil
}

func (helper *Helper) PutFile(ctx context.Context, localPath string, remotePath string, progressHandler ProgressHandler) error {
	fileInfo, err := os.Stat(localPath)
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()
	if progressHandler != nil {
		progressHandler(0, fileSize)
	}

	copyBufferSize := helper.copyBufferSize
	if copyBufferSize <= 0 {
		copyBufferSize = fileSize
	}

	f, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	var offset int64 = 0
	buffer := make([]byte, copyBufferSize)
	for {
		readLen, err := f.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		err = helper.client.PutFile(ctx, remotePath, offset, buffer[0:readLen])
		if err != nil {
			return err
		}

		offset += int64(readLen)
		if progressHandler != nil {
			progressHandler(offset, fileSize)
		}
	}

	return nil
}
