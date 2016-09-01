package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	npmblame "github.com/talend-glorieux/npm-blame"
)

func main() {
	var report = flag.Bool("report", false, `Report the issues to there owner
	(should always be used with the token flag)`)
	var token = flag.String("token", "", "GitHub token with public repo activated used for reporting")
	flag.Parse()

	if *report && *token == "" {
		fmt.Println("Please provide a token with public access for GitHub reporting by using the -token flag. https://help.github.com/articles/creating-an-access-token-for-command-line-use")
		os.Exit(-1)
	}

	np := npmblame.NewNpmPackages()
	if err := filepath.Walk(".", np.Blame); err != nil {
		fmt.Println("File system traversing error.", err)
		os.Exit(-1)
	}
	fmt.Print(np)

	if *report {
		fmt.Println("Do you want to report all of those issues? (Y/N)")
		var yn string
		fmt.Scanf("%s", &yn)
		if strings.ToLower(yn) == "y" || strings.ToLower(yn) == "yes" {
			fmt.Println("Reporting...")
			// TODO: Generate true report per packages
			report := npmblame.NewReport("talend-glorieux", "npm-blame", []int{42})
			issue, err := report.Send(npmblame.DefaultClient(*token))
			if err != nil {
				fmt.Println("ERROR", err)
				os.Exit(-1)
			}
			fmt.Printf("Created issue %d", issue.Number)
		}
	}
}
