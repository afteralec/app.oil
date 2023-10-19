package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"ariga.io/atlas-go-sdk/atlasexec"
)

func TestLintMigrations(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("pwd failed with err: %v", err)
	}
	client, err := atlasexec.NewClient(wd, "atlas")
	if err != nil {
		t.Errorf("Failed to create client with err: %v", err)
	}
	migrationCount, err := CountMigrations()
	if err != nil {
		t.Errorf("Failed counting migrations: %v", err)
	}
	p := atlasexec.LintParams{}
	p.DevURL = "mysql://root:pass@127.0.0.1:3306/oilatlasdev"
	p.DirURL = "file://migrations"
	p.Latest = uint64(migrationCount)
	report, err := client.Lint(context.Background(), &p)
	if err != nil {
		t.Errorf("Lint failed with error: %v", err)
	}
	for _, file := range report.Files {
		if file.Error != "" {
			t.Errorf("%v", file.Error)
		}
	}
}

func CountMigrations() (count int, err error) {
	pattern := filepath.Join("migrations", "*.sql")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return 0, err
	}
	count = len(files)
	return count, err
}
