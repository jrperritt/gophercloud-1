package snapshotcommands

import (
	"strings"

	"github.com/gophercloud/cli/vendor/github.com/fatih/structs"
	osSnapshots "github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack/blockstorage/v1/snapshots"
)

func snapshotSingle(snapshot *osSnapshots.Snapshot) map[string]interface{} {
	m := structs.Map(snapshot)
	if attachments := m["Attachments"].([]string); len(attachments) > 0 {
		m["Attachments"] = strings.Join(attachments, ",")
	} else {
		m["Attachments"] = ""
	}
	return m
}
