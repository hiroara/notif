package main

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

type Notifier struct{}

func (n *Notifier) Send(message, title, subtitle, sound string) (err error) {
	cmd := exec.Command("osascript")

	if err = pipeAll(cmd, os.Stdout, os.Stderr); err != nil {
		return
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}
	defer cmd.Wait()

	_, err = io.WriteString(stdin, strings.Join([]string{"display notification", escape(message), getOptions(title, subtitle, sound)}, " "))
	if err != nil {
		return
	}

	if err = stdin.Close(); err != nil {
		return
	}

	return err
}

func escape(s string) string {
	return "\"" + strings.Replace(s, "\"", "\\\"", -1) + "\""
}

func getOptions(title, subtitle, sound string) string {
	options := make([]string, 0, 5)
	options = append(options, "with", "title", escape(title))
	if subtitle != "" {
		options = append(options, "subtitle", escape(subtitle))
	}
	if sound != "" {
		options = append(options, "sound name", escape(sound))
	}
	return strings.Join(options, " ")
}

func pipeAll(cmd *exec.Cmd, stdout, stderr io.WriteCloser) (err error) {
	if err = pipe(cmd.StdoutPipe, stdout); err != nil {
		return
	}
	if err = pipe(cmd.StderrPipe, stderr); err != nil {
		return
	}
	return
}

func pipe(sourceGetter func() (io.ReadCloser, error), dist io.WriteCloser) (err error) {
	out, err := sourceGetter()
	if err != nil {
		return
	}
	go io.Copy(dist, out)
	return
}
