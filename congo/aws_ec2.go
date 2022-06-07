package congo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	tlog "github.com/Alana-Research/terminal-app-log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type Ec2Data struct {
	Name      string
	Id        string
	PrivateIp string
	Vpc       string
	Type      string
	Platform  string
}

func AwsListInstances(creds *AwsCredentials) (*[]Ec2Data, error) {
	tlog.Info("Listing EC2 instances...")

	awsClient := ec2.NewFromConfig(creds.Config)

	running_filter := types.Filter{
		Name:   aws.String("instance-state-name"),
		Values: []string{string(types.InstanceStateNameRunning)},
	}

	filters := []types.Filter{running_filter}

	awsInput := &ec2.DescribeInstancesInput{
		Filters:    filters,
		MaxResults: aws.Int32(1000),
	}

	paginator := ec2.NewDescribeInstancesPaginator(awsClient, awsInput, func(o *ec2.DescribeInstancesPaginatorOptions) {
		o.Limit = 1000
	})

	var instancesResponse []types.Reservation
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, errors.New(err.Error())
		}
		instancesResponse = append(instancesResponse, output.Reservations...)
	}

	var running_instances []Ec2Data
	for _, instances := range instancesResponse {
		for _, instance := range instances.Instances {
			instanceData := Ec2Data{
				Id:        *instance.InstanceId,
				Type:      string(instance.InstanceType),
				Platform:  *instance.PlatformDetails,
				PrivateIp: *instance.PrivateIpAddress,
				Vpc:       *instance.VpcId,
			}
			for _, i := range instance.Tags {
				if *i.Key == "Name" {
					instanceData.Name = *i.Value
					break
				}
			}
			running_instances = append(running_instances, instanceData)
		}
	}

	return &running_instances, nil
}

func AwsEC2Connect(creds *AwsCredentials, id string) error {
	fmt.Println("")
	tlog.BigInfo("Connecting to EC2: ", id)

	cmd := exec.Command("aws", "ssm", "start-session", "--target", id)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	tlog.BigInfo("Closed connection with EC2 ", id)
	return nil
}
