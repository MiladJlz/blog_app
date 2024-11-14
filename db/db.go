package db

const MongoDBNameEnvName = "MONGO_DB_NAME"

type Store struct {
	User UserStore
	Post PostStore
}
