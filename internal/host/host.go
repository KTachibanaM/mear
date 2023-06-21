package host

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/KTachibanaM/mear/internal/agent"
	"github.com/KTachibanaM/mear/internal/bucket"
	"github.com/KTachibanaM/mear/internal/engine"

	log "github.com/sirupsen/logrus"
)

func Host() error {
	// 1. Get agent binary url
	log.Println("getting agent binary url...")
	agent_binary_url := "http://minio-agent-binary:9000/bin/mear-agent"

	// 2. Provision buckets
	log.Println("provisioning buckets...")
	bucket_provisioner := bucket.NewDevcontainerBucketProvisioner()
	s3_source, s3_destination, s3_logs, err := bucket_provisioner.Provision()
	if err != nil {
		return err
	}

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
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(agent_args_json)

	// 4. Provision engine
	log.Println("provisioning engine...")
	engine_provisioner := engine.NewDevcontainerEngineProvisioner()
	engine_id, err := engine_provisioner.Provision(agent_binary_url, encoded)
	if err != nil {
		engine_teardown_err := engine_provisioner.Teardown(engine_id)
		bucket_teardown_err := bucket_provisioner.Teardown()
		return fmt.Errorf("failed to provision engine: %v; engine teardown error :%v; bucket teardown err: %v", err, engine_teardown_err, bucket_teardown_err)
	}

	// 5. Tail for logs and result
	log.Println("tailing for logs and result...")
	result := TailS3Logs(s3_logs)
	if result {
		log.Println("agent run succeeded")
	} else {
		log.Println("agent run failed")
	}

	// 6. Deprovision engine
	log.Println("deprovisioning engine...")
	err = engine_provisioner.Teardown(engine_id)
	if err != nil {
		bucket_teardown_err := bucket_provisioner.Teardown()
		return fmt.Errorf("failed to teardown engine: %v; bucket tear down error: %v", err, bucket_teardown_err)
	}

	// 7. Deprovision buckets
	err = bucket_provisioner.Teardown()
	if err != nil {
		return fmt.Errorf("failed to teardown buckets: %v", err)
	}

	return nil
}
