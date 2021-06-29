package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/codeship/codeship-go"
)

var (
	timeout = flag.Int("source", 300, "Timeout in seconds.")
)

func main() {
	flag.Parse()
	organization := os.Getenv("CODESHIP_ORGANIZATION")
	username := os.Getenv("CODESHIP_USERNAME")
	password := os.Getenv("CODESHIP_PASSWORD")
	projectID := os.Getenv("CI_PROJECT_ID")
	buildID := os.Getenv("CI_BUILD_ID")

	endTime := time.Now().Add(time.Duration(*timeout) * time.Second)

	ctx := context.Background()
	auth := codeship.NewBasicAuth(username, password)
	client, err := codeship.New(auth)
	if err != nil {
		log.Fatalf("Codeship authentication, %v", err)
	}
	org, err := client.Organization(ctx, organization)
	if err != nil {
		log.Fatalf("Codeship can not find organization, %v", err)
	}
	for time.Now().Before(endTime) {
		builds, _, err := org.ListBuilds(ctx, projectID, codeship.PerPage(50))
		if err != nil {
			log.Fatalf("Codeship cannot list builds, %v", err)
		}
		var ref string
		var start time.Time
		for _, build := range builds.Builds {
			if build.UUID == buildID {
				ref = build.Ref
				start = build.QueuedAt
			}
		}
		for _, build := range builds.Builds {
			if build.Ref == ref && build.UUID != buildID && build.Status != "success" && build.Status != "skipped" && build.Status != "stopped" && build.Status != "error" {
				if build.QueuedAt.After(start) {
					log.Printf("The is a newer build: %s", build.Links.Steps)
					_, _, err := org.StopBuild(ctx, projectID, buildID)
					if err != nil {
						log.Fatalf("Cannot stop build, %v", err)
					}
					os.Exit(0)
				}
			}
		}
		done := true
		for _, build := range builds.Builds {
			if build.Ref == ref && build.UUID != buildID && build.Status != "success" && build.Status != "skipped" && build.Status != "stopped" && build.Status != "error" {
				if build.QueuedAt.Before(start) {
					log.Printf("Waiting for older build")
					done = false
					break
				}
			}
		}
		if done {
			os.Exit(0)
		}
		time.Sleep(1 * time.Second)
	}
	log.Fatalf("Timeout")
}
