package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestBashBootstrap(t *testing.T) {
	f, err := os.CreateTemp("", "rig.*.bashrc")
	defer func() {
		if e := f.Close(); e != nil {
			fmt.Println(e)
		}
		_ = os.Remove(f.Name())
	}()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	_, err = f.WriteString("\n# Some string here\n\n")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	thing := &BashBootstrapper{
		rcFile: f.Name(),
	}

	if thing.IsBootstrapped() {
		fmt.Println("This is not possible... it does not exist!")
		t.FailNow()
	}
	err = thing.DoBootstrap()

	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	hasInitialText := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Some string here") {
			hasInitialText = true
			break
		}
	}

	if !hasInitialText {
		fmt.Println("Clobbers contents of file")
		t.FailNow()
	}
}
