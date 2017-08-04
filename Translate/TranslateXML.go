package Translate

import "github.com/beevik/etree"

func TransXML(doc *etree.Document) string {

	var result string

	root := doc.SelectElement("odoo")

	println("Root element: " + root.Tag)

	for _, data := range root.SelectElements("data") {
		println("Child element: " + data.Tag)

		for _, record := range data.SelectElements("record") {
			println("\nChild element2: " + record.Tag)
			println("____________________________________________")

			for _, field := range record.SelectElements("field") {

				for _ , s := range field.Attr{

					println(s.Value)
				}
			}
		}
	}

	return result
}
