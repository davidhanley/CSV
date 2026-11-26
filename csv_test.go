package csv

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	if err := os.WriteFile(path, []byte("name,age\n\"Hanley, David\",52\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	records, err := ReadFile(path, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatalf("unexpected record count: %d", len(records))
	}
	if records[1][0] != "Hanley, David" || records[1][1] != "52" {
		t.Fatalf("unexpected quoted row: %#v", records[1])
	}
}

func TestReadLines(t *testing.T) {
	records, err := ReadLines([]string{
		"race,points",
		"\"2022 T2T TAMPA\",150",
		"\"Line, With, Commas\",200",
	}, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 3 {
		t.Fatalf("unexpected record count: %d", len(records))
	}
	if records[2][0] != "Line, With, Commas" || records[2][1] != "200" {
		t.Fatalf("unexpected parsed row: %#v", records[2])
	}
}

func TestReadAllToleratesBareQuotes(t *testing.T) {
	// Real-world race files contain bare double-quotes in unquoted fields.
	records, err := ReadLines([]string{
		"place,name,age,gender",
		",Cody O\"connell,13,M",
		"Big \"D\" 2017",
	}, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 3 {
		t.Fatalf("unexpected record count: %d", len(records))
	}
	if records[1][1] != "Cody O\"connell" {
		t.Fatalf("unexpected bare-quote field: %#v", records[1])
	}
	if records[2][0] != "Big \"D\" 2017" {
		t.Fatalf("unexpected bare-quote title: %#v", records[2])
	}
}

func TestReadAllFromReader(t *testing.T) {
	records, err := ReadAll(strings.NewReader("a,b\n\"c,d\",e\n"), DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatalf("unexpected record count: %d", len(records))
	}
	if records[1][0] != "c,d" || records[1][1] != "e" {
		t.Fatalf("unexpected parsed row: %#v", records[1])
	}
}

func TestReadAllSkipsBlankLinesByDefault(t *testing.T) {
	records, err := ReadString("name\n\n50\n1,bob\n", DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	want := [][]string{{"name"}, {"50"}, {"1", "bob"}}
	if len(records) != len(want) {
		t.Fatalf("record count: got %#v", records)
	}
	for i := range want {
		if strings.Join(records[i], "|") != strings.Join(want[i], "|") {
			t.Fatalf("record %d: got %#v want %#v", i, records[i], want[i])
		}
	}
}

func TestReadAllPreservesBlankLines(t *testing.T) {
	opt := DefaultOptions()
	opt.KeepBlankLines = true
	records, err := ReadString("name\n\n50\n1,bob\n", opt)
	if err != nil {
		t.Fatal(err)
	}
	want := [][]string{{"name"}, {""}, {"50"}, {"1", "bob"}}
	if len(records) != len(want) {
		t.Fatalf("record count: got %#v", records)
	}
	for i := range want {
		if strings.Join(records[i], "|") != strings.Join(want[i], "|") {
			t.Fatalf("record %d: got %#v want %#v", i, records[i], want[i])
		}
	}
}

func TestReadAllPreservesCRLFBlankLines(t *testing.T) {
	opt := DefaultOptions()
	opt.KeepBlankLines = true
	records, err := ReadString("date\r\n\r\n50\r\n", opt)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 3 || records[1][0] != "" || records[2][0] != "50" {
		t.Fatalf("unexpected records: %#v", records)
	}
}

func TestReadAllKeepsBlankInsideQuotedField(t *testing.T) {
	records, err := ReadString("\"hello\n\nworld\"\n", DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0][0] != "hello\n\nworld" {
		t.Fatalf("unexpected quoted multiline: %#v", records)
	}
}

func TestFieldsPerRecordOption(t *testing.T) {
	opt := DefaultOptions()
	opt.FieldsPerRecord = 2
	// Third field should cause an error with fixed width 2.
	_, err := ReadString("a,b,c\n", opt)
	if err == nil {
		t.Fatal("expected error for fields-per-record mismatch")
	}

	opt.FieldsPerRecord = -1
	rows, err := ReadString("a,b,c\n", opt)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || len(rows[0]) != 3 {
		t.Fatalf("unexpected rows: %#v", rows)
	}
}
