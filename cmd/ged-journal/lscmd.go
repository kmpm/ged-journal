package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type LsCmd struct {
}

func (cmd *LsCmd) Run(cc *clicontext) error {
	// list all files in logpath
	// for each file, print the filename
	// if the file is a directory, print the filename followed by a colon
	// and then list the files in the directory
	err := filepath.Walk(cc.BasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			fmt.Println(path + ":")
			return nil
		}
		fmt.Println(path)
		return nil
	})
	return err
}
