package handlers

import (
	"io"
	"net/http"
	"strings"
	"golang.org/x/net/html"
	"net/url"
	"strconv"
	"encoding/json"
	"errors"
	"fmt"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json" )

	URL := r.FormValue("url")
	if URL == "" {
		http.Error(w,"URL supply error", 400)
	}

	body, err := fetchHTML(URL)
	defer body.Close()
	if err != nil {
		http.Error(w, "URL fetch error", 400)
	}

	pageSummary, err := extractSummary(URL, body)
	if err != nil {
		http.Error(w, "extracting summary error", 400)
	}
	json.NewEncoder(w).Encode(pageSummary)
	return
}


func fetchHTML(pageURL string) (io.ReadCloser, error) {
	
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New("StatusCode error")
	}

	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		return nil, errors.New("not a valid content type")
	}

	return resp.Body, nil
}

func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	
	tokenizer := html.NewTokenizer(htmlStream)
	page := new(PageSummary)
	HTMLTitle := false
	description := false
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
			//log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
			return nil, fmt.Errorf("error tokenizing HTML: %v", err)
		}
		if tokenType == html.EndTagToken {
			token := tokenizer.Token()
			if "head" == token.Data {
				break
			}
		}
		token := tokenizer.Token()
		if "title" == token.Data {
			if HTMLTitle == false {
				tokenType  = tokenizer.Next()
				if tokenType == html.TextToken {
					page.Title = tokenizer.Token().Data
					HTMLTitle = true
				}
			}
		}

		if "link" == token.Data {
			iconLink, iconType, iconHeight, iconWidth, check := iconHelper(token, "icon")
			if check {
				iconLink = absoluteURL(iconLink)
				p := &PreviewImage{}
				p.URL = iconLink
				p.Type = iconType
				p.Width = iconHeight
				p.Height = iconWidth
				page.Icon = p
			}
		}	

		if "meta" == token.Data {
			propTitle, check := extractHelper(token, "og:title")
			if check {
				page.Title = propTitle
				HTMLTitle = true
			}
			propType, check := extractHelper(token, "og:type")
			if check {
				page.Type = propType
			}
			propSiteName, check := extractHelper(token, "og:site_name")
			if check {
				page.SiteName = propSiteName
			}
			propDesc, check := extractHelper(token, "og:description")
			if check {
				page.Description = propDesc
				description = true
			} 
			if !check && description == false {
				propDesc2, check2 := extractHelper(token, "description")
				if check2 {
					page.Description = propDesc2
				}
			}
			propAuthor, check := extractHelper(token, "author")
			if check{
				page.Author = propAuthor
			}
			propURL, check := extractHelper(token, "og:url")
			if check {
				page.URL = propURL
			}
			propKeywords, check := extractHelper(token, "keywords")
			if check {
				s := strings.Split(propKeywords, ",")
				for i, val := range s {
					s[i] = strings.TrimSpace(val)
				}
				page.Keywords = s
			}
			imageLink, check := extractHelper(token, "og:image")
			if check {
				imageLink = absoluteURL(imageLink)
				p := &PreviewImage{}
				p.URL = imageLink
				page.Images = append(page.Images, p)
			}
			secureLink, check := extractHelper(token, "og:image:secure_url")
			if check {
				secureLink = absoluteURL(secureLink)
				page.Images[len(page.Images) - 1].SecureURL = secureLink
			}
			imageType, check := extractHelper(token, "og:image:type")
			if check {
				page.Images[len(page.Images) - 1].Type = imageType
			}
			imageW, check := extractHelper(token, "og:image:width")
			if check {
				page.Images[len(page.Images) - 1].Width, _ = strconv.Atoi(imageW)
			}
			imageH, check := extractHelper(token, "og:image:height")
			if check {
				page.Images[len(page.Images) - 1].Height, _ = strconv.Atoi(imageH)
			}
			alt, check := extractHelper(token, "og:image:alt") 
			if check {
				page.Images[len(page.Images) - 1].Alt = alt
			}
		}
	}
	return page, nil
}
	
	
func extractHelper(t html.Token, prop string) (content string, check bool) {
	for _, attr := range t.Attr {
		if (attr.Key == "property" || attr.Key == "name") && attr.Val == prop {
			check = true
		}

		if attr.Key == "content" {
			content = attr.Val
		}
	}

	return
}

func iconHelper(t html.Token, prop string) (href string, typ string, h int, w int, check bool) {
	for _, attr := range t.Attr {
		if attr.Key =="rel" && attr.Val == prop {
			check = true
		}
		if attr.Key == "href" {
			href = attr.Val
		}
		if attr.Key == "type" {
			typ = attr.Val
		}
		if attr.Key == "sizes" {
			sizes := strings.Split(attr.Val, "x")
			intW, _ := strconv.Atoi(sizes[0])
			w = intW
			if len(sizes) > 1 {
				intH, _ := strconv.Atoi(sizes[1])
				h = intH
			}
		}
	}
	return
}

func absoluteURL (u string) (abs string) {
	relative, _ := url.Parse(u)
	base, _ := url.Parse("http://test.com")
	abs = base.ResolveReference(relative).String()
	return 
}