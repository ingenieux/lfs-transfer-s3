package lfstransfers3

import (
	"bufio"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

type TransferAgent struct {
	Parameters *Parameters
	in         *os.File
	out        *os.File
	session    *session.Session
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
}

func NewAgent(parameters *Parameters) (*TransferAgent, error) {
	log.Printf("Inicializing new Agent with Parameters %+v\n", parameters)

	if parameters.PartSize > MAXIMUM_UPLOAD_PART_SIZE {
		log.Errorf("Oops: Requested Size %d higher than MAXIMUM_UPLOAD_PART_SIZE (%d)", parameters.PartSize, MAXIMUM_UPLOAD_PART_SIZE)

		return nil, EPART_SIZE_LARGER_THAN_ALLOWED
	}

	dirStat, err := os.Stat(parameters.TempDir)

	if nil != err {
		log.Errorf("While stating tempDir: %s", err.Error())

		return nil, err
	}

	if !dirStat.IsDir() {
		log.Errorf("Not a directory: %s", parameters.TempDir)

		return nil, EINVALIDDIRECTORY
	}

	transferAgent := &TransferAgent{Parameters: parameters}

	transferAgent.in = os.Stdin
	transferAgent.out = os.Stdout

	return transferAgent, nil
}

func (a *TransferAgent) Execute() (int, error) {
	s := bufio.NewScanner(a.in)

	for s.Scan() {
		var genericEvent GenericEvent

		if err := json.Unmarshal(s.Bytes(), &genericEvent); err != nil {
			return 2, err
		}

		var marshalError error
		var rc int
		var err error

		switch genericEvent.EventType {
		case "init":
			{
				var targetEvent InitEvent

				marshalError = json.Unmarshal(s.Bytes(), &targetEvent)

				if nil != marshalError {
					return 2, marshalError
				}

				rc, err = a.ExecuteInit(targetEvent)

				if nil != err {
					return rc, err
				}
			}
		case "terminate":
			a.ExecuteTerminate()
		case "upload":
			{
				var targetEvent UploadEvent

				marshalError = json.Unmarshal(s.Bytes(), &targetEvent)

				if nil != marshalError {
					return 2, marshalError
				}

				go func() {
					fmt.Fprintf(os.Stderr, "Dispatching upload for %+v\n", targetEvent)

					a.ExecuteUpload(targetEvent)
				}()
			}
		case "download":
			{
				var targetEvent DownloadEvent

				marshalError = json.Unmarshal(s.Bytes(), &targetEvent)

				if nil != marshalError {
					return 2, marshalError
				}

				rc, err = a.ExecuteDownload(targetEvent)
			}

			if nil != err {
				return rc, err
			}
		}

	}

	return 0, nil
}
