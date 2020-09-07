package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var files []string
	root := "csvs"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, info.Name())
		return nil
	})
	if err != nil {
		panic(err)
	}

	bar := progressbar.Default(int64(len(files)))
	for _, filename := range files {
		fmt.Println(filename)
		mapConvert,err:=readFile(filename)
		if err != nil {
			return
		}
		writeFile(filename,mapConvert)
		bar.Add(1)
	}
}


func readFile(filename string) ([]map[string]string, error){
	file, err := os.Open("csvs/"+filename)
	if err != nil {
		fmt.Println("Couldnt found file",err)
		return nil,err
	}

	r := csv.NewReader(file)

	row, err := r.Read()
	if err != nil {
		log.Fatal(err)
		return nil,err
	}
	var headers []string
	for _, field := range row {
		headers = append(headers,field)
	}
	var mapConvert []map[string]string

	mapConvert=make([]map[string]string,0)
	for {
		row, err = r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var rowConverted map[string]string
		rowConverted=make(map[string]string)

		for i, field := range row {
			pos:=headers[i]
			rowConverted[pos]=field
		}
		mapConvert=append(mapConvert,rowConverted)
	}

	closeReaderErr := file.Close()
	if closeReaderErr != nil {
		return nil,err
	}

	return mapConvert,nil

}

func writeFile(filename string,mapConvert []map[string]string){
	convertedJson, _ := json.Marshal(mapConvert)
	cleanFilename:=strings.Split(filename, ".")
	jsonFile, err := os.Create("jsons/"+cleanFilename[0]+".json")
	if err != nil {
		return
	}

	_,err=jsonFile.WriteString(string(convertedJson))
	if err != nil {
		return
	}

	closeWriterErr:=jsonFile.Close()
	if closeWriterErr != nil {
		return
	}

}