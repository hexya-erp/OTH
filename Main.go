package main

import (
	"fmt"
	"io/ioutil"
	"github.com/hexya-erp/OTH/Translate"
	"github.com/beevik/etree"
)

func main() {

	read, errr := ioutil.ReadFile("SourcePython/file.py")
	if errr != nil {
		fmt.Print(errr)
	}

	gocode := string(read)
	gocode = Translate.TransPyToGo(gocode)

	errw := ioutil.WriteFile("ResultGo/file.go", []byte(gocode), 0644)
	if errw != nil {
		fmt.Print(errw)
	}



	doc := etree.NewDocument()
	if err := doc.ReadFromFile("SourceXML/file.xml");err != nil{
		fmt.Print(err)
	}

	 Translate.TransXML(doc)

	//errxml := ioutil.WriteFile("ResultXML/file.xml", []byte(xml), 0644)
	//if errxml != nil {
	//	fmt.Print(errxml)
	//}


}