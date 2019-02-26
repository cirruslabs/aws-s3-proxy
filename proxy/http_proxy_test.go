package proxy

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const TestBucketName = "some-bucket"

func createFakeS3(t *testing.T) (*http.Server, *session.Session) {
	backend := s3mem.New()
	err := backend.CreateBucket(TestBucketName)
	if err != nil {
		t.Errorf("Failed to start fake S3: '%d'", err)
	}
	fakeS3 := gofakes3.New(backend)
	addr := "localhost:8123"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		t.Errorf("Failed to start fake S3: '%d'", err)
	}
	server := &http.Server{Addr: addr, Handler: fakeS3.Server()}
	go server.Serve(listener)
	if err != nil {
		t.Errorf("Failed to start servince fake S3: '%d'", err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		Endpoint:         &addr,
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials("test", "test", "test"),
	})

	if err != nil {
		t.Errorf("Failed to create fake AWS session: '%d'", err)
	}

	return server, sess
}

func Test_Blob_Exists(t *testing.T) {
	server, sess := createFakeS3(t)
	defer server.Close()

	storageProxy := NewStorageProxy(sess, TestBucketName, "")

	uploadFile(storageProxy, t, "some/object/file", "test")

	response := httptest.NewRecorder()
	storageProxy.checkBlobExists(response, "some/object/file")

	if response.Code == http.StatusOK {
		t.Log("Passed")
	} else {
		t.Errorf("Wrong status: '%d'", response.Code)
	}
}

func Test_Default_Prefix(t *testing.T) {
	server, sess := createFakeS3(t)
	defer server.Close()

	storageProxy := NewStorageProxy(sess, TestBucketName, "some/object/")

	uploadFile(storageProxy, t, "file", "test")

	response := httptest.NewRecorder()
	storageProxy.checkBlobExists(response, "file")

	if response.Code == http.StatusOK {
		t.Log("Passed")
	} else {
		t.Errorf("Wrong status: '%d'", response.Code)
	}
}

func Test_Blob_Download(t *testing.T) {
	expectedBlobContent := "my content"
	server, sess := createFakeS3(t)
	defer server.Close()

	storageProxy := NewStorageProxy(sess, TestBucketName, "")

	uploadFile(storageProxy, t, "some/file", expectedBlobContent)

	response := httptest.NewRecorder()
	storageProxy.downloadBlob(response, "some/file")

	if response.Code == http.StatusOK {
		t.Log("Passed")
	} else {
		t.Errorf("Wrong status: '%d'", response.Code)
	}

	downloadedBlobContent := response.Body.String()
	if downloadedBlobContent == expectedBlobContent {
		t.Log("Passed")
	} else {
		t.Errorf("Wrong content: '%s'", downloadedBlobContent)
	}
}

func Test_Blob_Upload(t *testing.T) {
	content := "my content"
	server, sess := createFakeS3(t)
	defer server.Close()

	storageProxy := NewStorageProxy(sess, TestBucketName, "")

	uploadFile(storageProxy, t, "test", content)
}

func uploadFile(storageProxy *StorageProxy, t *testing.T, name string, content string) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/"+name, strings.NewReader(content))
	storageProxy.uploadBlob(response, request, name)
	if response.Code == http.StatusCreated {
		t.Log("Passed")
	} else {
		t.Errorf("Wrong status: '%d' %v", response.Code, response.Body)
	}
}
