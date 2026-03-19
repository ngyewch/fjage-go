package shell

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ngyewch/fjage-go"
	"github.com/ngyewch/fjage-go/gateway"
)

type GetFileRsp struct {
	fjage.Message

	Directory bool           `json:"dir"` // TODO directory
	Filename  string         `json:"filename"`
	Offset    int64          `json:"ofs"` // TODO offset
	Contents  *gateway.Array `json:"contents"`
}

type DirEntry struct {
	Name         string
	Size         int64
	LastModified time.Time
}

func (rsp GetFileRsp) DirEntries() ([]DirEntry, error) {
	if !rsp.Directory {
		return nil, fmt.Errorf("%s is not a directory", rsp.Filename)
	}

	s := string(rsp.Contents.Data)
	lines := strings.Split(s, "\n")

	var entries []DirEntry
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) != 3 {
			return nil, fmt.Errorf("malformed directory listing")
		}
		name := parts[0]
		size, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("malformed directory listing: %w", err)
		}
		lastModified, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("malformed directory listing: %w", err)
		}
		entries = append(entries, DirEntry{
			Name:         name,
			Size:         size,
			LastModified: time.UnixMilli(lastModified),
		})
	}

	return entries, nil
}
