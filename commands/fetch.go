package commands

import (
	"fmt"
	"os"
	"time"

	rss "github.com/mattn/go-pkg-rss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os/signal"
)


type Config struct {
	Feeds []string
	Port int
}

var fetchCmd = cobra.Command{
	Use: "fetch",
	Short: "Fetch feeds",
	Long: "Bastion will fetch all feeds listed in the config file",
	Run: fetchRun,
}

func init() {
	fetchCmd.Flags().Int("rsstimeout", 5, "Timeout (in min) for RSS retrieval")
	viper.BindPFlag("rsstimeout", fetchCmd.Flags().Lookup("rsstimeout"))
}

func fetchRun(cmd *cobra.Command, arg []string) {

	Fetcher()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}

func Fetcher() {
	var config Config

	if err := viper.Marshal(&config); err != nil {
		fmt.Println(err)
	}

	for _, feed := range config.Feeds {
		go PoolFeed(feed)
	}
}

func PoolFeed(uri string) {
	timeout := viper.GetInt("RSSTimeout")
	if timeout < 1 {
		timeout = 1
	}
	feed := rss.New(timeout, true, chanHandler, itemHandler)

	for {
		if err := feed.Fetch(uri, nil); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s: %s", uri, err)
			return
		}

		fmt.Printf("Sleeping for %d seconds on %s\n", feed.SecondsTillUpdate(), uri)
		time.Sleep(time.Duration(feed.SecondsTillUpdate() * 1e9))
	}
}