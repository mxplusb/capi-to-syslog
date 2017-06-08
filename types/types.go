/*
Mike's note: This is the weirdest code I've ever written. I wrote most of it, but it's fucking weird.
*/

package types

import (
	"encoding/json"
	"errors"
	"time"
	"fmt"
)

const (
	AppCreateKind EventKind = iota
	AppStopKind
	AppDeleteKind
	AppUpdateKind
	AppCrashKind
	AppSSHKind
)

// EventKind is the type of auditing event from Cloud Foundry.
type EventKind int

// UnmarshalJSON maps an event to it's appropriate type so it can be unmarshalled more sanely.
func (ev *EventKind) UnmarshalJSON(data []byte) error {
	var tmp interface{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	s, ok := tmp.(string)
	if !ok {
		return errors.New("expected string") // improve this error message if you really care
	}
	switch s {
	case "audit.app.create":
		*ev = AppCreateKind
	case "audit.app.stop":
		*ev = AppStopKind
	case "audit.app.delete-request":
		*ev = AppDeleteKind
	case "audit.app.update":
		*ev = AppUpdateKind
	case "audit.app.ssh-authorized":
		*ev = AppSSHKind
	case "audit.app.ssh-unauthorized":
		*ev = AppSSHKind
	case "app.crash":
		*ev = AppCrashKind
	}
	return nil
}

// TypeHandler is our interface for typing the events API with the entity metadata by exploiting the JSON interface.
var TypeHandler = map[EventKind]func() interface{}{
	AppCreateKind: func() interface{} {
		// Need the extra level of indirection here, so return an unnamed struct.
		// See (*AppEntityEnv).UnmarshalJSON for more details on why this is so.
		return &struct {
			Metadata AppCreateMetadata `json:"metadata"`
		}{}
	},
	AppStopKind: func() interface{} {
		return &struct {
			Metadata AppStopEvent `json:"metadata"`
		}{}
	},
	AppDeleteKind: func() interface{} {
		return &struct {
			Metadata AppDeleteEvent `json:"metadata"`
		}{}
	},
	AppUpdateKind: func() interface{} {
		return &struct {
			Metadata AppUpdateEvent `json:"metadata"`
		}{}
	},
	AppSSHKind: func() interface{} {
		return &struct {
			Metadata AppSSHEvent `json:"metadata"`
		}{}
	},
	AppCrashKind: func() interface{} {
		return &struct {
			Metadata AppCrashEvent `json:"metadata"`
		}{}
	},
}

// AppEvent is the top-level object.
type AppEvent struct {
	TotalResults int        `json:"total_results"`
	TotalPages   int        `json:"total_pages"`
	PrevURL      string     `json:"prev_url,omitempty"`
	NextURL      string     `json:"next_url,omitempty"`
	Resources    []Resource `json:"resources"`
}

type Resource struct {
	Metadata EventMetadata `json:"metadata"`
	Entity   AppEntityEnv  `json:"entity"`
}

type EventMetadata struct {
	GUID      string    `json:"guid"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AppEntityEnv struct {
	AppEntity // embed this here, so we don't recurse in (*AppEntityEnv).UnmarshalJSON.
	Metadata interface{} `json:"metadata"`
}

func (a *AppEntityEnv) String() string {
	return fmt.Sprintf("%#v", a.Metadata)
}

func (e *AppEntityEnv) UnmarshalJSON(data []byte) error {
	// First, unmarshal the common bits.
	var env AppEntity
	if err := json.Unmarshal(data, &env); err != nil {
		return err
	}
	// Now, we know the event type and we can delegate to
	// the type handler, to give us the appropriate thing
	// to unmarshal the metadata into.
	metaFn := TypeHandler[env.Type]
	// Check for nil, to be safe.
	if metaFn == nil {
		return errors.New("unhandled event type")
	}
	// Call the function which gives us something to unmarshal to,
	// then do the actual unmarshalling. I think this works because
	// those types have their
	meta := metaFn()
	if err := json.Unmarshal(data, &env); err != nil {
		return err
	}
	// We now have valid metadata and can construct the whole thing.
	*e = AppEntityEnv{
		AppEntity: env,
		Metadata:  meta,
	}
	return nil
}

// AppEntity contains the common bits of an entity. It is embedded in AppEntityEnv.
type AppEntity struct {
	Type             EventKind `json:"type"`
	Actor            string    `json:"actor"`
	ActorType        string    `json:"actor_type"`
	ActorName        string    `json:"actor_name"`
	Actee            string    `json:"actee"`
	ActeeType        string    `json:"actee_type"`
	ActeeName        string    `json:"actee_name"`
	Timestamp        string    `json:"timestamp"`
	SpaceGUID        string    `json:"space_guid"`
	OrganizationGUID string    `json:"organization_guid"`
}

type AppCreateMetadata struct {
	Request struct {
		Name                  string `json:"name"`
		Instances             int    `json:"instances"`
		Memory                int    `json:"memory"`
		State                 string `json:"state"`
		EnvironmentJSON       string `json:"environment_json"`
		DockerCredentialsJSON string `json:"docker_credentials_json"`
	} `json:"request"`
}

type Request struct {
	Name             string `json:"name"`
	OrganizationGUID string `json:"organization_guid"`
	AllowSSH         bool   `json:"allow_ssh"`
}

type AppDeleteEvent struct {
	Recursive bool `json:"recursive"`
}

type AppCrashEvent struct {
	Instance        int    `json:"instance"`
	Index           int    `json:"index"`
	ExitStatus      string `json:"exit_status"`
	ExitDescription string `json:"exit_description"`
	Reason          string `json:"reason"`
}

type AppSSHEvent struct {
	Index int `json:"index"`
}

type AppUpdateEvent struct {
	Name                  string `json:"name"`
	Instances             int    `json:"instances"`
	Memory                int    `json:"memory"`
	State                 string `json:"state"`
	EnvironmentJSON       string `json:"environment_json"`
	DockerCredentialsJSON string `json:"docker_credentials_json"`
}

type AppStopEvent struct{}
