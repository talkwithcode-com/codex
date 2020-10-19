package engine

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

type mockCommander struct {
}

func (mc mockCommander) StdinPipe() (io.WriteCloser, error) {
	return nil, errors.New("Error pipe stdin")
}

func (mc mockCommander) StderrPipe() (io.ReadCloser, error) {
	return nil, errors.New("Error pipe stderr")
}

func (mc mockCommander) Output() ([]byte, error) {
	return []byte("test"), nil
}

type mockCommanderOutputError struct {
}

func (mc mockCommanderOutputError) StdinPipe() (io.WriteCloser, error) {
	_, p := io.Pipe()
	return p, nil
}

func (mc mockCommanderOutputError) StderrPipe() (io.ReadCloser, error) {
	p, _ := io.Pipe()
	return p, nil
}

func (mc mockCommanderOutputError) Output() ([]byte, error) {
	return nil, errors.New("output error")
}

type mockExec struct {
}

func (me mockExec) CommandContext(ctx context.Context, name string, arg ...string) lib.Commander {
	return &mockCommander{}
}

type mockExecOutput struct {
}

func (me mockExecOutput) CommandContext(ctx context.Context, name string, arg ...string) lib.Commander {
	return &mockCommanderOutputError{}
}

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

func TestEngine_Run(t *testing.T) {
	type fields struct {
		tempDir  string
		exec     lib.Executor
		fm       lib.FileCreateRemover
		filename string
	}
	type args struct {
		ctx   context.Context
		code  *lib.Code
		input []byte
	}

	eg := New(os.TempDir(), "main")

	goSourceCode := `
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

	goCodeWithStdin := `
			package main

		import (
			"bufio"
			"fmt"
			"os"
		)

		func main() {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			fmt.Fprint(os.Stdout, text)

		}
	`

	egWithStdin := &Engine{
		filename: "app",
		exec:     new(Exec),
		tempDir:  os.TempDir(),
		fm:       new(FileManger),
	}

	codeWithStd, _ := egWithStdin.WriteCode(lang.Go, goCodeWithStdin)

	code, _ := eg.WriteCode(lang.Go, goSourceCode)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *lib.Output
		wantErr bool
	}{
		{
			name: "Run golang code (without stdin provided)",
			args: args{
				ctx:  context.Background(),
				code: code,
			},
			fields: fields(*eg),
			want: &lib.Output{
				Stderr: "error from go\n",
				Stdout: "Hello\n",
			},
			wantErr: false,
		},
		{
			name: "Run golang code (stdin provided)",
			args: args{
				ctx:   context.Background(),
				code:  codeWithStd,
				input: []byte("hei"),
			},
			fields: fields(*egWithStdin),
			want: &lib.Output{
				Stderr: "",
				Stdout: "hei",
			},
			wantErr: false,
		}, {
			name: "Run golang: stdin should be error",
			args: args{
				ctx:   context.Background(),
				code:  codeWithStd,
				input: []byte("hei"),
			},
			fields: fields{
				filename: "app",
				exec:     new(mockExec),
				fm:       new(FileManger),
				tempDir:  os.TempDir(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Run golang: stderr should be error",
			args: args{
				ctx:   context.Background(),
				code:  codeWithStd,
				input: nil,
			},
			fields: fields{
				filename: "app",
				exec:     new(mockExec),
				fm:       new(FileManger),
				tempDir:  os.TempDir(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Run golang: output should be error",
			args: args{
				ctx:   context.Background(),
				code:  codeWithStd,
				input: nil,
			},
			fields: fields{
				filename: "app",
				exec:     new(mockExecOutput),
				fm:       new(FileManger),
				tempDir:  os.TempDir(),
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
			got, err := e.Run(tt.args.ctx, tt.args.code, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Engine.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Engine.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEngine_DeleteCode(t *testing.T) {
	os.Mkdir("temp", 0777)
	eg := New("temp", "main")

	eg.DeleteCode(&lib.Code{
		Path:     "temp",
		Language: lang.Go,
	})

	_, err := ioutil.ReadDir("temp")

	if err == nil {
		t.Errorf("Expected nil but got %v", err)
	}

}
