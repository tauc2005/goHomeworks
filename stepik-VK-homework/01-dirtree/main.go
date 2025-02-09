package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

const firstPrefix = "├───"
const lastPrefix = "└───"

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

func dirTree(out io.Writer, path string, printfiles bool) error {
	return readDirTree(out, path, printfiles, "")
}

func readDirTree(out io.Writer, path string, printfiles bool, levelPrefix string) error {
	items, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("Can't read path")
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name() < items[j].Name()
	})

	var cntFolder int
	for _, item := range items {
		if item.IsDir() {
			cntFolder++
			continue
		}
	}

	itemPrefix := firstPrefix
	last := false
	for i, item := range items {
		if item.IsDir() {
			cntFolder--
		}
		if i == len(items)-1 || (!printfiles && cntFolder == 0) {
			last = true
			itemPrefix = lastPrefix
		}

		if item.IsDir() {
			fmt.Fprint(out, levelPrefix, itemPrefix, item.Name(), "\n")
			path := fmt.Sprintf("%s%s%s", path, string(os.PathSeparator), item.Name())
			if last {
				readDirTree(out, path, printfiles, levelPrefix+"\t")
				continue
			}
			readDirTree(out, path, printfiles, levelPrefix+"│\t")
			continue
		}

		if !printfiles {
			continue
		}

		info, err := item.Info()
		if err != nil {
			return err
		}
		size := fmt.Sprintf("%db", info.Size())
		if info.Size() == 0 {
			size = "empty"
		}
		fmt.Fprintf(out, "%s%s%s (%s)\n", levelPrefix, itemPrefix, item.Name(), size)
	}

	return nil
}
