package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jackc/pgconn"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "pgtocsv",
	Short: "pgtocsv executes a query and returns the results in CSV format.",
	Long: `pgtocsv executes a query and returns the results in CSV format

If neither --file or --sql is specified the SQL will be read from STDIN.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		var sql string
		if viper.GetString("file") != "" {
			sqlBytes, err := ioutil.ReadFile(viper.GetString("file"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to read SQL from file: %v\n", err)
				os.Exit(1)
			}
			sql = string(sqlBytes)

		} else if viper.GetString("sql") != "" {
			sql = viper.GetString("sql")
		} else {
			sqlBytes, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to read SQL from stdin: %v\n", err)
				os.Exit(1)
			}
			sql = string(sqlBytes)
		}

		conn, err := pgconn.Connect(ctx, viper.GetString("database_url"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(ctx)

		w := csv.NewWriter(os.Stdout)

		mrr := conn.Exec(ctx, sql)

		for mrr.NextResult() {
			rr := mrr.ResultReader()
			fieldDescriptions := rr.FieldDescriptions()
			columnNames := make([]string, len(fieldDescriptions))
			for i := range fieldDescriptions {
				columnNames[i] = string(fieldDescriptions[i].Name)
			}

			w.Write(columnNames)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to write CSV header: %v\n", err)
				os.Exit(1)
			}

			for rr.NextRow() {
				values := rr.Values()
				row := make([]string, len(values))
				for i := range values {
					row[i] = string(values[i])
				}

				w.Write(row)
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to write CSV row: %v\n", err)
					os.Exit(1)
				}
			}

			_, err = rr.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error querying database: %v\n", err)
				os.Exit(1)
			}
		}

		err = mrr.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error querying database: %v\n", err)
			os.Exit(1)
		}
		w.Flush()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pgtocsv.yaml)")

	rootCmd.Flags().StringP("database-url", "d", "", "Database URL or DSN")
	viper.BindPFlag("database_url", rootCmd.Flags().Lookup("database-url"))

	rootCmd.Flags().StringP("file", "f", "", "File containing SQL to execute")
	viper.BindPFlag("file", rootCmd.Flags().Lookup("file"))

	rootCmd.Flags().StringP("sql", "s", "", "SQL to execute")
	viper.BindPFlag("sql", rootCmd.Flags().Lookup("sql"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".pgtocsv" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".pgtocsv")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
