package main

import (
	"fmt"
	"os"
	"path/filepath"

	npmblame "github.com/talend-glorieux/npm-blame"
)

func main() {
	np := npmblame.NewNpmPackages()
	if err := filepath.Walk(".", np.Blame); err != nil {
		fmt.Println("File system traversing error.", err)
		os.Exit(-1)
	}
	np.String()
	fmt.Print(np)
}
