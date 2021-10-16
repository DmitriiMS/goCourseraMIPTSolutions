//hw1_tree is the implementation of tree util in go
//it builds directory tree for given path
// if -f parameter is provided, it will also print files ant their sizes
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
)

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

//dirTree is an entry point which starts running an actual algorithm
func dirTree(out io.Writer, path string, printFiles bool) error {
	err := worker(out, path, printFiles, "")
	if err != nil { //check for errors
		return err
	}
	return nil
}

//getSortedDirList creates a list of objects in the given direcory, sorts it and returns.
//Also, if -f parameter is not provided, it prunes the list, deleting the file objects
func getSortedDirList(path string, printFiles bool) ([]os.DirEntry, error) {
	f, err := os.Open(path) //open file and check for errors
	if err != nil {
		return nil, errors.New("cant open path")
	}
	dirs, _ := f.ReadDir(0) //read contents of a dir and ignore errors because of files
	f.Close()
	if !printFiles { //if no -f parameter
		newDirs := []os.DirEntry{} //temp slice
		for _, d := range dirs {
			if d.IsDir() {
				newDirs = append(newDirs, d) //only if it's a directory, append it to temp
			}
		}
		dirs = newDirs // new slice of only dirs
	}
	//simple sort
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	return dirs, nil
}

//worker utilizes recursion to construct and print each obhect in a tree.
func worker(out io.Writer, path string, printFiles bool, prefix string) error {
	dirs, err := getSortedDirList(path, printFiles) //get a list of objects in path
	if err != nil {                                 //check for errors
		return err
	}
	for i, dir := range dirs {
		info, err := dir.Info() //get info of an object
		if err != nil {
			return err
		}
		//init some vars with empty values
		size := ""
		lines := ""
		newPrefix := ""
		if !dir.IsDir() && printFiles && info.Size() != 0 {
			size = fmt.Sprintf(" (%db)", info.Size()) // if it's not a dir and -f is present, get the size of a file
		} else if !dir.IsDir() && printFiles && info.Size() == 0 {
			size = " (empty)" // if size is 0, print (empty)
		}
		switch i { // check if it's the last element in a slice
		case len(dirs) - 1:
			lines = "└───" // it it is, make this line and prefix
			newPrefix = prefix + "	"
		default:
			lines = "├───" // if not, use this
			newPrefix = prefix + "│	"
		}
		fmt.Fprintln(out, prefix+lines+dir.Name()+size)               // print the string
		err = worker(out, path+"/"+dir.Name(), printFiles, newPrefix) // go deeper
		if err != nil {
			return err
		}
	}

	return nil
}
