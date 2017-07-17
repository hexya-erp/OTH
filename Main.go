package main

import (
	"fmt"
	"io/ioutil"
	"OTH/Translate"
)

func main() {

	read, err := ioutil.ReadFile("SourcePython/file.py")
	if err != nil {
	fmt.Print(err)
	}

	str:= string(read)
	str = Translate.TransPyToGo(str)

	 error := ioutil.WriteFile("ResultGo/file.go" , []byte(str) , 0644 )
	if error != nil {
		fmt.Print(error)
	}

}