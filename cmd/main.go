package cmd

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger/v3"
	"github.com/spf13/cobra"
	"github.com/stevehebert/frontsidecache/internal/badgerstore"
)

var rootCmd = &cobra.Command{
	Use:   "serve",
	Short: "Populate cache from kafka and host",
	Long: `A fast edge cache solution that drives simple key value
		lookup to the edge`,
	Run: func(cmd *cobra.Command, args []string) {
		b, _ := badgerstore.Connect()

		if b != nil {
			fmt.Println("connected")
		}

		err := b.Update(func(txn *badger.Txn) error {
			fmt.Println("writing")
			return txn.Set([]byte("answer"), []byte("42"))
		})

		fmt.Println(err)

		err = b.Update(func(txn *badger.Txn) error {
			fmt.Println("writing")
			return txn.Set([]byte("answer"), []byte("43"))
		})

		fmt.Println(err)

		if err != nil {
			fmt.Println(err)
		}

		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
