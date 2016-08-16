package main

import (
	"fmt"
	"os"
	"path/filepath"

	"talend.com/npmblame"
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
