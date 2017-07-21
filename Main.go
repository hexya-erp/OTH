package main

import (
	"fmt"
	"io/ioutil"
	"github.com/hexya-erp/OTH/Translate"
)

func main() {

	read, errr := ioutil.ReadFile("SourcePython/file.py")
	if errr != nil {
		fmt.Print(errr)
	}

	str := string(read)
	str = Translate.TransPyToGo(str)

	errw := ioutil.WriteFile("ResultGo/file.go", []byte(str), 0644)
	if errw != nil {
		fmt.Print(errw)
	}
}