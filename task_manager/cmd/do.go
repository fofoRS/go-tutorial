/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"strconv"
	"unicode"

	"github.com/fofoRS/go-tutorial/task_manager/db"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "marks a task as complete",
	Long:  "marks a task as complete",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires the number of the task you want to mark as complete")
		} else if !unicode.IsDigit([]rune(args[0])[0]) {
			return errors.New("values must be a digit")
		} else if v, _ := strconv.Atoi(args[0]); v < 1 {
			return errors.New("value must be greater than zero")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := strconv.Atoi(args[0])
		executeDoCommand(key)
	},
}

func init() {
	rootCmd.AddCommand(doCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// doCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func executeDoCommand(index int) {
	dbClient := db.NewConnection()
	defer dbClient.CloseConnection()
	completeTaskFunc := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("todo_new"))
		return bucket.Delete(itob(uint64(index)))
	}
	err := dbClient.Db.Update(completeTaskFunc)
	if err != nil {
		panic(err)
	}
}
