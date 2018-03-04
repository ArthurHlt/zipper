package main

import (
	"crypto/tls"
	"fmt"
	"github.com/ArthurHlt/zipper"
	"github.com/urfave/cli"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Name = "zipper"
	app.Usage = "use zipper in cli"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "type, t",
			Usage: "Choose source type",
		},
		cli.BoolFlag{
			Name:  "insecure, k",
			Usage: "Ignore certificate validation",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "zip",
			Aliases:   []string{"z"},
			Usage:     "create zip from a source",
			ArgsUsage: "<source uri>",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output, o",
					Value: "content.zip",
					Usage: "zip file in another path (you can set to - to write in stdout)",
				},
			},
			Action: zip,
		},
		{
			Name:      "sha1",
			Aliases:   []string{"s"},
			Usage:     "Get sha1 signature for the file from source",
			ArgsUsage: "<source uri>",
			Action:    sha1,
		},
		{
			Name:      "diff",
			Aliases:   []string{"s"},
			Usage:     "Check if file from source is different from your stored sha1",
			ArgsUsage: "<source uri> <stored sha1>",
			Action:    diff,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
func checkPath(c *cli.Context) error {
	if c.Args().First() == "" {
		return fmt.Errorf("You must pass an uri as first argument.")
	}
	return nil
}
func createSession(c *cli.Context) (*zipper.Session, error) {
	err := checkPath(c)
	if err != nil {
		return nil, err
	}
	path := c.Args().First()
	handlerType := c.GlobalString("type")
	zipper.SetHttpClient(&http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: c.GlobalBool("insecure"),
			},
		},
	})
	s, err := zipper.CreateSession(path, handlerType)
	if err != nil {
		return nil, err
	}
	return s, nil
}
func diff(c *cli.Context) error {
	s, err := createSession(c)
	if err != nil {
		return err
	}
	currentSha1 := c.Args().Get(1)
	if currentSha1 == "" {
		return fmt.Errorf("You must pass sha1 you've stored from sha1 command.")
	}
	diff, _, err := s.IsDiff(currentSha1)
	if err != nil {
		return err
	}
	if diff {
		fmt.Println("file from source is different")
		os.Exit(1)
	}
	fmt.Println("no change from source")
	return nil
}
func sha1(c *cli.Context) error {
	s, err := createSession(c)
	if err != nil {
		return err
	}
	sig, err := s.Sha1()
	if err != nil {
		return err
	}
	fmt.Print(sig)
	return nil
}
func zip(c *cli.Context) error {
	s, err := createSession(c)
	if err != nil {
		return err
	}
	z, err := s.Zip()
	if err != nil {
		return err
	}
	defer z.Close()
	output := c.String("output")
	if output == "-" {
		_, err = io.Copy(os.Stdout, z)
		return err
	}
	bar := pb.New(int(z.Size())).SetUnits(pb.U_BYTES).Prefix(filepath.Base(output))
	bar.Start()

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, bar.NewProxyReader(z))
	if err != nil {
		return err
	}
	fmt.Println("Downloaded and zipped at " + f.Name())
	return nil
}
