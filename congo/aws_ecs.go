package congo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tlog "github.com/Alana-Research/terminal-app-log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type EcsTaskData struct {
	Id         string
	StartedAt  string
	Containers []string
}

func NewECSClient(creds *AwsCredentials) *ecs.Client {
	ecsClient := ecs.NewFromConfig(creds.Config)
	return ecsClient
}

func GetClusters(ecsClient *ecs.Client) ([]string, error) {
	input := &ecs.ListClustersInput{
		MaxResults: aws.Int32(100),
	}

	paginator := ecs.NewListClustersPaginator(ecsClient, input, func(o *ecs.ListClustersPaginatorOptions) {
		o.Limit = 100 //Limited by AWS API
	})

	var clusters []string
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, errors.New(err.Error())
		}
		clusters = append(clusters, output.ClusterArns...)
	}

	if len(clusters) == 0 {
		return nil, errors.New("No clusters found.")
	}

	var clusterNames []string
	for _, arn := range clusters {
		arnSplit := strings.Split(arn, "/")
		name := arnSplit[len(arnSplit)-1]
		clusterNames = append(clusterNames, name)
	}
	return clusterNames, nil
}

func GetServices(ecsClient *ecs.Client, clusterName string) ([]string, error) {
	input := &ecs.ListServicesInput{
		Cluster:    aws.String(clusterName),
		MaxResults: aws.Int32(100),
	}

	paginator := ecs.NewListServicesPaginator(ecsClient, input, func(o *ecs.ListServicesPaginatorOptions) {
		o.Limit = 100 //limited by aws
	})

	var services []string
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, errors.New(err.Error())
		}
		services = append(services, output.ServiceArns...)
	}

	if len(services) == 0 {
		return nil, errors.New("No services found.")
	}

	var servicesNames []string
	for _, arn := range services {
		arnSplit := strings.Split(arn, "/")
		name := arnSplit[len(arnSplit)-1]
		servicesNames = append(servicesNames, name)
	}
	return servicesNames, nil
}

func GetTasks(ecsClient *ecs.Client, clusterName string, serviceName string) ([]EcsTaskData, error) {
	input := &ecs.ListTasksInput{
		Cluster:       aws.String(clusterName),
		ServiceName:   aws.String(serviceName),
		DesiredStatus: types.DesiredStatusRunning,
		MaxResults:    aws.Int32(100),
	}

	paginator := ecs.NewListTasksPaginator(ecsClient, input, func(o *ecs.ListTasksPaginatorOptions) {
		o.Limit = 100 //AWS limits to 100 the DescribeTasksInput slice lmao
	})

	var tasksData []types.Task
	for paginator.HasMorePages() {
		var tasks []string
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, errors.New(err.Error())
		}
		tasks = append(tasks, output.TaskArns...)

		if len(tasks) == 0 {
			return nil, errors.New("No tasks found.")
		}

		describeInput := &ecs.DescribeTasksInput{
			Cluster: aws.String(clusterName),
			Tasks:   tasks, //AWS limits to 100 that slice
		}

		tasksAux, err := ecsClient.DescribeTasks(context.TODO(), describeInput)
		if err != nil {
			return nil, err
		}

		if len(tasksAux.Tasks) == 0 {
			return nil, errors.New("No tasks found.")
		}

		tasksData = append(tasksData, tasksAux.Tasks...)
	}

	var tasksObject []EcsTaskData
	for _, taskData := range tasksData {
		newTask := EcsTaskData{}
		taskDay := strings.Split((*taskData.StartedAt).String(), " ")[0]
		taskHour := strings.Split((*taskData.StartedAt).String(), " ")[1]
		taskTime := []string{taskDay, taskHour}
		newTask.StartedAt = strings.Join(taskTime, " ")
		newTask.Id = strings.Split(*taskData.TaskArn, "/")[2]

		containerNames := func() []string {
			var containerNames []string
			for _, container := range taskData.Containers {
				containerNames = append(containerNames, *container.Name)
			}
			return containerNames
		}
		newTask.Containers = containerNames()
		tasksObject = append(tasksObject, newTask)
	}
	return tasksObject, nil
}

func GetContainers(task EcsTaskData) ([]string, error) {
	if len(task.Containers) == 0 {
		return nil, errors.New(fmt.Sprintf("No containers running for task %s", task.Id))
	}

	return task.Containers, nil
}

func AwsECSConnect(creds *AwsCredentials, task_id string, cluster_id string, container string) error {
	tlog.BigInfo("Connecting to ECS... ")

	cmd := exec.Command("aws", "ecs", "execute-command",
		"--cluster", cluster_id,
		"--task", task_id,
		"--container", container,
		"--command", "/bin/sh",
		"--interactive",
	)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	tlog.BigInfo("Closed connection with ECS container ")
	return nil
}
