package cmd

import (
	"fmt"
	"os"

	"github.com/PauSabatesC/congo/congo"

	"github.com/PauSabatesC/congo/tui"

	tlog "github.com/Alana-Research/terminal-app-log"
	"github.com/spf13/cobra"
)

var id []string

func init() {
	rootCmd.AddCommand(ec2)
	rootCmd.AddCommand(ecs)

	ec2.Flags().StringArrayVar(
		&id,
		"id",
		nil,
		"Specify the EC2 ID to connect to.",
	)
}

var ec2 = &cobra.Command{
	Use:   "ec2",
	Short: "Options to AWS EC2 connect.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && id == nil {
			creds, err := loginAWS()
			if err != nil {
				loginError(err)
			}

			instances, err := congo.AwsListInstances(creds)
			if err != nil {
				tlog.BigError("Couldn't list ec2 instances:", err)
				os.Exit(1)
			}

			instanceId, err := tui.SelectElementFromEC2(instances, "Select an EC2 instance")
			if err != nil {
				tlog.BigError("An error occurred getting the instances list value:", err)
				os.Exit(1)
			}

			connectToEc2(creds, instanceId)
			os.Exit(0)
		} else if id != nil {
			creds, err := loginAWS()
			if err != nil {
				loginError(err)
			}
			connectToEc2(creds, id[0])
		}
	},
}

var ecs = &cobra.Command{
	Use:   "ecs",
	Short: "Connect to an ECS container.",
	Run: func(cmd *cobra.Command, args []string) {
		creds, err := loginAWS()
		if err != nil {
			loginError(err)
		}

		ecsClient := congo.NewECSClient(creds)

		clusters, err := congo.GetClusters(ecsClient)
		if err != nil {
			tlog.BigError("Couldn't get list of clusters:", err)
			os.Exit(1)
		}
		clusterName, err := tui.SelectElementFromList(clusters, "Select an ECS cluster")
		if err != nil {
			tlog.BigError("An error occurred getting the cluster list value:", err)
			os.Exit(1)
		}

		if tui.Exited {
			os.Exit(0)
		}

		services, err := congo.GetServices(ecsClient, clusterName)
		if err != nil {
			tlog.BigError(fmt.Sprintf("Couldn't get list of ECS services for cluster %s:", clusterName), err)
			os.Exit(1)
		}
		serviceName, err := tui.SelectElementFromList(services, "Select an ECS service")
		if err != nil {
			tlog.BigError("An error occurred getting the service list value:", err)
			os.Exit(1)
		}

		if tui.Exited {
			os.Exit(0)
		}

		tasks, err := congo.GetTasks(ecsClient, clusterName, serviceName)
		if err != nil {
			tlog.BigError("Couldn't get list of tasks:", err)
			os.Exit(1)
		}
		taskObject, err := tui.SelectElementFromEcsTask(tasks, "Select an ECS task")
		if err != nil {
			tlog.BigError("An error occurred getting the task list value:", err)
			os.Exit(1)
		}

		if tui.Exited {
			os.Exit(0)
		}

		containers, err := congo.GetContainers(*taskObject)
		if err != nil {
			tlog.BigError("Couldn't get list of containers:", err)
			os.Exit(1)
		}
		containerName, err := tui.SelectElementFromList(containers, "Select an ECS container")
		if err != nil {
			tlog.BigError("An error occurred getting the container list value:", err)
			os.Exit(1)
		}

		if tui.Exited {
			os.Exit(0)
		}

		err = congo.AwsECSConnect(creds, taskObject.Id, clusterName, containerName)
		if err != nil {
			tlog.BigError("Couldn't connect to ECS container:", err)
			os.Exit(1)
		}
	},
}

func loginAWS() (*congo.AwsCredentials, error) {
	creds := congo.NewAwsCredentials()
	err := creds.LoginWithEnvProfile()
	if err != nil {
		return nil, err
	}

	err = creds.GetLoginInfo()
	if err != nil {
		return nil, err
	}

	tlog.Info("Successfully logged:")
	fmt.Println("	Account:", creds.Account)
	fmt.Println("	Region:", creds.Region)
	fmt.Println("	Role:", creds.Sts_arn)

	return creds, nil
}

func loginError(err error) {
	tlog.BigError("Couldn't login to AWS:", err)
	tlog.Info("Remember to export the AWS region. E.g.-> export AWS_REGION=\"us-east-1\"\n")
	os.Exit(1)
}

func connectToEc2(creds *congo.AwsCredentials, id string) {
	err := congo.AwsEC2Connect(creds, id)
	if err != nil {
		tlog.BigError("Couldn't connect to EC2 isntance:", err)
		tlog.Info("Remember to:")
		fmt.Println("    -Have \"aws\" cli correctly installed and added on PATH.")
		fmt.Println("    -Have SSM plugin installed")
		fmt.Println("    -EC2 has the correct IAM profile and network so it's possible to connect using AWS SSM.")
		os.Exit(1)
	}
}
