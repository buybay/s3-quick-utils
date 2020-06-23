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

package cmd

import (
	"fmt"
	"os"

	"github.com/buybay/s3-quick-utils/internal/pkg/s3Utils"
	"github.com/spf13/cobra"
)

var (
	region  string
	profile string
)

var rootCmd = &cobra.Command{
	Use:   "s3-quick-utils",
	Short: "S3 utilities focus on reduce execution times.",
}

var objCounterCmd = &cobra.Command{
	Use:   "counter [bucket_name]",
	Short: "Count objects in a S3 bucket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bucket := args[0]

		count, err := s3Utils.ObjectCounter(bucket, region, profile)
		if err != nil {
			fmt.Printf("Error: %v\n\n", err)
			os.Exit(1)
		}
		fmt.Printf("\nTotal objects in bucket %s: %d\n\n", bucket, count)
		os.Exit(0)
	},
	Example: "counter xxx-eu-central-1-production",
}

var deleteCmd = &cobra.Command{
	Use:   "delete [bucket_name]",
	Short: "Delete all the objects in a S3 bucket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bucket := args[0]

		count, err := s3Utils.DeleteBucketObjects(bucket, region, profile)
		if err != nil {
			fmt.Printf("Error: %v\n\n", err)
			os.Exit(1)
		}
		fmt.Printf("\nTotal objects in deleted from bucket %s: %d\n\n", bucket, count)
		os.Exit(0)
	},
	Example: "delete xxx-eu-central-1-production",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&region, "region", "eu-central-1", "AWS region for the S3 bucket")
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "AWS local profile")

	rootCmd.AddCommand(objCounterCmd)
	rootCmd.AddCommand(deleteCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
