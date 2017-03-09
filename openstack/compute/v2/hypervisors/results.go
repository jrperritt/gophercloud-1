package hypervisors

import (
	"encoding/json"
	"fmt"

	"github.com/gophercloud/gophercloud/pagination"
)

type Topology struct {
	Sockets int `json:"sockets"`
	Cores   int `json:"cores"`
	Threads int `json:"threads"`
}

type CPUInfo struct {
	Vendor   string   `json:"vendor"`
	Arch     string   `json:"arch"`
	Model    string   `json:"model"`
	Features []string `json:"features"`
	Topology Topology `json:"topology"`
}

type Service struct {
	Host           string `json:"host"`
	ID             int    `json:"id"`
	DisabledReason string `json:"disabled_reason"`
}

type Hypervisor struct {
	// A structure that contains cpu information like arch, model, vendor, features and topology
	CPUInfo CPUInfo `json:"-"`
	// The current_workload is the number of tasks the hypervisor is responsible for.
	// This will be equal or greater than the number of active VMs on the system
	// (it can be greater when VMs are being deleted and the hypervisor is still cleaning up).
	CurrentWorkload int `json:"current_workload"`
	// Status of the hypervisor, either "enabled" or "disabled"
	Status string `json:"status"`
	// State of the hypervisor, either "up" or "down"
	State string `json:"state"`
	// Actual free disk on this hypervisor in GB
	DiskAvailableLeast int `json:"disk_available_least"`
	// The hypervisor's IP address
	HostIP string `json:"host_ip"`
	// The free disk remaining on this hypervisor in GB
	FreeDiskGB int `json:"free_disk_gb"`
	// The free RAM in this hypervisor in MB
	FreeRamMB int `json:"free_ram_mb"`
	// The hypervisor host name
	HypervisorHostname string `json:"hypervisor_hostname"`
	// The hypervisor type
	HypervisorType string `json:"hypervisor_type"`
	// The hypervisor version
	HypervisorVersion int `json:"-"`
	// Unique ID of the hypervisor
	ID int `json:"id"`
	// The disk in this hypervisor in GB
	LocalGB int `json:"local_gb"`
	// The disk used in this hypervisor in GB
	LocalGBUsed int `json:"local_gb_used"`
	// The memory of this hypervisor in MB
	MemoryMB int `json:"memory_mb"`
	// The memory used in this hypervisor in MB
	MemoryMBUsed int `json:"memory_mb_used"`
	// The number of running vms on this hypervisor
	RunningVMs int `json:"running_vms"`
	// The hypervisor service object
	Service Service `json:"service"`
	// The number of vcpu in this hypervisor
	VCPUs int `json:"vcpus"`
	// The number of vcpu used in this hypervisor
	VCPUsUsed int `json:"vcpus_used"`
}

func (r *Hypervisor) UnmarshalJSON(b []byte) error {

	type tmp Hypervisor
	var s struct {
		tmp
		CPUInfo           interface{} `json:"cpu_info"`
		HypervisorVersion interface{} `json:"hypervisor_version"`
	}

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = Hypervisor(s.tmp)

	var tmpb []byte

	switch t := s.CPUInfo.(type) {
	case string:
		tmpb = []byte(t)
	case map[string]interface{}:
		tmpb, err = json.Marshal(t)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("CPUInfo has unexpected type: %T", t)
	}

	err = json.Unmarshal(tmpb, &r.CPUInfo)
	if err != nil {
		return err
	}

	// A feature in OpenStack may return this value as floating point
	switch t := s.HypervisorVersion.(type) {
	case int:
		r.HypervisorVersion = t
	case float64:
		r.HypervisorVersion = int(t)
	default:
		return fmt.Errorf("HypervisorVersion has unexpected type: %T", t)
	}

	return nil
}

type HypervisorPage struct {
	pagination.LinkedPageBase
}

func (page HypervisorPage) IsEmpty() (bool, error) {
	hypervisors, err := ExtractHypervisors(page)
	return len(hypervisors) == 0, err
}

func ExtractHypervisors(p pagination.Page) ([]Hypervisor, error) {
	var h struct {
		Hypervisors []Hypervisor `json:"hypervisors"`
	}
	err := (p.(HypervisorPage)).ExtractInto(&h)
	return h.Hypervisors, err
}
