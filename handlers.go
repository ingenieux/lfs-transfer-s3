package lfstransfers3

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

func (a *TransferAgent) reportError(err error) {
	eventError := EventError{Code: 1, Message: err.Error()}

	bytesToReturn, _ := json.Marshal(eventError)

	a.out.Write(bytesToReturn)
	a.out.Write([]byte("\n"))
}

func (a *TransferAgent) ExecuteInit(e InitEvent) (int, error) {
	a.out.Write([]byte("{}\n"))

	if newSession, err := session.NewSession(&aws.Config{Region: aws.String(a.Parameters.Region)}); nil != err {
		//a.out.Write(json.Marshal(err))

		a.reportError(err)

		return 1, err
	} else {
		a.session = newSession
	}

	// Create a downloader with the session and custom options
	a.downloader = s3manager.NewDownloader(a.session, func(d *s3manager.Downloader) {
		d.Concurrency = e.ConcurrentTransfers
		d.PartSize = a.Parameters.PartSize
	})

	// ...and an uploader
	a.uploader = s3manager.NewUploader(a.session, func(d *s3manager.Uploader) {
		d.Concurrency = e.ConcurrentTransfers
		d.PartSize = a.Parameters.PartSize
		d.LeavePartsOnError = true
	})

	return 0, nil
}

func (a *TransferAgent) ExecuteUpload(e UploadEvent) (int, error) {
	urlAsString, _ := e.Action["href"].(string)

	log.Debugf("urlAsString: %s", urlAsString)

	url, err := url.Parse(urlAsString)

	log.Debugf("Parsed url: %+v", url)

	if nil != err {
		a.reportError(err)

		return 1, err
	}

	pathElements := strings.Split(url.Path, "/")

	log.Debugf("pathElements: %s\n", pathElements)

	bucket := pathElements[1]

	key := strings.Join(pathElements[2:], "/")

	fileToUpload, err := os.OpenFile(e.Path, os.O_RDONLY, 0660)

	if nil != err {
		a.reportError(err)

		return 1, err
	}

	defer func() {
		fileToUpload.Close()
	}()

	uploadInput := &s3manager.UploadInput{
		Bucket:       aws.String(bucket),
		Key:          aws.String(key),
		StorageClass: aws.String("STANDARD_IA"),
		ContentType:  aws.String("application/binary"),
		Body:         fileToUpload,
	}

	log.Debugf("Dispatching upload for %+v\n", uploadInput)

	if _, err := a.uploader.Upload(uploadInput); nil != err {
		fmt.Fprintf(os.Stderr, "F*deu\n")

		completeEvent := &CompleteEvent{
			EventType: "complete",
			Oid:       e.Oid,
			Error: &EventError{
				Code:    500,
				Message: err.Error(),
			},
		}

		bytesToSend, _ := json.Marshal(completeEvent)

		a.out.Write(bytesToSend)
		a.out.Write([]byte("\n"))
	} else {
		log.Debugf("Complete\n")

		completeEvent := &CompleteEvent{
			EventType: "complete",
			Oid:       e.Oid,
			Path:      &urlAsString,
		}

		bytesToSend, _ := json.Marshal(completeEvent)

		a.out.Write(bytesToSend)
		a.out.Write([]byte("\n"))
	}

	return 0, nil
}

func (a *TransferAgent) ExecuteDownload(e DownloadEvent) (int, error) {
	urlAsString, _ := e.Action["href"].(string)

	log.Debugf("urlAsString: %s", urlAsString)

	url, err := url.Parse(urlAsString)

	log.Debugf("Parsed url: %+v", url)

	if nil != err {
		a.reportError(err)

		return 1, err
	}

	pathElements := strings.Split(url.Path, "/")

	log.Debugf("pathElements: %s\n", pathElements)

	bucket := pathElements[1]

	key := strings.Join(pathElements[2:], "/")

	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	var tempFile *os.File

	if tempFile, err = ioutil.TempFile("", "lfs-"); nil != err {
		a.reportError(err)

		return 0, err
	}

	absPath := tempFile.Name()

	defer func() {
		tempFile.Close()
	}()

	_, err = a.downloader.Download(tempFile, getObjectInput)

	if nil != err {
		completeEvent := &CompleteEvent{
			EventType: "complete",
			Oid:       e.Oid,
			Error: &EventError{
				Code:    500,
				Message: err.Error(),
			},
		}

		bytesToSend, _ := json.Marshal(completeEvent)

		a.out.Write(bytesToSend)
		a.out.Write([]byte("\n"))
	} else {
		fmt.Fprintf(os.Stderr, "Complete\n")

		completeEvent := &CompleteEvent{
			EventType: "complete",
			Oid:       e.Oid,
			Path:      &absPath,
		}

		bytesToSend, _ := json.Marshal(completeEvent)

		a.out.Write(bytesToSend)
		a.out.Write([]byte("\n"))
	}

	return 0, nil
}

func (a *TransferAgent) ExecuteTerminate() {
	os.Exit(0)
}
