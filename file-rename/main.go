package main

import (
	"flag"
	"fmt"
	"os"
	// "path/filepath"
	"io/ioutil"
	"regexp"
	"strings"
	"strconv"
)

const (
	currentNum        = `##NN`
	totalNum          = `##TT`
	currentNumRegex   = `(?P<currentNum>\d+)`
	totalNumRegex     = `(?P<totalNum>\d+)`
	currentNumReplace = `${currentNum}`
	totalNumReplace   = `${totalNum}`
)

func main() {
	var from, to string
	var verbose bool

	flag.StringVar(&from, "from", "", "regex `pattern` to match files to rename")
	flag.StringVar(&to, "to", "", "non-regex `pattern` to rename each file to")
	flag.BoolVar(&verbose, "v", false, "verbose mode")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Renames files in current directory only as per pattern\nUsage:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "In patterns,\n\tuse", currentNum, "to denote current counter value in rename\n\tuse", totalNum, "to denote total count of files to be renamed")
		fmt.Fprintln(os.Stderr, "\t`from` and `to` should be different\n\t`currentNum` and `totalNum` named regex groups should not be present in patterns")
	}
	flag.Parse()

	if from == "" || to == "" || from == to {
		flag.Usage()
		return
	}

	fromPattern := strings.Replace(from, currentNum, currentNumRegex, 1)
	fromPattern = strings.Replace(fromPattern, totalNum, totalNumRegex, 1)
	fromPattern = "^" + fromPattern + "$"
	fmt.Println("fromPattern", fromPattern)
	if strings.Contains(fromPattern, currentNum) || strings.Contains(fromPattern, totalNum) {
		fatalError("Invaild `from` pattern: Maximum one", currentNum, "and", totalNum, "are allowed in `from` pattern")
	}

	fromRegex, err := regexp.Compile(fromPattern)
	if err != nil {
		fatalError("from pattern:", err.Error())
	}

	var matchedFiles int
	// matchedFiles = CountMatchingWalk(fromRegex, verbose)
	///*
	files, err := ioutil.ReadDir(".")
	if err != nil {
		fatalError(err.Error())
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if fromRegex.MatchString(file.Name()) {
			if verbose {
				fmt.Println("matched", file.Name())
			}
			matchedFiles++
		}
	}
	//*/
	var proceed rune
	fmt.Println("Found", matchedFiles, "files to rename.")
	if matchedFiles == 0 {
		return
	}
	fmt.Print("Proceed? (y/n)\t")
	fmt.Scanf("%c", &proceed)
	if proceed != 'y' && proceed != 'Y' {
		return
	}
	fmt.Println("Renaming...")
	var renamedFiles int
	// renamedFiles = RenameMatchingWalk(fromRegex, to, matchedFiles)
	///*
	newName, err := renamer(fromRegex, to, matchedFiles)
	if err != nil {
		fatalError(err.Error())
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := newName(file.Name())
		if file.Name() != name {
			os.Rename(file.Name(), name)
			renamedFiles++
		}
	}
	//*/
	fmt.Println("Done")
	if matchedFiles != renamedFiles {
		fatalError("Unexpected mismatch: Matched files-", matchedFiles, "Renamed files-", renamedFiles)
	}
}

/*
func CountMatchingWalk(fromPattern *regexp.Regexp, verbose bool) int {
	var filesMatchedCount int

	err := filepath.Walk(".", func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if info.IsDir() && info.Name() != "." {
			fmt.Println("skipping dir", info.Name())
			return filepath.SkipDir
		}

		// fmt.Println("matching", info.Name())
		if fromPattern.MatchString(info.Name()) {
			if verbose {
				fmt.Println("matched", file.Name())
			}
			filesMatchedCount++
		}
		return nil
	})

	if err != nil {
		fatalError(err.Error())
	}

	return filesMatchedCount
}

func RenameMatchingWalk(fromRegex *regexp.Regexp, to string, totalCount int) int {
	var filesRenameCount int

	newName, err := renamer(fromRegex, to, totalCount)

	if err != nil {
		fatalError(err.Error())
	}

	err = filepath.Walk(".", func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if info.IsDir() && info.Name() != "." {
			// fmt.Println("skipping dir", info.Name())
			return filepath.SkipDir
		}

		name := newName(info.Name())
		if name != info.Name() {
			err = os.Rename(info.Name(), name)
			if err != nil {
				fatalError(err.Error())
			}
			filesRenameCount++
		}
		return nil
	})

	if err != nil {
		fatalError(err.Error())
	}

	return filesRenameCount
}
*/

func renamer(fromRegex *regexp.Regexp, to string, totalCount int) (func(string) string, error) {	

	toStr := strings.Replace(to, currentNum, currentNumReplace, 1)
	if strings.Contains(fromRegex.String(), totalNumRegex) {
		toStr = strings.Replace(toStr, totalNum, totalNumReplace, 1)
	} else {
		toStr = strings.Replace(toStr, totalNum, strconv.Itoa(totalCount), 1)
	}

	return func(oldName string) string {
		return fromRegex.ReplaceAllString(oldName, toStr)
	}, nil
}

func fatalError(s ...interface{}) {
	fmt.Fprintln(os.Stderr, s...)
	os.Exit(1)
}

