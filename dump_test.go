package main

import (
	"reflect"
	"testing"
)

func TestGetProjects(t *testing.T) {
	tests := []struct {
		projects []string
		want     []string
	}{
		{},
		{[]string{"foo", "bar"}, []string{"foo", "bar"}},
		{[]string{"foo", "bar", ""}, []string{"foo", "bar"}},
		// TODO: Add tests involving error cases.
	}
	for _, tt := range tests {
		cmd := helperCommand("echo", tt.projects...)
		got, err := getProjects(cmd)
		if err != nil {
			t.Errorf("getProjects(%v) returned non-nil error: %v", cmd.Args, err)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("getProjects(%v) = %v, want %v", cmd.Args, got, tt.want)
		}
	}
}

func TestGetSpaceSeparated(t *testing.T) {
	tests := []struct {
		projects []string
		want     []string
	}{
		{},
		{[]string{"foo", "bar"}, []string{"foo", "bar"}},
		{[]string{"foo", "bar", ""}, []string{"foo", "bar"}},
		// TODO: Add tests involving error cases.
	}
	for _, tt := range tests {
		cmd := helperCommand("echo", tt.projects...)
		got, err := getSpaceSeparated(cmd)
		if err != nil {
			t.Errorf("getSpaceSeparated(%v) returned non-nil error: %v", cmd.Args, err)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("getSpaceSeparated(%v) = %v, want %v", cmd.Args, got, tt.want)
		}
	}
}
