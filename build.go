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
	parseTemplates(templates)
}
