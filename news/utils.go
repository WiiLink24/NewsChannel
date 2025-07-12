package news

import (
	"bufio"
	"bytes"
	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/image/draw"
	"image/jpeg"
	"html"
	"regexp"
	"strings"

	"image"
	_ "image/jpeg"
	_ "image/png"
)

func IsDuplicateArticle(previousArticles []string, currentArticle string) bool {
	for _, previousArticle := range previousArticles {
		diff := difflib.NewMatcher([]string{currentArticle}, []string{previousArticle})
		if diff.QuickRatio() >= 0.85 {
			return true
		}
	}

	return false
}

func ConvertImage(data []byte) []byte {
	origImage, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil
	}

	newImage := image.NewRGBA(image.Rect(0, 0, 200, 200))
	draw.BiLinear.Scale(newImage, newImage.Bounds(), origImage, origImage.Bounds(), draw.Over, nil)

	var outputImgWriter bytes.Buffer
	err = jpeg.Encode(bufio.NewWriter(&outputImgWriter), newImage, nil)
	if err != nil {
		return nil
	}

	return outputImgWriter.Bytes()
}

func CleanHTMLEntities(content string) string {
	content = html.UnescapeString(content)

	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	content = htmlTagRegex.ReplaceAllString(content, "")
	
	// RTVE specific tag that calls on another article
	content = strings.ReplaceAll(content, "@@NOTICIA[16657506,IMAGEN,FIRMA]", "")

	replacements := map[string]string{
		"&nbsp;":   " ",
		"&lt;":     "<",
		"&gt;":     ">",
		"&amp;":    "&",
		"&quot;":   "\"",
		"&apos;":   "'",
		"&uacute;": "ú",
		"&iacute;": "í",
		"&oacute;": "ó",
		"&aacute;": "á",
		"&eacute;": "é",
		"&ntilde;": "ñ",
		"&Uacute;": "Ú",
		"&Iacute;": "Í",
		"&Oacute;": "Ó",
		"&Aacute;": "Á",
		"&Eacute;": "É",
		"&Ntilde;": "Ñ",
		"&uuml;":   "ü",
		"&Uuml;":   "Ü",
	}

	for entity, char := range replacements {
		content = strings.ReplaceAll(content, entity, char)
	}

	return strings.TrimSpace(content)
}
