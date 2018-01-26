package translate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var packagename string
var rawcode [][][]string
var imports string
var defs []string

//call the translates rules to build a string corresponding to the generated go code
func TransPyToGo(str string, pack string) string {

	var content string

	packagename = pack
	rawcode = nil
	imports = ""
	defs = nil

	GenerateSlices(str)
	GenerateDefs(str)
	content = TransRules()

	var result string = "package " + packagename + " \n\n  import (\n" + imports + "\n) \n\n func init() { \n\n " + content + " \n }"

	return result

}

//Put the initial code into slice of slice of slice  so it split  the code like this : words < lines < class
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

// build a slice of string containing the python functions
func GenerateDefs(str string) {

	cut := strings.Split(str, "@api.")

	for c := range cut {
		cut[c] = "@api." + cut[c]
		def := strings.Split(cut[c], "def ")

		for d := range def {
			defcut := strings.Split(def[d], "class")
			for f := range defcut {
				finalcut := strings.Split(defcut[f], "fields.")
				defs = append(defs, finalcut[0])
			}

		}
	}
}

// return a string in go code corresponding to the original python code
func TransRules() string {

	var result string
	var selectionimportset bool = false
	var classname string

	for class := range rawcode {

		var inherit = false

		for line := range rawcode[class] {

			if rawcode[class][line][0] == "_inherit" {
				classname = CheckBuitInNames(CamelCase(TrimString(string(rawcode[class][line][2]))))
				inherit = true
				break
			} else if rawcode[class][line][0] == "_name" {
				classname = CheckBuitInNames(CamelCase(TrimString(string(rawcode[class][line][2]))))
				break
			}
		}

		var fieldDeclarationStarted bool

		for line := range rawcode[class] {

			typemodel := strings.Split(rawcode[class][line][0], ".")

			if len(typemodel) > 1 && typemodel[1] == "TransientModel):" && !inherit {
				result += "\n\npool." + classname + "().DeclareTransientModel()\n"

			} else if len(typemodel) > 1 && typemodel[1] == "Model):" && !inherit {
				result += "\n\npool." + classname + "().DeclareModel()\n"
			}

			if len(rawcode[class][line]) >= 3 && len(rawcode[class][line][2]) > 7 && rawcode[class][line][2][:7] == "fields." {

				if !fieldDeclarationStarted {
					result += fmt.Sprintf("pool.%s().AddFields(map[string]models.FieldDefinition{\n", classname)
					fieldDeclarationStarted = true
				}

				cut := strings.Split(rawcode[class][line][2], "(")
				fieldtype := cut[0][7:]
				odoofieldname := rawcode[class][line][0]
				fieldname := CamelCase(odoofieldname)
				fieldname = "\"" + fieldname + "\""

				switch fieldtype {

				case "Char":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])
					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					} else {
						name = fieldname
					}

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "string":
							name = "\"" + TrimString(strings.TrimSpace(value[1])) + "\""
						case "required":
							if len(value[1]) == 1 {
								if value[1] == "0" {
									body += ", Required: false"
								} else {
									body += ", Required: true"
								}
							} else {
								body += ", Required : " + strings.ToLower(value[1])
							}
						case "translate":
							body += ", Translate: " + strings.ToLower(value[1])
						case "compute":
							body += ", Compute: pool." + classname + "().Methods()." + CamelCase(strings.Trim(strings.Trim(value[1], "'"), "_")) + "()"
						case "help":
							help := GetHelpText(class, line)
							body += " ,Help :\"" + help + "\""
						case "index":
							if len(value[1]) == 1 {
								if value[1] == "0" {
									body += ", Index: false"
								} else {
									body += ", Index: true"
								}
							} else {
								body += ", Index: " + strings.ToLower(value[1])
							}
						case "copy":
							if value[1] == "True" {
								body += ", NoCopy: false"
							} else {
								body += ", NoCopy: true"
							}
						case "related":
							var result string
							var i int = 0
							related := CamelCaseForRelated(strings.Trim(value[1], "'"))
							split := strings.Split(related, ".")

							for s := range split {

								if (split[s][len(split[s])-2:]) == "Id" {
									if i > 0 {
										result += "."
									}
									result += split[s][:len(split[s])-2]

								} else if (split[s][len(split[s])-3:]) == "Ids" {
									if i > 0 {
										result += "."
									}
									result += split[s][:len(split[s])-3] + "s"

								}

								i++
							}
						case "inverse":
							body += ", Inverse: pool." + classname + "().Methods().Inverse" + CamelCase(strings.Trim(fieldname, "'\"")) + "()"
						case "store":
							body += ", Stored: " + strings.ToLower(TrimString(strings.TrimSpace(value[1])))
						default:
							body += fmt.Sprintf("/*%s*/", value)
						}
					}
					if name != "" {
						body = "String :" + name + ", " + body
					}
					result += fmt.Sprintf("%s: models.CharField{%s},\n", fieldname, body)

				case "Many2one":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					if len(args) < 2 {
						body += fmt.Sprintf("/* %s */", args)
						result += fmt.Sprintf("%s: models.Many2OneField{%s},\n", fieldname, body)
						break
					}
					args[1] = strings.TrimSpace(args[1])
					if string(args[1][0]) == "'" || string(args[1][0]) == "\"" {
						name = "\"" + TrimString(strings.Trim(args[1], ")")) + "\""
					}
					foreignkey := CheckBuitInNames(CamelCase(strings.Trim(strings.TrimSpace(args[0]), "'")))

					if (foreignkey[len(foreignkey)-2:]) == "Id" {
						foreignkey = foreignkey[:len(foreignkey)-2]

					} else if (foreignkey[len(foreignkey)-3:]) == "Ids" {
						foreignkey = foreignkey[:len(foreignkey)-3] + "s\""
					}

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "required":
							if len(value[1]) == 1 {
								if value[1] == "0" {
									body += ", Required: false"
								} else {
									body += ", Required: true"
								}
							} else {
								body += ", Required : " + strings.ToLower(value[1])
							}
						case "default":

							def := value[1]

							for d := range defs {
								if len(defs[d]) > len(value[1]) && defs[d][:len(value[1])] == value[1] {

									def = defs[d]
								}
							}

							body += ", Default: func(env models.Environment) interface{}{\n" +
								"/*" + def + "*/\n" +
								"return 0}"

						case "ondelete":
							body += ", OnDelete : models." + CamelCase(TrimString(value[1]))
						case "help":
							help := GetHelpText(class, line)
							body += " ,Help :\"" + help + "\""
						case "index":
							if len(value[1]) == 1 {
								if value[1] == "0" {
									body += ", Index: false"
								} else {
									body += ", Index: true"
								}
							} else {
								body += ", Index: " + strings.ToLower(value[1])
							}
						case "readonly":
							body += "/* readonly=true */ "
						case "related":
							var result string
							var i int = 0
							related := CamelCaseForRelated(strings.Trim(value[1], "'"))
							split := strings.Split(related, ".")

							for s := range split {

								if (split[s][len(split[s])-2:]) == "Id" {
									if i > 0 {
										result += "."
									}
									result += split[s][:len(split[s])-2]

								} else if (split[s][len(split[s])-3:]) == "Ids" {
									if i > 0 {
										result += "."
									}
									result += split[s][:len(split[s])-3] + "s"

								}

								i++
							}
							body += ", Related: \"" + result + "\""
						case "store":
							body += ", Stored: " + strings.ToLower(TrimString(strings.TrimSpace(value[1])))
						case "string":
							name = "\"" + TrimString(strings.TrimSpace(value[1])) + "\""
						case "domain":
							body += "/*, Filter: " + value[1] + "*/"
						case "compute":
							body += ", Compute : pool." + classname + "().Methods()." + CamelCase(strings.Trim(TrimString(value[1]), "_")) + "()"
						case "inverse":
							body += ", Inverse: pool." + classname + "().Methods().Inverse" + CamelCase(strings.Trim(fieldname, "'\"")) + "()"
						case "company_dependent":
							body += "/*, CompanyDependent : " + strings.ToLower(value[1]) + "*/"
						default:
							body += fmt.Sprintf("/*%s*/", value)
						}
					}
					body = "String :" + name + " , RelationModel: pool." + foreignkey + "()" + body
					result += fmt.Sprintf("%s: models.Many2OneField{%s},\n", fieldname, body)

				case "One2many":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					foreignkey := CamelCase(TrimString(strings.TrimSpace(args[1])))
					args[2] = strings.TrimSpace(args[2])
					if string(args[2][0]) == "'" || string(args[2][0]) == "\"" {
						name = strings.Trim(strings.Trim(args[2], ")"), "(")
						name = "\"" + TrimString(name) + "\""
					} else {
						name = fieldname
					}
					body += "String :" + name
					body += " ,RelationModel : pool." + CheckBuitInNames(CamelCase(TrimString(strings.TrimSpace(args[0])))) + "()"

					if (foreignkey[len(foreignkey)-2:]) == "Id" {
						foreignkey = foreignkey[:len(foreignkey)-2]

					} else if (foreignkey[len(foreignkey)-3:]) == "Ids" {
						foreignkey = foreignkey[:len(foreignkey)-3] + "s\""
					}

					body += " ,ReverseFK : \"" + foreignkey + "\""

					if TrimString(odoofieldname[len(odoofieldname)-3:]) == "_id" {
						body += " , JSON : \"" + odoofieldname + "\""
						fieldname = fieldname[:len(fieldname)-3] + "\""

					} else if TrimString(odoofieldname[len(odoofieldname)-4:]) == "_ids" {
						body += " , JSON : \"" + odoofieldname + "\""
						fieldname = fieldname[:len(fieldname)-4] + "s\""
					}

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "copy":
							if value[1] == "True" {
								body += ", NoCopy: false"
							} else {
								body += ", NoCopy: true"
							}
						case "default":
							if string(value[1][0]) == "_" {
								var def string

								for d := range defs {
									if len(defs[d]) > len(value[1]) && defs[d][:len(value[1])] == value[1] {

										def = defs[d]
									}
								}

								body += ", Default: func(env models.Environment) interface{}{\n" +
									"/*" + def + "*/return 0}"
							} else {
								body += ", Default: models.DefaultValue(" + value[1] + ")"
							}
						case "help":
							help := GetHelpText(class, line)
							body += " ,Help :\"" + help + "\""
						case "required":
							if len(value[1]) == 1 {
								if value[1] == "0" {
									body += ", Required: false"
								} else {
									body += ", Required: true"
								}
							} else {
								body += ", Required : " + strings.ToLower(value[1])
							}
						case "readonly":
							body += "/* readonly */"
						case "string":
							name = "\"" + TrimString(strings.TrimSpace(value[1])) + "\""
						default:
							body += fmt.Sprintf("/*%s*/", value)
						}
					}
					if name != "" {
						body = "String :" + name + ", " + body
					}
					result += fmt.Sprintf("%s: models.One2ManyField{%s},\n", fieldname, body)

				case "Selection":
					if selectionimportset == false {
						imports += "\"github.com/hexya-erp/hexya/hexya/models/types\"\n"
						selectionimportset = true
					}

					var body string
					var s string
					var count = 0
					ok := false

					for ok == false {
						for w := range rawcode[class][line+count] {
							s += " " + rawcode[class][line+count][w]

						}
						if string(s[len(s)-1]) != ")" {
							count += 1
						} else {
							ok = true
						}
					}

					s = strings.TrimSpace(s)
					cut := strings.Split(s, "[")
					if len(cut) < 2 {
						body += fmt.Sprintf("/*%s*/", s)
						result += fmt.Sprintf("%s: models.SelectionField{%s},\n", fieldname, body)
						break
					}
					cut2 := strings.Split(cut[1], "]")
					if len(cut2) < 2 {
						body += fmt.Sprintf("/*%s*/", s)
						result += fmt.Sprintf("%s: models.SelectionField{%s},\n", fieldname, body)
						break
					}
					args := strings.Split(cut2[1], ",")
					if len(args) < 2 {
						body += fmt.Sprintf("/*%s*/", s)
						result += fmt.Sprintf("%s: models.SelectionField{%s},\n", fieldname, body)
						break
					}
					selectable := strings.Split(cut2[0], "),")

					name := ""
					args[1] = strings.TrimSpace(args[1])
					if string(args[1][0]) == "'" || string(args[1][0]) == "\"" {
						name = "\"" + TrimString(args[1]) + "\""
					}
					body += ", Selection : types.Selection{\n"

					for sub := range selectable {
						sec := strings.Split(selectable[sub], ",")
						if len(sec) > 2 {
							i := 1
							for len(sec)-2 >= i {
								sec[1] += sec[1+i]

								i++
							}
							sec = sec[:len(sec)-(len(sec)-2)]
						}
						if len(sec) != 2 {
							body += fmt.Sprintf("/*%s*/", selectable)
							break
						}
						sec0 := TrimString(strings.Trim(strings.TrimSpace(sec[0]), "("))
						sec1 := TrimString(strings.Trim(strings.TrimSpace(sec[1]), ")"))
						body += "\"" + sec0 + "\" : \"" + sec1 + "\",\n"

					}

					body += "}"

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")
						switch strings.TrimSpace(value[0]) {
						case "help":
							body += ", Help : \"" + TrimString(value[1]) + "\""
						case "default":
							body += ", Default: models.DefaultValue(\"" + TrimString(value[1]) + "\")"
						case "required":
							if len(value[1]) == 1 {
								if value[1] == "0" {
									body += ", Required: false"
								} else {
									body += ", Required: true"
								}
							} else {
								body += ", Required : " + strings.ToLower(value[1])
							}
						case "index":
							if len(value[1]) == 1 {
								if value[1] == "0" {
									body += ", Index: false"
								} else {
									body += ", Index: true"
								}
							} else {
								body += ", Index: " + strings.ToLower(value[1])
							}

						case "implied_group":
							body += "/*, ImpliedGroup :\"" + TrimString(value[1]) + "\"*/"
						case "string":
							name = "\"" + TrimString(strings.TrimSpace(value[1])) + "\""

						default:
							body += fmt.Sprintf("/*%s*/", value)
						}
					}

					if name != "" {
						body = "String :" + name + ", " + body
					}
					result += fmt.Sprintf("%s: models.SelectionField{%s},\n", fieldname, body)

				case "Integer":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])
					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					}
					body += "String :" + name

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "help":
							help := GetHelpText(class, line)
							body += " ,Help :\"" + help + "\""
						case "default":
							if string(value[1][0]) == "_" {
								var def string

								for d := range defs {
									if len(defs[d]) > len(value[1]) && defs[d][:len(value[1])] == value[1] {

										def = defs[d]
									}
								}

								body += ", Default: func(env models.Environment) interface{}{\n" +
									"/*" + def + "*/return 0}"
							} else {
								body += ", Default: models.DefaultValue(" + value[1] + ")"
							}
						case "required":
							if len(value[1]) == 1 {
								if value[1] == "0" {
									body += ", Required: false"
								} else {
									body += ", Required: true"
								}
							} else {
								body += ", Required : " + strings.ToLower(value[1])
							}
						case "index":
							if len(value[1]) == 1 {
								if value[1] == "0" {
									body += ", Index: false"
								} else {
									body += ", Index: true"
								}
							} else {
								body += ", Index: " + strings.ToLower(value[1])
							}
						case "compute":
							body += ", Compute : pool." + classname + "().Methods()." + CamelCase(strings.Trim(TrimString(value[1]), "_")) + "()"

						case "string":
							name = "\"" + TrimString(strings.TrimSpace(value[1])) + "\""
						default:
							body += fmt.Sprintf("/*%s*/", value)
						}
					}
					if name != "" {
						body = "String :" + name + ", " + body
					}
					result += fmt.Sprintf("%s: models.IntegerField{%s},\n", fieldname, body)

				case "Float", "Monetary":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])

					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					}
					body += "String :" + name

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "default":
							if string(value[1][0]) == "_" {
								var def string

								for d := range defs {
									if len(defs[d]) > len(value[1]) && defs[d][:len(value[1])] == value[1] {

										def = defs[d]
									}
								}

								body += ", Default: func(env models.Environment) interface{}{\n" +
									"/*" + def + "*/return 0}"
							} else {
								body += ", Default: models.DefaultValue(" + value[1] + ")"
							}
						case "digits":
							if _, err := strconv.Atoi(value[1]); err == nil {
								body += ",  Digits: nbutils.Digits{" + value[1] + "," + value[1] + "}"
							} else {
								body += "/*,  Digits:" + value[1] + "*/"
							}
						case "compute":
							body += ", Compute: pool." + classname + "().Methods()." + CamelCase(strings.Trim(strings.Trim(value[1], "'"), "_")) + "()"
						case "help":
							help := GetHelpText(class, line)
							body += " ,Help :\"" + help + "\""
						case "required":
							if len(value[1]) == 1 {
								if value[1] == "0" {
									body += ", Required: false"
								} else {
									body += ", Required: true"
								}
							} else {
								body += ", Required : " + strings.ToLower(value[1])
							}
						case "company_dependent":
							body += "/*, CompanyDependent : " + strings.ToLower(value[1]) + "*/"
						case "inverse":
							body += ", Inverse: pool." + classname + "().Methods().Inverse" + CamelCase(strings.Trim(fieldname, "'\"")) + "()"
						case "store":
							body += ", Stored: " + strings.ToLower(TrimString(strings.TrimSpace(value[1])))
						case "related":
							var result string
							var i int = 0
							related := CamelCaseForRelated(strings.Trim(value[1], "'"))
							split := strings.Split(related, ".")

							for s := range split {

								if (split[s][len(split[s])-2:]) == "Id" {
									if i > 0 {
										result += "."
									}
									result += split[s][:len(split[s])-2]

								} else if (split[s][len(split[s])-3:]) == "Ids" {
									if i > 0 {
										result += "."
									}
									result += split[s][:len(split[s])-3] + "s"

								}

								i++
							}
							body += ", Related: \"" + result + "\""
						case "search":
							body += "/*, Search: \"" + TrimString(value[1]) + "\"*/"

						case "string":
							name = "\"" + TrimString(strings.TrimSpace(value[1])) + "\""
						default:
							body += fmt.Sprintf("/*%s*/", value)
						}
					}
					if name != "" {
						body = "String :" + name + ", " + body
					}
					result += fmt.Sprintf("%s: models.FloatField{%s},\n", fieldname, body)

				case "Boolean":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])
					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					}
					body += "String :" + name

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "default":
							if len(value[1]) == 1 {
								if value[1] == "0" {
									body += ", Default: models.DefaultValue(false)"
								} else {
									body += ", Default: models.DefaultValue(true)"
								}
							} else {
								body += ", Default: models.DefaultValue(" + strings.ToLower(value[1]) + ")"
							}

						case "string":
							name = "\"" + TrimString(strings.TrimSpace(value[1])) + "\""
						case "help":
							help := GetHelpText(class, line)
							body += " ,Help :\"" + help + "\""
						default:
							body += fmt.Sprintf("/*%s*/", value)
						}
					}
					if name != "" {
						body = "String :" + name + ", " + body
					}
					result += fmt.Sprintf("%s: models.BooleanField{%s},\n", fieldname, body)

				case "Many2many":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[1] = strings.TrimSpace(args[1])
					if string(args[1][0]) == "'" || string(args[1][0]) == "\"" {
						name = "\"" + TrimString(args[1]) + "\""
					}

					foreignkey := CheckBuitInNames(CamelCase(strings.Trim(strings.TrimSpace(args[0]), "'")))

					if (foreignkey[len(foreignkey)-2:]) == "Id" {
						foreignkey = foreignkey[:len(foreignkey)-2]

					} else if (foreignkey[len(foreignkey)-3:]) == "Ids" {
						foreignkey = foreignkey[:len(foreignkey)-3] + "s\""
					}

					if TrimString(odoofieldname[len(odoofieldname)-3:]) == "_id" {
						body += " , JSON : \"" + odoofieldname + "\""
						fieldname = fieldname[:len(fieldname)-3] + "\""

					} else if TrimString(odoofieldname[len(odoofieldname)-4:]) == "_ids" {
						body += " , JSON : \"" + odoofieldname + "\""
						fieldname = fieldname[:len(fieldname)-4] + "s\""
					}

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "string":
							name = "\"" + TrimString(strings.TrimSpace(value[1])) + "\""
						case "compute":
							body += ", Compute: pool." + classname + "().Methods()." + CamelCase(strings.Trim(strings.Trim(value[1], "'"), "_")) + "()"
						case "ondelete":
							body += "/*, OnDelete : models." + CamelCase(TrimString(value[1])) + "*/"
						default:
							body += fmt.Sprintf("/*%s*/", value)
						}
					}
					body = "String :" + name + " , RelationModel: pool." + foreignkey + "()" + body
					result += fmt.Sprintf("%s: models.Many2ManyField{%s},\n", fieldname, body)

				case "Binary":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])
					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					}
					body += "String :" + name

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {

						case "string":
							name = "\"" + TrimString(strings.TrimSpace(value[1])) + "\""
						case "help":
							help := GetHelpText(class, line)
							body += " ,Help :\"" + help + "\""
						case "compute":
							body += ", Compute: pool." + classname + "().Methods()." + CamelCase(strings.Trim(strings.Trim(value[1], "'"), "_")) + "()"
						case "inverse":
							body += ", Inverse: pool." + classname + "().Methods().Inverse" + CamelCase(strings.Trim(fieldname, "'\"")) + "()"
						case "attachment":
							body += "/*, Attachment: " + strings.ToLower(TrimString(strings.TrimSpace(value[1]))) + "*/"
						default:
							body += fmt.Sprintf("/*%s*/", value)
						}
					}
					if name != "" {
						body = "String :" + name + ", " + body
					}
					result += fmt.Sprintf("%s: models.BinaryField{%s},\n", fieldname, body)

				case "Date":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])
					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					}
					body += "String :" + name

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {

						case "string":
							name = "\"" + TrimString(strings.TrimSpace(value[1])) + "\""
						case "help":
							help := GetHelpText(class, line)
							body += " ,Help :\"" + help + "\""
						default:
							body += fmt.Sprintf("/*%s*/", value)
						}
					}
					if name != "" {
						body = "String :" + name + ", " + body
					}
					result += fmt.Sprintf("%s: models.DateField{%s},\n", fieldname, body)

				case "Datetime":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])
					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					}
					body += "String :" + name

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "string":
							name = "\"" + TrimString(strings.TrimSpace(value[1])) + "\""
						case "help":
							help := GetHelpText(class, line)
							body += " ,Help :\"" + help + "\""
						default:
							body += fmt.Sprintf("/*%s*/", value)
						}
					}
					if name != "" {
						body = "String :" + name + ", " + body
					}
					result += fmt.Sprintf("%s: models.DateTimeField{%s},\n", fieldname, body)

				case "Text":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])
					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					}
					body += "String :" + name

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "translate":
							body += ", Translate: " + strings.ToLower(value[1])
						case "string":
							name = "\"" + TrimString(strings.TrimSpace(value[1])) + "\""
						case "help":
							help := GetHelpText(class, line)
							body += " ,Help :\"" + help + "\""
						default:
							body += fmt.Sprintf("/*%s*/", value)
						}
					}
					if name != "" {
						body = "String :" + name + ", " + body
					}
					result += fmt.Sprintf("%s: models.TextField{%s},\n", fieldname, body)

				case "Html":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])
					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					}
					body += "String :" + name

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {

						case "string":
							name = "\"" + TrimString(strings.TrimSpace(value[1])) + "\""
						case "help":
							help := GetHelpText(class, line)
							body += " ,Help :\"" + help + "\""

						default:
							body += fmt.Sprintf("/*%s*/", value)
						}
					}
					if name != "" {
						body = "String :" + name + ", " + body
					}
					result += fmt.Sprintf("%s: models.HTMLField{%s},\n", fieldname, body)

				default:
					fmt.Println("Unknown fieldType: ", fieldtype)

				}

			} else if rawcode[class][line][0] == "_sql_constraints" {
				if fieldDeclarationStarted {
					result += "\n})\n"
					fieldDeclarationStarted = false
				}

				var count int = 1

				for rawcode[class][line+count][0] != "]" {
					var thisline string
					for w := range rawcode[class][line+count] {
						thisline += rawcode[class][line+count][w]
						thisline += " "
					}
					args := strings.Split(thisline, ",")

					if len(args) > 3 && args[3] != " " {
						i := 2
						var sqlstring = (args[1])
						for string(sqlstring[len(sqlstring)-2]) != ")" {
							sqlstring += "," + (args[i])
							i++
						}

						args[1] = sqlstring
						args[2] = strings.TrimRight(strings.TrimLeft(args[len(args)-1], "("), ")")
					}

					if len(args) == 3 {
						name := CamelCase(GetArgsSqlConstraint(args[0]))
						sql := GetArgsSqlConstraint(args[1])
						errorstring := strings.Trim(GetArgsSqlConstraint(args[2]), "\"")
						result += fmt.Sprintf("pool.%s().AddSQLConstraint(\"%s\", \"%s\", \"%s\")\n", classname, name, sql, errorstring)
					} else {
						result += fmt.Sprintf("pool.%s().AddSQLConstraint(/* %v */)\n", classname, args)
					}

					count += 1
				}

			} else if rawcode[class][line][0] == "def" {
				if fieldDeclarationStarted {
					result += "\n})\n"
					fieldDeclarationStarted = false
				}

				var body string
				var args string
				var getargs []string
				var def string

				for d := range defs {
					if len(defs[d]) >= len(rawcode[class][line][1]) && rawcode[class][line][1] == defs[d][:len(rawcode[class][line][1])] {
						def = defs[d]
						break
					}
				}

				cut := strings.Split(rawcode[class][line][1], "(")
				name := cut[0]
				getargs = GetArgsFunc(def)

				if len(name) > 4 && name[:4] == "_set" {
					name = CamelCase("inverse" + name[4:])
				} else {
					name = CamelCase(strings.Trim(cut[0], "_"))
				}

				if len(rawcode[class][line-1][0]) > 5 && string(rawcode[class][line-1][0][:5]) == "@api." {
					body += "  //"
					for w := range rawcode[class][line-1] {

						body += rawcode[class][line-1][w]
					}
					body += "\n"
				}

				body += "  /*def " + def + "*/"

				if len(getargs) > 1 {
					args = " ,"
					for g := range getargs {
						if strings.HasPrefix(getargs[g], "Self") {
							continue
						}
						s := strings.Split(getargs[g], "=")
						args += strings.ToLower(TrimString(s[0])) + " interface{},"
					}
				}
				args = strings.TrimRight(args, ",")

				result += "pool." + classname + "().Methods()." + name + "().DeclareMethod(" +
					"\n`" + name + "` ," +
					"\nfunc (rs pool." + classname + "Set" + args + "){\n" +
					body + "})\n"
			}

		}
		if fieldDeclarationStarted {
			result += "\n})\n"
			fieldDeclarationStarted = false
		}

	}
	return result
}

