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
	"github.com/aws/aws-sdk-go/aws/session"
)

// Helper to build a session with region or profile as parameter. Both
// optionals
func NewSession(region, profile string) *session.Session {

	// session configuration
	if region != "" && profile != "" {
		return session.Must(session.NewSessionWithOptions(session.Options{
			Config:  aws.Config{Region: aws.String(region)},
			Profile: profile,
		}))
	} else if region != "" {
		return session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{Region: aws.String(region)},
		}))
	} else if profile != "" {
		return session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
		}))
	} else {
		return session.Must(session.NewSession())
	}
}
