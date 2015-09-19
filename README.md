Slideshare downloader
===============

A simple Golang program which is able to download private and public slideshare presentation.

Following a simple hack explained in this presentation: http://www.slideshare.net/MalteLandwehr/slideshare-download, I've implemented this simple program
in order to let my dear friends to download all the presentation that they want.

I know, I'm a very good man :)


**Dependecies**

Obviously this programs needs the Golang compiler available on the official website for every operating system.

If you are a developer, you need to install *gofpdf* (https://godoc.org/code.google.com/p/gofpdf) and *goquery* (https://github.com/PuerkitoBio/goquery).


**HOW TO**
First of all you need to build the source code using the standard `go build`

After that you need to run the file **main.go** using `go run` specifing the suggested parameters specified in the USAGE; here there are:

-dir string
	Directory in which all the files will be saved
-pdf_path string
    Path of the generated PDF (default "presentation.pdf")
-url string
    Insert slide presentation URL

**TODO**

It is a really simple program but works pretty well. 
There is something not well implemented (I'm a little gopher) and I think that I'll update my code in order to improve it.
If you find something wrong or something that could be improved feel free to make a pull request or put it in the issue section. 
