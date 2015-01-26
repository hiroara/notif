package main

import (
	"github.com/codegangsta/cli"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "notif"
	app.Usage = "send messages to the Notification Center on OSX."
	app.Author = "Hiroki Arai"
	app.Email = "hiroara62@gmail.com"
	app.Version = "0.0.1"
	app.Action = notif
	app.Flags = flags

	app.Run(os.Args)
}

var flags = []cli.Flag{
	cli.StringFlag{
		Name:  "title, t",
		Value: "notif",
		Usage: "title of notification",
	},
	cli.StringFlag{
		Name:  "subtitle, s",
		Usage: "subtitle of notification",
	},
}

func notif(c *cli.Context) {
	cmd := exec.Command("osascript")

	handleError(pipe(cmd.StdoutPipe, os.Stdout))

	handleError(pipe(cmd.StderrPipe, os.Stderr))

	stdin, err := cmd.StdinPipe()
	handleError(err)

	input, err := getInput(c)
	handleError(err)
	_, err = io.WriteString(stdin, "display notification"+input+getOptions(c))
	handleError(err)
	handleError(stdin.Close())

	handleError(cmd.Start())
	defer cmd.Wait()
}

func pipe(sourceGetter func() (io.ReadCloser, error), dist io.WriteCloser) (err error) {
	out, err := sourceGetter()
	if err != nil {
		return
	}
	go io.Copy(dist, out)
	return
}

func escape(s string) string {
	return "\"" + strings.Replace(s, "\"", "\\\"", -1) + "\""
}

func getInput(c *cli.Context) (string, error) {
	if len(c.Args()) > 0 {
		return escape(strings.Join(c.Args(), " ")), nil
	}
	bytes, err := ioutil.ReadAll(os.Stdin)
	return escape(string(bytes)), err
}

func getOptions(c *cli.Context) string {
	return strings.Join([]string{"with", getTitle(c), getSubTitle(c)}, " ")
}

func getTitle(c *cli.Context) string {
	title := c.String("title")
	return "title " + escape(title)
}

func getSubTitle(c *cli.Context) string {
	subtitle := c.String("subtitle")
	if subtitle == "" {
		return ""
	}
	return "subtitle " + escape(subtitle)
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
