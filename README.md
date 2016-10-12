mosaique
========

Simple Go tool for creating mosaique pictures from a set of photos.

![mosaique-sample](https://github.com/chmllr/mosaique/raw/master/mosaique.jpg)

How to use
----------

### Compile

First compile the tools runing `make` in the root directory.

### Create Mosaiques

1. Run fetcher on your folder with fotos to be used as mosaique tiles:

       ./bin/fetcher /path/to/photos

    This will create the file `color.txt` in the caller's directory containing color information for the mosaique.
2. Run generator specifying the color file and the photo to be turned into mosaique:

       ./bin/generator color.txt /path/tp/photo.jpg

    This will create a file `/path/tp/photo.jpg_mosaique.jpg` in the directory of the original photo.

Restrictions
------------

- Only JPGs are supported.
- Tile size is currently 28x28 pixels only.

License
-------

MIT.