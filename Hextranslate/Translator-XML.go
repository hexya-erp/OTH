package Hextranslate

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ToXmlField(input []OdooField) []XmlDataField {
	out := []XmlDataField{}
	for _, inputEntry := range input {
		outputEntry := XmlDataField{
			XMLNAME:  inputEntry.XMLNAME,
			InnerXML: inputEntry.InnerXML,
			Name:     inputEntry.Name,
			Ref:      inputEntry.Ref,
		}
		out = append(out, outputEntry)
	}
	return out
}

func ToXmlRecord(input []OdooRecord) []XmlDataRecord {
	out := []XmlDataRecord{}
	for _, inputEntry := range input {
		outputEntry := XmlDataRecord{
			XMLNAME: inputEntry.XMLNAME,
			Field:   ToXmlField(inputEntry.Field),
			Id:      inputEntry.Id,
			Model:   inputEntry.Model,
		}
		out = append(out, outputEntry)
	}
	return out
}

func ReadOdooRecord(rec OdooRecord) GenericRecord {
	genericRecord := GenericRecord{
		ID:   rec.Id,
		Type: rec.Model,
	}
	if rec.Model != "ir.ui.view" && rec.Model != "ir.actions.act_window" && rec.Model != "ir.ui.menu" {
		fmt.Printf("Unknown model %s. record considered as future csv data : %v\n", rec.Model, rec)
		laterCSVs = append(laterCSVs, rec)
	}
	for _, field := range rec.Field {
		switch field.Name {
		case "name":
			genericRecord.Name = field.InnerXML
		case "model", "res_model":
			genericRecord.Model = field.InnerXML
		case "arch":
			genericRecord.Data = field.InnerXML
		case "view_id":
			genericRecord.ViewId = field.Ref
		case "view_mode":
			genericRecord.ViewMode = field.InnerXML
		}
	}
	return genericRecord
}

func WriteHexyaRecord(inrec GenericRecord) (HexyaRecord, int) {
	outrec := HexyaRecord{
		Id:       inrec.ID,
		Type:     inrec.Type,
		Name:     inrec.Name,
		Model:    CamelCase(strings.Split(inrec.Model, ".")),
		ViewId:   inrec.ViewId,
		ViewMode: inrec.ViewMode,
		InnerXML: inrec.Data,
	}
	var outtype int
	switch inrec.Type {
	case "ir.ui.view":
		outtype = RECTYPEVIEW
		outrec.Name = ""
		outrec.Type = ""
	case "ir.actions.act_window":
		outtype = RECTYPEACTION
	case "ir.ui.menu":
		outtype = RECTYPEMENU
	default:
		outtype = 0
	}
	return outrec, outtype
}

func XMLPostTranslate(input []byte, data []int) []byte {
	//WARNING - POORLY OPTIMIZED
	//replace record placeholders
	output := bytes.Replace(input, []byte("HexyaHexyaXml>"), []byte("hexya>"), -1)
	var rectype string
	for _, val := range data {
		switch val {
		case RECTYPEVIEW:
			rectype = "view"
		case RECTYPEACTION:
			rectype = "action"
		case RECTYPEMENU:
			rectype = "menu"
		default:
			rectype = "Unknown"
		}
		output = bytes.Replace(output, []byte("HEXYA_HEXYA_TMP_RECORD"), []byte(rectype), 2)
	}
	//remove empty attrs
	lines := strings.Split(string(output), "\n")
	output = []byte{}
	output = append(output, "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n"...)
	for _, line := range lines {
		if strings.Contains(line, "=\"\"") {
			words := strings.Split(line, " ")
			line = string([]byte{})
			for _, word := range words {
				if strings.Contains(word, "=\"\"") {
					if strings.HasSuffix(word, ">") {
						word = ">"
					} else if strings.HasSuffix(word, "></action>") {
						word = "></action>"
					} else {
						word = ""
					}
				}
				line = string(append(append([]byte(line), []byte(word)...), ' '))
			}
			line = string(append([]byte("\t"), []byte(strings.Join(strings.Fields(line), " "))...))
		}
		output = append(append(output, []byte(line)...), '\n')
	}
	return output
}

