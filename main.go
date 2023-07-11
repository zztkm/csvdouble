package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const version = "0.0.1"

var revision = "HEAD"

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(1)
}

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func writeWrapper(w *csv.Writer, record []string) error {
	for i, field := range record {
		record[i] = `"` + field + `"`
	}
	return w.Write(record)
}

func main() {
	var showVersion bool
	var showHelp bool

	// options
	fs := flag.NewFlagSet("csvdouble", flag.ExitOnError)
	fs.BoolVar(&showVersion, "version", false, "show version")
	fs.BoolVar(&showVersion, "v", false, "show version")
	fs.BoolVar(&showHelp, "help", false, "show help")
	fs.BoolVar(&showHelp, "h", false, "show help")
	fs.Usage = func() {
		fmt.Println(`Usage:
  csvdouble <path>

Flags:`)
		fs.PrintDefaults()
		fmt.Println(`Repository:
  https://github.com/zztkm/csvdouble`)
	}

	fs.Parse(os.Args[1:])

	if showVersion {
		fmt.Printf("version: %s, revision: %s\n", version, revision)
		return
	}
	if showHelp {
		fs.Usage()
		return
	}

	if fs.NArg() < 1 {
		fatal(errors.New("please specify csv file path. `excel2csv -h` for more details"))
	}

	filename := fs.Arg(0)
	// 読み込みファイル
	rf, err := os.Open(filename)
	if err != nil {
		fatal(err)
	}
	defer rf.Close()
	r := csv.NewReader(rf)

	// 書き込みファイル
	wf, err := os.Create(getFileNameWithoutExt(filename) + "-double.csv")
	if err != nil {
		fatal(err)
	}
	defer wf.Close()
	w := csv.NewWriter(wf)
	defer w.Flush()

	for {
		records, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fatal(err)
		}
		writeWrapper(w, records)
	}
}
