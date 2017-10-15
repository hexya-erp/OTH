package main

import (
	"fmt"
	"io/ioutil"

	"os"

	"path/filepath"

	"github.com/beevik/etree"
	"github.com/hexya-erp/OTH/translate"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println(`
OTH usage
---------
Scaffold an Hexya module from an Odoo module

OTH <python-src-dir> <go-dest-dir>

python-src-dir:  directory of an Odoo module
go-dest-dir:     target directory for the Hexya module`)
		os.Exit(1)
	}

	sourceDir, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}
	destDir, err := filepath.Abs(os.Args[2])
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat(destDir); err == nil {
		fmt.Println("Destination directory already exists", destDir)
		os.Exit(1)
	}

	packageName := filepath.Base(destDir)

	var filename string

	content, _ := ioutil.ReadDir(filepath.Join(sourceDir))

	os.Mkdir(filepath.Join(destDir), 0755)
	os.Mkdir(filepath.Join(destDir, "resources"), 0755)

	for _, c := range content {

		if c.IsDir() {

			switch c.Name() {

			case "models":
				modelsDir := filepath.Join(sourceDir, "models")

				filespython, _ := ioutil.ReadDir(modelsDir)

				for _, fp := range filespython {

					read, errr := ioutil.ReadFile(filepath.Join(modelsDir, fp.Name()))
					if errr != nil {
						fmt.Print(errr)
					}

					if fp.Name() != "__init__.py" {
						gocode := translate.TransPyToGo(string(read), packageName)
						filename = fp.Name()

						errw := ioutil.WriteFile(filepath.Join(destDir, filename[:len(filename)-2]+"go"), []byte(gocode), 0644)
						if errw != nil {
							fmt.Print(errw)
						}
					}
				}

			case "views":
				viewsDir := filepath.Join(sourceDir, "views")
				filesxml, _ := ioutil.ReadDir(viewsDir)
				for _, fx := range filesxml {

					doc := etree.NewDocument()
					if err := doc.ReadFromFile(filepath.Join(viewsDir, fx.Name())); err != nil {
						fmt.Print(err)
					}

					xml := translate.TransXML(doc, packageName)

					xml.WriteToFile(filepath.Join(destDir, "resources", fx.Name()))
				}

			case "wizard":
				wizardDir := filepath.Join(sourceDir, "wizard")
				fileswiz, _ := ioutil.ReadDir(wizardDir)
				for _, wiz := range fileswiz {

					if wiz.Name()[len(wiz.Name())-2:] == "py" {

						read, errr := ioutil.ReadFile(filepath.Join(wizardDir, wiz.Name()))
						if errr != nil {
							fmt.Print(errr)
						}

						if wiz.Name() != "__init__.py" {
							gocode := translate.TransPyToGo(string(read), packageName)
							filename = "wizard_" + wiz.Name()

							errw := ioutil.WriteFile(filepath.Join(destDir, filename[:len(filename)-2]+"go"), []byte(gocode), 0644)
							if errw != nil {
								fmt.Print(errw)
							}
						}

					} else if wiz.Name()[len(wiz.Name())-3:] == "xml" {

						doc := etree.NewDocument()
						if err := doc.ReadFromFile(filepath.Join(wizardDir, wiz.Name())); err != nil {
							fmt.Print(err)
						}

						xml := translate.TransXML(doc, packageName)

						xml.WriteToFile(filepath.Join(destDir, "resources", "wizard_"+wiz.Name()))
					}
				}

			case "security":
				securityDir := filepath.Join(sourceDir, "security")
				filessec, _ := ioutil.ReadDir(securityDir)

				for _, sec := range filessec {

					if sec.Name()[len(sec.Name())-3:] == "csv" {

						read, errr := ioutil.ReadFile(filepath.Join(securityDir, sec.Name()))
						if errr != nil {
							fmt.Print(errr)
						}

						csv := translate.TransCSV(string(read), packageName)

						errw := ioutil.WriteFile(filepath.Join(destDir, "security.go"), []byte(csv), 0644)
						if errw != nil {
							fmt.Print(errw)
						}
					}
				}

			}

		}

	}

	hexya := translate.GenerateHexya()

	err = ioutil.WriteFile(filepath.Join(destDir, "000hexya.go"), []byte(hexya), 0644)
	if err != nil {
		fmt.Print(err)
	}

}
