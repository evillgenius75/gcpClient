package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/servicemanagement/v1"
)

type Config struct {
	projectID string
	jsonpath  string
}

// Credentials

//var projectID = flags                                                        // Your ProjectID
//var jsonPath = filepath.Join(os.Getenv("HOME"), "mypersonalgcpjsonkey.json") // path to your JSON file

// GoogleCloudClient is a generic wrapper for talking with individual services inside Google Cloud
// such as Cloud Resource Manager, IAM, Services, Billing and DNS
type GoogleCloudClient struct {
	// Structs from Google library
	Resource *cloudresourcemanager.Service
	IAM      *iam.Service
	Service  *servicemanagement.APIService
	Billing  *cloudbilling.APIService
	DNS      *dns.Service

	// Required user input
	ProjectID string
	JSONPath  string
}

// NewGoogleCloudClient returns a pointer to the `GoogleCloudClient` instance
func NewGoogleCloudClient(projectID string, jsonPath string) (*GoogleCloudClient, error) {
	ctx := context.Background()

	// Client for Cloud Resource Manager
	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		return nil, fmt.Errorf("Error with Cloud Resource Manager Service: %v", err)
	}

	// Client for IAM
	iamService, err := iam.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		return nil, fmt.Errorf("Error with the IAM Service: %v", err)
	}

	// Client for Service Infrastructure Manager
	servicemanagementService, err := servicemanagement.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		return nil, fmt.Errorf("Error with the Service Management Service: %v", err)
	}

	// Client for Cloud Billing
	cloudbillingService, err := cloudbilling.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		return nil, fmt.Errorf("Error with the Cloud Billing Account: %v", err)
	}

	// Client for Google Cloud DNS API
	dnsService, err := dns.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		return nil, fmt.Errorf("Error with the Cloud DNS: %v", err)
	}

	return &GoogleCloudClient{
		Resource:  cloudresourcemanagerService,
		IAM:       iamService,
		Service:   servicemanagementService,
		Billing:   cloudbillingService,
		DNS:       dnsService,
		ProjectID: projectID,
		JSONPath:  jsonPath,
	}, nil
}

// ListProjects lists the Projects of a GCP service account and returns an error
func (c *GoogleCloudClient) ListProjects() (*cloudresourcemanager.ListProjectsResponse, error) {
	projectsList, err := c.Resource.Projects.List().Do()
	if err != nil {
		return nil, err
	}
	return projectsList, nil
}

// GetProject returns a project from GCP
func (c *GoogleCloudClient) GetProject(projectID string) (*cloudresourcemanager.Project, error) {
	project, err := c.Resource.Projects.Get(projectID).Do()
	if err != nil {
		return nil, err
	}
	return project, nil
}

func main() {
	var cfg Config

	flag.StringVar(&cfg.projectID, "project-id", "", "ProjectID for Client Access")
	flag.StringVar(&cfg.jsonpath, "jsonKey", "", "The path to the JSON key for Client Access")
	flag.Parse()

	projectID := cfg.projectID
	jsonPath := filepath.Join(os.Getenv("HOME"), cfg.jsonpath)

	fmt.Printf("Project ID: %s jsonPath: %s", projectID, jsonPath)
	gcpClient, err := NewGoogleCloudClient(projectID, jsonPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gcpClient.Resource.BasePath)
	fmt.Println(gcpClient.IAM.BasePath)
	fmt.Println(gcpClient.Service.BasePath)
	fmt.Println(gcpClient.Billing.BasePath)
	fmt.Println(gcpClient.DNS.BasePath)

	resp, err := gcpClient.ListProjects()
	if err != nil {
		log.Fatal(err)
	}
	for _, project := range resp.Projects {
		fmt.Printf("Project Name: %s\tProjectID: %s\tProject Number: %d\r\n", project.Name, project.ProjectId, project.ProjectNumber)
	}

	// gcpClient.Resource.Projects.Delete("hidden-howl-252922")
	resp1, err := gcpClient.GetProject(projectID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp1.ProjectNumber)

}
