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
var imports string
var defs []string

func TransPyToGo(str string) string {

	var content string

	packagename = ""
	packagenameset = false
	rawcode = nil
	imports = ""
	defs = nil

	GenerateSlices(str)
	GenerateDefs(str)
	content = TransRules()

	var result string = "package " + packagename + " \n\n  import (\n" + imports + "\n) \n\n func init() { \n\n " + content + " \n }"

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

func GenerateDefs(str string) {

	cut := strings.Split(str, "@api.")

	for c := range cut {
		cut[c] = "@api." + cut[c]
		def := strings.Split(cut[c], "def ")

		for d := range def {
			finaldef := strings.Split(def[d], "class")
			for f := range finaldef {

				defs = append(defs, finaldef[f])
			}
		}
	}

}

// return a string in go code corresponding to the original python code
func TransRules() string {

	var result string
	var selectionimportset bool = false
	for class := range rawcode {
		getclassname := strings.Split(rawcode[class][0][0], "(")
		classname := getclassname[0]
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
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])
					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					} else {
						name = fieldname
					}
					body += "String :" + name

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "required":
							body += ", Required: " + strings.ToLower(value[1])
						case "translate":
							body += ", Translate: " + strings.ToLower(value[1])
						case "compute":
							body += ", Compute: \"" + CamelCase(strings.Trim(strings.Trim(value[1], "'"), "_")) + "\""
						case "help":
							help := GetHelpText(class, line)

							for i := range help {
								if help[i][len(help[i])-4:] == "help" {

									regex, err := regexp.Compile("\"")
									if err != nil {
										return err.Error()
									}
									cut := help[i+1]
									cut = regex.ReplaceAllString(cut, "")
									body += " ,Help :\"" + cut + "\""
									body += " ,Help :\"" + cut + "\""
								}
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
						case "copy":
							if value[1] == "True" {
								body += ", NoCopy: false"
							} else {
								body += ", NoCopy: true"
							}
						case "oldname":
							//TODO

						default:
							//println("Char: " + value[0])
						}
					}
					result += "pool." + classname + "().AddCharField(" + fieldname + ", models.StringFieldParams{" + body + "})\n"

				case "Many2one":
					var body string
					var readonly string
					args := GetArgsFields(class, line)
					name := ""
					args[1] = strings.TrimSpace(args[1])
					if string(args[1][0]) == "'" || string(args[1][0]) == "\"" {
						name = "\"" + TrimString(args[1]) + "\""
					} else {
						name = fieldname
					}
					foreignkey := CamelCase(strings.Trim(strings.TrimSpace(args[0]), "'"))

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "required":
							body += ", Required: " + strings.ToLower(value[1])
						case "default":
							if string(value[1][0]) == "_" {
								body += ", Default : pool." + classname + "." + CamelCase(strings.Trim(value[1], "_")) + "()"
							} else {
								body += ", Default: func(models.Environment, models.FieldMap) interface{} {return " + value[1] + "}"
							}
						case "ondelete":
							body += ", OnDelete : models." + CamelCase(TrimString(value[1]))
						case "help":
							help := GetHelpText(class, line)

							for i := range help {
								if help[i][len(help[i])-4:] == "help" {

									regex, err := regexp.Compile("\"")
									if err != nil {
										return err.Error()
									}
									cut := help[i+1]
									cut = regex.ReplaceAllString(cut, "")
									body += " , Help :\"" + cut + "\""
								}
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
						case "readonly":
							readonly = "pool." + classname + "().Fields()." + strings.Trim(fieldname, "\"") + "().RevokeAccess(security.GroupEveryone, security.Write)\n"
						case "related":
							body += ", Related: \"" + CamelCase(strings.Trim(value[1], "'")) + "\""
						case "store":
							body += ", Stored: " + strings.ToLower(value[1])
						case "string":
							name = "\"" + TrimString(strings.TrimSpace(value[1])) + "\""
						case "auto_join":
							//TODO
						case "domain":
							//TODO
						case "oldname":
							//TODO

						default:
							//println("Many2One: " + value[0])
						}
					}
					body = "String :" + name + " , RelationModel: pool." + foreignkey + "()" + body
					result += "pool." + classname + "().AddMany2OneField(" + fieldname + ",models.ForeignKeyFieldParams{" + body + "})\n"
					result += readonly

				case "One2many":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[2] = strings.TrimSpace(args[2])
					if string(args[2][0]) == "'" || string(args[2][0]) == "\"" {
						name = "\"" + TrimString(args[2]) + "\""
					} else {
						name = fieldname
					}
					body += "String :" + name
					body += " ,RelationModel : pool." + CamelCase(TrimString(strings.TrimSpace(args[0]))) + "()"
					body += " ,ReverseFK : \"" + CamelCase(TrimString(strings.TrimSpace(args[1]))) + "\""

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
								body += ", Default : pool." + classname + "." + CamelCase(strings.Trim(value[1], "_")) + "()"
							} else {
								body += ", Default: func(models.Environment, models.FieldMap) interface{} {return " + value[1] + "}"
							}
						default:
							//println("One2Many: " + value[0])
						}
					}
					result += "pool." + classname + "().AddOne2ManyField(" + fieldname + ", models.ReverseFieldParams{" + body + "})\n"

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
					cut2 := strings.Split(cut[1], "]")
					args := strings.Split(cut2[1], ",")
					selectable := strings.Split(cut2[0], "),")

					name := ""
					args[1] = strings.TrimSpace(args[1])
					if string(args[1][0]) == "'" || string(args[1][0]) == "\"" {
						name = "\"" + TrimString(args[1]) + "\""
					} else {
						name = fieldname
					}
					body += "String :" + name

					body += ", Selection : types.Selection{\n"

					for s := range selectable {
						sec := strings.Split(selectable[s], ",")
						if len(sec) > 2 {
							i := 1
							for len(sec)-2 >= i {
								sec[1] += sec[1+i]

								i++
							}
							sec = sec[:len(sec)-(len(sec)-2)]
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
							body += ", Default: func(models.Environment, models.FieldMap) interface{} {return \"" + TrimString(value[1]) + "\"}"
						case "required":
							body += ", Required: " + strings.ToLower(value[1])
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
						default:
							//println("selection: " + value[0])
						}
					}
					result += "pool." + classname + "().AddSelectionField(" + fieldname + ", models.SelectionFieldParams{" + body + "})\n"
				case "Integer":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])
					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					} else {
						name = fieldname
					}
					body += "String :" + name

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "help":
							help := GetHelpText(class, line)

							for i := range help {
								if help[i][len(help[i])-4:] == "help" {

									regex, err := regexp.Compile("\"")
									if err != nil {
										return err.Error()
									}
									cut := help[i+1]
									cut = regex.ReplaceAllString(cut, "")
									body += " ,Help :\"" + cut + "\""
								}
							}
						case "default":
							if string(value[1][0]) == "_" {
								body += ", Default : pool." + classname + "." + CamelCase(strings.Trim(value[1], "_")) + "()"
							} else {
								body += ", Default: func(models.Environment, models.FieldMap) interface{} {return " + value[1] + "}"
							}
						case "required":
							body += ", Required: " + strings.ToLower(value[1])
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
							body += ", Compute : pool." + CamelCase(strings.Trim(TrimString(value[1]), "_")) + "()"

						default:
							//println("Integer: " + value[0])
						}
					}
					result += "pool." + classname + "().AddIntegerField(" + fieldname + ", models.SimpleFieldParams{" + body + "})\n"

				case "Datetime":
					result += "pool." + classname + "().AddDateTimeField(" + fieldname + ", models.SimpleFieldParams{})\n"

				case "Float":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])

					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					} else {
						name = fieldname
					}
					body += "String :" + name

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "default":
							if string(value[1][0]) == "_" {
								body += ", Default : pool." + classname + "." + CamelCase(strings.Trim(value[1], "_")) + "()"
							} else {
								body += ", Default: func(models.Environment, models.FieldMap) interface{} {return " + value[1] + "}"
							}
						case "digits":
							//TODO
						case "compute":
							body += ", Compute: \"" + CamelCase(strings.Trim(strings.Trim(value[1], "'"), "_")) + "\""
						case "help":
							help := GetHelpText(class, line)

							for i := range help {
								if help[i][len(help[i])-4:] == "help" {

									regex, err := regexp.Compile("\"")
									if err != nil {
										return err.Error()
									}
									cut := help[i+1]
									cut = regex.ReplaceAllString(cut, "")
									body += " ,Help :\"" + cut + "\""
								}
							}
						case "required":
							body += ", Required: " + strings.ToLower(value[1])
						case "company_dependent":
							//TODO

						default:
							//println("Float: " + value[0])
						}
					}
					result += "pool." + classname + "().AddFloatField(" + fieldname + ", models.FloatFieldParams{" + body + "})\n"

				case "Boolean":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])
					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					} else {
						name = fieldname
					}
					body += "String :" + name

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "default":
							if len(value[1]) == 1 {
								if value[1] == "0" {
									body += ", Default: func(models.Environment, models.FieldMap) interface{} {return false}"
								} else {
									body += ", Default: func(models.Environment, models.FieldMap) interface{} {return true}"
								}
							} else {
								body += ", Default: func(models.Environment, models.FieldMap) interface{} {return " + strings.ToLower(value[1]) + "}"
							}

						case "help":
							help := GetHelpText(class, line)

							for i := range help {
								if help[i][len(help[i])-4:] == "help" {

									regex, err := regexp.Compile("\"")
									if err != nil {
										return err.Error()
									}
									cut := help[i+1]
									cut = regex.ReplaceAllString(cut, "")
									body += " ,Help :\"" + cut + "\""
								}
							}
						default:
							//println("Boolean: " + value[0])
						}
					}
					result += "pool." + classname + "().AddBooleanField(" + fieldname + ", models.SimpleFieldParams{" + body + "})\n"

				case "Many2many":
					var body string
					var readonly string
					args := GetArgsFields(class, line)
					name := ""
					args[1] = strings.TrimSpace(args[1])
					if string(args[1][0]) == "'" || string(args[1][0]) == "\"" {
						name = "\"" + TrimString(args[1]) + "\""
					} else {
						name = fieldname
					}

					foreignkey := CamelCase(strings.Trim(strings.TrimSpace(args[0]), "'"))

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "string":
							name = "\"" + TrimString(strings.TrimSpace(value[1])) + "\""
						case "ondelete":
							body += ", OnDelete : models." + CamelCase(TrimString(value[1]))
						case "compute":
							body += ", Compute: \"" + CamelCase(strings.Trim(strings.Trim(value[1], "'"), "_")) + "\""
						case "readonly":
							readonly = "pool." + classname + "().Fields()." + strings.Trim(fieldname, "\"") + "().RevokeAccess(security.GroupEveryone, security.Write)\n"
						default:
							//println("Many2Many: " + value[0])
						}
					}
					body = "String :" + name + " , RelationModel: pool." + foreignkey + "()" + body
					result += "pool." + classname + "().AddMany2ManyField(" + fieldname + ", models.Many2ManyFieldParams{" + body + "})\n"
					result += readonly

				case "Binary":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])
					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					} else {
						name = fieldname
					}
					body += "String :" + name

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "attachment":
							//TODO
						case "help":
							help := GetHelpText(class, line)

							for i := range help {
								if help[i][len(help[i])-4:] == "help" {

									regex, err := regexp.Compile("\"")
									if err != nil {
										return err.Error()
									}
									cut := help[i+1]
									cut = regex.ReplaceAllString(cut, "")
									body += " ,Help :\"" + cut + "\""
								}
							}
						case "compute":
							body += ", Compute: \"" + CamelCase(strings.Trim(strings.Trim(value[1], "'"), "_")) + "\""
						case "inverse":
							//TODO

						default:
							//println("Binary: " + value[0])
						}
					}
					result += "pool." + classname + "().AddBinaryField(" + fieldname + ", models.SimpleFieldParams{" + body + "})\n"

				case "Date":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])
					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					} else {
						name = fieldname
					}
					body += "String :" + name

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {

						case "help":
							help := GetHelpText(class, line)

							for i := range help {
								if help[i][len(help[i])-4:] == "help" {

									regex, err := regexp.Compile("\"")
									if err != nil {
										return err.Error()
									}
									cut := help[i+1]
									cut = regex.ReplaceAllString(cut, "")
									body += " ,Help :\"" + cut + "\""
								}
							}
						default:
							//println("Date: " + value[0])
						}
					}
					result += "pool." + classname + "().AddDateField(" + fieldname + ", models.SimpleFieldParams{" + body + "})\n"

				case "DateTime":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])
					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					} else {
						name = fieldname
					}
					body += "String :" + name

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "help":
							help := GetHelpText(class, line)

							for i := range help {
								if help[i][len(help[i])-4:] == "help" {

									regex, err := regexp.Compile("\"")
									if err != nil {
										return err.Error()
									}
									cut := help[i+1]
									cut = regex.ReplaceAllString(cut, "")
									body += " ,Help :\"" + cut + "\""
								}
							}
						default:
							//println("DateTime: " + value[0])
						}
					}
					result += "pool." + classname + "().AddDateTimeField(" + fieldname + ", models.SimpleFieldParams{" + body + "})\n"

				case "Text":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])
					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					} else {
						name = fieldname
					}
					body += "String :" + name

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {
						case "translate":
							body += ", Translate: " + strings.ToLower(value[1])
						case "help":
							help := GetHelpText(class, line)

							for i := range help {
								if help[i][len(help[i])-4:] == "help" {

									regex, err := regexp.Compile("\"")
									if err != nil {
										return err.Error()
									}
									cut := help[i+1]
									cut = regex.ReplaceAllString(cut, "")
									body += " ,Help :\"" + cut + "\""
								}
							}
						default:
							//println("Text: " + value[0])
						}
					}

					result += "pool." + classname + "().AddTextField(" + fieldname + " , models.StringFieldParams{" + body + "})\n"

				case "Html":
					var body string
					args := GetArgsFields(class, line)
					name := ""
					args[0] = strings.TrimSpace(args[0])
					if string(args[0][0]) == "'" || string(args[0][0]) == "\"" {
						name = "\"" + TrimString(args[0]) + "\""
					} else {
						name = fieldname
					}
					body += "String :" + name

					for i := range args {
						arg := strings.Trim(args[i], ")")
						value := strings.Split(arg, "=")

						switch strings.TrimSpace(value[0]) {

						case "help":
							help := GetHelpText(class, line)

							for i := range help {
								if help[i][len(help[i])-4:] == "help" {

									regex, err := regexp.Compile("\"")
									if err != nil {
										return err.Error()
									}
									cut := help[i+1]
									cut = regex.ReplaceAllString(cut, "")
									body += " ,Help :\"" + cut + "\""
								}
							}

						default:
							//println("Html: " + value[0])
						}
					}
					result += "pool." + classname + "().AddHTMLField(" + fieldname + " , models.StringFieldParams{})\n"

				default:
					//println(fieldtype)

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

					if len(args) > 3 && args[3] != " " {
						i := 2
						var sqlstring = CamelCase(args[1])
						for string(sqlstring[len(sqlstring)-2]) != ")" {
							sqlstring += "," + CamelCase(args[i])
							i++
						}
						args[1] = sqlstring
						args[2] = strings.TrimRight(strings.TrimLeft(args[len(args)-1], "("), ")")
					}

					name := GetArgsSqlConstraint(args[0])
					sql := GetArgsSqlConstraint(args[1])
					errorstring := GetArgsSqlConstraint(args[2])

					result += "pool." + classname + "().AddSQLConstraint(" + name + " , " + sql + " , " + errorstring + ")\n"

					count += 1
				}

			} else if rawcode[class][line][0] == "def" {

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
				name := CamelCase(strings.Trim(cut[0], "_"))
				getargs = GetArgsFunc(def)

				if len(rawcode[class][line-1][0]) > 5 && string(rawcode[class][line-1][0][:5]) == "@api." {
					body += "  //"
					for w := range rawcode[class][line-1] {

						body += rawcode[class][line-1][w]
					}
					body += "\n"
				}

				body += "  /*def " + def + "*/"

				if len(getargs) > 1 {
					args = " , args struct{"
					for g := range getargs {
						if len(getargs[g]) > 4 {
							if getargs[g][:4] != "Self" {
								s := strings.Split(getargs[g], "=")
								args += TrimString(s[0]) + " interface{}\n"
							}
						}

					}
					args += "}"
				}

				result += "pool." + classname + "().Methods()." + name + "().DeclareMethod(" +
					"\n`" + name + "` ," +
					"\nfunc (rs pool." + classname + "Set" + args + "){\n" +
					body + "})\n"
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
		} else if string(onestring[c]) != "_" && string(onestring[c]) != "." {

			result += string(onestring[c])
		} else {
			uppercase = true
		}

	}

	return result
}

func GetArgsSqlConstraint(arg string) string {

	var result string
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
	result = "\"" + CamelCase(res[1]) + "\""

	return result
}

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
	cut := strings.Split(s, "(")
	result = strings.Split(cut[1], ",")

	return result
}

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

func GetHelpText(c int, l int) []string {
	var result []string
	var s string
	var count = 0
	ok := false

	for ok == false {
		for w := range rawcode[c][l+count] {
			s += " "
			s += rawcode[c][l+count][w]

		}
		if string(s[len(s)-1]) != ")" && string(s[len(s)-2]) != "\"" {
			count += 1
		} else {
			ok = true
		}
	}

	s = strings.TrimSpace(s)
	s = strings.Trim(s, ")")

	cut := strings.Split(s, "(")
	result = strings.Split(cut[1], "=")

	return result
}

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