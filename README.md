Slideshare downloader
===============

A simple Golang program which is able to download private and public slideshare presentation.

Following a simple hack explained in this presentation: http://www.slideshare.net/MalteLandwehr/slideshare-download, I've implemented this simple program in order to let my dear friends to download all the presentation that they want.

I know, I'm a very good man :)

**Dependecies**

Obviously this programs needs the Golang compiler available on the official website for every operating system.
In order to convert SWF file into PNG image I've used a simple tool knows as *swftools* which contains a very useful utility that is able to make automatically the conversion (http://www.swftools.org/).

If you are a developer, you need to install *gofpdf* (https://godoc.org/code.google.com/p/gofpdf) and *go.net/html* (https://godoc.org/code.google.com/p/go.net/html).

**TODO**

It is a really simple program but works pretty well. 
There is something not well implemented (I'm a little gopher) and I think that I'll update my code in order to improve it.
If you find something wrong or something that could be improved feel free to make a pull request or put it in the issue section. 
