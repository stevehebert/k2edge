package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stevehebert/frontsidecache/http"
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

		err := b.Set([]byte("answer"), []byte("42"))

		fmt.Println(err)

		err = b.Set([]byte("answer"), []byte("43"))

		fmt.Println(err)

		if err != nil {
			fmt.Println(err)
		}
		s := http.New(b, ":http")
		s.Start(context.Background())

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
