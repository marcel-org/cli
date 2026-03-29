package cli

import (
	"flag"
	"reflect"
	"testing"
)

func TestNormalizeArgsMovesFlagsAheadOfPositionals(t *testing.T) {
	fs := newFlagSet("event add")
	fs.String("date", "", "event date")
	fs.String("time", "", "event time")

	got := normalizeArgs(fs, []string{"deadline backoffice", "--date", "2026-04-01", "--time", "09:00"})
	want := []string{"--date", "2026-04-01", "--time", "09:00", "deadline backoffice"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizeArgs mismatch\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestNormalizeArgsKeepsBoolFlagsWithoutConsumingPositionals(t *testing.T) {
	fs := newFlagSet("event list")
	fs.Bool("all", false, "include past events")
	fs.Bool("json", false, "output JSON")

	got := normalizeArgs(fs, []string{"today", "--all", "--json"})
	want := []string{"--all", "--json", "today"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizeArgs mismatch\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestParseArgsAllowsPositionalBeforeFlags(t *testing.T) {
	fs := newFlagSet("event add")
	date := fs.String("date", "", "event date")

	if err := parseArgs(fs, []string{"deadline backoffice", "--date", "2026-04-01"}); err != nil {
		t.Fatalf("parseArgs returned error: %v", err)
	}

	if *date != "2026-04-01" {
		t.Fatalf("expected date to be parsed, got %q", *date)
	}

	if fs.NArg() != 1 || fs.Arg(0) != "deadline backoffice" {
		t.Fatalf("expected positional title to be preserved, got %q", fs.Arg(0))
	}
}

func TestFlagConsumesValueRecognizesBoolFlags(t *testing.T) {
	fs := flag.NewFlagSet("event list", flag.ContinueOnError)
	fs.Bool("all", false, "include past events")

	if flagConsumesValue(fs, "--all") {
		t.Fatal("expected --all to be treated as bool flag")
	}
}