//Camel case for standard string
func CamelCase(onestring string) string {

	var result string

	uppercase := true
	for c := range onestring {

		if uppercase == true {

			result += string(unicode.ToUpper(rune(onestring[c])))
			uppercase = false
		} else if string(onestring[c]) != "_" && string(onestring[c]) != "." {

			result += string(onestring[c])
		} else {
			uppercase = true
		}

	}

	return result
}

//Camel case for string in related fields
func CamelCaseForRelated(onestring string) string {

	var result string

	uppercase := true
	for c := range onestring {

		if uppercase == true {

			result += string(unicode.ToUpper(rune(onestring[c])))
			uppercase = false
		} else if string(onestring[c]) == "." {

			result += string(onestring[c])
			uppercase = true
		} else if string(onestring[c]) != "_" {

			result += string(onestring[c])
		} else {
			uppercase = true
		}

	}

	return result
}

//return a string corresponding to the args of a sql constraint
func GetArgsSqlConstraint(arg string) string {

	r := regexp.MustCompile(`'(.+)?'`)

	for c := range arg {
		if string(arg[c]) == "'" {
			r = regexp.MustCompile(`'(.+)?'`)
			break
		}
		if string(arg[c]) == "\"" {
			r = regexp.MustCompile(`"(.+)?"`)
			break
		}
	}

	res := r.FindStringSubmatch(arg)
	if len(res) > 1 {
		return res[1]
	}

	return strings.Join(res, " ")
}

