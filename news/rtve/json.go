package rtve

import (
	"NewsChannel/news"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var htmlRegex = regexp.MustCompile("<.*?>")

func httpGet(url string) ([]byte, error) {
    c := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
	
    resp, err := c.Do(req)
    if err != nil || resp.StatusCode != http.StatusOK {
        return nil, err
    }

    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return body, nil
}

// RTVEResponse represents the structure of RTVE API response
type RTVEResponse struct {
    Page struct {
        Items       []RTVEArticle `json:"items"`
        Number      int           `json:"number"`
        Size        int           `json:"size"`
        Offset      int           `json:"offset"`
        Total       int           `json:"total"`
        TotalPages  int           `json:"totalPages"`
        NumElements int           `json:"numElements"`
    } `json:"page"`
}

type RTVEArticle struct {
    ID                       string  `json:"id"`
    Title                    string  `json:"title"`
    AnteTitle                *string `json:"anteTitle"`
    LongTitle                string  `json:"longTitle"`
    ShortTitle               *string `json:"shortTitle"`
    Summary                  string  `json:"summary"`
    Text                     string  `json:"text"`
    Image                    string  `json:"image"`
    ImageSEO                 string  `json:"imageSEO"`
    HTMLUrl                  string  `json:"htmlUrl"`
    HTMLShortUrl             string  `json:"htmlShortUrl"`
    PublicationDate          string  `json:"publicationDate"`
    ModificationDate         string  `json:"modificationDate"`
    PublicationDateTimestamp int64   `json:"publicationDateTimestamp"`
    ContentType              string  `json:"contentType"`
    Language                 string  `json:"language"`
    MainCategory             string  `json:"mainCategory"`
    PubState                 struct {
        Code        string `json:"code"`
        Description string `json:"description"`
    } `json:"pubState"`
}

func (r *RTVE) getArticles(url string, topic news.Topic) ([]news.Article, error) {
    data, err := httpGet(url)
    if err != nil {
        return nil, err
    }

    var response RTVEResponse
    err = json.Unmarshal(data, &response)
    if err != nil {
        return nil, err
    }

	log.Printf("Fetched %d articles from RTVE for topic %d", len(response.Page.Items), topic)

    var articles []news.Article
    for i, rtveArticle := range response.Page.Items {
        if i >= 1 { // Limit to 1 article per category like Reuters
            break
        }

        // Use longTitle if available, otherwise fall back to title
        title := rtveArticle.Title
        if rtveArticle.LongTitle != "" {
            title = rtveArticle.LongTitle
        }

        // Check for duplicates
        if news.IsDuplicateArticle(r.oldArticleTitles, title) {
            continue
        }

        // Use the text field as content, clean HTML tags
        content := rtveArticle.Text
        if content == "" {
            // Fall back to summary if no text
            content = rtveArticle.Summary
        }

        // Clean HTML tags from content
        content = htmlRegex.ReplaceAllString(content, "")
        
        // Skip if no content
        if len(strings.TrimSpace(content)) == 0 {
            continue
        }

        // Get thumbnail - try imageSEO first, then image
        var thumbnail *news.Thumbnail
        imageURL := rtveArticle.ImageSEO
        if imageURL == "" {
            imageURL = rtveArticle.Image
        }
        
        if imageURL != "" {
            thumbnail, _ = r.getThumbnail(imageURL) // Ignore errors, continue without thumbnail
        }

        // Parse location from content and category
        location := r.extractLocation(content, rtveArticle.MainCategory)

        article := news.Article{
            Title:     title,
            Content:   &content,
            Topic:     topic,
            Location:  location,
            Thumbnail: thumbnail,
        }

        articles = append(articles, article)
    }

    return articles, nil
}

func (r *RTVE) getThumbnail(imageURL string) (*news.Thumbnail, error) {
    if imageURL == "" {
        return nil, nil
    }

    // Ensure URL is absolute
    if !strings.HasPrefix(imageURL, "http") {
        if strings.HasPrefix(imageURL, "//") {
            imageURL = "https:" + imageURL
        } else if strings.HasPrefix(imageURL, "/") {
            imageURL = "https://img.rtve.es" + imageURL
        } else {
            imageURL = "https://img.rtve.es/" + imageURL
        }
    }

    data, err := httpGet(imageURL)
    if err != nil {
        return nil, err
    }

    if len(data) == 0 {
        return nil, nil
    }

    return &news.Thumbnail{
        Image:   news.ConvertImage(data),
        Caption: "",
    }, nil
}

func (r *RTVE) extractLocation(text, category string) *news.Location {
    // Extract the location from the category path
    // For instance "Noticias/España/Cataluña/Barcelona" -> "Barcelona"
    if category != "" {
        parts := strings.Split(category, "/")
        if len(parts) > 0 {
            lastPart := strings.TrimSpace(parts[len(parts)-1])
            if lastPart != "" && lastPart != "Noticias" && lastPart != "España" {
                // Clean up common category names that aren't locations
                if lastPart != "Especiales" && lastPart != "Nacional" && lastPart != "Internacional" {
                    return &news.Location{Name: strings.ToUpper(lastPart)}
                }
            }
        }
    }

    // Fallback to text-based location detection
    searchText := strings.ToUpper(text + " " + category)

    // Common Spanish locations
    commonSpanishLocations := map[string]news.Location{
        "MADRID":     {Name: "MADRID"},
        "BARCELONA":  {Name: "BARCELONA"},
        "VALENCIA":   {Name: "VALENCIA"},
        "SEVILLA":    {Name: "SEVILLA"},
        "SEVILLE":    {Name: "SEVILLA"},
        "BILBAO":     {Name: "BILBAO"},
        "ZARAGOZA":   {Name: "ZARAGOZA"},
        "MÁLAGA":     {Name: "MÁLAGA"},
        "MALAGA":     {Name: "MÁLAGA"},
        "MURCIA":     {Name: "MURCIA"},
        "PALMA":      {Name: "PALMA"},
        "SANTANDER":  {Name: "SANTANDER"},
        "CÓRDOBA":    {Name: "CÓRDOBA"},
        "CORDOBA":    {Name: "CÓRDOBA"},
        "VALLADOLID": {Name: "VALLADOLID"},
        "TOLEDO":     {Name: "TOLEDO"},
        "PAMPLONA":   {Name: "PAMPLONA"},
        "BURGOS":     {Name: "BURGOS"},
        "SALAMANCA":  {Name: "SALAMANCA"},
        "CÁDIZ":      {Name: "CÁDIZ"},
        "CADIZ":      {Name: "CÁDIZ"},
        "HUELVA":     {Name: "HUELVA"},
        "BADAJOZ":    {Name: "BADAJOZ"},
        "CÁCERES":    {Name: "CÁCERES"},
        "CACERES":    {Name: "CÁCERES"},
        "LEÓN":       {Name: "LEÓN"},
        "LEON":       {Name: "LEÓN"},
        "PONTEVEDRA": {Name: "PONTEVEDRA"},
        "VIGO":       {Name: "VIGO"},
        "CORUÑA":     {Name: "A CORUÑA"},
        "SANTIAGO":   {Name: "SANTIAGO"},
        "OVIEDO":     {Name: "OVIEDO"},
        "GIJÓN":      {Name: "GIJÓN"},
        "GIJON":      {Name: "GIJÓN"},
        "VITORIA":    {Name: "VITORIA"},
        "DONOSTIA":   {Name: "DONOSTIA"},
        "SAN SEBASTIÁN": {Name: "DONOSTIA"},
        "NAVARRA":    {Name: "NAVARRA"},
        "CATALUÑA":   {Name: "CATALUÑA"},
        "CATALUNYA":  {Name: "CATALUÑA"},
        "ANDALUCÍA":  {Name: "ANDALUCÍA"},
        "ANDALUCIA":  {Name: "ANDALUCÍA"},
    }

    for city, location := range commonSpanishLocations {
        if strings.Contains(searchText, city) {
            return &location
        }
    }

    return nil
}