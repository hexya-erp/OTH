package Hextranslate

import (
	"bufio"
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hexya-erp/hexya/src/tools/fileutils"
)

func GetAllFiles(iPath string) *list.List {
	out := list.New()
	fmt.Printf("Starting walk though input dir. Path: %s\n", iPath)
	err := filepath.Walk(iPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Could not access %q : %v\n", path, err)
			if path == iPath {
				return err
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}
		out.PushBack(path)
		return nil
	})
	if err != nil {
		FatalError("Error walking the input path %q: %v", iPath, err)
	}
	return out
}

func WriteRawOutputFile(filePath string, content *list.List) {
	filePath = string(append(append([]byte(genArgs.outputPath), []byte("/unchanged")...), []byte(strings.TrimPrefix(filePath, genArgs.inputPath))...))
	CreateDir(filepath.Dir(filePath))
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Could not create file %s. %v\n", filePath, err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for line := content.Front(); line != nil; line = line.Next() {
		writer.WriteString(line.Value.(string))
		writer.WriteString("\n")
	}
	writer.Flush()
	fmt.Printf("%s written.\n", filePath)
}

func TranslateDataDemoFile(file string) {
	fileExtension := path.Ext(file)
	switch fileExtension {
	case ".xml":
		TranslateXML(file)
	case ".csv":
		if IsNotInStringArray(strings.Replace(file, genArgs.inputPath+"/", "", -1), append(genArgs.manifest["data"], genArgs.manifest["demo"]...)) {
			NoTranslation(file)
			return
		}
		TranslateCSV(file)
	}
}

func Translator() {
	// DATA FILES
	laterCSVs = []OdooRecord{}
	fmt.Println("stating Data")
	GCurentFileConsideration := "data"
	for _, fp := range genArgs.FilesData {
		TranslateDataDemoFile(fp)
	}
	fmt.Println("stating remaining data csv translations")
	data := CsvDecodeXml(XmlDataOdoo{Data: []XmlDataData{{Record: ToXmlRecord(laterCSVs)}}})
	CsvWriteCsv(data, GCurentFileConsideration)
	// DEMO FILES
	laterCSVs = []OdooRecord{}
	fmt.Println("stating Demo")
	GCurentFileConsideration = "demo"
	for _, fp := range genArgs.FilesDemo {
		TranslateDataDemoFile(fp)
	}
	fmt.Println("stating remaining demo csv translations")
	data = CsvDecodeXml(XmlDataOdoo{Data: []XmlDataData{{Record: ToXmlRecord(laterCSVs)}}})
	CsvWriteCsv(data, GCurentFileConsideration)
	// PYTHONFILES
	for _, fp := range genArgs.FilesPython {
		TranslatePy2(fp)
	}
	fmt.Println("Done!")
}

func StartTranslator() {
	allfiles := GetAllFiles(genArgs.inputPath)
	for file := allfiles.Front(); file != nil; file = file.Next() {
		fp := file.Value.(string)
		switch {
		case !IsNotInStringArray(strings.Replace(fp, genArgs.inputPath, "", -1), genArgs.pythonFiles):
			genArgs.FilesPython = append(genArgs.FilesPython, fp)
		case !IsNotInStringArray(strings.Replace(fp, genArgs.inputPath+"/", "", -1), genArgs.manifest["data"]):
			genArgs.FilesData = append(genArgs.FilesData, fp)
		case !IsNotInStringArray(strings.Replace(fp, genArgs.inputPath+"/", "", -1), genArgs.manifest["demo"]):
			genArgs.FilesDemo = append(genArgs.FilesDemo, fp)
		case strings.Contains(fp, "/static/"):
			fileutils.Copy(fp, strings.Replace(fp, genArgs.inputPath, genArgs.outputPath, -1))
		default:
			NoTranslation(fp)
		}
	}
	Translator()
}

func GetDefaultManifest() Manifest {
	out := make(Manifest)
	out["name"] = []string{"Module Name"}
	out["version"] = []string{"0.0.0"}
	out["category"] = []string{"none"}
	out["installable"] = []string{"False"}
	out["auto_install"] = []string{"False"}
	return out
}

func ReadManifestEntryList(input string) []string {
	out := []string{}
	for _, str := range strings.Split(input, ",") {
		if len(str) > 0 {
			out = append(out, strings.TrimSpace(str))
		}
	}
	return out
}

func ReadManifestEntry(key, val string) {
	switch key {
	case "depends", "data", "demo":
		genArgs.manifest[key] = ReadManifestEntryList(strings.TrimSpace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(val, "'", "", -1), "\n", "", -1), "[", "", -1), "]", "", -1)))
	case "description":
		genArgs.manifest[key] = []string{strings.TrimSpace(strings.Replace(val, "\"\"\"", "", -1))}
	default:
		genArgs.manifest[key] = []string{strings.TrimSpace(strings.Replace(val, "'", "", -1))}
	}
}
func ReadManifest() {
	filePath := genArgs.inputPath + "/__manifest__.py"
	fileContent, err := ioutil.ReadFile(filePath)
	genArgs.manifest = GetDefaultManifest()
	if err != nil {
		fmt.Println("could not read/find manifest file")
		fmt.Println("No data files will be translated")
		return
	}
	rxTest := regexp.MustCompile("'([\\w]*)'[\\s]*:")
	rxGet := regexp.MustCompile("(.*)('([\\w]*)'[\\s]*:)(.*)")
	var entryLabel string
	var entryValue string
	rawEntries := make(map[string]string)
	for _, csv := range strings.Split(strings.Replace(string(fileContent), "\n", "~", -1), ",") {
		if match := rxTest.MatchString(csv); match {
			matches := rxGet.FindStringSubmatch(csv)
			if len(entryValue) > 0 && len(entryLabel) > 0 {
				rawEntries[entryLabel] = entryValue
			}
			entryLabel = matches[3]
			entryValue = matches[4]
		} else {
			entryValue += "," + csv
		}
	}
	for key, val := range rawEntries {
		ReadManifestEntry(key, strings.Replace(val, "~", "\n", -1))
	}
}

func ReadInitFile(path string) {
	content, err := ioutil.ReadFile(path + "/__init__.py")
	if err != nil {
		fmt.Printf("could not read %s. %v\n", path, err)
		return
	}
	rx := regexp.MustCompile("(?m)^.*import[\\s]+(\\w*)")
	fmt.Println("Reading ", path+"/__init__.py")
	for _, match := range rx.FindAllStringSubmatch(string(content), -1) {
		if _, err := os.Stat(path + "/" + match[1]); err == nil {
			ReadInitFile(path + "/" + match[1])
		} else {
			newPath := strings.Replace(path, genArgs.inputPath, "", -1) + "/" + match[1] + ".py"
			fmt.Println("Found python file", newPath)
			genArgs.pythonFiles = append(genArgs.pythonFiles, newPath)
		}
	}
}

func ReadInit() {
	filePath := genArgs.inputPath + "/__init__.py"
	genArgs.pythonFiles = []string{}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("could not read/find init file")
		fmt.Println("No python files will be translated")
		return
	}
	ReadInitFile(genArgs.inputPath)
}
