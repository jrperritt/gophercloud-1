package snapshotcommands

import (
	"strings"

	"github.com/fatih/structs"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/snapshots"
)

func snapshotSingle(snapshot *snapshots.Snapshot) map[string]interface{} {
	m := structs.Map(snapshot)
	if attachments := m["Attachments"].([]string); len(attachments) > 0 {
		m["Attachments"] = strings.Join(attachments, ",")
	} else {
		m["Attachments"] = ""
	}
	return m
}
