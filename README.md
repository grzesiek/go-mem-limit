# go-mem-limit

Simple and naive library making it possible to limit the amount of dynamically
allocated memory on heap in Go applications.

When memory allocation in a given moment reaches a defined threshold a provided
callback function is going to be called.

This samples allocated memory for entire process.
