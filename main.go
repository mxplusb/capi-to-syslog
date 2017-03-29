package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"gopkg.in/robfig/cron.v2"
)

var CapiClientID, CapiClientSecret, CapiSystemURI string
var InsecureSkipVerify bool
var AuditableEvents []string

type Events struct {
	TotalResults int         `json:"total_results"`
	TotalPages   int         `json:"total_pages"`
	PrevURL      string      `json:"prev_url,omitempty"`
	NextURL      string      `json:"next_url,omitempty"`
	Resources    []Resources `json:"resources"`
}

type Resources struct {
	Metadata Metadata
	Entity   Entity
}

type Metadata struct {
	GUID      string    `json:"guid"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Entity struct {
	Type             string    `json:"type"`
	Actor            string    `json:"actor"`
	ActorType        string    `json:"actor_type"`
	ActorName        string    `json:"actor_name"`
	Actee            string    `json:"actee"`
	ActeeType        string    `json:"actee_type"`
	ActeeName        string    `json:"actee_name"`
	Timestamp        time.Time `json:"timestamp"`
	EntityMetadata   Request   `json:"metadata"`
	SpaceGUID        string    `json:"space_guid"`
	OrganizationGUID string    `json:"organization_guid"`
}

type Request struct {
	Name             string `json:"name"`
	OrganizationGUID string `json:"organization_guid"`
	AllowSSH         bool   `json:"allow_ssh"`
}

func GetLogs(client *http.Client, r *http.Request, ch chan bool) {

	resp, err := client.Do(r)
	if err != nil {
		ch <- false
		panic(err)
	}
	defer resp.Body.Close()

	var events Events
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		ch <- false
		panic(err)
	}

	for resource := range events.Resources {
		fmt.Printf("%#v\n", events.Resources[resource].Entity)
	}
}

// type:audit.app.ssh-authorized,type:audit.app.ssh-unauthorized,type:audit.app.create,type:audit.app.start,type:audit.app.stop,type:audit.app.update,type:audit.app.delete-request,type:audit.service_key.create,type:audit.service_key.delete,type:audit.space.create

func RequestBuilder(idx int, listenChan chan bool, client *http.Client) {
	req, err := http.NewRequest("GET", "https://api."+CapiSystemURI+"/v2/events", nil)
	if err != nil {
		listenChan <- false
		panic(err)
	}

	now := time.Now()
	then := now.Add(time.Duration(time.Minute * -1))

	q := req.URL.Query()
	q.Add("q", AuditableEvents[idx])
	q.Add("q", "timestamp>"+then.Format(time.RFC3339))
	req.URL.RawQuery = q.Encode()

	fmt.Printf("looking for events here: %s\n", AuditableEvents[idx])

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	GetLogs(client, req, listenChan)
}

func SetEnvVars() {
	var err error
	CapiClientID = os.Getenv("CAPI_CLIENT_ID")
	CapiClientSecret = os.Getenv("CAPI_CLIENT_SECRET")
	CapiSystemURI = os.Getenv("CAPI_SYSTEM_URI")

	if CapiSystemURI == "" || CapiClientSecret == "" || CapiClientID == "" {
		panic("cannot continue with the CAPI configs!")
	}

	InsecureSkipVerifyString := os.Getenv("INSECURE_SKIP_VERIFY")

	if InsecureSkipVerifyString == "true" {
		InsecureSkipVerify, err = strconv.ParseBool(InsecureSkipVerifyString)
		if err != nil {
			panic(err)
		}
	}

	localAudits := os.Getenv("CAPI_EVENTS")
	events := strings.Split(localAudits, ",")
	AuditableEvents = append(AuditableEvents, events...)

	if len(AuditableEvents) == 0 {
		panic("can't watch nothing.")
	}
}

func main() {

	SetEnvVars()

	// create out Oauth2 config.
	config := clientcredentials.Config{
		ClientID:     CapiClientID,
		ClientSecret: CapiClientSecret,
		TokenURL:     "https://uaa." + CapiSystemURI + "/oauth/token",
	}

	// create the Oauth2 http client.
	var client *http.Client

	if InsecureSkipVerify {
		tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		sslcli := &http.Client{Transport: tr}
		ctx := context.TODO()
		ctx = context.WithValue(ctx, oauth2.HTTPClient, sslcli)
		client = config.Client(ctx)
	} else {
		ctx := context.Background()
		client = config.Client(ctx)
	}

	listenChan := make(chan bool)

	bleh := func() {
		for idx := range AuditableEvents {
			RequestBuilder(idx, listenChan, client)
		}
	}

	bleh()

	c := cron.New()
	c.AddFunc("59 * * * * *", bleh)
	c.Start()

	for i := range listenChan {
		if !i {
			return
		}
	}
	close(listenChan)
}
