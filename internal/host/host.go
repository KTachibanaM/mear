package host

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/KTachibanaM/mear/internal/agent"
	"github.com/KTachibanaM/mear/internal/bucket"
	"github.com/KTachibanaM/mear/internal/engine"
	"github.com/KTachibanaM/mear/internal/utils"

	log "github.com/sirupsen/logrus"
)

func Host() error {
	// 1. Get agent binary url
	log.Println("getting agent binary url...")
	agent_binary_url := "http://minio-agent-binary:9000/bin/mear-agent"

	// 2. Provision buckets
	log.Println("provisioning buckets...")
	source_bucket_provisioner := bucket.NewNoOpBucketProvisioner(bucket.DevContainerSource, false)
	destination_bucket_provisioner := bucket.NewNoOpBucketProvisioner(bucket.DevContainerDestination, true)
	logs_bucket_provisioner := bucket.NewNoOpBucketProvisioner(bucket.DevContainerLogs, true)
	bucket_provisioner := bucket.NewMultiBucketProvisioner(
		source_bucket_provisioner,
		destination_bucket_provisioner,
		logs_bucket_provisioner,
	)
	s3_targets, err := bucket_provisioner.Provision()
	if err != nil {
		bucket_teardown_err := bucket_provisioner.Teardown()
		return utils.CombineErrors(err, bucket_teardown_err)
	}
	s3_source := s3_targets[0]
	s3_destination := s3_targets[1]
	s3_logs := s3_targets[2]

	// 3. Gather agent args
	log.Println("gathering agent args...")
	agent_args := agent.NewAgentArgs(
		s3_source,
		s3_destination,
		s3_logs,
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
	result := NewS3LogsTailer(s3_logs).Tail()
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
	err = bucket_provisioner.Teardown()
	if err != nil {
		return fmt.Errorf("failed to teardown buckets: %v", err)
	}

	return nil
}
