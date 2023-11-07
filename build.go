// 2. Traverse templates folder to find files and convert using html/template and array of structs

package main

import (
	"html/template"
)

var singlePaths []string
var singleFiles []MarkdownFile

var contentPaths map[string][]string
var contentFiles map[string]ContentType

var templatePaths []string
var templates map[string]*template.Template

var assetPaths []string

func build() {

	singlePaths = getSingle()
	contentPaths = getContent()
	templatePaths = getTemplates()
	assetPaths = getAssets()

	singleFiles = parseMarkdown(singlePaths)
	contentFiles = parseContent(contentPaths)
	templates = parseTemplates(templatePaths)
	parseAssets(assetPaths)

	buildFiles(singleFiles, templates)
	buildContent(contentFiles, templates)
}
