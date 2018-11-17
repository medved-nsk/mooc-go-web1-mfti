package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

const (
	tab = string(os.PathSeparator)
	TRI = "├───"
	CNR = "└───"
	VBR = "│	"
	SPC = "	"
)

/* Sorting entries */
type FI []os.FileInfo

func(fi FI) Swap(i, j int) {
	fi[i], fi[j] = fi[j], fi[i]
}

func(fi FI) Len() int {
	return len(fi)
}

func(fi FI) Less(i, j int) bool {
	return fi[i].Name() < fi[j].Name()
}
/* ^^^ Sorting ^^^ */

func dirTree(out io.Writer, path string, showFiles bool) error {
	baumSpazieren(out, path, showFiles, 0, "")
	return nil
}

func isLast(count, printed int) bool {
	return printed == count-1
}

func levelPrefix(count, printed int) string {
	if isLast(count, printed) { // last (or the only) level's entry
		return SPC
	} else { // count > 1:
		return VBR
	}
}

func printEntry(out io.Writer, entry string) {
	fmt.Fprintf(out, "%s\n", entry)
}

func entryAnchor(fLast bool) string {
	if fLast {
		return CNR
	} else {
		return TRI
	}
}

func printDir(out io.Writer, f os.FileInfo, pref string) {
	fmt.Fprintf(out, "%s%s\n", pref, f.Name())
}

func printFile(out io.Writer, f os.FileInfo, pref string) {
	fmt.Fprintf(out, "%s%s\n", pref, f.Name()+" ("+sizeToStr(f.Size(), "b")+")")
}

func baumSpazieren(out io.Writer, path string, showFiles bool, level int, pref string) {
	flist, err := ioutil.ReadDir(path)
	if err != nil {
		// TODO: handle error
	}

	var count int
	sort.Sort(FI(flist))
	if (showFiles) {
		count = len(flist)
	} else {
		for _, fn:= range flist {
			if fn.IsDir() { count++ }
		}
	}
	printed := 0

	for _, fn := range flist {
		if fn.IsDir() {
			printDir(out, fn, pref + entryAnchor(isLast(count, printed)))
			baumSpazieren(out, path+tab+fn.Name(), showFiles, level+1, pref + levelPrefix(count, printed))
			printed++
		} else { // ! IsDir
			if showFiles {
				printFile(out, fn, pref + entryAnchor(isLast(count, printed)))
				printed++
			}
		}
	} // for
}

func sizeToStr(sz int64, suf string) string {
	if sz > 0 {
		return strconv.FormatInt(sz, 10) + suf
	} else {
		return "empty"
	}
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
