package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type AccTool struct {
	idsToName map[string]string
	nameToIDs map[string]string
	print     printer

	strictMode bool
}

const accountIDRegexp = "\\b\\d{12}\\b"
const onlyDigitsRegexp = "^\\d+$"
const unknownItem = "unknown"

func NewTool(c *Config, noColor, verbose, strictMode bool) *AccTool {
	return &AccTool{
		idsToName:  MapIDsToName(c),
		nameToIDs:  MapNameToIDs(c),
		print:      newPrinter(noColor, verbose),
		strictMode: strictMode,
	}
}

func (a *AccTool) ReplaceAccountIDsFromFiles(filepaths []string) error {
	for _, path := range filepaths {
		path, _ := filepath.Abs(path)
		err := a.ReplaceAccountIDsFromFile(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AccTool) ReplaceAccountIDsFromFilesGlob(regexp string) error {
	matches, _ := filepath.Glob(regexp)
	for _, match := range matches {
		path, _ := filepath.Abs(match)
		err := a.ReplaceAccountIDsFromFile(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AccTool) ReplaceAccountIDsFromFile(filepath string) error {
	f, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return ErrUnableLoadFile
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		a.replaceWithAccountIDs(sc.Text())
	}
	if err := sc.Err(); err != nil {
		return ErrUnableReadFile
	}
	return nil
}

func (a *AccTool) ReplaceAccountIDsFromStdin() error {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		a.replaceWithAccountIDs(sc.Text())
	}
	if err := sc.Err(); err != nil {
		return ErrUnableReadFile
	}
	return nil
}

func (a *AccTool) replaceWithAccountIDs(currentLine string) {
	if !hasAccountID(currentLine) {
		fmt.Println(currentLine)
		return
	}

	containsUnknownItems := false
	newLine := currentLine

	accountIDs := findAllAccountIDs(currentLine)
	for _, accountID := range accountIDs {
		found := a.lookupNameForID(accountID)
		if found == unknownItem {
			containsUnknownItems = true
		}
		newLine = strings.Replace(newLine, accountID, fmt.Sprintf("%s", found), -1)
	}

	if containsUnknownItems && a.strictMode {
		a.print.Diff(currentLine, newLine)
		os.Exit(1)
	}

	a.print.Diff(currentLine, newLine)
}

func (a *AccTool) SearchAll(patterns []string) {
	hasMatches := false
	for _, pattern := range patterns {
		found := a.Search(pattern)
		if hasMatches != true && found {
			hasMatches = true
		}
	}

	if !hasMatches && a.strictMode {
		os.Exit(1)
	}
}

func (a *AccTool) Search(pattern string) bool {
	hasMatches := false

	if isOnlyDigits(pattern) {
		for k, v := range a.idsToName {
			match, _ := regexp.MatchString(toRegexp(pattern), k)
			if match {
				a.print.Success(k, v)
				hasMatches = true
			}
		}
		return hasMatches
	}

	for k, v := range a.nameToIDs {
		match, _ := regexp.MatchString(toRegexp(pattern), k)
		if match {
			a.print.Success(k, v)
			hasMatches = true
		}
	}
	return hasMatches
}

func toRegexp(pattern string) string {
	var result strings.Builder
	for i, literal := range strings.Split(pattern, "*") {
		if i > 0 {
			result.WriteString(".*")
		}
		result.WriteString(regexp.QuoteMeta(literal))
	}
	return result.String()
}

func (a *AccTool) lookupNameForID(id string) string {
	found := a.idsToName[id]
	if found == "" {
		return unknownItem
	}
	return found
}

func (a *AccTool) lookupIDForName(name string) string {
	found := a.nameToIDs[name]
	if found == "" {
		return unknownItem
	}
	return found
}

func hasAccountID(line string) bool {
	match, _ := regexp.MatchString(accountIDRegexp, line)
	return match
}

func isOnlyDigits(line string) bool {
	match, _ := regexp.MatchString(onlyDigitsRegexp, line)
	return match
}

func findAllAccountIDs(line string) []string {
	r, _ := regexp.Compile(accountIDRegexp)
	matches := r.FindAllString(line, -1)
	return matches
}

type printer struct {
	verbose bool
}

func newPrinter(noColor bool, verbose bool) printer {
	color.NoColor = !noColor
	return printer{
		verbose: verbose,
	}
}

func (p *printer) Success(present, found string) {
	if p.verbose {
		color.Green("%s (%s)", present, found)
		return
	}
	color.Green("%s", found)
}

func (p *printer) Fail(present, found string) {
	if p.verbose {
		color.Red("%s (%s)", present, found)
		return
	}
	color.Red("%s", found)
}

func (p *printer) Diff(present, found string) {
	if p.verbose {
		color.Green("%s", found)
		color.Red("%s", present)
		return
	}
	color.Green("%s", found)
}
