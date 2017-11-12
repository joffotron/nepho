package cloudformation

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/k0kubun/pp"
)

type StackDefinition string

type Cfn struct {
	stackName string
	cfn       cloudformationiface.CloudFormationAPI
}

func New(stackName string) *Cfn {
	conn := session.Must(session.NewSession())
	return &Cfn{
		stackName: stackName,
		cfn:       cloudformation.New(conn),
	}
}

func (c *Cfn) Create(template string) error {
	capabilities := []*string{
		aws.String("CAPABILITY_IAM"),
		aws.String("CAPABILITY_NAMED_IAM"),
	}

	input := &cloudformation.CreateStackInput{
		StackName:    &c.stackName,
		TemplateBody: &template,
		Capabilities: capabilities,
	}
	out, err := c.cfn.CreateStack(input)
	if err != nil {
		return err
	}
	pp.Println(out)
	return nil
}
