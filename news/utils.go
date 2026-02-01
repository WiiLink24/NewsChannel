package news

import (
	"bufio"
	"bytes"
	"fmt"
	"html"
	"image/jpeg"
	"io"
	"math"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/image/draw"

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
	req.Header.Set("User-Agent", "WiiLink News Channel File Generator")

	if len(userAgent) > 0 && userAgent[0] != "" {
		req.Header.Set("User-Agent", userAgent[0])
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request to %v failed: %v", url, err)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request to %v failed: Status Code %v", url, resp.StatusCode)
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
		diff := difflib.NewMatcher(strings.Split(currentArticle, ""), strings.Split(previousArticle, ""))
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

	origBounds := origImage.Bounds()
	origWidth := origBounds.Dx()
	origHeight := origBounds.Dy()

	scaleFactor := math.Max(
		float64(200)/float64(origWidth),
		float64(200)/float64(origHeight),
	)

	resizeWidth := int(math.Ceil(float64(origWidth) * scaleFactor))
	resizeHeight := int(math.Ceil(float64(origHeight) * scaleFactor))

	resizedImage := image.NewRGBA(image.Rect(0, 0, resizeWidth, resizeHeight))
	draw.BiLinear.Scale(resizedImage, resizedImage.Bounds(), origImage, origBounds, draw.Over, nil)

	var outputImgWriter bytes.Buffer
	err = jpeg.Encode(bufio.NewWriter(&outputImgWriter), resizedImage, nil)
	if err != nil {
		return nil
	}

	return outputImgWriter.Bytes()
}

func SanitizeText(content string) string {
	content = html.UnescapeString(content)

	iframeRegex := regexp.MustCompile(`(?s)<iframe.*?>.*?</iframe>`)
	content = iframeRegex.ReplaceAllString(content, "")

	scriptRegex := regexp.MustCompile(`(?s)<script.*?>.*?</script>`)
	content = scriptRegex.ReplaceAllString(content, "")

	// Remove all HTML tags
	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	content = htmlTagRegex.ReplaceAllString(content, "")

	// RTVE specific tag that calls on another article
	noticiaRegex := regexp.MustCompile(`@@NOTICIA\[[^\]]*\]`)
	content = noticiaRegex.ReplaceAllString(content, "")

	fotoRegex := regexp.MustCompile(`@@FOTO\[[^\]]*\]`)
	content = fotoRegex.ReplaceAllString(content, "")

	mediaRegex := regexp.MustCompile(`@@MEDIA\[[^\]]*\]`)
	content = mediaRegex.ReplaceAllString(content, "")

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
		"​":        "",
		"‌":        "",
		"‍":        "",
		"⁠":        "",
		"‑":        "-",
		"\t":       "",
		"ᵉ":        "e",
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
