package engine

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/talkwithcode-com/codex/lib"
	"github.com/talkwithcode-com/codex/lib/lang"
)

const golangCode = `
	package main

	import (
		"fmt"
		"os"
	)

	func main() {
		fmt.Fprintln(os.Stdout, "Hello")
		fmt.Fprintln(os.Stderr, "error from go")
	}
`

const javaScriptCode = `
	console.log("Hello")
	console.error("error from javascript)
`

type mockFileManager struct{}

func (mfm *mockFileManager) Create(name string) (*os.File, error) {
	return nil, errors.New("failed_create_file")
}

func (mfm *mockFileManager) Remove(path string) error {
	return nil
}

func TestEngine_WriteCode(t *testing.T) {
	type fields struct {
		tempDir  string
		exec     lib.Executor
		fm       lib.FileCreateRemover
		filename string
	}
	type args struct {
		language   lang.Language
		sourceCode string
	}

	f :=
		fields{
			tempDir:  os.TempDir(),
			fm:       new(FileManger),
			filename: "main",
		}

	fMockFm := f
	fMockFm.fm = new(mockFileManager)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *lib.Code
		wantErr bool
	}{
		{
			name:   "Write Golang code into tempDir",
			fields: f,
			args: args{
				language:   lang.Go,
				sourceCode: golangCode,
			},
			want: &lib.Code{
				Path:     fmt.Sprintf("%s/main.%s", os.TempDir(), lang.LanguageConfig[lang.Go].Extension),
				Language: lang.Go,
			},
			wantErr: false,
		},
		{
			name:   "Write JavaScript code into tempDir",
			fields: f,
			args: args{
				language:   lang.JavaScript,
				sourceCode: javaScriptCode,
			},
			want: &lib.Code{
				Path:     fmt.Sprintf("%s/main.%s", os.TempDir(), lang.LanguageConfig[lang.JavaScript].Extension),
				Language: lang.JavaScript,
			},
			wantErr: false,
		},
		{
			name:   "Should be error to create code",
			fields: fMockFm,
			args: args{
				language:   lang.JavaScript,
				sourceCode: javaScriptCode,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Engine{
				tempDir:  tt.fields.tempDir,
				exec:     tt.fields.exec,
				fm:       tt.fields.fm,
				filename: tt.fields.filename,
			}
			got, err := e.WriteCode(tt.args.language, tt.args.sourceCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("Engine.WriteCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Engine.WriteCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
