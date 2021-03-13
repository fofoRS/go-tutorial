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
	"os"
	"strconv"
	"strings"

	"github.com/fofoRS/go-tutorial/task_manager/db"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists all of our incomplete tasks",
	Long:  "lists all of our incomplete tasks",
	Run: func(cmd *cobra.Command, args []string) {
		executeListCommand()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func executeListCommand() {
	dbClient := db.NewConnection()
	defer dbClient.CloseConnection()

	getListOfTaskFunc := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("todo_new"))
		tasks := make([]string, 10)
		bucket.ForEach(func(k, v []byte) error {
			var task strings.Builder
			task.WriteString(strconv.Itoa(btoi(k)))
			task.WriteRune('.')
			task.WriteRune(' ')
			task.WriteString(string(v))
			task.WriteString("\n")
			tasks = append(tasks, task.String())
			return nil
		})
		consoleWriter := os.Stdout
		for _, v := range tasks {
			consoleWriter.WriteString(v)
		}
		return nil
	}

	dbClient.Db.View(getListOfTaskFunc)
}

func btoi(numberByte []byte) int {
	value := binary.BigEndian.Uint64(numberByte)
	return int(value)
}
