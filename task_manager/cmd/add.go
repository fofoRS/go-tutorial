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
	"encoding/binary"
	"fmt"

	"github.com/fofoRS/go-tutorial/task_manager/db"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new Task to the list",
	Long:  "Adds a new Task to the list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		executeCommand(args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func executeCommand(args []string) {
	dbClient := db.NewConnection()
	defer dbClient.CloseConnection()
	insertTaskFunc := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("todo_new"))
		nextID, _ := bucket.NextSequence()
		return bucket.Put(itob(nextID), []byte(args[0]))
	}
	dbClient.Db.Update(insertTaskFunc)
}

func itob(v uint64) []byte {
	fmt.Println(v)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
