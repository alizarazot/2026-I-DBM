package constants

const (
	ENV_BASE             = "LINEA_"
	ENV_ADDR             = ENV_BASE + "ADDR"
	ENV_JWT_SECRET       = ENV_BASE + "JWT_SECRET"
	ENV_MONGODB_URI      = ENV_BASE + "MONGODB_URI"
	ENV_MONGODB_DATABASE = ENV_BASE + "MONGODB_DATABASE"
)

const (
	DB_COLLECTION_USER       = "user"
	DB_COLLECTION_CFC        = "cfc"
	DB_COLLECTION_CFC_ANSWER = "cfc_answer"
)
