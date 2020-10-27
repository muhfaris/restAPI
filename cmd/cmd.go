package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/muhfaris/restAPI/internal/pkg/logging"
	"github.com/muhfaris/restAPI/router"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/globalsign/mgo"
)

var (
	cfgFile string
	dbPool  *mgo.Database
	portApp = 3030
	logger  *logging.Logger
)

type Person struct {
	Name  string
	Phone string
}

var rootCmd = &cobra.Command{
	Use:   "restAPI",
	Short: "Boostrap of restAPI",
	Long:  "Prototype of restAPI with Golang",
	PreRun: func(cmd *cobra.Command, args []string) {
		initDatabase()
		router.Init(logger, dbPool)
	},
	Run: func(cmd *cobra.Command, args []string) {
		r := mux.NewRouter()
		port := viper.GetInt("app.port")

		r.HandleFunc("/", router.Home)

		apiRouter := r.PathPrefix("/api/v1").Subrouter()
		apiRouter.Handle("/articles", router.HandlerFunc(router.HandlerArticleList)).Methods(http.MethodGet)
		apiRouter.Handle("/articles", router.HandlerFunc(router.HandlerArticleCreate)).Methods(http.MethodPost)
		apiRouter.Handle("/articles/{id}", router.HandlerFunc(router.HandlerArticleDetail)).Methods(http.MethodGet)
		apiRouter.Handle("/articles/{id}", router.HandlerFunc(router.HandlerArticleUpdate)).Methods(http.MethodPut)
		apiRouter.Handle("/articles/{id}", router.HandlerFunc(router.HandlerArticleDelete)).Methods(http.MethodDelete)

		log.Println("Listen api at :", port)
		http.ListenAndServe(fmt.Sprintf(":%d", port), r)

		srv := &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			ReadTimeout:  time.Duration(viper.GetInt("app.read_timeout")) * time.Second,
			WriteTimeout: time.Duration(viper.GetInt("app.write_timeout")) * time.Second,
			Handler:      r,
		}

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Err.WithError(err).Fatalln("Listen and serve error.")
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is config.toml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("toml")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		//Search config with name ".name"
		viper.SetConfigName("config")
		viper.AddConfigPath("./config")
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Println(err)
		}
	}

	//Read env
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initDatabase() {
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	//	username := viper.GetString("database.username")
	//	password := viper.GetString("database.password")
	dbName := viper.GetString("database.name")

	uri := fmt.Sprintf("%s:%s", host, port)

	session, err := mgo.Dial(uri)
	if err != nil {
		logger.Err.Fatalf("can not connect to database %s", err)
	}

	dbPool = session.DB(dbName)
	session.SetMode(mgo.Monotonic, true)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	//rootCmd.AddCommand(cronCmd)

	if err := rootCmd.Execute(); err != nil {
		logger.Out.WithError(err).Fatalln("Execute command error.")
	}
}
