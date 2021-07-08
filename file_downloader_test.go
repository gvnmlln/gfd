package main

import (
	"os"
	"reflect"
	"testing"
)

const testResourcesDir = "test-resources/"

var wd, _ = os.Getwd()

func TestParseFileName(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{name: "basic", args: args{url: "directory/filename.ext"}, want: "filename.ext", want1: "filename"},
		{name: "no extension", args: args{url: "directory/filepath"}, want: "filepath", want1: "filepath"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ParseFileName(tt.args.url)
			if got != tt.want {
				t.Errorf("ParseFileName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseFileName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestParseURLsFromTextFile(t *testing.T) {
	type args struct {
		filepath string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "valid filepath, only valid urls",
			args: args{filepath: testResourcesDir + "only-valid-urls.txt"},
			want: []string{
				"https://upload.wikimedia.org/wikipedia/commons/f/f1/Paul_E._Patton_2013.jpg",
				"https://upload.wikimedia.org/wikipedia/commons/8/84/Dmitri_N_Smirnov_%C2%A9Kompozitor.jpg"},
			wantErr: false,
		},

		{name: "valid filepath, valid and invalid urls",
			args: args{filepath: testResourcesDir + "valid-and-invalid-urls.txt"},
			want: []string{
				"https://upload.wikimedia.org/wikipedia/commons/f/f1/Paul_E._Patton_2013.jpg",
				"https://upload.wikimedia.org/wikipedia/commons/8/84/Dmitri_N_Smirnov_%C2%A9Kompozitor.jpg"},
			wantErr: false,
		},

		{name: "valid filepath, only invalid urls",
			args:    args{filepath: testResourcesDir + "only-invalid-urls.txt"},
			want:    nil,
			wantErr: false,
		},

		{name: "invalid filepath",
			args:    args{filepath: "fish"},
			want:    nil,
			wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseURLsFromTextFile(tt.args.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseURLsFromTextFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseURLsFromTextFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExists(t *testing.T) {
	type args struct {
		filepath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "test if this file exist", args: args{filepath: "file_downloader_test.go"}, want: true},
		{name: "test if home directory exists", args: args{filepath: wd}, want: true},
		{name: "invalid filepath", args: args{filepath: "hello world"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExists(tt.args.filepath); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURLIsValid(t *testing.T) {
	type args struct {
		u string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "valid url", args: args{u: "https://www.google.com/"}, want: true},
		{name: "invalid url", args: args{u: "fish"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := URLIsValid(tt.args.u); got != tt.want {
				t.Errorf("URLIsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
