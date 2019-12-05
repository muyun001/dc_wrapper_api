package settings

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var DbConnection string
var DbUsername string
var DbPassword string
var DbHost string
var DbPort string
var DbDatabase string

var Debug bool

var DcFixApi string
var DcApi string
var DcUserId string
var DcIsOuter bool

func init() {
	checkEnv()
	LoadSetting()
}

func checkEnv() {
	_ = godotenv.Load()
	needChecks := []string{
		"DB_CONNECTION", "DB_HOST", "DB_PORT", "DB_DATABASE", "DB_USERNAME", "DB_PASSWORD",
		"DC_USER_ID", "DC_FIX_API",
	}

	for _, envKey := range needChecks {
		if os.Getenv(envKey) == "" {
			log.Fatalf("env %s missed", envKey)
		}
	}
}

func LoadSetting() {
	DbConnection = os.Getenv("DB_CONNECTION")
	DbUsername = os.Getenv("DB_USERNAME")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
	DbDatabase = os.Getenv("DB_DATABASE")

	debug := os.Getenv("DEBUG")
	if debug != "" && debug != "false" && debug != "0" {
		Debug = true
	}

	DcFixApi = os.Getenv("DC_FIX_API")
	DcUserId = os.Getenv("DC_USER_ID")
	isOuter := os.Getenv("DC_IS_OUTER")
	if isOuter != "" && isOuter != "false" && isOuter != "0" {
		DcIsOuter = true
	}
}
