package base

import "os"

var (
	clusterNameKey = "GOGO_CLUSTER"
	envNameKey     = "GOGO_ENV"
)

func SetGlobalClusterName(name string) {
	clusterNameKey = name
}

func SetGlobalEnvName(name string) {
	envNameKey = name
}

func GetEnv() string {
	return os.Getenv(envNameKey)
}

func GetCluster() string {
	return os.Getenv(clusterNameKey)
}

func IsProduction() bool {
	return GetCluster() == "prod"
}

func IsUAT() bool {
	return GetCluster() == "uat"
}

func IsDev() bool {
	return GetCluster() == "dev"
}

func InCloud() bool {
	return GetCluster() != ""
}

func InLocal() bool {
	return GetCluster() == ""
}
