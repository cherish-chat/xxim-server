package mgmtmodel

type AppLineConfig struct {
	Config  string               `json:"config"`  // json格式
	AesIv   string               `json:"aesIv"`   // 32 bytes
	AesKey  string               `json:"aesKey"`  // 32 bytes
	Storage AppLineConfigStorage `json:"storage"` // 存储
}
type AppLineConfigStorage struct {
	Type     string                     `json:"type"`     // cos oss minio kodo
	ObjectId string                     `json:"objectId"` // 用于更新
	CdnAddr  string                     `json:"cdnAddr"`  // cdn地址 http://cdn.xxx.com
	Cos      *AppLineConfigStorageCos   `json:"cos,omitempty"`
	Oss      *AppLineConfigStorageOss   `json:"oss,omitempty"`
	Minio    *AppLineConfigStorageMinio `json:"minio,omitempty"`
	Kodo     *AppLineConfigStorageKodo  `json:"kodo,omitempty"`
}
type AppLineConfigStorageCos struct {
	AppId      string `json:"appId"`
	SecretId   string `json:"secretId"`
	SecretKey  string `json:"secretKey"`
	BucketName string `json:"bucketName"`
	Region     string `json:"region"`
}
type AppLineConfigStorageOss struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	BucketName      string `json:"bucketName"`
}
type AppLineConfigStorageMinio struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	BucketName      string `json:"bucketName"`
	SSL             bool   `json:"ssl"`
}
type AppLineConfigStorageKodo struct {
	AccessKey  string `json:"accessKey"`
	SecretKey  string `json:"secretKey"`
	BucketName string `json:"bucketName"`
}
