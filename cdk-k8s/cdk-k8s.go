package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/joho/godotenv"
	"os"
)

type CdkK8SStackProps struct {
	awscdk.StackProps
}

var (
	AccountRegion string
	VpcID         string
	AmiID         string
	SubnetID      string
	EC2KeyPair    string
	AccountID     string
)

func init() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		return
	}

	AccountID = os.Getenv("ACCOUNT_ID")
	AccountRegion = os.Getenv("ACCOUNT_REGION")
	VpcID = os.Getenv("VPC_ID")
	SubnetID = os.Getenv("SUBNET_ID")
	EC2KeyPair = os.Getenv("EC2_KEYPAIR")
	AmiID = os.Getenv("AMI_ID")
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(AccountID),
		Region:  jsii.String(AccountRegion),
	}
}

func GetUserData() string {
	script, _ := os.ReadFile("userData.sh")

	return string(script)
}

func NewRole(stack awscdk.Stack) awsiam.Role {
	k8sRole := awsiam.NewRole(stack, jsii.String("K8sRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewCompositePrincipal(
			awsiam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), nil),
			awsiam.NewServicePrincipal(jsii.String("cloudformation.amazonaws.com"), nil),
		),
	})

	k8sRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect:    awsiam.Effect_ALLOW,
		Actions:   jsii.Strings("sts:AssumeRole"),
		Resources: jsii.Strings("*"),
	}))

	return k8sRole
}

func CreateSecurityGroup(stack awscdk.Stack) awsec2.CfnSecurityGroup {
	props := awsec2.CfnSecurityGroupProps{
		GroupDescription: jsii.String("Desc"),
		VpcId:            jsii.String(VpcID),
		GroupName:        jsii.String("K8sSecurityGroup"),
		SecurityGroupIngress: []awsec2.CfnSecurityGroup_IngressProperty{
			{
				IpProtocol: jsii.String("TCP"),
				FromPort:   jsii.Number(22),
				ToPort:     jsii.Number(22),
				CidrIp:     jsii.String("0.0.0.0/0"),
			},
			{
				IpProtocol: jsii.String("TCP"),
				FromPort:   jsii.Number(80),
				ToPort:     jsii.Number(80),
				CidrIp:     jsii.String("0.0.0.0/0"),
			},
			{
				IpProtocol: jsii.String("TCP"),
				FromPort:   jsii.Number(443),
				ToPort:     jsii.Number(443),
				CidrIp:     jsii.String("0.0.0.0/0"),
			},
			{
				IpProtocol: jsii.String("TCP"),
				FromPort:   jsii.Number(6443),
				ToPort:     jsii.Number(6443),
				CidrIp:     jsii.String("0.0.0.0/0"),
			},
			{
				IpProtocol: jsii.String("TCP"),
				FromPort:   jsii.Number(2379),
				ToPort:     jsii.Number(2380),
				CidrIp:     jsii.String("0.0.0.0/0"),
			},
			{
				IpProtocol: jsii.String("TCP"),
				FromPort:   jsii.Number(10250),
				ToPort:     jsii.Number(10252),
				CidrIp:     jsii.String("0.0.0.0/0"),
			},
			{
				IpProtocol: jsii.String("TCP"),
				FromPort:   jsii.Number(6783),
				ToPort:     jsii.Number(6784),
				CidrIp:     jsii.String("0.0.0.0/0"),
			},
			{
				IpProtocol: jsii.String("TCP"),
				FromPort:   jsii.Number(2049),
				ToPort:     jsii.Number(2049),
				CidrIp:     jsii.String("0.0.0.0/0"),
			},
		},
	}

	K8sSecurityGroup := awsec2.NewCfnSecurityGroup(stack, jsii.String("SEC"), &props)

	return K8sSecurityGroup
}

func CreateEC2Instance(stack awscdk.Stack, secGroup awsec2.CfnSecurityGroup) awsec2.CfnInstance {
	userData := GetUserData()

	props := awsec2.CfnInstanceProps{
		ImageId:          jsii.String(AmiID),
		InstanceType:     jsii.String("t2.medium"),
		SecurityGroupIds: jsii.Strings(*secGroup.AttrGroupId()),
		SubnetId:         jsii.String(SubnetID),
		KeyName:          jsii.String(EC2KeyPair),
		UserData:         awscdk.Fn_Base64(&userData),
	}

	K8sInstance := awsec2.NewCfnInstance(stack, jsii.String("K8sInstance"), &props)

	return K8sInstance
}

func NewCdkK8SStack(scope constructs.Construct, id string, props *CdkK8SStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	stack := NewCdkK8SStack(app, "CdkK8SStack", &CdkK8SStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	NewRole(stack)
	secGroup := CreateSecurityGroup(stack)
	CreateEC2Instance(stack, secGroup)

	app.Synth(nil)
}
