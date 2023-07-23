package do

import "fmt"

var DropletSizesToSlug = map[int]map[int]string{
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
}

func PickDropletSlug(ram, cpu int) (string, error) {
	if _, ok := DropletSizesToSlug[ram]; !ok {
		return "", fmt.Errorf("invalid droplet ram size: %d", ram)
	}

	if _, ok := DropletSizesToSlug[ram][cpu]; !ok {
		return "", fmt.Errorf("invalid droplet cpu count: %d", cpu)
	}

	return DropletSizesToSlug[ram][cpu], nil
}
