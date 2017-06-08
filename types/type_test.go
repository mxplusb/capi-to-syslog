package types

import (
	"encoding/json"
	"testing"
)

var supportedEvents = []string{
	"audit.app.create",
	"audit.app.delete-request",
	"app.crash",
	"audit.app.ssh-authorized",
	"audit.app.ssh-unauthorized",
	"audit.app.start",
	"audit.app.stop",
	"audit.app.update",
}

var typeMap = map[string]string{
	"audit.app.create":           `{ "total_results": 1, "total_pages": 1, "prev_url": null, "next_url": null, "resources": [ { "metadata": { "guid": "f7ace6b1-7abc-40b1-aef4-155d0e2725f3", "url": "/v2/events/f7ace6b1-7abc-40b1-aef4-155d0e2725f3", "created_at": "2016-06-08T16:41:23Z", "updated_at": "2016-06-08T16:41:26Z" }, "entity": { "type": "audit.app.create", "actor": "uaa-id-40", "actor_type": "user", "actor_name": "user@example.com", "actee": "49275867-0e4d-46b6-b9c4-62e0e4a14eb2", "actee_type": "app", "actee_name": "name-265", "timestamp": "2016-06-08T16:41:23Z", "metadata": { "request": { "name": "new", "instances": 1, "memory": 84, "state": "STOPPED", "environment_json": "PRIVATE DATA HIDDEN", "docker_credentials_json": "PRIVATE DATA HIDDEN" } }, "space_guid": "c21a3373-a7b6-4db1-9f73-f57ecc3e8383", "organization_guid": "ed942a38-b29c-4659-bdb1-4fdc65903dff" } } ] }`,
	"audit.app.delete-request":   `{ "total_results": 1, "total_pages": 1, "prev_url": null, "next_url": null, "resources": [ { "metadata": { "guid": "c1d9970f-e8a9-46a1-bffd-89a6dca00e6b", "url": "/v2/events/c1d9970f-e8a9-46a1-bffd-89a6dca00e6b", "created_at": "2016-06-08T16:41:26Z", "updated_at": "2016-06-08T16:41:26Z" }, "entity": { "type": "audit.app.delete-request", "actor": "uaa-id-92", "actor_type": "user", "actor_name": "user@example.com", "actee": "5fadf60f-6463-4341-9bfd-c569501e2781", "actee_type": "app", "actee_name": "name-990", "timestamp": "2016-06-08T16:41:26Z", "metadata": { "request": { "recursive": false } }, "space_guid": "6934f00f-b2c5-427d-a9cf-43c56fdfcbe8", "organization_guid": "3166cf2e-916f-4026-8528-b6594a219a45" } } ] }`,
	"app.crash":                  `{ "total_results": 1, "total_pages": 1, "prev_url": null, "next_url": null, "resources": [ { "metadata": { "guid": "14bb5650-4fe7-42b2-a085-6337056a2088", "url": "/v2/events/14bb5650-4fe7-42b2-a085-6337056a2088", "created_at": "2016-06-08T16:41:26Z", "updated_at": "2016-06-08T16:41:26Z" }, "entity": { "type": "app.crash", "actor": "32994dbf-9378-4a53-8c05-ad232a1e0137", "actor_type": "app", "actor_name": "name-1051", "actee": "32994dbf-9378-4a53-8c05-ad232a1e0137", "actee_type": "app", "actee_name": "name-1051", "timestamp": "2016-06-08T16:41:26Z", "metadata": { "instance": 0, "index": 1, "exit_status": "1", "exit_description": "out of memory", "reason": "crashed" }, "space_guid": "f9f47800-e974-40fa-b078-b76614928562", "organization_guid": "fd6f03a5-cab7-4d84-bd4d-f73505ff10da" } } ] }`,
	"audit.app.ssh-authorized":   `{ "total_results": 1, "total_pages": 1, "prev_url": null, "next_url": null, "resources": [ { "metadata": { "guid": "20eec8ac-20a2-4775-8926-54ff9f3f1776", "url": "/v2/events/20eec8ac-20a2-4775-8926-54ff9f3f1776", "created_at": "2016-06-08T16:41:24Z", "updated_at": "2016-06-08T16:41:26Z" }, "entity": { "type": "audit.app.ssh-authorized", "actor": "uaa-id-60", "actor_type": "user", "actor_name": "user@example.com", "actee": "95fe1035-300b-4acd-b092-3e8a07b67648", "actee_type": "app", "actee_name": "name-550", "timestamp": "2016-06-08T16:41:24Z", "metadata": { "index": 1 }, "space_guid": "dc008557-9c91-4e5c-9fec-5a46c49c9c7c", "organization_guid": "4e9f2aed-cb5e-42b7-b337-e3df877668ca" } } ] }`,
	"audit.app.ssh-unauthorized": `{ "total_results": 1, "total_pages": 1, "prev_url": null, "next_url": null, "resources": [ { "metadata": { "guid": "ccd2dbe9-356a-46b7-8c15-9946b16ff467", "url": "/v2/events/ccd2dbe9-356a-46b7-8c15-9946b16ff467", "created_at": "2016-06-08T16:41:27Z", "updated_at": "2016-06-08T16:41:26Z" }, "entity": { "type": "audit.app.ssh-unauthorized", "actor": "uaa-id-107", "actor_type": "user", "actor_name": "user@example.com", "actee": "b7acfab0-f241-49bc-9562-0d801d47aa52", "actee_type": "app", "actee_name": "name-1221", "timestamp": "2016-06-08T16:41:27Z", "metadata": { "index": 1 }, "space_guid": "95cbac34-c6ee-4efc-a963-9f6194044a4d", "organization_guid": "48e7ab86-b057-4e82-84c9-8c43766fae34" } } ] }`,
	"audit.app.start":            `{ "total_results": 1, "total_pages": 1, "prev_url": null, "next_url": null, "resources": [ { "metadata": { "guid": "6d5aa0c7-dcb4-4fa2-807e-f3c884f694a8", "url": "/v2/events/6d5aa0c7-dcb4-4fa2-807e-f3c884f694a8", "created_at": "2016-06-08T16:41:24Z", "updated_at": "2016-06-08T16:41:26Z" }, "entity": { "type": "audit.app.start", "actor": "uaa-id-48", "actor_type": "user", "actor_name": "user@example.com", "actee": "guid-29b0e354-599a-4582-a1b2-4a03d25b5f86", "actee_type": "v3-app", "actee_name": "name-383", "timestamp": "2016-06-08T16:41:24Z", "metadata": { }, "space_guid": "148a25e3-f4c7-4d22-af7f-490717c23e17", "organization_guid": "cdafb6a0-95b0-4426-8454-482685dd6073" } } ] }`,
	"audit.app.stop":             `{ "total_results": 1, "total_pages": 1, "prev_url": null, "next_url": null, "resources": [ { "metadata": { "guid": "1698d26a-39bb-4210-8036-3296ed93fed1", "url": "/v2/events/1698d26a-39bb-4210-8036-3296ed93fed1", "created_at": "2016-06-08T16:41:27Z", "updated_at": "2016-06-08T16:41:26Z" }, "entity": { "type": "audit.app.stop", "actor": "uaa-id-109", "actor_type": "user", "actor_name": "user@example.com", "actee": "guid-278aa447-6947-4a8b-a4ac-270bf5edec84", "actee_type": "v3-app", "actee_name": "name-1251", "timestamp": "2016-06-08T16:41:27Z", "metadata": { }, "space_guid": "4a10b249-1b42-4351-85de-a96ece28191a", "organization_guid": "67ee124e-0ad4-4715-9751-1d6ecafb7199" } } ] }`,
	"audit.app.update":           `{ "total_results": 1, "total_pages": 1, "prev_url": null, "next_url": null, "resources": [ { "metadata": { "guid": "c21ef683-a8ff-417c-a551-3d267bfc23f4", "url": "/v2/events/c21ef683-a8ff-417c-a551-3d267bfc23f4", "created_at": "2016-06-08T16:41:27Z", "updated_at": "2016-06-08T16:41:26Z" }, "entity": { "type": "audit.app.update", "actor": "uaa-id-115", "actor_type": "user", "actor_name": "user@example.com", "actee": "3eba670b-d6af-4d4a-ad61-5763df7e0400", "actee_type": "app", "actee_name": "name-1331", "timestamp": "2016-06-08T16:41:27Z", "metadata": { "request": { "name": "new", "instances": 1, "memory": 84, "state": "STOPPED", "environment_json": "PRIVATE DATA HIDDEN", "docker_credentials_json": "PRIVATE DATA HIDDEN" } }, "space_guid": "c5080860-98dd-4ffe-b3ea-3260d408a3d3", "organization_guid": "43442266-5e2a-46ca-9f85-53a5b44c2b9b" } } ] }`,
}

func TestTypeMap(t *testing.T) {
	for idx := range supportedEvents {
		t.Run(supportedEvents[idx], func(t *testing.T) {
			t.Logf("testing %s json extraction", supportedEvents[idx])
			var ev AppEvent
			if err := json.Unmarshal([]byte(typeMap[supportedEvents[idx]]), &ev); err != nil {
				t.Log(err)
				t.Fail()
			}
		})
	}
}
