package types

type (
	DatabaseConnector interface {
		Save(entity interface{}) error
		Find(output interface{}) error
	}
)
