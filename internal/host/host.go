package host

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/KTachibanaM/mear/internal/agent"
	"github.com/KTachibanaM/mear/internal/agent_bin"
	"github.com/KTachibanaM/mear/internal/bucket"
	"github.com/KTachibanaM/mear/internal/do"
	"github.com/KTachibanaM/mear/internal/engine"
	"github.com/KTachibanaM/mear/internal/s3"
	"github.com/KTachibanaM/mear/internal/ssh"
	"github.com/KTachibanaM/mear/internal/utils"

	log "github.com/sirupsen/logrus"
)

type JobOutput struct {
	DestinationTarget *s3.S3Target `json:"destination_target"`
	ExtraFfmpegArgs   []string     `json:"extra_ffmpeg_args"`
}

type Job struct {
	InputFile string       `json:"input_file"`
	Outputs   []*JobOutput `json:"outputs"`
}

func Host(
	jobs []*Job,
	agent_execution_timeout_minutes int,
	stack string,
	retain_engine,
	retain_buckets bool,
	droplet_ram,
	droplet_cpu int,
	do_access_key_id,
	do_secret_access_key,
	do_token string,
) (*agent.AgentFailure, error) {
	// 0. Pre-validations
	private_key, public_key, err := ssh.Keygen()
	if err != nil {
		return nil, fmt.Errorf("could not generate ssh key pair: %v", err)
	}
	var input_exts []string
	for _, job := range jobs {
		ext, err := utils.InferExt(job.InputFile)
		if err != nil {
			return nil, fmt.Errorf("could not infer ext from input filename %v: %v", job.InputFile, err)
		}
		input_exts = append(input_exts, ext)
	}

	var do_bucket_name string
	var droplet_name string
	var droplet_slug string
	if stack == "do" {
		do_bucket_name, err = utils.GetRandomName("mear-s3", bucket.DigitalOceanSpacesBucketSuffixLength, bucket.DigitalOceanSpacesBucketNameMaxLength)
		if err != nil {
			return nil, fmt.Errorf("could not generate random string for bucket name: %v", err)
		}
		droplet_name, err = utils.GetRandomName("mear-engine", engine.DigitalOceanDropletSuffixLength, engine.DigitalOceanDropletNameMaxLength)
		if err != nil {
			return nil, fmt.Errorf("could not generate random string for droplet name: %v", err)
		}
		droplet_slug, err = do.PickDropletSlug(droplet_ram, droplet_cpu)
		if err != nil {
			return nil, fmt.Errorf("could not pick droplet slug: %v", err)
		}
	}

	// 1. Get agent binary url
	log.Println("getting agent binary url...")
	var ab agent_bin.AgentBinary
	if stack == "dev" {
		ab = agent_bin.NewDevContainerAgentBinary()
	} else if stack == "do" {
		ab = agent_bin.NewGithubAgentBinary()
	} else {
		return nil, fmt.Errorf("unknown stack name %v", stack)
	}
	agent_binary_url, err := ab.RetrieveUrl()
	if err != nil {
		return nil, fmt.Errorf("could not get agent binary url: %v", err)
	}

	// 2. Provision buckets
	log.Println("provisioning buckets...")
	var s3_session *s3.S3Session
	var bucket_name string
	if stack == "dev" {
		s3_session = bucket.DevContainerS3Session
		bucket_name = "mear-dev"
	} else if stack == "do" {
		do_dc_picker := do.NewStaticDigitalOceanDataCenterPicker("sfo3")
		s3_session = bucket.NewDigitalOceanSpacesS3Session(
			do_dc_picker,
			do_access_key_id,
			do_secret_access_key,
		)
		bucket_name = do_bucket_name
	} else {
		return nil, fmt.Errorf("unknown stack name %v", stack)
	}

	bucket_provisioner := bucket.NewMultiBucketProvisioner()
	s3_bucket := s3.NewS3Bucket(s3_session, bucket_name)
	err = bucket_provisioner.Provision(
		[]*s3.S3Bucket{s3_bucket},
	)
	if err != nil {
		bucket_teardown_err := bucket_provisioner.Teardown()
		return nil, utils.CombineErrors(err, bucket_teardown_err)
	}
	var source_targets []*s3.S3Target
	for i, input_ext := range input_exts {
		source_targets = append(source_targets, s3.NewS3Target(s3_bucket, fmt.Sprintf("%v.in.%v", i, input_ext)))
	}

	// 3. Upload videos
	for i, job := range jobs {
		log.Printf("uploading video %v ...", job.InputFile)
		err = s3.UploadFile(job.InputFile, source_targets[i], true)
		if err != nil {
			bucket_teardown_err := bucket_provisioner.Teardown()
			return nil, utils.CombineErrors(err, bucket_teardown_err)
		}
	}

	// 4. Gather agent args
	log.Println("gathering agent args...")
	agent_args := agent.NewAgentArgs()
	for i, job := range jobs {
		agent_job := agent.NewAgentJob(source_targets[i])
		for _, output := range job.Outputs {
			agent_job.AddJobDestination(
				output.DestinationTarget,
				output.ExtraFfmpegArgs,
			)
		}
		agent_args.AddJob(agent_job)
	}
	agent_args_json, err := json.MarshalIndent(agent_args, "", "")
	if err != nil {
		bucket_teardown_err := bucket_provisioner.Teardown()
		return nil, utils.CombineErrors(err, bucket_teardown_err)
	}
	encoded_agent_args := base64.StdEncoding.EncodeToString(agent_args_json)

	// 5. Provision engine
	log.Println("provisioning engine...")
	var engine_provisioner engine.EngineProvisioner
	if stack == "dev" {
		engine_provisioner = engine.NewDevcontainerEngineProvisioner()
	} else if stack == "do" {
		do_dc_picker := do.NewStaticDigitalOceanDataCenterPicker("sfo3")
		engine_provisioner = engine.NewDigitalOceanEngineProvisioner(
			do_token,
			do_dc_picker,
			droplet_name,
			droplet_slug,
			"debian-11-x64",
		)
	} else {
		return nil, fmt.Errorf("unknown stack name %v", stack)
	}

	ip_address, err := engine_provisioner.Provision(agent_binary_url, public_key)
	if err != nil {
		engine_teardown_err := engine_provisioner.Teardown()
		bucket_teardown_err := bucket_provisioner.Teardown()
		return nil, utils.CombineErrors(err, engine_teardown_err, bucket_teardown_err)
	}

	// 6. Run agent on engine
	var ssh_status *ssh.SshStatus
	commands := []string{
		"apt update",
		"apt install -y curl",
		"curl --fail -sL " + agent_binary_url + " -o /root/mear-agent",
		"chmod +x /root/mear-agent",
		"/root/mear-agent " + encoded_agent_args,
	}
	for _, command := range commands {
		if !strings.Contains(command, encoded_agent_args) {
			log.Printf("ssh executing command: '%v'", command)
		} else {
			log.Printf("ssh executing command: '%v'", strings.Replace(command, encoded_agent_args, "<agent args redacted>", -1))
		}
		is_mear_agent := strings.Contains(command, "/root/mear-agent")
		timeout := 1 * time.Minute
		if is_mear_agent {
			timeout = time.Duration(agent_execution_timeout_minutes) * time.Minute
		}
		ssh_status = ssh.SshExec(ip_address, "root", private_key, command, timeout)
		if ssh_status.Status != ssh.SshSuccess {
			log.Errorf("failed to ssh execute command: %v", err)
			break
		}
	}

	// 7. Deprovision engine
	if !retain_engine {
		log.Println("deprovisioning engine...")
		err = engine_provisioner.Teardown()
		if err != nil {
			bucket_teardown_err := bucket_provisioner.Teardown()
			return nil, utils.CombineErrors(err, bucket_teardown_err)
		}
	} else {
		log.Warnln("retaining engine. you might want to deprovision manually.")
	}

	// 8. Deprovision buckets
	if !retain_buckets {
		log.Println("deprovisioning buckets...")
		err = bucket_provisioner.Teardown()
		if err != nil {
			return nil, fmt.Errorf("failed to teardown buckets: %v", err)
		}
	} else {
		log.Warnln("retaining buckets. you might want to deprovision manually.")
	}

	if ssh_status.Status == ssh.SshSuccess {
		return nil, nil
	} else if ssh_status.Status == ssh.SshError {
		return nil, fmt.Errorf("failed to ssh execute command: %v", ssh_status.Err)
	} else {
		return ssh_status.AgentFailure, nil
	}
}
