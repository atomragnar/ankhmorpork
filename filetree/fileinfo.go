package filetree

import (
	"io/fs"
	"time"
)

type SerializableFileInfo struct {
	Name    string
	Size    int64
	Mode    fs.FileMode
	ModTime time.Time
	IsDir   bool
}

func NewSerializableFileInfo(info fs.FileInfo) SerializableFileInfo {
	return SerializableFileInfo{
		Name:    info.Name(),
		Size:    info.Size(),
		Mode:    info.Mode(),
		ModTime: info.ModTime(),
		IsDir:   info.IsDir(),
	}
}
