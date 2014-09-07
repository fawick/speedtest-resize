speedtest-resize
================

Compare various Image resize algorithms for the Go language

I am writing a web gallery called gonagall in Go
(https://github.com/fawick/gonagall). For that, I need a efficient solution for
scaling and resizing a lot of images (mostly JPGs) to generate thumbnails and
bandwidth-friendly sized copies from high-resolution original photo files.

In this project I compare the speed of a few selected image resizing algorithms
with each other as well as with ImageMagick and GraphicsMagick. The competitors
are 

- https://github.com/nfnt/resize, Pure golang image resizing, more precisely
  only Nearest-Neighbor interpolation that comes with that Go package.
- https://github.com/disintegration/gift Again, I use one of the fastest
  algorithms of the package. Here, it's called 'Box'
- https://github.com/disintegration/imaging Again, I use one of the fastest
  algorithms of the package. Here, it's called 'Box'
- [ImageMagick convert](http://www.imagemagick.org/script/convert.php) with the options `-resize 150x150>`
- [ImageMagick convert](http://www.imagemagick.org/script/convert.php) with the
  options `-define "jpeg:size=300x300 -thumbnail 150x150>`. `-thumbnail` is
considered to be faster than resize, and the `-define` will reduce the size (in
terms of memory footprint) of the original image on reading.
- [GraphicsMagick convert](http://www.graphicsmagick.org/convert.html) with the
  options `-define "jpeg:size=300x300 -thumbnail 150x150>`.
- https://github.com/gographics/imagick Go wrapper for the MagickWand API,
  again the Box algorithm is used for the sake of comparing the results.
- https://github.com/lazywei/go-opencv Go binding for OpenCV, using the fastest
  algorithm.
- https://github.com/bamiauz/rez, pure go resizer, using bilinear interpolation
  in these tests
- https://github.com/DAddYE/vips, bindings for libvips
  (https://github.com/jcupitt/libvips)

Some implementations are no longer available and have been removed:

- [resize.go](https://code.google.com/p/appengine-go/source/browse/example/moustachio/resize/resize.go)
  by the Go authors, which performs a Nearest-Neighbor scaling

### Installation

To run the tests `go get` the source and compile/run it:

    $ go get -u github.com/fawick/speedtest-resize
    $ cd $GOPATH/src/speedtest-resize
    $ go run main.go <jpg file folder>

Alternatively, call the go command (or the compiled binary) from the image folder without 
supplying a parameter

    $ cd <jpg file folder>
    $ go run $GOPATH/src/speedtest-resize/main.go

### Benchmark

Im my test scenario all of these tools/packages are unleashed on a directory
containing JPG photo files, all of which have a resolution of 5616x3744 pixels
(aspect ratio 2:1, both landscape and portrait). 

For each tool/package and for all files, the total time for loading the
original file, scaling the image to a thumbnail of 150x100 pixels, and writing
it to a new JPG file is measured. In the end, the total runtime for processing
the 10 first files and the average time per file is printed for each
tool/package.

The scenario is run on a Intel(R) Pentium(R) Dual T2390 @ 1.86GHz running
Ubuntu 14.04. Here are the results:

| Tables                            | Time (file avg.)  | Pure Go | 
| --------------------------------- | -----------------:| ------- |
| ImageMagick -thumbnail            | 0.325s            |         |
| VIPS                              | 0.354s            |         |
| GraphicsMagick -thumbnail         | 0.752s            |         |
| github.com/gographics/imagick     | 1.100s            |         |  
| ImageMagick -resize               | 2.286s            |         |
| rez_bilinear                      | 2.942s            |    X    |
| github.com/nfnt/resize            | 3.509s            |    X    |
| github.com/disintegration/imaging | 4.694s            |    X    |

--------

Yet another scenario ran by [lazywei](https://github.com/lazywei), 2.5GHz Intel Core i5, Mac OS X 10.9.1:

| Tables               | Average time per file  |
| -------------------- | ----------------------:|
| magickwand_box       |  155.371531ms          |
| imaging_Box          |  463.459339ms          |
| moustaschio_resize   |  719.051885ms          |
| Nfnt_NearestNeighbor |  1.436507946s          |
| OpenCv               |   97.353041ms          |

--------

Yet another scenario ran by [bamiaux](https://github.com/bamiaux), 3.3GHz Intel Core i5, win 7:

| Tables               | Average time per file  |
| -------------------- | ----------------------:|
| rez_bilinear         |  148ms                 |
| imaging_Box          |  243ms                 |
| moustaschio_resize   |  407ms                 |
| Nfnt_NearestNeighbor |  233ms                 |

--------

So, what is to learn from that? While all of the currently existing
pure-Go-language solutions do a pretty good job in generating good-looking
thumbnails, they are much slower than the veteran dedicated image processing
toolboxes. That is hardly surprising, given that both ImageMagick and
GraphicsMagick have been around for decades and have been optimized to work as
efficient as possible. Go and its image processing packages are still the new
kids on the block, and while they work pretty neat for the occasional tweak of
an image or two, I rather not use them as the default image processor in
[gonagall](http://github.com/fawick/gonagall) yet. 

I was surprised to find that GraphicsMagick was slower than ImageMagick in my
test scenario, as I expected it to be exactly the other way around with
GraphicsMagick's fancy multi-processor algorithms.

While the imagick Wrapper is written in Go, it uses CGO bindings of the C
MagickWand API. It outperforms the pure-Go approaches (five times faster than
http://github.com/disintegration/imaging) but it still slower than calling
ImageMagick in an external process. Of the above 1.13 seconds, only around 275
millisecs were used for resizing and saving an individual file, while over 850
ms were used by simply loading the file. I wonder how much optimization can
still be done in the imagick loading routines.


