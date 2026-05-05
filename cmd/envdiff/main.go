package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/envdiff/envdiff/internal/parser"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff <file1.env> <file2.env>\n\n")
		fmt.Fprintf(os.Stderr, "Compares two .env files and reports missing or mismatched keys.\n")
	}
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	fileA, fileB := flag.Arg(0), flag.Arg(1)

	envA, err := parser.ParseFile(fileA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	envB, err := parser.ParseFile(fileB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	issues := diff(fileA, envA, fileB, envB)
	if len(issues) == 0 {
		fmt.Println("No differences found.")
		return
	}

	for _, issue := range issues {
		fmt.Println(issue)
	}
	os.Exit(2)
}

// diff compares two EnvMaps and returns a list of human-readable issue strings.
func diff(nameA string, envA parser.EnvMap, nameB string, envB parser.EnvMap) []string {
	var issues []string

	for key := range envA {
		if _, ok := envB[key]; !ok {
			issues = append(issues, fmt.Sprintf("MISSING in %s: %s", nameB, key))
		}
	}

	for key := range envB {
		if _, ok := envA[key]; !ok {
			issues = append(issues, fmt.Sprintf("MISSING in %s: %s", nameA, key))
		}
	}

	return issues
}
