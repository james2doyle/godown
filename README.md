godown
======

A Markdown parser written in Go

### Usage

    godown test.md

This will compile a new file called `test.html`.

    godown --stdout=true test.md

This will print the compiled HTML to the stdout so you can use it as you wish.

### Supported Features

* Paragraphs (new lines)
* Headers
* Links (no title attribute)
* Images
* Bold
* Italic
* Deletion with the \~\~
* Blockquotes

### Working on

* Unordered Lists
* Ordered Lists
* Links with Title support
* Code Blocks (with classes and indents)
* CLI Flags for in/out and options

### Testing

There is a `test.md` file in the project. It shows all the supported features.

#### Benchmarks

* [sundown](https://github.com/vmg/sundown) test.md  0.00s user 0.00s system 45% cpu 0.007 total
* godown test.md  0.01s user 0.00s system 85% cpu 0.013 total
* [marked](https://github.com/chjj/marked) test.md  0.05s user 0.02s system 91% cpu 0.077 total
