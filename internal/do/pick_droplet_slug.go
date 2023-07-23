package do

import "fmt"

var DropletRamAndCpuToSlug = map[int]map[int]string{
	1: {
		1: "s-1vcpu-1gb",
	},
	2: {
		1: "s-1vcpu-2gb",
		2: "s-2vcpu-2gb",
	},
	4: {
		2: "s-2vcpu-4gb",
	},
	8: {
		4: "s-4vcpu-8gb",
	},
}

func PickDropletSlug(ram, cpu int) (string, error) {
	if _, ok := DropletRamAndCpuToSlug[ram]; !ok {
		return "", fmt.Errorf("invalid droplet ram size: %d", ram)
	}

	if _, ok := DropletRamAndCpuToSlug[ram][cpu]; !ok {
		return "", fmt.Errorf("invalid droplet cpu count: %d", cpu)
	}

	return DropletRamAndCpuToSlug[ram][cpu], nil
}
