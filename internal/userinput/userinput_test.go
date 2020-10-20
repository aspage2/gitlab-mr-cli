package userinput

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func catStrategy(filename string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte("This is the description"))
	return err
}

func Test_LargeInput(t *testing.T) {
	result, err := LargeInput("This is the first line\n", PureLargeInputStrategy(catStrategy))
	assert.NoError(t, err)
	assert.Equal(t, "This is the first line\nThis is the description", result)
}

func TestUseEditor_Error(t *testing.T) {
	_, err := UseEditor("pycharm")
	assert.Error(t, err)
}

func TestYesOrNo(t *testing.T) {
	for _, test := range [...]struct {
		in  string
		def bool
		exp bool
	}{
		{"y", true, true},
		{"y", false, true},
		{"n", true, false},
		{"n", false, false},
		{"", true, true},
		{"", false, false},
	} {
		t.Run(test.in, func(t *testing.T) {
			inBuf := bytes.NewBuffer([]byte(test.in + "\n"))
			r, w, _ := os.Pipe()
			oldStdin := os.Stdin
			os.Stdin = r
			defer func() {
				_ = r.Close()
				_ = w.Close()
				os.Stdin = oldStdin
			}()
			iC := make(chan struct{})
			go func() {
				_, _ = io.Copy(w, inBuf)
				close(iC)
			}()
			res, err := YesOrNo("", test.def)
			assert.NoError(t, err)
			assert.Equal(t, test.exp, res)
			<-iC
		})
	}
}
