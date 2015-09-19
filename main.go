package main

import (
	"flag"
	"fmt"
	"github.com/aleSuglia/slideshare_down/converter"
	"github.com/aleSuglia/slideshare_down/info"
	"os"
	"regexp"
	"sort"
	"strconv"
	"sync"
)

var (
	slideOriginalURL string
	slideDir         string
	pdfPath          string
)

type slideSorter []string

func (a slideSorter) Len() int      { return len(a) }
func (a slideSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Sorts the slice according to the slide number contained in each string
// Each slide name is in the form "slide-$num_slide.jpeg"
func (a slideSorter) Less(i, j int) bool {
	r := regexp.MustCompile(`\D+-([0-9]+).jpg`)
	firstNum, _ := strconv.ParseInt(r.FindStringSubmatch(a[i])[1], 10, 64)
	secNum, _ := strconv.ParseInt(r.FindStringSubmatch(a[j])[1], 10, 64)

	return firstNum < secNum

}

func init() {
	flag.StringVar(&slideOriginalURL, "url", "", "Insert slide presentation URL")
	flag.StringVar(&slideDir, "dir", "", "Directory in which all the files will be saved")
	flag.StringVar(&pdfPath, "pdf_path", "presentation.pdf", "Path of the generated PDF")
}

func downloadService(slideList []string) (*chan string, *sync.WaitGroup) {
	computedSlides := make(chan string)
	var wg sync.WaitGroup

	wg.Add(len(slideList))
	for _, slide := range slideList {
		go func(slideUrl string) {
			defer wg.Done()
			fmt.Println("Downloading: " + slideUrl)

			slideFileName, err := info.DownloadSlideImage(slideUrl, slideDir)

			if err != nil {
				return
			}

			// Sending result to the channel
			fmt.Println("Completed slide: ", slideUrl)
			computedSlides <- slideFileName

		}(slide)
	}

	return &computedSlides, &wg
}

func main() {

	flag.Parse()
	if !flag.Parsed() || flag.NFlag() != 3 {
		flag.Usage()
		os.Exit(-1)
	}

	fmt.Println("Started reading page: " + slideOriginalURL)

	slideList, err := info.GetSlideList(slideOriginalURL)

	if err != nil {
		fmt.Println(err)
	} else {

		gen := converter.NewPdf()

		// Starts the download service
		computedSlides, wg := downloadService(slideList)

		// A separated goroutine will wait all the other
		// in order to close the channel (see after)
		go func() {
			wg.Wait()
			close(*computedSlides)
		}()

		orderedSlides := make([]string, 0, 0)
		// The for-loop will wait until someone closes the channel
		for slide := range *computedSlides {
			orderedSlides = append(orderedSlides, slide)
		}

		// TODO: Use an appropriate data structure
		// which grants insert in order (ordered set?)
		// Sorts them
		sort.Sort(slideSorter(orderedSlides))

		for _, val := range orderedSlides {
			fmt.Println("Adding slide: ", val)
			gen.AddJPEGImage(val)
		}

		// Serialize the whole pdf on disk
		gen.ClosePdf(pdfPath)
	}
}
