package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

var headerRegex = "(#+)(.*)"
var blockquoteRegex = "\n(&gt;|\\>)(.*)"
var newlineRegex = "([^\n]+)"

// dont wrap tags that begin with <h <b <u <o <l
var newlineSkipRegex = "\\<(h|b|u|o|l)(.*)"

func GenerateMap() map[string]string {
	// the simple replacements
	x := map[string]string{
		// inline code
		"<code>$1</code>":               "`(.*?)`",
		// images
		"<img src=\"$2\" alt=\"$1\" />": "!\\[([^\\[]+)\\]\\(([^\\)]+)\\)",
		// link
		"<a href=\"$2\">$1</a>":         "\\[([^\\[]+)\\]\\(([^\\)]+)\\)",
		// bold
		"<strong>$2</strong>":           "(\\*\\*|__)(.*?)\\*\\*",
		// italics
		"<em>$2</em>":                   "(\\*|_)(.*?)\\*",
		// del
		"<del>$1</del>":                 "\\~\\~(.*?)\\~\\~",
	}
	return x
}

// parses the simple wrappers from GenerateMap
func SimpleParser(contents string) string {
	arr := GenerateMap() // get the map
	for key, value := range arr {
		re := regexp.MustCompile(value)
		contents = re.ReplaceAllString(contents, key)
	}
	return contents
}

func HandleHeaders(contents string) string {
	re := regexp.MustCompile(headerRegex)
	return re.ReplaceAllStringFunc(contents, func(input string) string {
		// count the number of # to find the size
		count := strings.Count(input, "#")
		// wrap the html and trim the #
		return fmt.Sprintf("<h%d>%s</h%d>", count, strings.TrimSpace(strings.Trim(input, "#")), count)
	})
}

// placeholders for lists
// func HandleUl(contents string) string {
// 	re := regexp.MustCompile("\n\\*(.*)")
// 	return re.ReplaceAllStringFunc(contents, func(input string) string {
// 		return fmt.Sprintf("\n<ul>\n\t<li>%s</li>\n</ul>", strings.TrimSpace(input))
// 	})
// }

// func HandleOl(contents string) string {
// 	re := regexp.MustCompile("\n[0-9]+\\.(.*)")
// 	return re.ReplaceAllStringFunc(contents, func(input string) string {
// 		return fmt.Sprintf("\n<ol>\n\t<li>%s</li>\n</ol>", strings.TrimSpace(input))
// 	})
// }

func HandleBlockquotes(contents string) string {
	re := regexp.MustCompile(blockquoteRegex)
	return re.ReplaceAllStringFunc(contents, func(input string) string {
		// trim out the "> " in the blockquote
		input = strings.Replace(input, "> ", "", -1)
		return fmt.Sprintf("\n<blockquote><p>%s</p></blockquote>", strings.TrimSpace(input))
	})
}

func HandleNewLines(contents string) string {
	re := regexp.MustCompile(newlineRegex)
	return re.ReplaceAllStringFunc(contents, func(input string) string {
		re2 := regexp.MustCompile(newlineSkipRegex)
		// see if this is already wrapped with some html
		index := re2.FindStringIndex(input)
		if index != nil {
			// she is wrapped! return
			return input
		} else {
			// handle this B
			return fmt.Sprintf("<p>%s</p>", input)
		}
	})
}

// the string of the file
var s = ""

func ReadFile(filename string) {
	// read whole the file
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file")
		log.Fatal(err)
	}
	// assign the contens of the file
	s = string(contents)
}

func WriteFile(filename string) {
	// get the file extension
	ext := filepath.Ext(filename)
	// we want a .html file written
	newfilename := strings.Replace(filename, ext, ".html", 1)
	// add some lines to top and bottom
	s = fmt.Sprintf("\n%s\n", s)
	err := ioutil.WriteFile(newfilename, []byte(s), 0700)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("compiled", newfilename)
	}
}

var stdFlag = flag.Bool("stdout", false, "Print the output to stdout instead of a .html file")

func main() {
	flag.Parse()
	var arg = flag.Arg(0)
	if arg != "" {
		ReadFile(arg)
		s = SimpleParser(s)
		s = HandleHeaders(s)
		s = HandleBlockquotes(s)
		// s = HandleUl(s)
		// s = HandleOl(s)
		s = HandleNewLines(s)
		// check to see if this is for the stdout
		if *stdFlag == true {
			fmt.Println(s)
		} else {
			WriteFile(arg)
		}
	} else {
		fmt.Println("file argument required")
	}
}
