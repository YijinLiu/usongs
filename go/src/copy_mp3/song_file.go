package main

import(
	"hash/fnv"
	"log"
	"os"
	"path/filepath"
	"strings"
)
	
type SongFile struct {
	path string
	size int64
	id   uint64
}

type SongFileList []*SongFile

func (l SongFileList) Len() int { return len(l) }
func (l SongFileList) Less(i, j int) bool { return l[i].id < l[j].id }
func (l SongFileList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

func NewSongFile(name string, size int64) *SongFile {
	fullPath, err := filepath.Abs(name)
	if err != nil {
		log.Fatal(err)
	}
	h := fnv.New64a()
	if _, err := h.Write([]byte(fullPath)); err != nil {
		log.Fatal(err)
	}
	return &SongFile{
		path: fullPath,
		size: size,
		id:   h.Sum64(),
	}
}

func ListMp3Files(f *os.File, sfl []*SongFile) []*SongFile {
	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if fi.Mode().IsDir() {
		fis, err := f.Readdir(0)
		if err != nil {
			log.Fatal(err)
		}
		for _, fi := range fis {
			name := filepath.Join(f.Name(), fi.Name())
			if f, err := os.Open(name); err != nil {
				log.Printf("Failed to open '%s': %v", name, err)
			} else {
				sfl = ListMp3Files(f, sfl)
			}
		}
	} else if fi.Mode().IsRegular() {
		if strings.ToLower(filepath.Ext(f.Name())) == ".mp3" {
			sfl = append(sfl, NewSongFile(f.Name(), fi.Size()))
		}
	}
	return sfl
}