func WriteHexyaData(genRecords []GenericRecord) []byte {
	hexya := HexyaHexyaXml{}
	var postTranslateData []int
	for _, rec := range genRecords {
		hexrec, ptd := WriteHexyaRecord(rec)
		hexya.Data.Records = append(hexya.Data.Records, hexrec)
		postTranslateData = append(postTranslateData, ptd)
	}
	out, err := xml.MarshalIndent(&hexya, "", "    ")
	if err != nil {
		fmt.Printf("Could not Marshal. %v\n", err)
	}
	out = XMLPostTranslate(out, postTranslateData)
	return out
}

func WriteHexyaXML(path string, genRecords []GenericRecord) {
	path = string(append(append([]byte(genArgs.outputPath), []byte("/resources")...), []byte(strings.TrimPrefix(path, genArgs.inputPath))...))
	CreateDir(filepath.Dir(path))
	err := ioutil.WriteFile(path, WriteHexyaData(genRecords), 0777)
	if err != nil {
		fmt.Printf("Could not write file. %v\n", err)
		return
	}
	fmt.Printf("XML file %s is written\n", path)
}

func StartXmltoXml(file io.Reader, filePath string) {
	byteValue, _ := ioutil.ReadAll(file)
	var odoo OdooOdoo
	var genericRecordArray []GenericRecord
	xml.Unmarshal(byteValue, &odoo)
	for _, rec := range odoo.Record {
		genericRecordArray = append(genericRecordArray, ReadOdooRecord(rec))
	}
	for _, data := range odoo.Data {
		for _, rec := range data.Record {
			genericRecordArray = append(genericRecordArray, ReadOdooRecord(rec))
		}
	}
	ReadTemplates(odoo.Template)
	for _, data := range odoo.Data {
		ReadTemplates(data.Template)
	}
	WriteHexyaXML(filePath, genericRecordArray)
}

func WriteSecurityData(filePath string, odoo SecuOdoo) {
	file, _ := os.Open(filePath)
	defer file.Close()
	output := list.New()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		output.PushBack(scanner.Text())
		if scanner.Text() == "//vars" {
			output.PushBack("\nvar (")
			for _, data := range odoo.Data {
				for _, rec := range data.Record {
					if rec.Model == "res.groups" {
						for _, field := range rec.Field {
							switch field.Name {
							case "name":
								rec.Name = field.InnerXML
							}
						}
						groupName := CamelCase(strings.Split(rec.Id, "_"))
						output.PushBack(fmt.Sprintf("\t// %s %s\n\t%s *security.Group", groupName, rec.Name, groupName))
					}
				}
			}
			output.PushBack(")\n")
		}
		if strings.HasSuffix(scanner.Text(), "//group_init") {
			for _, data := range odoo.Data {
				for _, rec := range data.Record {
					if rec.Model == "res.groups" {
						for _, field := range rec.Field {
							switch field.Name {
							case "name":
								rec.Name = field.InnerXML
							}
						}
						groupName := CamelCase(strings.Split(rec.Id, "_"))
						output.PushBack(fmt.Sprintf("\t%s = security.Registry.NewGroup(\"%s_%s\", \"%s\")",
							groupName, genArgs.moduleName, rec.Id, rec.Name))
					}
				}
			}
		}
	}
	WriteFile(filePath, output)
}

func StartXmltoSecurity(file io.Reader, filePath string) {
	filePath = strings.Replace(filePath, genArgs.inputPath, genArgs.outputPath, -1)
	InitSecurityGoFile(filePath)
	byteValue, _ := ioutil.ReadAll(file)
	var odoo SecuOdoo
	err := xml.Unmarshal(byteValue, &odoo)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	WriteSecurityData(filePath, odoo)
}

