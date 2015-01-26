package main

import (
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
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
	cli.StringFlag{
		Name:  "sound, S",
		Value: "Default",
		Usage: "sonund of notification",
	},
}

func notif(c *cli.Context) {
	n := &Notifier{}
	message, err := getMessage(c)
	if err != nil {
		log.Fatal(err)
	}

	err = n.Send(message, c.String("title"), c.String("subtitle"), c.String("sound"))
	if err != nil {
		log.Fatal(err)
	}
}

func getMessage(c *cli.Context) (string, error) {
	if len(c.Args()) > 0 {
		return strings.Join(c.Args(), " "), nil
	}
	bytes, err := ioutil.ReadAll(os.Stdin)
	return string(bytes), err
}
