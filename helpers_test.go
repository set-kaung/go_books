package main

import "testing"

func TestCreateLocalSoruces(t *testing.T) {
	_, err := createLocalSources()
	if err != nil {
		t.Fatalf("file to check local sources %+v", err)
	}
}
