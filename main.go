package main

import (
	"fmt"
	"io/ioutil"

	"os"

	"github.com/beevik/etree"
	"github.com/hexya-erp/OTH/translate"
)

func main() {

	var gocode string
	var filename string
	root, _ := ioutil.ReadDir("OTH/source")

	for _, r := range root {

		content, _ := ioutil.ReadDir("OTH/source/" + r.Name())

		os.Mkdir("OTH/result/"+r.Name(), os.FileMode(0775))
		os.Mkdir("OTH/result/"+r.Name()+"/resources", os.FileMode(0775))

		for _, c := range content {

			if c.IsDir() {

				switch c.Name() {

				case "models":

					filespython, _ := ioutil.ReadDir("OTH/source/" + r.Name() + "/models")

					for _, fp := range filespython {

						read, errr := ioutil.ReadFile("OTH/source/" + r.Name() + "/models/" + fp.Name())
						if errr != nil {
							fmt.Print(errr)
						}

						if fp.Name() != "__init__.py" {
							gocode = string(read)
							gocode = translate.TransPyToGo(gocode, r.Name())
							filename = fp.Name()

							errw := ioutil.WriteFile("OTH/result/"+r.Name()+"/"+filename[:len(filename)-2]+"go", []byte(gocode), 0644)
							if errw != nil {
								fmt.Print(errw)
							}
						}
					}

				case "views":

					filesxml, _ := ioutil.ReadDir("OTH/source/" + r.Name() + "/views")
					for _, fx := range filesxml {

						doc := etree.NewDocument()
						if err := doc.ReadFromFile("OTH/source/" + r.Name() + "/views/" + fx.Name()); err != nil {
							fmt.Print(err)
						}

						xml := translate.TransXML(doc)

						xml.WriteToFile("OTH/result/" + r.Name() + "/resources/" + fx.Name())
					}

				case "wizard":

					fileswiz, _ := ioutil.ReadDir("OTH/source/" + r.Name() + "/wizard")
					for _, wiz := range fileswiz {

						if wiz.Name()[len(wiz.Name())-2:] == "py" {

							read, errr := ioutil.ReadFile("OTH/source/" + r.Name() + "/wizard/" + wiz.Name())
							if errr != nil {
								fmt.Print(errr)
							}

							if wiz.Name() != "__init__.py" {
								gocode = string(read)
								gocode = translate.TransPyToGo(gocode, r.Name())
								filename = "wizard_" + wiz.Name()

								errw := ioutil.WriteFile("OTH/result/"+r.Name()+"/"+filename[:len(filename)-2]+"go", []byte(gocode), 0644)
								if errw != nil {
									fmt.Print(errw)
								}
							}

						} else if wiz.Name()[len(wiz.Name())-3:] == "xml" {

							doc := etree.NewDocument()
							if err := doc.ReadFromFile("OTH/source/" + r.Name() + "/wizard/" + wiz.Name()); err != nil {
								fmt.Print(err)
							}

							xml := translate.TransXML(doc)

							xml.WriteToFile("OTH/result/" + r.Name() + "/resources/wizard_" + wiz.Name())
						}
					}

				}

			}

		}

		hexya := translate.GenerateHexya()

		err := ioutil.WriteFile("OTH/result/"+r.Name()+"/000-hexya.go", []byte(hexya), 0644)
		if err != nil {
			fmt.Print(err)
		}

	}
}
