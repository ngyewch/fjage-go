package shell

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"slices"

	"github.com/ngyewch/fjage-go/services/shell"
)

type ProgressHandler func(current int64, total int64)

type Helper struct {
	client  *Client
	options HelperOptions
	httpUrl *url.URL
}

var (
	DefaultHelperOptions = HelperOptions{
		CopyBufferSize:         16384,
		HttpClient:             http.DefaultClient,
		HttpGetFileDirectories: []string{"logs"},
	}
)

type HelperOptions struct {
	CopyBufferSize         int64
	HttpClient             *http.Client
	HttpGetFileDirectories []string
}

func NewHelper(client *Client, options *HelperOptions) (*Helper, error) {
	if options == nil {
		options = &DefaultHelperOptions
	}
	transportUrl, err := url.Parse(client.gw.Transport().Url())
	if err != nil {
		return nil, err
	}
	httpScheme := ""
	if transportUrl.Scheme == "ws" {
		httpScheme = "http"
	} else if transportUrl.Scheme == "wss" {
		httpScheme = "https"
	}
	var httpUrl *url.URL
	if httpScheme != "" {
		httpUrl, err = url.Parse(fmt.Sprintf("%s://%s", httpScheme, transportUrl.Host))
		if err != nil {
			return nil, err
		}
	}
	return &Helper{
		client:  client,
		options: *options,
		httpUrl: httpUrl,
	}, nil
}

func (helper *Helper) ListFiles(ctx context.Context, remotePath string) ([]shell.DirEntry, error) {
	response, err := helper.client.GetFile(ctx, remotePath, 0, 0)
	if err != nil {
		return nil, err
	}

	return response.DirEntries()
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

	if (helper.httpUrl != nil) && slices.Contains(helper.options.HttpGetFileDirectories, path.Dir(remotePath)) {
		relativeUrl, err := url.Parse(remotePath)
		if err != nil {
			return err
		}
		remoteUrl := helper.httpUrl.ResolveReference(relativeUrl)
		httpRequest, err := http.NewRequest("GET", remoteUrl.String(), nil)
		if err != nil {
			return err
		}
		httpResponse, err := helper.options.HttpClient.Do(httpRequest)
		if err != nil {
			return err
		}
		defer func(body io.ReadCloser) {
			_ = body.Close()
		}(httpResponse.Body)
		if httpResponse.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected HTTP status: %s", httpResponse.Status)
		}

		f, err := os.Create(localPath)
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)
		var offset int64
		buffer := make([]byte, 32*1024)
		for {
			readLen, err := httpResponse.Body.Read(buffer)
			if readLen > 0 {
				_, err = f.Write(buffer[0:readLen])
				if err != nil {
					return err
				}
				offset += int64(readLen)
				if progressHandler != nil {
					progressHandler(offset, fileSize)
				}
			}
			if err != nil {
				if err == io.EOF {
					break
				} else {
					return err
				}
			}
		}
	} else {
		copyBufferSize := helper.options.CopyBufferSize
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

	copyBufferSize := helper.options.CopyBufferSize
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

func (helper *Helper) DeleteFile(ctx context.Context, remotePath string) error {
	err := helper.client.PutFile(ctx, remotePath, 0, nil)
	if err != nil {
		return err
	}
	return nil
}

func (helper *Helper) ExecuteCommand(ctx context.Context, command string) error {
	return helper.client.ExecuteCommand(ctx, command)
}

func (helper *Helper) ExecuteScript(ctx context.Context, scriptFile string, scriptArgs []string) error {
	return helper.client.ExecuteScript(ctx, scriptFile, scriptArgs)
}
