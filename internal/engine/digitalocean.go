package engine

import (
	"context"
	"fmt"

	"github.com/KTachibanaM/mear/internal/do"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"

	log "github.com/sirupsen/logrus"
)

var DigitalOceanDropletSuffixLength = 8
var DigitalOceanDropletNameMaxLength = 20

type DigitalOceanEngineProvisioner struct {
	token              string
	do_dc              string
	droplet_name       string
	droplet_size       string
	droplet_image_slug string
	droplet_id         int
}

func NewDigitalOceanEngineProvisioner(token string, dc_picker do.DigitalOceanDataCenterPicker, droplet_name, droplet_size, droplet_image_slug string) *DigitalOceanEngineProvisioner {
	return &DigitalOceanEngineProvisioner{
		token:              token,
		do_dc:              dc_picker.Pick(),
		droplet_name:       droplet_name,
		droplet_size:       droplet_size,
		droplet_image_slug: droplet_image_slug,
	}
}

func (p *DigitalOceanEngineProvisioner) createClient() (*godo.Client, *context.Context) {
	ctx := context.TODO()
	token_source := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: p.token},
	)
	oauth_client := oauth2.NewClient(ctx, token_source)
	client := godo.NewClient(oauth_client)

	return client, &ctx
}

func (p *DigitalOceanEngineProvisioner) Provision(agent_binary_url, encoded_agent_args string) error {
	client, ctx := p.createClient()

	log.Printf("creating droplet %v...", p.droplet_name)
	create_request := &godo.DropletCreateRequest{
		Name:   p.droplet_name,
		Region: p.do_dc,
		Size:   p.droplet_size,
		Image: godo.DropletCreateImage{
			Slug: p.droplet_image_slug,
		},
		UserData: fmt.Sprintf(`#cloud-config
package_update: true
package_upgrade: true

packages:
    - curl

runcmd:
    - curl -sL %v -o /root/mear-agent
    - chmod +x /root/mear-agent
    - /root/mear-agent %v`, agent_binary_url, encoded_agent_args),
	}

	droplet, _, err := client.Droplets.Create(*ctx, create_request)
	if err != nil {
		return fmt.Errorf("failed to create droplet: %v", err)
	}

	p.droplet_id = droplet.ID

	return nil
}

func (p *DigitalOceanEngineProvisioner) Teardown() error {
	if p.droplet_id == 0 {
		return fmt.Errorf("droplet was never provisioned")
	}

	client, ctx := p.createClient()

	log.Printf("deleting droplet %v...", p.droplet_name)
	_, err := client.Droplets.Delete(*ctx, p.droplet_id)
	if err != nil {
		return fmt.Errorf("failed to request deleting droplet: %v", err)
	}

	return nil
}