//return a slice of arguments from the given field
func GetArgsFields(c int, l int) []string {
	var result []string
	var s string
	var count = 0
	ok := false

	for ok == false {
		for w := range rawcode[c][l+count] {
			s += " " + rawcode[c][l+count][w]
		}
		if string(s[len(s)-1]) != ")" {
			count += 1
		} else {
			ok = true
		}
	}

	s = strings.TrimSpace(s)
	cut := strings.SplitN(s, "(", 2)
	result = strings.Split(cut[1], ",")

	return result

}

//return a slice of the arguments from the given function
func GetArgsFunc(s string) []string {
	var args []string

	cut := strings.Split(s, ")")
	cut1 := strings.Split(cut[0], "(")
	if len(cut1) > 1 {
		args = strings.Split(cut1[1], ",")
	}

	for a := range args {

		args[a] = CamelCase(TrimString(strings.TrimSpace(args[a])))
		args[a] = args[a]
	}

	return args
}

//Get the text from an help argument by ignoring certain characters
func GetHelpText(c int, l int) string {
	var s string
	var count = 0

	for {
		for w := range rawcode[c][l+count] {
			s += " "
			s += rawcode[c][l+count][w]

		}
		if string(s[len(s)-1]) == ")" {
			break
		}
		s = strings.TrimSpace(s) + "#~#~#"
		count += 1
	}

	s = strings.TrimSpace(s)
	s = strings.Trim(s, ")")

	cut := strings.SplitN(s, "(", 2)
	cut1 := strings.SplitN(cut[1], "help=", 2)
	if len(cut1) <= 1 {
		return ""
	}
	cut2 := strings.SplitN(cut1[1], "\",", 2)
	help := cut2[0]
	help = strings.Replace(help, "\"#~#~# \"", " ", -1)

	return TrimString(help)
}

// trim ' or " characters from a given string
func TrimString(s string) string {

	var result string

	for x := range s {
		if string(s[x]) == "'" {
			result = strings.Trim(s, "'")
			break
		} else {
			result += strings.Trim(s, "\"")
			break
		}
	}
	return result
}

// write the body of hexya.go
func GenerateHexya() string {

	var result string

	result += "package " + packagename + "\n\n"
	result += "import(\n\"github.com/hexya-erp/hexya/hexya/server\"\n)\n\n"
	result += "const MODULE_NAME string = \"" + packagename + "\"\n\n"
	result += "func init() {\nserver.RegisterModule(&server.Module{\nName:     MODULE_NAME,\nPostInit: func() {},\n})\n}"

	return result
}

//verify the class name given and replace it by it's hexya equivalent if necessary
func CheckBuitInNames(classname string) string {

	var result string

	switch classname {

	case "ResCompany":
		result = "Company"
	case "ResPartner":
		result = "Partner"
	case "ResCurrency":
		result = "Currency"
	case "ResCountryGroup":
		result = "CountryGroup"
	default:
		result = classname
	}

	return result
}
