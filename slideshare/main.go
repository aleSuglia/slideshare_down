package main

import (
	"fmt"
	"os"
	"slideshare/converter"
	"slideshare/info"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type SlideSorter []string

func (a SlideSorter) Len() int      { return len(a) }
func (a SlideSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Sorts the slice according to the slide number contained in each string
// Each slide name is in the form "slide-$num_slide.png"
func (a SlideSorter) Less(i, j int) bool {
	firstSlidePart, secSlidePart := a[i][strings.Index(a[i], string(os.PathSeparator)):],
		a[j][strings.Index(a[j], string(os.PathSeparator)):]

	firstNum, _ := strconv.Atoi(firstSlidePart[strings.Index(firstSlidePart, "-")+1 : strings.Index(firstSlidePart, ".")])

	secNum, _ := strconv.Atoi(secSlidePart[strings.Index(secSlidePart, "-")+1 : strings.Index(secSlidePart, ".")])

	return firstNum < secNum

}

// Downloads the slide with the specified name in the specified directory
func downloadSlide(slideName, slideDir string, buffer *[]string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	fmt.Println("Downloading: " + slideName)
	slideFileName, errConv := converter.SaveSWFPage(slideName, slideDir)

	if errConv != nil {
		fmt.Println(errConv)
	} else {
		imageName, errExec := converter.ExecuteConversion(slideFileName)

		if errExec != nil {
			fmt.Println(errExec)
		} else {
			slice := *buffer
			*buffer = append(slice, imageName)
		}
	}

}

// Needs four command line parameters:
// $slide_url: The presentation's URL on slideshare.net
// $output_dir: directory in which all the files will be saved
// $output_pdf: Name of the produced pdf
func main() {
	if len(os.Args) != 4 {
		fmt.Println("Incorrect number of parameters")
		fmt.Println("< usage >: slideshare $slide_url $output_dir $output_pdf")
		os.Exit(-1)
	}

	slideOriginalURL := os.Args[1] // original slide URL
	slideDir := os.Args[2]         // directory in which will be saved all the files
	pdfName := os.Args[3]          // pdf name

	fmt.Println("Started reading page: " + slideOriginalURL)

	slideImageURL, err := info.GetSlideImageURL(slideOriginalURL)
	if err != nil {
		fmt.Println(err)
	} else {

		fmt.Println("Processing presentation URL")
		trimmedStr := info.TransformImageURL(slideImageURL)
		//fmt.Println(trimmedStr)
		slideList, err := info.GetSlideList(trimmedStr)
		if err != nil {
			fmt.Println(err)
		} else {
			gen := converter.NewPdf()
			computedSlides := make([]string, 0, len(slideList))
			var waitGroup sync.WaitGroup

			for _, slide := range slideList {
				waitGroup.Add(1)
				go downloadSlide(slide, slideDir, &computedSlides, &waitGroup)
			}

			// Waits until all the goroutine finished
			waitGroup.Wait()

			// TODO: Use an appropriate data structure
			// which grants insert in order (ordered set?)
			// Sorts them
			sort.Sort(SlideSorter(computedSlides))

			for _, val := range computedSlides {
				gen.AddPngImage(val)
			}

			// Serialize the whole pdf on disk
			gen.ClosePdf(pdfName)
		}
	}
}
