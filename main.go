package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func InsertLineInTranslFile(d *bufio.Writer, s string) {
	d.WriteString(s + "\n")
	d.Flush()
}
func CheckForTargetAndInsertLineInTranslFile(d *bufio.Writer, targ string, curr string) {
	if strings.Contains(curr, "target") {
		// the line has the target tag, so we just rewrite on that line
		InsertLineInTranslFile(d, targ)
	} else {
		// the line has no target tag, insert target and current line
		InsertLineInTranslFile(d, targ+"\n"+curr)
	}
}
func FindSourceLineInTranslationFile(t *os.File, d *bufio.Writer, prewSourceLine *string, currSourceLine string) {
	t.Seek(0, 0)
	var tScanner = bufio.NewScanner(t)
	var isFoundTranslationInTransFile = false
	for tScanner.Scan() {
		var targetLine = tScanner.Text()
		if strings.Contains(targetLine, *prewSourceLine) && *prewSourceLine != "" {
			// Found target <source> in tranlation file, take next line with target
			isFoundTranslationInTransFile = true
			*prewSourceLine = ""
			continue
		}
		if isFoundTranslationInTransFile {
			CheckForTargetAndInsertLineInTranslFile(d, targetLine, currSourceLine)
			isFoundTranslationInTransFile = false
			break
		}
	}
}
func GoLineByLineInSourceFile(fi *os.File, fo *os.File, t *os.File) {
	var prewSourceLine = ""
	scanner := bufio.NewScanner(fi)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	for scanner.Scan() {
		var currSourceLine = scanner.Text()
		datawriter := bufio.NewWriter(fo)
		if strings.Contains(currSourceLine, "<source>") {
			prewSourceLine = currSourceLine
			InsertLineInTranslFile(datawriter, currSourceLine)
			continue
		}
		if prewSourceLine == "" {
			InsertLineInTranslFile(datawriter, currSourceLine)
			continue
		}
		if strings.Contains(currSourceLine, "translated") {
			// translated skip serching
			prewSourceLine = ""
			InsertLineInTranslFile(datawriter, currSourceLine)
			continue
		}
		if (strings.Contains(currSourceLine, "<target>") && strings.Contains(currSourceLine, "</target>")) &&
			currSourceLine[strings.Index(currSourceLine, ">")+1] != 60 {
			// <target>transaltion</target> - has transaltion, skip serching
			// <target></target> 60 = <, if after > is not < - skip serching
			// <target> </target> - has translation even it is empty
			prewSourceLine = ""
			InsertLineInTranslFile(datawriter, currSourceLine)
			continue
		}
		FindSourceLineInTranslationFile(t, datawriter, &prewSourceLine, currSourceLine)
	}
}

func main() {
	var toTranslateFileName string
	var transaltionFileName string
	var newTranslationFileName = "translated_rename.xlf"

	fmt.Println("Enter the name of 'source.xlf', need-translation file")
	fmt.Scan(&toTranslateFileName)
	fmt.Println("Enter the name of 'translation.xlf', translated file")
	fmt.Scan(&transaltionFileName)

	fi, err := os.Open(toTranslateFileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	fo, err := os.Create(newTranslationFileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	t, err := os.Open(transaltionFileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := t.Close(); err != nil {
			panic(err)
		}
	}()
	GoLineByLineInSourceFile(fi, fo, t)
}
