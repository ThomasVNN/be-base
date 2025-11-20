package postgresql

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"git-codecommit.ap-southeast-1.amazonaws.com/v1/repos/be-base/storage/local"
	"github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

type DatabaseAuth struct {
	Host     string
	Port     int
	UserName string
	Password string
}

func GetConnection() (db *gorm.DB, err error) {
	databas_eAuth, err := getDatabaseAuth()
	if err != nil{
		log.Printf("Cannot get AWS ENV:", err)
		return  nil, err
	}
	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		databas_eAuth.Host, databas_eAuth.Port, databas_eAuth.UserName, databas_eAuth.Password, local.Getenv("PG_DB_NAME"))
	log.Print(psql)
	DB3, err := gorm.Open("postgres", psql)
	if err != nil {
		log.Printf("Cannot connect to %s database", "postgres")
	}
	log.Printf("We are connected to the %s database", "postgres")
	return DB3, err
}

func getDatabaseAuth() (*DatabaseAuth, error) {
	var databaseAuth = DatabaseAuth{}
	if local.Getenv("ENVIRONMENT") == "dev" {
		databaseAuth = DatabaseAuth{
			Host:     local.Getenv("PG_DB_HOST"),
			Port:     5432,
			UserName: local.Getenv("PG_DB_USER"),
			Password: local.Getenv("PG_DB_PASS"),
		}
		return &databaseAuth, nil

	}
	secretName := local.Getenv("AWS_SECRET_NAME")
	region := local.Getenv("AWS_REGION")
	conf := aws.Config{
		Region: aws.String(region),
	}
	svc := secretsmanager.New(awssession.New(&conf))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}
	result, err := svc.GetSecretValue(input)
	if err != nil {
		return nil, err
	}
	var secretString, decodedBinarySecret string

	if result.SecretString != nil {
		secretString = *result.SecretString
		json.Unmarshal([]byte(secretString), &databaseAuth)
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			fmt.Println("Base64 Decode Error:", err)
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])
		json.Unmarshal([]byte(decodedBinarySecret), &databaseAuth)
	}
	return &databaseAuth, nil

}
