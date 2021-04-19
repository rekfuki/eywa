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
	AverageObjectSize float64          `bson:"avgObjSize" json:"average_object_size"`
	DataSize          float64          `bson:"dataSize" json:"data_size"`
	StorageSize       float64          `bson:"storageSize" json:"storage_size"`
	Indexes           int              `bson:"indexes" json:"indexes"`
	IndexSize         float64          `bson:"indexSize" json:"index_size"`
	TotalSize         float64          `bson:"totalSize" json:"total_size"`
	CollectionsInfo   []CollectionInfo `json:"collections_info"`
}

// CollectionInfo holds information a collection
type CollectionInfo struct {
	Namespace         string         `bson:"ns" json:"namespace"`
	AverageObjectSize float64        `bson:"avgObjSize" json:"average_object_size"`
	StorageSize       float64        `bson:"storageSize" json:"storage_size"`
	IndexCount        int            `bson:"nindexes" json:"index_count"`
	TotalSize         float64        `bson:"totalSize" json:"total_size"`
	IndexSizes        map[string]int `bson:"indexSizes" json:"index_sizes"`
	TotalIndexSize    float64        `bson:"totalIndexSize" json:"total_index_size"`
}
