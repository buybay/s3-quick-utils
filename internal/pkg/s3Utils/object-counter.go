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
	"sync/atomic"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cheggaaa/pb/v3"
)

const (
	prefixLen = 1
)

func ObjectCounter(
	bucket string,
	region string,
	profile string) (uint64, error) {

	statusBar := pb.Start64(0)
	defer statusBar.Finish()

	prefixJobs := genPrefix(prefixLen)

	counterChan := make(chan uint64)
	errorsChan := make(chan error)
	statusChan := make(chan string)

	svc := s3.New(NewSession(region, profile))

	for _, prefix := range prefixJobs {
		go counterListObjectsPrefix(svc, bucket, prefix, nil, counterChan, statusChan, errorsChan)
	}

	var prefixDone uint64
	var objectCounter uint64
	for {
		select {
		case count := <-counterChan:
			atomic.AddUint64(&objectCounter, count)
			statusBar.Add64(int64(count))
		case <-statusChan:
			atomic.AddUint64(&prefixDone, 1)
			if prefixDone >= uint64(len(prefixJobs)) {
				return objectCounter, nil
			}
		case err := <-errorsChan:
			return objectCounter, err
		}
	}

}