func CsvDecodeXml(odoo XmlDataOdoo) CsvData {
	out := CsvData{
		ModelEntryCount: make(map[string]int),
		ModelHeads:      make(map[string][]string),
		ModelContent:    make(map[string][]CsvValue)}
	for _, data := range odoo.Data {
		for _, rec := range data.Record {
			out.Models = AddUniqueEntry(out.Models, rec.Model)
			out.ModelHeads[rec.Model] = AddUniqueEntry(out.ModelHeads[rec.Model], "ID")
			out.ModelContent[rec.Model] = append(out.ModelContent[rec.Model], CsvValue{})
			out.ModelEntryCount[rec.Model] += 1
			out.ModelContent[rec.Model][out.ModelEntryCount[rec.Model]-1].Value = make(map[string]string)
			out.ModelContent[rec.Model][out.ModelEntryCount[rec.Model]-1].Value["ID"] = rec.Id
			for _, field := range rec.Field {
				out.ModelHeads[rec.Model] = AddUniqueEntry(out.ModelHeads[rec.Model], field.Name)
				out.ModelContent[rec.Model][out.ModelEntryCount[rec.Model]-1].Value[field.Name] = BestString(field.Ref, field.InnerXML)
			}
		}
	}
	return out
}

func CsvReadXml(file io.Reader) CsvData {
	odoo := XmlDataOdoo{}
	byteValue, _ := ioutil.ReadAll(file)
	xml.Unmarshal(byteValue, &odoo)
	return CsvDecodeXml(odoo)
}

func CsvWriteCsv(data CsvData, fileConsideration string) {
	for _, model := range data.Models {
		outPath := genArgs.outputPath + "/" + fileConsideration + "/" + CamelCase(strings.Split(model, ".")) + ".csv"
		var output *list.List
		output, data.ModelHeads[model] = InitCsvDataFile(outPath, data.ModelHeads[model])
		for _, entry := range data.ModelContent[model] {
			outline := entry.Value[data.ModelHeads[model][0]]
			for _, fieldName := range data.ModelHeads[model][1:] {
				outline = outline + "," + entry.Value[fieldName]
			}
			output.PushBack(outline)
		}
		WriteFile(outPath, output)
	}
}

func StartXmltoCsv(file io.Reader, filePath string) {
	data := CsvReadXml(file)
	fileConsideration := "demo"
	if IsNotInStringArray(strings.Replace(filePath, genArgs.inputPath+"/", "", -1), genArgs.manifest["demo"]) {
		fileConsideration = "data"
	}
	CsvWriteCsv(data, fileConsideration)
}

func TranslateXML(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Could not open file. %v\n", err)
		return
	}
	defer file.Close()
	fmt.Printf("Translating XML file: %s\n", filePath)
	folders := strings.TrimPrefix(filepath.Dir(filePath), genArgs.inputPath)
	switch {
	case strings.Contains(folders, "views"):
		StartXmltoXml(file, strings.Replace(filePath, "/views", "", -1))
	case strings.Contains(folders, "wizard"):
		changedName := append([]byte("_wizard_"), []byte(filepath.Base(filePath))...)
		filePath = string(append([]byte(strings.TrimSuffix(filePath, filepath.Base(filePath))), changedName...))
		StartXmltoXml(file, strings.Replace(filePath, "/wizard", "", -1))
	case strings.Contains(folders, "data"):
		StartXmltoCsv(file, filePath)
	case strings.Contains(folders, "security"):
		changedName := "security.go"
		filePath = string(append([]byte(strings.TrimSuffix(filePath, filepath.Base(filePath))), changedName...))
		StartXmltoSecurity(file, strings.Replace(filePath, "/security/", "/", -1))
	default:
		NoTranslation(filePath)
	}
}
