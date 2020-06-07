package Hextranslate

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/hexya-addons/web/odooproxy"
)

func FatalError(str string, args ...interface{}) {
	fmt.Println(str)
	for _, a := range args {
		fmt.Println(a)
	}
	FtExit()
}

func FtExit() {
	fmt.Println("Exiting program...")
	os.Exit(0)
}

func PrintVersion() {
	fmt.Println("Starting HexTranslate")
}

func CheckRequestedHelp() {
	if genArgs.help {
		out, err := exec.Command("man", "./hextranslate_manual").Output()
		if err != nil {
			FatalError("could not execute command. Error: %v", err)
		}
		fmt.Println(string(out))
		FtExit()
	}
}

func InitSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\n\n\n\n\n\n\n\nCalled program termination. Program will stop\n")
		FtExit()
	}()
}

func GetRawFileContent(output *list.List, filepath string) {
	output.Init()
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("could not open file %s. content is considered to be empty.\n", filepath)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		output.PushBack(scanner.Text())
	}
}

func CamelCase(words []string) string {
	var out []byte
	for _, word := range words {
		out = append(out, []byte(strings.Title(strings.ToLower(word)))...)
	}
	return string(out)
}

func CreateDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			fmt.Printf("Could not create directory. %v\n", err)
			return
		}
	}
}

func InitSecurityGoFile(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		CreateDir(filepath.Dir(path))
		file, _ := os.Create(path)
		file.WriteString("package " + genArgs.moduleName + "\n\nimport (\n\t\"github.com/hexya-addons/base\"\n\t\"github.com/hexya-erp/hexya/src/models/security\"\n\t\"github.com/hexya-erp/hexya/src/models/security\"\n\t\"github.com/hexya-erp/pool/h\"\n)\n\n//vars\nfunc init() {\n//group_init\n\n//rights\n}\n")
		file.Close()
	}
}

func WriteFile(path string, content *list.List) {
	CreateDir(filepath.Dir(path))
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Could not create file %s. %v\n", path, err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for line := content.Front(); line != nil; line = line.Next() {
		writer.WriteString(line.Value.(string))
		writer.WriteString("\n")
	}
	writer.Flush()
	fmt.Printf("%s written.\n", path)
}

func AddUniqueEntry(array []string, value string) []string {
	for _, val := range array {
		if val == value {
			return array
		}
	}
	array = append(array, value)
	return array
}

func BestString(a, b string) string {
	if len(a) >= len(b) {
		return a
	} else {
		return b
	}
}

func WriteHeadCsv(heads []string) *list.List {
	out := list.New()
	line := heads[0]
	for _, head := range heads[1:] {
		line = line + "," + head
	}
	out.PushBack(line)
	return out
}

func InitCsvDataFile(path string, heads []string) (*list.List, []string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		out := WriteHeadCsv(heads)
		return out, heads
	} else {
		file, _ := os.Open(path)
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		var heads2 []string
		for _, head := range strings.Split(scanner.Text(), ",") {
			heads2 = append(heads2, head)
		}
		for _, head := range heads {
			heads2 = AddUniqueEntry(heads2, head)
		}
		out := WriteHeadCsv(heads2)
		for scanner.Scan() {
			out.PushBack(scanner.Text())
		}
		return out, heads2
	}
}

func NoTranslation(path string) {
	output := list.New()
	fmt.Printf("No translation found for %s. File will be copied without modifications\n", path)
	GetRawFileContent(output, path)
	WriteRawOutputFile(path, output)
}

func IsNotInStringArray(val string, array []string) bool {
	for _, str := range array {
		if str == val {
			return false
		}
	}
	return true
}

func CleanName(str string) string {
	str = strings.Replace(strings.Replace(strings.TrimSpace(str), "\"", "", -1), "'", "", -1)
	str = odooproxy.ConvertModelName(str)
	return str
}

func RemoveExtraneousSpaces(str string) string {
	str = strings.TrimSpace(str)
	oldstr := str
	str = strings.Replace(str, "\t", " ", -1)
	str = strings.Replace(str, "  ", " ", -1)
	for str != oldstr {
		oldstr = str
		str = strings.Replace(str, "\t", " ", -1)
		str = strings.Replace(str, "  ", " ", -1)
	}
	return str
}
