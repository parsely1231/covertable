package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

type csvExporter struct {
	profileTree profileTree
}

func makeCSVExporter(profileTree profileTree) csvExporter {
	return csvExporter{
		profileTree: profileTree,
	}
}

func (v csvExporter) Save(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	w := makeWriter(f)

	err = w.Write(csvHeader())
	if err != nil {
		return err
	}

	for _, node := range v.profileTree.Walk() {
		record, err := csvRecordFromNode(node)
		if err != nil {
			return err
		}

		err = w.Write(record)
		if err != nil {
			return errors.Errorf("failed to write csv record %w", err)
		}
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

type writer struct {
	innerWriter *csv.Writer
}

func makeWriter(w io.Writer) writer {
	return writer{
		innerWriter: csv.NewWriter(w),
	}
}

func (w writer) Write(record csvRecord) error {
	stringRecord := []string{
		record.FilePath,
		record.Coverage,
		record.TotalStatements,
		record.TestedStatements,
		record.NotTestedStatements,
	}
	err := w.innerWriter.Write(stringRecord)
	return err
}

func (w writer) Flush() error {
	w.innerWriter.Flush()
	return w.innerWriter.Error()
}

type csvRecord struct {
	FilePath            string
	Coverage            string
	TotalStatements     string
	TestedStatements    string
	NotTestedStatements string
}

func csvRecordFromNode(node profileNode) (csvRecord, error) {
	var coverage float64
	if node.TotalStatements() == 0 {
		coverage = 0
	} else {
		coverage = float64(node.TestedStatements()) / float64(node.TotalStatements()) * 100

	}
	notTested := node.TotalStatements() - node.TestedStatements()

	return csvRecord{
		FilePath:            node.FilePath(),
		Coverage:            strconv.FormatFloat(coverage, 'f', 2, 64),
		TotalStatements:     strconv.Itoa(node.TotalStatements()),
		TestedStatements:    strconv.Itoa(node.TestedStatements()),
		NotTestedStatements: strconv.Itoa(notTested),
	}, nil
}

func csvHeader() csvRecord {
	return csvRecord{
		FilePath:            "filepath",
		Coverage:            "coverage",
		TotalStatements:     "total statements",
		TestedStatements:    "tested",
		NotTestedStatements: "not tested",
	}
}
