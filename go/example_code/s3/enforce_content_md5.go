/*
   Copyright 2010-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

   This file is licensed under the Apache License, Version 2.0 (the "License").
   You may not use this file except in compliance with the License. A copy of
   the License is located at

    http://aws.amazon.com/apache2.0/

   This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
   CONDITIONS OF ANY KIND, either express or implied. See the License for the
   specific language governing permissions and limitations under the License.
*/
package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Downloads an item from an S3 Bucket in the region configured in the shared config
// or AWS_REGION environment variable.
//
// Usage:
//    go run s3_download.go BUCKET ITEM
func main() {
	h := md5.New()
	content := strings.NewReader("")
	content.WriteTo(h)
	svc := s3.New(
		session.New(),
		&aws.Config{Region: aws.String("us-west-2")},
	)

	r, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String("testBucket"),
		Key:    aws.String("testKey"),
	})

	md5s := base64.StdEncoding.EncodeToString(h.Sum(nil))
	r.HTTPRequest.Header.Set("Content-MD5", md5s)
	url, err := r.Presign(15 * time.Minute)
	if err != nil {
		fmt.Println("error presigning request", err)
		return
	}

	req, err := http.NewRequest("PUT", url, strings.NewReader(""))
	req.Header.Set("Content-MD5", md5s)
	if err != nil {
		fmt.Println("error creating request", url)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	fmt.Println(resp, err)
}
