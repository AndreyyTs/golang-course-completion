package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	var prefix string = "│"
	if printFiles {
		err = drawTreeWithFiles(out, path, prefix)
	} else {
		err = drawTreeWithoutFiles(out, path, prefix)
	}

	return err
}

func drawTreeWithFiles(out io.Writer, path string, prefix string) (err error) {
	files, err1 := ioutil.ReadDir(path)
	if err1 != nil {
		return err1
	}

	//сортировка по имени
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for i, file := range files {

		var str string
		str = strings.TrimSuffix(prefix, "│")
		if i == len(files)-1 {
			str = str + "└───"
		} else {
			str = str + "├───"
		}

		str = str + file.Name()
		if !file.IsDir() {
			if file.Size() == 0 {
				str = str + " (empty)"
			} else {
				str = str + " (" + strconv.FormatInt(file.Size(), 10) + "b)"
			}
		}
		str = str + "\n"
		fmt.Fprint(out, str)
		if file.IsDir() {
			if i == len(files)-1 {
				prefix = strings.TrimSuffix(prefix, "│")
			}
			drawTreeWithFiles(out, path+"/"+file.Name(), prefix+"\t│")
		}

	}
	return nil
}

func drawTreeWithoutFiles(out io.Writer, path string, prefix string) (err error) {
	files, err1 := ioutil.ReadDir(path)
	if err1 != nil {
		return err1
	}

	//сортировка по имени
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	var countDir int = func(files []os.FileInfo) (countDir int) {
		countDir = 0
		for _, file := range files {
			if file.IsDir() {
				countDir++
			}
		}
		return countDir
	}(files)
	var numberDir = 0

	for _, file := range files {

		var str string
		str = strings.TrimSuffix(prefix, "│")
		if numberDir == countDir-1 {
			str = str + "└───"
		} else {
			str = str + "├───"
		}
		str = str + file.Name() + "\n"

		if file.IsDir() {
			if numberDir == countDir-1 {
				prefix = strings.TrimSuffix(prefix, "│")
			}
			fmt.Fprint(out, str)
			drawTreeWithoutFiles(out, path+"/"+file.Name(), prefix+"\t│")
			numberDir++
		}

	}
	return nil
}
