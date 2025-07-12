package news

import (
	"bufio"
	"bytes"
	"html"
	"image/jpeg"
	"regexp"
	"strings"
	"fmt"
	"io"
	"net/http"

	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/image/draw"
	"github.com/PuerkitoBio/goquery"

	"image"
	_ "image/jpeg"
	_ "image/png"
)

func HttpGet(url string, userAgent ...string) ([]byte, error) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    if len(userAgent) > 0 && userAgent[0] != "" {
        req.Header.Set("User-Agent", userAgent[0])
    }

    resp, err := client.Do(req)
    if err != nil || resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP request failed: %v", err)
    }

    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return body, nil
}

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
		"\n":        " ",
	}

	for entity, char := range replacements {
		content = strings.ReplaceAll(content, entity, char)
	}

	return strings.TrimSpace(content)
}

func ExtractImageCaption(articleURL string, find string) string {
    if articleURL == "" {
        return ""
    }

    data, err := HttpGet(articleURL)
    if err != nil {
        return ""
    }

    doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
    if err != nil {
        return ""
    }

    caption := doc.Find(find).Text()
    if caption != "" {
        caption = strings.ReplaceAll(caption, "Quelle:", "")
        caption = strings.TrimSpace(caption)
        return caption
    }

    return ""
}