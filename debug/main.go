package main

//import (
//	"context"
//	"fmt"
//
//	"github.com/aws/aws-sdk-go-v2/config"
//	"github.com/aws/aws-sdk-go-v2/service/bedrock"
//)
//
//const region = "us-west-2"
//
//// main uses the AWS SDK for Go (v2) to create an Amazon Bedrock client and
//// list the available foundation models in your account and the chosen region.
//// This example uses the default settings specified in your shared credentials
//// and config files.
//func main() {
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	sdkConfig, err := config.LoadDefaultConfig(
//		ctx,
//		config.WithRegion(region),
//	)
//	if err != nil {
//		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
//		fmt.Println(err)
//		return
//	}
//	bedrockClient := bedrock.NewFromConfig(sdkConfig)
//
//	modelId := "anthropic.claude-3-5-sonnet-20240620-v1:0"
//
//	modeloutput, err := bedrockClient.GetFoundationModel(ctx, &bedrock.GetFoundationModelInput{
//		ModelIdentifier: &modelId,
//	})
//	if err != nil {
//		fmt.Printf("Couldn't get foundation model. Here's why: %v\n", err)
//		return
//	}
//
//	fmt.Println(*modeloutput.ModelDetails.ModelId)
//
//	//result, err := bedrockClient.ListFoundationModels(ctx, &bedrock.ListFoundationModelsInput{})
//	//if err != nil {
//	//	fmt.Printf("Couldn't list foundation models. Here's why: %v\n", err)
//	//	return
//	//}
//	//if len(result.ModelSummaries) == 0 {
//	//	fmt.Println("There are no foundation models.")
//	//}
//	//for _, modelSummary := range result.ModelSummaries {
//	//	fmt.Println(*modelSummary.ModelId)
//	//}
//}
