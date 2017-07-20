package Translate

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"
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

//Put the initial code into slice of slice of slice so it's split by classes > lines > words
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

// return a string in go code corresponding to the original python code
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

			if len(rawcode[class][line]) >= 3 && len(rawcode[class][line][2]) > 7 && rawcode[class][line][2][:7] == "fields." {

				cut := strings.Split(rawcode[class][line][2], "(")
				fieldtype := cut[0][7:]
				fieldname := CamelCase(rawcode[class][line][0])
				fieldname = "\"" + fieldname + "\""

				switch fieldtype {

				case "Char":
					result += "pool." + classname + "().AddCharField(" + fieldname + ", models.StringFieldParams{})\n"

				case "Many2one":
					result += "pool." + classname + "().AddMany2OneField(" + fieldname + ",models.ForeignKeyFieldParams{})\n"

				case "One2many":
					result += "pool." + classname + "().AddOne2ManyField(" + fieldname + ", models.ReverseFieldParams{})\n"

				case "Selection":
					result += "pool." + classname + "().AddSelectionField(" + fieldname + ", models.SelectionFieldParams{})\n"

				case "Integer":
					result += "pool." + classname + "().AddIntegerField(" + fieldname + ", models.SimpleFieldParams{})\n"

				case "Datetime":
					result += "pool." + classname + "().AddDateTimeField(" + fieldname + ", models.SimpleFieldParams{})\n"

				case "Float":
					result += "pool." + classname + "().AddFloatField(" + fieldname + ", models.FloatFieldParams{})\n"

				case "Boolean":
					result += "pool." + classname + "().AddBooleanField(" + fieldname + ", models.SimpleFieldParams{})\n"

				case "Many2many":
					result += "pool." + classname + "().AddMany2ManyField(" + fieldname + ", models.Many2ManyFieldParams{})\n"

				case "Binary":
					result += "pool." + classname + "().AddBinaryField(" + fieldname + ", models.SimpleFieldParams{})\n"

				case "Date":
					result += "pool." + classname + "().AddDateField(" + fieldname + ", models.SimpleFieldParams{})\n"

				case "Text":
					result += "pool." + classname + "().AddTextField(" + fieldname + " , models.StringFieldParams{})\n"

				case "Html":
					result += "pool." + classname + "().AddHTMLField(" + fieldname + " , models.StringFieldParams{})\n"

				default:
					println(fieldtype)

				}

			} else if rawcode[class][line][0] == "_sql_constraints" {
				var count int = 1

				for rawcode[class][line+count][0] != "]" {
					var thisline string
					for w := range rawcode[class][line+count] {
						thisline += rawcode[class][line+count][w]
						thisline += " "
					}
					args := strings.Split(thisline, ",")

					name := getArgSqlConstraint(args[0])
					sql := getArgSqlConstraint(args[1])
					errorstring := getArgSqlConstraint(args[2])
					println(args[2])

					result += "pool." + classname + "().AddSQLConstraint(" + name + " , " + sql + " , (" + errorstring + "))\n"

					count += 1
				}

			} else if rawcode[class][line][0] == "def" {

				var body string

				var i int = 1
				for ok := true; ok; ok = rawcode[class][line+i][0] != "def" {

					body += "//"
					for w := range rawcode[class][line+i] {

						body += rawcode[class][line+i][w]
						body += " "
					}
					body += "\n"
					i++

					if line+i == len(rawcode[class]) {
						break
					}
				}

				cut := strings.Split(rawcode[class][line][1], "(")
				name := CamelCase(strings.Trim(cut[0], "_"))

				result += "pool." + classname + "().Method()." + name + "().DeclareMethod(" +
					"\n`" + name + "` ," +
					"\nfunc (){" +
					body+"})\n"
			}

		}

	}
	return result
}

func CamelCase(onestring string) string {

	var result string

	uppercase := true
	for c := range onestring {

		if uppercase == true {

			result += string(unicode.ToUpper(rune(onestring[c])))
			uppercase = false
		} else if string(onestring[c]) == "_" {
			uppercase = true
		} else {
			result += string(onestring[c])
		}
	}

	return result
}

func getArgSqlConstraint(arg string) string {

	var result string
	r := regexp.MustCompile(`'(.+)?'`)

	for i := range arg {
		if string(arg[i]) == "'" {
			r = regexp.MustCompile(`'(.+)?'`)
			break
		}
		if string(arg[i]) == "\"" {
			r = regexp.MustCompile(`"(.+)?"`)
			break
		}
	}

	res := r.FindStringSubmatch(arg)
	result = CamelCase(res[1])
	result = "\"" + result + "\""

	return result

}
