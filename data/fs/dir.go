package fs

import (
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/zurrty/cranky/util"
)

var FORMATS = []string{"gif", "jpg", "jpeg", "png", "webp"}

type Directory struct {
	RootPath string
	Files    []string
	Index    int
}

func (dir *Directory) GetFile() string {
	if dir.Index < len(dir.Files) && dir.Index >= 0 {
		return dir.Files[dir.Index]
	}
	return ""
}

func (dir *Directory) SetIndex(newIdx int) {
	if newIdx >= len(dir.Files) {
		dir.Index = newIdx - dir.Index - 1
	} else if newIdx < 0 {
		dir.Index = len(dir.Files) + newIdx // adding a negative aka subtracting
	} else {
		dir.Index = newIdx
	}
}
func (dir *Directory) IndexOf(filename string) int {
	for i, f := range dir.Files {
		if filename == f {
			return i
		}
	}
	return 0 // didnt find it... TOO BAD!!! YOU GET ZERO!!!
}

func OpenDirectory(dirPath string) Directory {
	files, err := ioutil.ReadDir(dirPath)
	util.Assert(err)
	dir := Directory{
		RootPath: dirPath,
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			fileSplit := strings.Split(file.Name(), ".")
			ext := fileSplit[len(fileSplit)-1]
			for _, FMT := range FORMATS {
				if FMT == strings.ToLower(ext) {
					dir.Files = append(dir.Files, file.Name())
					break
				}
			}
		}
	}
	sort.Strings(dir.Files)
	return dir
}

func IsDir(path string) bool {
	f, err := os.Open(path)
	util.Assert(err)
	stat, err := f.Stat()
	util.Assert(err)
	return stat.IsDir()
}
