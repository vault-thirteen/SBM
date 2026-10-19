# SBM

## Simple Bit Map

This package provides a model and methods to work with the SBM format. "SBM" 
is an acronym for "simple bit map". The SBM format is a format for two-level 
(monochrome, black-&-white, bi-level) raster graphical images.

The SBM format is used to store a raw uncompressed array of binary pixels 
together with basic meta-data describing the size of the pixel array. SBM 
format is stored in a mixed encoding. This means that meta-data is encoded as 
plain _ASCII_ text symbols and pixel array is encoded using the binary format.

The array of pixels is composed of pixels row by row, starting with the top 
row, and having the bottom row at the end. Each row is composed of `W` pixels, 
where `W` is the array's width. Total number of rows is `H`, where `H` is the 
array's height. The array contains `A` pixels total, where `A` is the array's 
area, the multiple of `W` and `H`. Each pixel in the array is a separate bit, 
where zero bit (`0`) is black (dark colour) and one bit (`1`) is white (light 
colour). Due to the limitations of current hardware, the order of bits in each 
byte is not controlled by this library (package). The least significant bit is 
considered to be the first bit, the most significant bit is the last bit.

The internal SBM model type stores both fields: an array of bits (as separate 
objects) and an array of bytes which are created by the concatenation of all 
the bits as real machine's bits (not as objects). This is done for convenience 
of various internal manipulations with either bits or bytes. When the object 
of the SBM model type is stored into the stream, it stores only the bytes array 
and omits the internal array of bit objects. The same happens when the SBM 
model type is read from the stream – bytes are received from the stream, not 
the bits.
