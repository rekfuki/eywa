package types

const MongoDBCredentialsSecretName = "mongodb-credentials"

// CreateUserDatabaseResponse ...
type CreateUserDatabaseResponse struct {
	Username string `json:"username"`
	Database string `json:"database"`
	Password string `json:"password"`
}

// UserDatabaseInfo holds information about the users database
type UserDatabaseInfo struct {
	UserID            string           `json:"user_id"`
	Name              string           `bson:"db" json:"name"`
	CollectionCount   int              `bson:"collections" json:"collection_count"`
	Objects           int              `bson:"objects" json:"objects"`
	AverageObjectSize int              `bson:"avgObjSize" json:"average_object_size"`
	DataSize          int              `bson:"dataSize" json:"data_size"`
	StorageSize       int              `bson:"storageSize" json:"storage_size"`
	Indexes           int              `bson:"indexes" json:"indexes"`
	IndexSize         int              `bson:"indexSize" json:"index_size"`
	TotalSize         int              `bson:"totalSize" json:"total_size"`
	CollectionsInfo   []CollectionInfo `json:"collections_info"`
}

// CollectionInfo holds information a collection
type CollectionInfo struct {
	Namespace         string         `bson:"ns" json:"namespace"`
	AverageObjectSize int            `bson:"avgObjSize" json:"average_object_size"`
	StorageSize       int            `bson:"storageSize" json:"storage_size"`
	IndexCount        int            `bson:"nindexes" json:"index_count"`
	TotalSize         int            `bson:"totalSize" json:"total_size"`
	IndexSizes        map[string]int `bson:"indexSizes" json:"index_sizes"`
	TotalIndexSize    int            `bson:"totalIndexSize" json:"total_index_size"`
}
