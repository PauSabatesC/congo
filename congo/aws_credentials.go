package congo

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AwsCredentials struct {
	Config  aws.Config
	Ctx     context.Context
	Account string
	Sts_arn string
	Region  string
}

func (cred *AwsCredentials) LoginWithEnvProfile() error {
	cred.Ctx = context.TODO()
	cnf, err := config.LoadDefaultConfig(cred.Ctx)
	if err != nil {
		return errors.New(
			fmt.Sprintf("ERROR while trying to get AWS credentials: %s", err.Error()))
	}

	cred.Config = cnf
	return nil
}

func (cred *AwsCredentials) GetLoginInfo() error {
	client := sts.NewFromConfig(cred.Config)
	identity, err := client.GetCallerIdentity(cred.Ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return errors.New(
			fmt.Sprintf("ERROR while trying to connect to AWS: %s", err.Error()))
	}

	cred.Account = aws.ToString(identity.Account)
	cred.Sts_arn = aws.ToString(identity.Arn)
	cred.Region = cred.Config.Region
	return nil
}

func NewAwsCredentials() *AwsCredentials {
	cred := &AwsCredentials{}
	return cred
}
