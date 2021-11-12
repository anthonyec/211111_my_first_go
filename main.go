package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"gitlab.com/golang-commonmark/markdown"
)

type page struct {
	slug    string
	path    string
	assets  string
	title   string
	content string
}

func getTitleFromHTML(html string) string {
	headingRegex, _ := regexp.Compile("<h1>.*</h1>")
	anchorRegex, _ := regexp.Compile("<h1>.*</h1>")

	headingMatch := headingRegex.FindString(html)
	anchorMatch := anchorRegex.FindString(html)

	headingMatch = strings.Replace(headingMatch, "<h1>", "", 1)
	headingMatch = strings.Replace(headingMatch, "</h1>", "", 1)
	headingMatch = strings.Replace(headingMatch, anchorMatch, "", 1)

	return headingMatch
}

func getDateFromFileName(fileName string) string {
	regex, _ := regexp.Compile("^(19[0-9]{2}|2[0-9]{3})-(0[1-9]|1[012])-([123]0|[012][1-9]|31)")
	return regex.FindString(fileName)
}

func parseCollectionFromFilesystem(name string, source string) *[]page {
	md := markdown.New(markdown.XHTMLOutput(true))
	files, err := ioutil.ReadDir(source)
	pages := []page{}
	// pages := make([]page, len(files))

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fileName := file.Name()

		if fileName == ".DS_Store" {
			continue
		}

		filePath := path.Join(source, fileName)
		fileExtension := path.Ext(filePath)

		isDirectory := file.IsDir()
		markdownPath := filePath

		if isDirectory {
			markdownPath = path.Join(source, fileName, "index.md")
		}

		if _, err := os.Stat(markdownPath); errors.Is(err, os.ErrNotExist) {
			// fmt.Println("Markdown does not exist!")
			continue
		}

		fileContents, err := ioutil.ReadFile(markdownPath)

		if err != nil {
			panic(err)
		}

		content := md.RenderToString(fileContents)
		title := getTitleFromHTML(content)
		date := getDateFromFileName(fileName)

		slug := strings.Replace(fileName, fileExtension, "", 1)
		slug = strings.Replace(slug, date+"-", "", 1)

		cwd, _ := os.Getwd()
		assets := path.Join(cwd, filePath)

		if !isDirectory {
			assets = ""
		}

		// const [headers, contentWithoutHeaders] =
		//   getCommentHeadersFromContent(content);

		newPage := page{
			slug:    slug,
			path:    "./dist/" + slug,
			title:   title,
			content: content,
			assets:  assets,
		}

		pages = append(pages, newPage)
	}

	return &pages
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func main() {
	defer timeTrack(time.Now(), "posts")
	parseCollectionFromFilesystem("posts", "./content/_posts")

	// fmt.Println(collectionPages)
}
