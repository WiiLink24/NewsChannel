package news

import (
	"bufio"
	"bytes"
	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/image/draw"
	"image/jpeg"

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
