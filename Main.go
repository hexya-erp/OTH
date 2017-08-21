package main

import (
	"fmt"
	"io/ioutil"

	"github.com/beevik/etree"
	"github.com/hexya-erp/OTH/Translate"
)

func main() {

	filespy, _ := ioutil.ReadDir("OTH/SourcePython")
	for _, f := range filespy {

		read, errr := ioutil.ReadFile("OTH/SourcePython/" + f.Name())
		if errr != nil {
			fmt.Print(errr)
		}

		gocode := string(read)
		gocode = Translate.TransPyToGo(gocode)

		errw := ioutil.WriteFile("OTH/ResultGo/"+f.Name()[:len(f.Name())-2]+"go", []byte(gocode), 0644)
		if errw != nil {
			fmt.Print(errw)
		}

	}

	filesxml, _ := ioutil.ReadDir("OTH/SourceXML")
	for _, f := range filesxml {

		doc := etree.NewDocument()
		if err := doc.ReadFromFile("OTH/SourceXML/" + f.Name()); err != nil {
			fmt.Print(err)
		}

		Translate.TransXML(doc, f.Name())
	}

}
