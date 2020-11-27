package main

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"leapbeyond.ai/models"
	"log"
)

// makeResourcesClient creates a resource client using the current session.
func makeResourcesClient(session *models.AzureSession) resources.GroupsClient {
	client := resources.NewGroupsClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer

	return client
}

// createResourceGroup tries to create a resource group, using the base name in the configuration.
// It will return an non-nil error on error.
func createResourceGroup(session *models.AzureSession) error {
	name:=makeResourceGroupName(session)
	log.Printf("Begin creating resource group %s", name)

	client := makeResourcesClient(session)

	resourceType := "Microsoft.Resources/resourceGroups"

	group, err := client.CreateOrUpdate(ctx, name, resources.Group{
		Name:     &name,
		Type:     &resourceType,
		Location: &session.TargetLocation,
		Tags: map[string]*string{
			"Name":    &name,
			"Client":  &clientTag,
			"Owner":   &ownerTag,
			"Project": &projectTag,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create resource group %s: %v", name, err)
	}

	session.ResourceGroupName = name
	log.Printf("resource group created: %s", *group.ID)
	return nil
}

func destroyResourceGroup(session *models.AzureSession) error {
	client := makeResourcesClient(session)
	future, err := client.Delete(ctx, session.ResourceGroupName)
	if err != nil {
		return fmt.Errorf("failed to delete resource group %s: %v", session.ResourceGroupName, err)
	}

	log.Printf("Waiting to delete resource group %s", session.ResourceGroupName)
	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("failed on waiting to delete resource group %s: %v", session.ResourceGroupName, err)
	}

	return nil
}
