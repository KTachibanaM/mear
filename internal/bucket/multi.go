package bucket

import (
	"github.com/KTachibanaM/mear/internal/s3"
	"github.com/hashicorp/go-multierror"
)

type MultiBucketProvisioner struct {
	provisioners       []BucketProvisioner
	provisioned_states []bool
}

func NewMultiBucketProvisioner(provisioners ...BucketProvisioner) *MultiBucketProvisioner {
	return &MultiBucketProvisioner{
		provisioners:       provisioners,
		provisioned_states: make([]bool, len(provisioners)),
	}
}

func (p *MultiBucketProvisioner) Provision() ([]*s3.S3Target, error) {
	var s3_targets []*s3.S3Target
	var errors error
	for i, provisioner := range p.provisioners {
		s3_target, err := provisioner.Provision()
		p.provisioned_states[i] = err != nil
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			s3_targets = append(s3_targets, s3_target)
		}
	}
	if errors != nil {
		return nil, errors
	}
	return s3_targets, nil
}

func (p *MultiBucketProvisioner) Teardown() error {
	var errors error
	for i, provisioner := range p.provisioners {
		if p.provisioned_states[i] {
			if err := provisioner.Teardown(); err != nil {
				errors = multierror.Append(errors, err)
			}
		}
	}
	return errors
}
