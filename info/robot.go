// package which contains all the function used
// in order to obtain the presentation's information
package info

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

// Used when a specific HTML tag isn't found
type TagNotFoundError struct {
	errorMsg string
}

func (e TagNotFoundError) Error() string {
	return e.errorMsg
}

func GetSlideList(presentationURL string) ([]string, error) {
	//var slideContSelector string = ".slide_container"
	var slideImgSelector string = ".slide_image"
	var imgURLAttribute string = "data-full"

	doc, err := goquery.NewDocument(presentationURL)

	if err != nil {
		return nil, err
	}

	// allocate for a single slide, than extend it
	// for each slide that we find in the HTML page
	slideList := make([]string, 0, 0)

	// find the slide container in the web page
	// for each section in it, retrieve the img tag that contains the images' URL
	doc.Find(slideImgSelector).Each(func(i int, s *goquery.Selection) {
		// each children of the slide container is a section
		// each sections' children is an "img" tag
		if url, ok := s.Attr(imgURLAttribute); ok {
			slideList = append(slideList, url)
		}
	})

	if len(slideList) == 0 {
		return nil, TagNotFoundError{"No slide sections in the HTML page!"}
	}

	return slideList, nil

}

func DownloadSlideImage(slideURL, slideDir string) (string, error) {
	pageReq, errRequest := http.Get(slideURL)

	if errRequest != nil {
		return "", errRequest
	}

	defer pageReq.Body.Close()

	image, errImage := ioutil.ReadAll(pageReq.Body)

	if errImage != nil {
		return "", errImage
	}

	re := regexp.MustCompile(`http:\/\/image\.slidesharecdn\.com\/*.*\/.*\/(\D+)-([0-9]+)-.*`)
	matched := re.FindStringSubmatch(slideURL)
	slideFileName := slideDir + string(os.PathSeparator) + matched[1] + "-" + matched[2] + ".jpg"

	ioutil.WriteFile(slideFileName, image, os.ModePerm)

	return slideFileName, nil
}
