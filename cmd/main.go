package main

import (
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/cirruslabs/aws-s3-proxy/proxy"
	"log"
)

func main() {
	var port int64
	flag.Int64Var(&port, "port", 8080, "Port to serve")
	var bucketName string
	flag.StringVar(&bucketName, "bucket", "", "S3 Name")
	var region string
	flag.StringVar(&region, "region", "", "S3 region")
	var defaultPrefix string
	flag.StringVar(&defaultPrefix, "prefix", "", "Optional prefix for all objects. For example, use --prefix=foo/ to work under foo directory in a bucket.")
	flag.Parse()

	if bucketName == "" {
		log.Fatal("Please specify S3 Bucket")
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region:          aws.String(region),
		S3UseAccelerate: aws.Bool(true),
	}))

	storageProxy := proxy.NewStorageProxy(sess, bucketName, defaultPrefix)

	err := storageProxy.Serve(port)
	if err != nil {
		log.Fatalf("Failed to start proxy: %s", err)
	}
}
