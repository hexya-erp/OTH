package translate

import (
	"encoding/csv"
	"io"
	"log"
	"strings"
)

//create authorisations  in go from a csv file
func TransCSV(csvfile string, packagename string) string {

	var imports string
	var content string
	var result string

	r := csv.NewReader(strings.NewReader(string(csvfile)))

	for {
		var model = ""
		var group = ""
		var read = false
		var write = false
		var create = false
		var unlink = false

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		for re := range record {

			if re == 2 {
				model = CamelCase(record[re])
			} else if re == 3 {
				group = CamelCase(record[re])
			} else if re == 4 && record[re] == "1" {
				read = true
			} else if re == 5 && record[re] == "1" {
				write = true
			} else if re == 6 && record[re] == "1" {
				create = true
			} else if re == 7 && record[re] == "1" {
				unlink = true
			}
		}

		if group == "" {
			group = "security.GroupEveryone"
		}
		if group[:4] == "Base" {
			group = "base." + group[4:]
		}

		if model[:5] == "Model" {
			model = model[5:]
		}

		if read == true && write == true {
			content += "pool." + model + "().Methods().AllowAllToGroup(" + group + ")\n"

		} else {

			if read == true {
				content += "pool." + model + "().Methods().Load().AllowGroup(" + group + ")\n"
			}
			if write == true {
				content += "pool." + model + "().Methods().Write().AllowGroup(" + group + ")\n"
			}
			if create == true {
				content += "pool." + model + "().Methods().Create().AllowGroup(" + group + ")\n"
			}
			if unlink == true {
				content += "pool." + model + "().Methods().Unlink().AllowGroup(" + group + ")\n"
			}
		}

	}

	result = "package " + packagename + " \n\n  import (\n" + imports + "\n) \n\n func init() { \n\n " + content + " \n }"

	return result
}
