// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package s3Utils

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cheggaaa/pb/v3"
	"golang.org/x/sync/semaphore"
)

var (
	maxWorkers   = runtime.GOMAXPROCS(0)
	semDelete    = semaphore.NewWeighted(int64(maxWorkers / 2))
	semList      = semaphore.NewWeighted(int64(maxWorkers / 2))
	deleteBlocks = 1000
	ctx          = context.Background()
)

func deleteObject(svc *s3.S3,
	bucket string,
	items []string,
	ctx context.Context,
	semDelete *semaphore.Weighted,
	deleteCounterChan chan uint64,
	errorChan chan error) error {

	if len(items) == 0 {
		return nil
	}

	if err := semDelete.Acquire(ctx, 1); err != nil {
		errorChan <- err
		return err
	}
	defer semDelete.Release(1)

	var deleteItems []*s3.ObjectIdentifier
	for _, elem := range items {
		delItem := s3.ObjectIdentifier{Key: aws.String(elem)}
		deleteItems = append(deleteItems, &delItem)
	}

	input := &s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &s3.Delete{
			Objects: deleteItems,
			Quiet:   aws.Bool(false),
		},
	}

	for {
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_, err := svc.DeleteObjectsWithContext(ctxTimeout, input)
		if err != nil {
			if strings.Contains(err.Error(), "SlowDown") ||
				strings.Contains(err.Error(), "RequestCanceled") {
				// Manage SlowDown with an Context timeout
				delay := rand.Intn(100) * 100
				time.Sleep(time.Duration(delay) * time.Millisecond)
			} else {
				errorChan <- err
				return err
			}
		} else {
			break
		}
	}

	if deleteCounterChan != nil {
		deleteCounterChan <- uint64(len(items))
	}
	return nil
}

func DeleteBucketObjects(
	bucket string,
	region string,
	profile string) (uint64, error) {

	fmt.Println("Deleting bucket ", bucket)

	statusBar := pb.Start64(0)
	defer statusBar.Finish()

	prefixJobs := genPrefix(1)

	deleteCounterChan := make(chan uint64)
	objectsChan := make(chan string)
	errorsChan := make(chan error)
	statusChan := make(chan string)

	var output []string

	svc := s3.New(NewSession(region, profile))

	for _, prefix := range prefixJobs {
		go func() {
			if err := semList.Acquire(ctx, 1); err != nil {
				errorsChan <- err
			}
			defer semList.Release(1)

			counterListObjectsPrefix(
				svc,
				bucket,
				prefix,
				objectsChan,
				nil,
				statusChan,
				errorsChan)
		}()
	}

	var prefixDone uint64
	var objectCounter, deletedCounter uint64
	var muxLastDelete sync.Mutex
	for {
		select {
		case count := <-deleteCounterChan:
			atomic.AddUint64(&deletedCounter, count)
			statusBar.SetCurrent(int64(deletedCounter))

			if (deletedCounter + uint64(deleteBlocks)) >= objectCounter {
				deletedCounter = objectCounter
				statusBar.SetCurrent(int64(deletedCounter))
				return deletedCounter, nil
			}

		case obj := <-objectsChan:
			muxLastDelete.Lock()
			output = append(output, obj)
			atomic.AddUint64(&objectCounter, uint64(1))
			statusBar.SetTotal(int64(objectCounter))
			if len(output) == deleteBlocks {
				go deleteObject(svc,
					bucket,
					output,
					ctx,
					semDelete,
					deleteCounterChan,
					errorsChan)
				output = []string{}
			}
			muxLastDelete.Unlock()
		case <-statusChan:
			atomic.AddUint64(&prefixDone, 1)
			if prefixDone >= uint64(len(prefixJobs)) {
				go deleteObject(
					svc,
					bucket,
					output,
					ctx,
					semDelete,
					nil,
					errorsChan)
				if deletedCounter == objectCounter {
					return deletedCounter, nil
				}
			}

		case err := <-errorsChan:
			return deletedCounter, err
		}
	}

}
