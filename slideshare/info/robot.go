// package which contains all the function used
// in order to obtain the presentation's information
package info

import (
	"code.google.com/p/go.net/html"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// Defines the structure of the XML page
// that is used by the Slideshare website to contain
// all the presentation's pages
type xmlShowNode struct {
	Show  string `xml:"Id,attr"`
	Slide []xmlSlideNode
}

// A single page element in the XML file
type xmlSlideNode struct {
	Slide string `xml:"Src,attr"`
}

// Used when a specific HTML tag isn't found
type TagNotFoundError struct {
	errorMsg string
}

func (e TagNotFoundError) Error() string {
	return e.errorMsg
}

// Transform the retrieved URL in order to get
// another URL which contains a list of reference
// to SWF file that represents the presentation's slides
func TransformImageURL(contentURL string) string {
	// first step: Remove thumbnail_ss
	replacer := strings.NewReplacer("ss_thumbnails/", "")
	noThumbString := replacer.Replace(contentURL)

	// second_step: Remove the text which starts with "-thumbnail"
	lastPartIndex := strings.Index(noThumbString, "-thumbnail")

	trimmedStr := noThumbString[:lastPartIndex]

	// append ".xml"
	return trimmedStr + ".xml"

}

// Retrieves the correct URL from the webpage's HTML code
// in order to start retrieving each slide from the original presentation
func GetSlideImageURL(slideURL string) (string, error) {
	var metaStr, propKey, propValue, nameKey, nameValue, thumbValue string = "meta", "property", "og:image", "name", "og_image", "thumbnail"

	resp, err := http.Get(slideURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	doc, err := html.Parse(io.Reader(resp.Body))
	if err != nil {
		return "", err
	}

	slideRealURL := ""
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == metaStr {
			attributeMap := mapAttributeList(n.Attr)
			if attributeMap[propKey] == propValue && attributeMap[nameKey] == nameValue {
				slideRealURL = attributeMap["content"]
			} else if attributeMap[nameKey] == thumbValue {
				slideRealURL = attributeMap["content"]
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if slideRealURL == "" {
		return slideRealURL, TagNotFoundError{"Unable to find 'meta content' tag in HTML page"}
	}

	return slideRealURL, nil

}

// Returns a map containing a pair <property, property_value>
// obtained from the HTML code of the webpage
func mapAttributeList(attributeList []html.Attribute) (attributeMap map[string]string) {
	attributeMap = make(map[string]string)

	for _, attribute := range attributeList {
		attributeMap[attribute.Key] = attribute.Val
	}

	return
}

// Returns the list of all the URL of the slide
// parsing the XML file which contains all of them
func GetSlideList(slideListPage string) ([]string, error) {
	resp, errRequest := http.Get(slideListPage)

	if errRequest != nil {
		return nil, errRequest
	}

	defer resp.Body.Close()

	page, errPage := ioutil.ReadAll(resp.Body)

	if errPage != nil {
		return nil, errPage
	}

	var resData xmlShowNode

	errXml := xml.Unmarshal(page, &resData)

	if errXml != nil {
		return nil, errXml
	}

	var slideList []string

	for _, value := range resData.Slide {
		slideList = append(slideList, value.Slide)
	}

	return slideList, nil

}
