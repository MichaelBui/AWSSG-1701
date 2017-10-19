package controllers

import (
	"github.com/michaelbui/AWSSG-1710/libraries"
	"github.com/michaelbui/AWSSG-1710/types"
)

func getDatabaseConnector(config *types.DBConfig) types.DatabaseConnector {
	var conn types.DatabaseConnector
	switch config.Type {
	case "sqlite3":
		conn = libraries.NewSqliteConnection(config)
	}
	return conn
}

func getCloudConnector(config *types.CloudConfig) types.CloudConnector {
	var conn types.CloudConnector
	switch config.Type {
	case "aws":
		conn = libraries.NewAwsConnector(config.Configs)
	}
	return conn
}
