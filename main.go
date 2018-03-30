package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"

	"github.com/ChimeraCoder/anaconda"
	"github.com/fatih/color"
)

//go:generate bash ./geninfo.sh

func mainLoop(sig <-chan os.Signal) os.Signal {
	log.Printf("Start main loop.")
	anaconda.SetConsumerKey(ConsumerKey)
	anaconda.SetConsumerSecret(ConsumerSecret)
	tw := anaconda.NewTwitterApi(OAuthToken, OAuthSecret)
	defer tw.Close()
	stream := tw.PublicStreamSample(nil)
	defer stream.Stop()
	for {
		select {
		case s, ok := <-sig:
			if !ok {
				log.Fatalf("Signal chan %v is closed.", s)
			}
			return s
		case msg, ok := <-stream.C:
			if !ok {
				log.Fatalf("Stream chan %v is closed.", stream.C)
			}
			switch tweet := msg.(type) {
			case anaconda.Tweet:
				user := tweet.User
				if user.Following {
					continue
				}
				text := tweet.Text
				for _, kw := range Keywords {
					if strings.Contains(text, kw) {
						log.Printf("%s(%s) blocked (by tweet): %s\n%s", user.ScreenName, user.Name, kw, text)
						_, err := tw.BlockUserId(user.Id, nil)
						if err != nil {
							log.Error("Failed to block", err)
						}
						break
					} else if strings.Contains(user.ScreenName, kw) {
						log.Printf("%s(%s) blocked (by screen name): %s\n%s", user.ScreenName, user.Name, kw, user.ScreenName)
						_, err := tw.BlockUserId(user.Id, nil)
						if err != nil {
							log.Error("Failed to block", err)
						}
						break
					} else if strings.Contains(user.Description, kw) {
						log.Printf("%s(%s) blocked (by description): %s\n %s", user.ScreenName, user.Name, kw, user.Description)
						_, err := tw.BlockUserId(user.Id, nil)
						if err != nil {
							log.Error("Failed to block", err)
						}
						break
					}
				}
			case anaconda.StatusDeletionNotice:
				// pass
			default:
				fmt.Printf("unknown type(%T) : %v \n", msg, msg)
			}
		}
	}
}

func printLogo() {
	log.Info("****************************************")
	log.Info(color.BlueString("  block-them-all  "))
	log.Info("****************************************")
	log.Infof("Build at: %s", color.MagentaString("%s", BuildAt()))
	log.Infof("Git Revision: \n%s", color.MagentaString("%s", DecodeGitRev()))
}

func main() {
	//var err error

	printLogo()
	flag.Parse()
	log.Info("----------------------------------------")
	log.Info("Initializing...")
	log.Info("----------------------------------------")

	log.Info(color.GreenString("                                    [OK]"))

	log.Info("----------------------------------------")
	log.Info("Initialized.")
	log.Info("----------------------------------------")

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := mainLoop(sig)
	log.Fatalf("Signal (%v) received, stopping\n", s)
}
