package Translate

import (
	"bytes"
	"regexp"
	"strings"
)

var packagename string
var packagenameset bool = false
var rawcode [][][]string

func TransPyToGo(str string) string {

	var content string

	GenerateSlices(str)
	content = TransRules()

	var result string = "package " + packagename + " \n\n func init() { \n\n " + content + " \n }"

	return result

}

func GenerateSlices(str string) {

	// preparing the document (delete space , etc..)
	regex, err := regexp.Compile("\n\n")
	if err != nil {
		return
	}
	str = regex.ReplaceAllString(str, "\n")
	str = strings.TrimSpace(str)

	classes := strings.Split(str, "class")
	var class = make([][][]string, len(classes))

	for c := range classes {

		classes[c] = strings.TrimSpace(classes[c])
		lines := strings.Split(classes[c], "\n")
		var line = make([][]string, len(lines))

		for l := range lines {

			lines[l] = strings.TrimSpace(lines[l])
			words := strings.Split(lines[l], " ")

			for w := range words {

				line[l] = append(line[l], words[w])

			}

			class[c] = append(class[c], line[l])
		}

		// dont add the empty class
		if c != 0 {
			rawcode = append(rawcode, class[c])
		}

	}

}

func TransRules() string {

	var result string

	for class := range rawcode {

		length := len(rawcode[class][0][0]) - 15
		classname := rawcode[class][0][0][:length]
		result += "\n\npool." + classname + "().DeclareModel()\n"

		for line := range rawcode[class] {

			if packagenameset == false && rawcode[class][line][0] == "_name" {

				var buffer bytes.Buffer

				for c := range rawcode[class][line][2] {

					if string(rawcode[class][line][2][c]) == "." {
						break
					} else {
						//dont write first character that is always " or '
						if c != 0 {
							buffer.WriteByte(rawcode[class][line][2][c])
						}
					}
				}
				packagename = buffer.String()
				packagenameset = true
			}

		}

	}

	//switch word[w] {
	//
	//case "class":
	//line[l+1]

	//

	//
	//case "def":
	//
	//case "fields":
	return result
}
