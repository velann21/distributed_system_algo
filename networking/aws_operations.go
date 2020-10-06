package main

import (
	//"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	//"reflect"

	//"reflect"

	//"log"
)


type CreateVpcStruct struct {
	Tag []*Tag `json:"tags"`
	Name string `json:"name"`
	CredLocation string `json:"credLocation"`
	CredName string `json:"credName"`
	Region string `json:"region"`
	CIDR string `json:"cidr"`

}
type Tag struct{
	Key string `json:"key"`
	Value string `json:"value"`

}


func CreateVpc (createVpc CreateVpcStruct) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewSharedCredentials(createVpc.CredLocation, createVpc.CredName),
		Region: aws.String(createVpc.Region)},
	)
	if err != nil{
		return
	}
	val, err := sess.Config.Credentials.Get()
	if err != nil{
		fmt.Println(err)
	}
    fmt.Println(val)

	input := &ec2.CreateVpcInput{
		CidrBlock: aws.String(createVpc.CIDR),
		TagSpecifications: []*ec2.TagSpecification{
			&ec2.TagSpecification{
				ResourceType: aws.String(ec2.ResourceTypeVpc),
				Tags: FormTags(createVpc.Tag),
			}},
		}

	svc := ec2.New(sess)
	resu, err := svc.CreateVpc(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}
    fmt.Println(resu.Vpc.VpcId)
}


func CreateSubNet(){
	svc := ec2.New(session.New())
	input := &ec2.CreateSubnetInput{
		CidrBlock: aws.String("10.0.1.0/24"),
		VpcId:     aws.String("vpc-a01106c2"),
	}

	result, err := svc.CreateSubnet(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}


func FormTags(tags []*Tag)[]*ec2.Tag{
	ec2tags := []*ec2.Tag{}
	for _, value := range tags{
		kv := ec2.Tag{Key:aws.String(value.Key), Value:aws.String(value.Value)}
		ec2tags = append(ec2tags, &kv)
	}
	return ec2tags
}

func ExampleEC2_CreateVolume_shared01() {
	svc := ec2.New(session.New())
	input := &ec2.CreateVolumeInput{
		AvailabilityZone: aws.String("us-east-1a"),
		Iops:             aws.Int64(1000),
		SnapshotId:       aws.String("snap-066877671789bd71b"),
		VolumeType:       aws.String("io1"),
	}

	result, err := svc.CreateVolume(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func ExampleEC2_CreateVolume_shared00() {
	svc := ec2.New(session.New())
	input := &ec2.CreateVolumeInput{
		AvailabilityZone: aws.String("us-east-1a"),
		Size:             aws.Int64(80),
		VolumeType:       aws.String("gp2"),
	}

	result, err := svc.CreateVolume(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func ExampleEC2_CreateTags_shared00() {
	svc := ec2.New(session.New())
	input := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String("ami-78a54011"),
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Stack"),
				Value: aws.String("production"),
			},
		},
	}

	result, err := svc.CreateTags(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func ExampleEC2_CreateSubnet_shared00() {
	svc := ec2.New(session.New())
	input := &ec2.CreateSubnetInput{
		CidrBlock: aws.String("10.0.1.0/24"),
		VpcId:     aws.String("vpc-a01106c2"),
	}

	result, err := svc.CreateSubnet(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func ExampleEC2_CreateSecurityGroup_shared00() {
	svc := ec2.New(session.New())
	input := &ec2.CreateSecurityGroupInput{
		Description: aws.String("My security group"),
		GroupName:   aws.String("my-security-group"),
		VpcId:       aws.String("vpc-1a2b3c4d"),
	}

	result, err := svc.CreateSecurityGroup(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func ExampleEC2_CreateRouteTable_shared00() {
	svc := ec2.New(session.New())
	input := &ec2.CreateRouteTableInput{
		VpcId: aws.String("vpc-a01106c2"),
	}

	result, err := svc.CreateRouteTable(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}
func ExampleEC2_CreateRoute_shared00() {
	svc := ec2.New(session.New())
	input := &ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		GatewayId:            aws.String("igw-c0a643a9"),
		RouteTableId:         aws.String("rtb-22574640"),
	}

	result, err := svc.CreateRoute(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}
func main() {
	tags := make([]*Tag, 0)
	tags = append(tags, &Tag{Key:"Name", Value:"velanVpc"})
	vpcReq := CreateVpcStruct{
		Tag:tags,
		Name: "velanVpc",
		CredLocation: "/Users/singaravelannandakumar/.aws/credentials",
		CredName: "default",
		Region: "eu-central-1",
		CIDR: "10.0.0.0/16",
	}
	CreateVpc(vpcReq)
}