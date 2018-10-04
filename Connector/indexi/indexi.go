package indexi

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type SingleMusic struct {
	Fname string `json:"filename"`
	Fpath string `json:"filepath"`
}

func GetRefreshedMusicList() []SingleMusic {
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
		return []SingleMusic{}
	case <-donec:
		os.Rename("allMusic.csv" , "allMusicFiles.csv")
		return GetMusicList()
	}
}

func GetMusicList() []SingleMusic {
	csvFile, _ := os.Open("allMusicFiles.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	var musicFiles []SingleMusic
	for {
		line, e := reader.Read()
		if e == io.EOF {
			break
		} else if e != nil {
			log.Fatal(e)
		}
		musicFiles = append(musicFiles, SingleMusic{
			Fname: line[0],
			Fpath: line[1],
		})
	}

	return musicFiles

}

func MusicListToJson (list []SingleMusic) string {
	filesJson, _ := json.Marshal(list)
	return string(filesJson)
}

func FileExt (filename string) string {
	if len(filename) == 0 {
		return ""
	}

	dotIndex := strings.LastIndex(filename , ".")
	return filename[dotIndex:]
}

func GetFileDetails (filepath string) (os.FileInfo , error){

	finfo , err := os.Stat(filepath)
	if err != nil {
		log.Fatal(err)
		return nil , err
	}

	return finfo , nil
}
