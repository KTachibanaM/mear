package host

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/KTachibanaM/mear/internal/agent"
	"github.com/KTachibanaM/mear/internal/bucket"
	"github.com/KTachibanaM/mear/internal/do"
	"github.com/KTachibanaM/mear/internal/engine"
	"github.com/KTachibanaM/mear/internal/s3"
	"github.com/KTachibanaM/mear/internal/utils"
	"github.com/joho/godotenv"

	log "github.com/sirupsen/logrus"
)

func Host(upload_filename, save_to_filename string, skip_deprovision_engine, skip_deprovision_buckets bool) error {
	input_ext, err := utils.InferExt(upload_filename)
	if err != nil {
		return fmt.Errorf("could not infer ext from upload filename: %v", err)
	}

	// 0. Load credentials
	err = godotenv.Load()
	if err != nil {
		return fmt.Errorf("could not load .env file: %v", err)
	}
	access_key_id, secret_access_key, err := bucket.GetDigitalOceanSpacesCredentialsFromEnv()
	if err != nil {
		return fmt.Errorf("could not get DigitalOcean Spaces credentials from env: %v", err)
	}
	do_token, err := do.GetDigitalOceanTokenFromEnv()
	if err != nil {
		return fmt.Errorf("could not get DigitalOcean token from env: %v", err)
	}

	// 1. Get agent binary url
	log.Println("getting agent binary url...")
	agent_binary_url, err := NewGithubAgentBinaryURLRetriever().RetrieveUrl()
	if err != nil {
		return fmt.Errorf("could not get agent binary url: %v", err)
	}

	// 2. Provision buckets
	log.Println("provisioning buckets...")
	bucket_name, err := utils.GetRandomName("mear-s3", bucket.DigitalOceanSpacesBucketSuffixLength, bucket.DigitalOceanSpacesBucketNameMaxLength)
	if err != nil {
		return fmt.Errorf("could not generate random string for bucket name: %v", err)
	}
	do_dc_picker := do.NewStaticDigitalOceanDataCenterPicker("sfo3")
	do_spaces_session := bucket.NewDigitalOceanSpacesS3Session(
		do_dc_picker,
		access_key_id,
		secret_access_key,
	)

	s3_bucket := s3.NewS3Bucket(do_spaces_session, bucket_name)
	source_target := s3.NewS3Target(s3_bucket, fmt.Sprintf("input.%s", input_ext))
	destination_target := s3.NewS3Target(s3_bucket, "output.mp4")

	bucket_provisioner := bucket.NewMultiBucketProvisioner()
	err = bucket_provisioner.Provision(
		[]*s3.S3Bucket{s3_bucket},
	)
	if err != nil {
		bucket_teardown_err := bucket_provisioner.Teardown()
		return utils.CombineErrors(err, bucket_teardown_err)
	}

	// 3. Upload file
	log.Println("uploading file...")
	err = s3.UploadFile(upload_filename, source_target, true)
	if err != nil {
		bucket_teardown_err := bucket_provisioner.Teardown()
		return utils.CombineErrors(err, bucket_teardown_err)
	}

	// 4. Gather agent args
	log.Println("gathering agent args...")
	agent_args := agent.NewAgentArgs(
		source_target,
		destination_target,
		[]string{},
	)
	agent_args_json, err := json.MarshalIndent(agent_args, "", "")
	if err != nil {
		bucket_teardown_err := bucket_provisioner.Teardown()
		return utils.CombineErrors(err, bucket_teardown_err)
	}
	encoded := base64.StdEncoding.EncodeToString(agent_args_json)

	// 5. Provision engine
	log.Println("provisioning engine...")
	droplet_name, err := utils.GetRandomName("mear-engine", engine.DigitalOceanDropletSuffixLength, engine.DigitalOceanDropletNameMaxLength)
	if err != nil {
		return fmt.Errorf("could not generate random string for droplet name: %v", err)
	}
	engine_provisioner := engine.NewDigitalOceanEngineProvisioner(
		do_token,
		do_dc_picker,
		droplet_name,
		"s-1vcpu-512mb-10gb",
		"ubuntu-22-04-x64",
	)
	err = engine_provisioner.Provision(agent_binary_url, encoded)
	if err != nil {
		engine_teardown_err := engine_provisioner.Teardown()
		bucket_teardown_err := bucket_provisioner.Teardown()
		return utils.CombineErrors(err, engine_teardown_err, bucket_teardown_err)
	}

	// 6. Run agent on engine
	// TODO
	result := false

	// 7. Deprovision engine
	if !skip_deprovision_engine {
		log.Println("deprovisioning engine...")
		err = engine_provisioner.Teardown()
		if err != nil {
			bucket_teardown_err := bucket_provisioner.Teardown()
			return utils.CombineErrors(err, bucket_teardown_err)
		}
	} else {
		log.Warnln("skipped deprovisioning engine. you might want to deprovision manually.")
	}

	// 8. Download file
	if result {
		log.Println("downloading file...")
		err = s3.DownloadFile(save_to_filename, destination_target, true)
		if err != nil {
			bucket_teardown_err := bucket_provisioner.Teardown()
			return utils.CombineErrors(err, bucket_teardown_err)
		}
	} else {
		log.Warnln("failed to run agent. skipped downloading file.")
	}

	// 9. Deprovision buckets
	if !skip_deprovision_buckets {
		log.Println("deprovisioning buckets...")
		err = bucket_provisioner.Teardown()
		if err != nil {
			return fmt.Errorf("failed to teardown buckets: %v", err)
		}
	} else {
		log.Warnln("skipped deprovisioning buckets. you might want to deprovision manually.")
	}

	return nil
}
