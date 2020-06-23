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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	paginationSize = 1000
)

func counterListObjectsPrefix(c *s3.S3,
	bucket string,
	prefix string,
	objectsChan chan string,
	counterChan chan uint64,
	resultChan chan string,
	errorChan chan error) error {

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(paginationSize),
		Prefix:  aws.String(prefix),
	}

	for {
		result, err := c.ListObjectsV2(input)
		if err != nil {
			errorChan <- err
			return err
		}

		if counterChan != nil {
			counterChan <- uint64(*result.KeyCount)
		}

		if objectsChan != nil {
			for _, elem := range result.Contents {
				objectsChan <- *elem.Key
			}
		}

		if !*result.IsTruncated {
			// no more results
			resultChan <- prefix
			return nil
		}

		input.SetContinuationToken(*result.NextContinuationToken)
	}
}
