package news

import (
	"bufio"
	"bytes"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
)

/*func IsDuplicateArticle(articles []Article, currentArticle string) bool {
	for _, article := range articles {
		diff := difflib.NewMatcher([]string{currentArticle}, []string{article.Article.Title})
		if diff.QuickRatio() >= 0.85 {
			return true
		}
	}

	return false
}*/

func ConvertImage(data []byte) []byte {
	origImage, err := jpeg.Decode(bytes.NewReader(data))
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
