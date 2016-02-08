package main

import(
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

var srcDirFlag = flag.String("src-dir", "", "The folder containing your mp3 files")
var dstDirFlag = flag.String("dst-dir", "", "The folder to copy file to")

func readLastId(dir string) uint64 {
	lastIdFile := filepath.Join(dir, "lastid.txt")
	content, err := ioutil.ReadFile(lastIdFile)
	if err != nil {
		log.Printf("Failed to read '%s': %v", lastIdFile, err)
		return 0
	}
	id, err := strconv.ParseUint(string(content), 16, 64)
	if err != nil {
		log.Print("Invalid lastid.txt file content: ", string(content))
		return 0
	}
	return id
}

func writeLastId(dir string, lastId uint64) {
	lastIdFile := filepath.Join(dir, "lastid.txt")
	str := strconv.FormatUint(lastId, 16)
	if err := ioutil.WriteFile(lastIdFile, []byte(str), os.FileMode(0644)); err != nil {
		log.Printf("Failed to write to '%s': %v", lastIdFile, err)
	}
}

func main() {
	flag.Parse()
	
	srcDir := *srcDirFlag
	dstDir := *dstDirFlag
	if len(srcDir) == 0 || len(dstDir) == 0 {
		log.Fatal("Both --src-dir and --dst-dir are required")
	}
	
	sd, err := os.Open(srcDir)
	if err != nil {
		log.Fatalf("Failed to open '%s': %v", srcDir, err)
	}	
	
	sfl := ListMp3Files(sd, nil)
	log.Printf("Total %d mp3 files found.", len(sfl))
	sort.Sort(SongFileList(sfl))
	
	lastId := readLastId(dstDir)
	start := 0
	for start < len(sfl) {
		if sfl[start].id > lastId {
			break
		}
		start++
	}
	
	i := start
	for {
		if _, free, err := DiskSpace(dstDir); err != nil {
			log.Fatal(err)
		} else if free < sfl[i].size {
			log.Printf("Only %d left, '%s' needs '%d'. Exiting..", free, sfl[i].path, sfl[i].size)
			break
		}
		dstPath := filepath.Join(dstDir, filepath.Base(sfl[i].path))
		if err := CopyFile(sfl[i].path, dstPath); err != nil {
			log.Fatalf("Failed to copy '%s' to '%s': %v", sfl[i].path, dstPath, err)
		}
		lastId = sfl[i].id

		i++
		if i == len(sfl) {
			i = 0
		}
		if i == start {
			break
		}
	}
	writeLastId(dstDir, lastId)
}