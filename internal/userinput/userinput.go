// package userinput provides utilities for getting input from the user.
package userinput

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func writeAndClose(f *os.File, data []byte) error {
	_, err := f.Write(data)
	f.Close()
	return err
}

// A LargeInputStrategy implements some means of populating user
// input into the given file. Strategies should take into account
// that the file may be pre-populated with a template/default value.
type LargeInputStrategy interface {
	GetInput(filename string) error
}

// PureLargeInputStrategy is a wrapper type for functions which already implement
// the LargeInputStrategy.GetInput signature so they may be used in functions
// which require a LargeInputStrategy-implementor.
type PureLargeInputStrategy func(string) error

func (lis PureLargeInputStrategy) GetInput(filename string) error {
	return lis(filename)
}

func VimStrategy(filename string) error {
	cmd := exec.Command("vim", filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func NanoStrategy(filename string) error {
	cmd := exec.Command("nano", filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func VsCodeStrategy(filename string) error {
	return exec.Command("code", "-w", filename).Run()
}

func TyporaStrategy(filename string) error {
	return exec.Command("typora", filename).Run()
}

// UseEditor chooses a LargeInputStrategy implementation based
// on a provided "enum" string. Accepted enums are "vim", "nano", "vscode", "typora"
func UseEditor(editor string) (LargeInputStrategy, error) {
	switch editor {
	case "vim":
		return PureLargeInputStrategy(VimStrategy), nil
	case "nano":
		return PureLargeInputStrategy(NanoStrategy), nil
	case "vscode":
		return PureLargeInputStrategy(VsCodeStrategy), nil
	case "typora":
		return PureLargeInputStrategy(TyporaStrategy), nil
	default:
		return nil, errors.New(fmt.Sprintf("not an editor option: %s", editor))
	}
}

// LargeInput gathers input from the user via a separate text editor.
// The template arg is what is initially displayed the editor. The
// LargeInputStrategy represents the editor to use.
func LargeInput(template string, editor LargeInputStrategy) (string, error) {
	tempfile, err := ioutil.TempFile("", "*.md")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempfile.Name())
	if err := writeAndClose(tempfile, []byte(template)); err != nil {
		return "", err
	}
	if err := editor.GetInput(tempfile.Name()); err != nil {
		return "", err
	}
	fileContents, err := ioutil.ReadFile(tempfile.Name())
	if err != nil {
		return "", err
	}
	return string(fileContents), nil
}

// YesOrNo returns a boolean value based on a yes/no input from the
// user (yes=true, no=false). If the user inputs nothing (i.e. presses ENTER)
// the value of yesIsDefault is returned.
func YesOrNo(prompt string, yesIsDefault bool) (bool, error) {
	var optionStr string
	if yesIsDefault {
		optionStr = "[Y/n]"
	} else {
		optionStr = "[y/N]"
	}
	for {
		choice := StdinPrompt(prompt + " " + optionStr)
		switch strings.ToLower(choice) {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		case "":
			return yesIsDefault, nil
		}
	}
}

// StdinPrompt receives a one-line response from the user via stdin.
func StdinPrompt(prompt string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt + " ")
	scanner.Scan()
	return strings.TrimRight(scanner.Text(), "\n")
}
