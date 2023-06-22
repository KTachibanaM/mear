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

func Host() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("could not load .env file: %v", err)
	}

	// 1. Get agent binary url
	log.Println("getting agent binary url...")
	agent_binary_url := "http://minio:9000/bin/mear-agent"

	// 2. Provision buckets
	logs_bucket_name, err := bucket.GetDigitalOceanSpacesBucketName("mear-logs")
	if err != nil {
		return fmt.Errorf("could not generate random string for logs bucket name: %v", err)
	}
	access_key_id, secret_access_key, err := bucket.GetDigitalOceanSpacesCredentialsFromEnv()
	if err != nil {
		return fmt.Errorf("could not get DigitalOcean Spaces credentials from env: %v", err)
	}
	log.Println("provisioning buckets...")
	source_target := s3.NewS3Target(
		s3.NewS3Bucket(
			bucket.DevContainerS3Session,
			"source",
		),
		"MakeMine1948_256kb.rm",
	)
	destination_session := bucket.DevContainerS3Session
	destination_bucket := s3.NewS3Bucket(destination_session, "destination")
	destination_target := s3.NewS3Target(destination_bucket, "output.mp4")
	logs_session := bucket.NewDigitalOceanSpacesS3Session(
		do.NewStaticDigitalOceanDataCenterPicker("nyc3"),
		access_key_id, secret_access_key,
	)
	logs_bucket := s3.NewS3Bucket(logs_session, logs_bucket_name)
	logs_target := s3.NewS3Target(logs_bucket, "agent.log")

	bucket_provisioner := bucket.NewMultiBucketProvisioner()
	err = bucket_provisioner.Provision(
		[]*s3.S3Bucket{destination_bucket, logs_bucket},
	)
	if err != nil {
		bucket_teardown_err := bucket_provisioner.Teardown()
		return utils.CombineErrors(err, bucket_teardown_err)
	}

	// 3. Gather agent args
	log.Println("gathering agent args...")
	agent_args := agent.NewAgentArgs(
		source_target,
		destination_target,
		logs_target,
		[]string{},
	)
	agent_args_json, err := json.MarshalIndent(agent_args, "", "")
	if err != nil {
		bucket_teardown_err := bucket_provisioner.Teardown()
		return utils.CombineErrors(err, bucket_teardown_err)
	}
	encoded := base64.StdEncoding.EncodeToString(agent_args_json)

	// 4. Provision engine
	log.Println("provisioning engine...")
	engine_provisioner := engine.NewDevcontainerEngineProvisioner()
	err = engine_provisioner.Provision(agent_binary_url, encoded)
	if err != nil {
		engine_teardown_err := engine_provisioner.Teardown()
		bucket_teardown_err := bucket_provisioner.Teardown()
		return utils.CombineErrors(err, engine_teardown_err, bucket_teardown_err)
	}

	// 5. Tail for logs and result
	log.Println("tailing for logs and result...")
	result := NewS3LogsTailer(logs_target).Tail()
	if result {
		log.Println("agent run succeeded")
	} else {
		log.Println("agent run failed")
	}

	// 6. Deprovision engine
	log.Println("deprovisioning engine...")
	err = engine_provisioner.Teardown()
	if err != nil {
		bucket_teardown_err := bucket_provisioner.Teardown()
		return utils.CombineErrors(err, bucket_teardown_err)
	}

	// 7. Deprovision buckets
	log.Println("deprovisioning buckets...")
	err = bucket_provisioner.Teardown()
	if err != nil {
		return fmt.Errorf("failed to teardown buckets: %v", err)
	}

	return nil
}
