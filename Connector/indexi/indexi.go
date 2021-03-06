package indexi

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func GetRefreshedMusicList(drives []string) MusicFiles {
	// start
	cmd := exec.Command("eCLink.exe")
	if err := cmd.Start(); err != nil {
		panic(err)
	}

	// wait or timeout
	donec := make(chan error, 1)
	go func() {
		donec <- cmd.Wait()
	}()
	select {
	case <-time.After(25 * time.Second):
		return MusicFiles{}
	case <-donec:
		os.Remove("allMusicFiles.csv")
		os.Rename("allMusic.csv" , "allMusicFiles.csv")
		return GetMusicList(drives)
	}
}

func GetMusicList(drives []string) MusicFiles {
	csvFile, _ := os.Open("allMusicFiles.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	var musicFiles MusicFiles
	for {
		line, e := reader.Read()
		if e == io.EOF {
			break
		} else if e != nil {
			log.Fatal(e)
		}

		if SliceContains(drives , string(strings.ToUpper(line[1])[0])) {
			musicFiles = append(musicFiles, MusicStruct{
				Fname: line[0],
				Fpath: line[1],
			})
		}

	}

	return musicFiles
}

func SliceContains (a []string , i string) bool {
	for _ , value := range a {
		if string(value) == i {
			return true
		}
	}
	return false
}

func FileExt (filename string) string {
	if len(filename) == 0 {
		return ""
	}

	dotIndex := strings.LastIndex(filename , ".")
	return filename[dotIndex:]
}

func AvaliableFileSystems () []string {
	return getDrives()
}

func getDrives() (r []string) {
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		_, err := os.Open(string(drive)+":\\")
		if err == nil {
			r = append(r, string(drive))
		}
	}
	return
}

func GetFileDetails (filepath string) (os.FileInfo , error){

	finfo , err := os.Stat(filepath)
	if err != nil {
		log.Fatal(err)
		return nil , err
	}

	return finfo , nil
}