package Hextranslate

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func WriteGoSecurityLine(path string, output *list.List) {
	infile, _ := os.Open(path)
	defer infile.Close()
	scanner := bufio.NewScanner(infile)
	scanner.Scan()
	for scanner.Scan() {
		entry := SecurityEntry{}
		words := strings.Split(scanner.Text(), ",")
		entry.ModelId = CamelCase(strings.Split(strings.TrimPrefix(words[2], "model_"), "_"))
		groupIdSplit := strings.Split(words[3], ".")
		entry.GroupId = CamelCase(strings.Split(groupIdSplit[len(groupIdSplit)-1], "_"))
		if words[3] == "" {
			entry.GroupId = "security.GroupEveryone"
		}
		if len(words[3]) >= 5 && words[3][:5] == "base." {
			entry.GroupId = "base." + entry.GroupId
		}
		if words[4] == "1" && words[5] == "1" && words[6] == "1" && words[7] == "1" {
			output.PushBack(fmt.Sprintf("\th.%s().Methods().AllowAllToGroup(%s)", entry.ModelId, entry.GroupId))
		} else {
			if words[4] == "1" {
				output.PushBack(fmt.Sprintf("\th.%s().Methods().%s().AllowGroup(%s)", entry.ModelId, "Load", entry.GroupId))
			}
			if words[5] == "1" {
				output.PushBack(fmt.Sprintf("\th.%s().Methods().%s().AllowGroup(%s)", entry.ModelId, "Write", entry.GroupId))
			}
			if words[6] == "1" {
				output.PushBack(fmt.Sprintf("\th.%s().Methods().%s().AllowGroup(%s)", entry.ModelId, "Create", entry.GroupId))
			}
			if words[7] == "1" {
				output.PushBack(fmt.Sprintf("\th.%s().Methods().%s().AllowGroup(%s)", entry.ModelId, "Unlink", entry.GroupId))
			}
		}
	}
}

func TranslateCSVSecurity(path string) {
	outpath := fmt.Sprintf("%s/security.go", genArgs.outputPath)
	InitSecurityGoFile(outpath)
	file, _ := os.Open(outpath)
	defer file.Close()
	output := list.New()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		output.PushBack(scanner.Text())
		if scanner.Text() == "//rights" {
			WriteGoSecurityLine(path, output)
		}
	}
	WriteFile(outpath, output)
	out, err := exec.Command("goimports", "-w", fmt.Sprintf("%s", outpath)).CombinedOutput()
	if err != nil {
		fmt.Printf("could not execute command. Error: %v\n%s\n", err, out)
	}
}

func TranslateCSV(file string) {
	switch filepath.Base(file) {
	case "ir.model.access.csv":
		TranslateCSVSecurity(file)
	default:
		NoTranslation(file)
	}
}
