package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"gopkg.in/robfig/cron.v2"
	"strings"
)

const (
	APIEndpoint = "https://api.run.haas-88.pez.pivotal.io"
)

var Events = []string{"type:audit.app.ssh-authorized", "type:audit.app.ssh-unauthorized"}

func getLogs(client *http.Client, r *http.Request, ch chan bool) {
	resp, err := client.Do(r)
	if err != nil {
		ch <- false
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- false
		panic(err)
	}

	fixed := strings.Replace(string(body), "\n", "", -1)
	fixed = strings.Replace(fixed, "    ", "" , -1)

	fmt.Print(fmt.Sprintf("%s\n", fixed))
}

func main() {
	config := clientcredentials.Config{
		ClientID:     "uaa-to-syslog",
		ClientSecret: "pivotal123!",
		TokenURL:     "https://uaa.run.haas-88.pez.pivotal.io/oauth/token",
	}

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	sslcli := &http.Client{Transport: tr}
	ctx := context.TODO()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, sslcli)
	client := config.Client(ctx)

	listenChan := make(chan bool)

	bleh := func() {
		for idx := range Events {
			req, err := http.NewRequest("GET", APIEndpoint+"/v2/events", nil)
			if err != nil {
				listenChan <- false
				panic(err)
			}

			q := req.URL.Query()
			q.Add("q", Events[idx])
			req.URL.RawQuery = q.Encode()

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			getLogs(client, req, listenChan)
		}
	}

	c := cron.New()
	c.AddFunc("5 * * * * *", bleh)
	c.Start()

	for i := range listenChan {
		if !i {
			return
		}
	}
	close(listenChan)
}