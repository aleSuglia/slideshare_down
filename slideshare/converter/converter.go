// Defines specific functions and types needed
// to make conversion operations
package converter

import (
	"code.google.com/p/gofpdf"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// Represents a way to generate a pdf
type PdfGenerator struct {
	pdfObj *gofpdf.Fpdf
}

// Constructs a pdf using the appropriate
// values for presentation
func NewPdf() *PdfGenerator {
	return &PdfGenerator{gofpdf.New("L", "mm", "A4", "")}
}

// Adds a PNG image whose name is contained in the imageName parameter,
// to the current pdf
func (pdfGen *PdfGenerator) AddPngImage(imageName string) {
	pdfGen.pdfObj.Image(imageName, 0, 0, 0, 0, true, "PNG", 0, "")
}

// Closes the current pdf and serialize it in a file whose
// name is specified as a parameter
func (pdfGen *PdfGenerator) ClosePdf(pdfName string) error {
	return pdfGen.pdfObj.OutputFileAndClose(pdfName)

}

// Saves the SWF image located at the specified url in a
// specified directory with the same name
// An error is returned if something was wrong during I/O operations
func SaveSWFPage(swfUrl, swfDir string) (string, error) {
	pageReq, errRequest := http.Get(swfUrl)

	if errRequest != nil {
		return "", errRequest
	}

	defer pageReq.Body.Close()

	swfImage, errImage := ioutil.ReadAll(pageReq.Body)

	if errImage != nil {
		return "", errImage
	}

	slideIndex := strings.Index(swfUrl, "-slide")

	slideFileName := swfDir + string(os.PathSeparator) + swfUrl[slideIndex+1:]

	ioutil.WriteFile(slideFileName, swfImage, os.ModePerm)

	return slideFileName, nil

}

// Converts the specified SWF file into a PNG file (with the same name)
// using the swfrender tool (http://www.swftools.org/)
func ExecuteConversion(slideFileName string) (string, error) {
	replacer := strings.NewReplacer(".swf", ".png")
	imageName := replacer.Replace(slideFileName)

	swfCommand := exec.Command("swfrender", slideFileName, "-o", imageName)

	return imageName, swfCommand.Run()
}
