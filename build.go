// 2. Traverse templates folder to find files and convert using html/template and array of structs

package main

var markdownFiles []MarkdownFile

func build() {
	files := getContent()
	// assets := getAssets()
	templates := getTemplates()

	// log.Println(files)
	// log.Println(assets)
	// log.Println(templates)

	parseMarkdown(files)
	//parseAssets()

	// ts, err := template.ParseFiles("./templates/baseof.html", "./templates/home.html")
	// if err != nil {
	// 	panic(err.Error())
	// }

	// log.Println(ts.Name())
	// log.Println(ts.Tree)
	// log.Println(ts.Tree.Root.Nodes)

	// outputFile, err := os.Create("test.html")
	// if err != nil {
	// 	panic(err)
	// }
	// defer outputFile.Close()

	// data := MarkdownFile{
	// 	Title:       "Hello",
	// 	Path:        "",
	// 	Destination: "",
	// 	Content:     "This is some content",
	// 	Metadata:    nil,
	// }

	// err = ts.ExecuteTemplate(outputFile, "baseof", data)
	parseTemplates(templates)

}
