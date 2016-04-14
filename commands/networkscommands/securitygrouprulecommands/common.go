package securitygrouprulecommands

import (
	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	osSecurityGroupRules "github.com/rackspace/rack/internal/github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
)

func securityGroupRuleSingle(rule *osSecurityGroupRules.SecGroupRule) map[string]interface{} {
	m := structs.Map(rule)

	m["SecurityGroupID"] = m["SecGroupID"]

	return m
}
