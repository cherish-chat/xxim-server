package utils

import (
	"log"
	"os"
)

func GetPodIp() string {
	podIp := os.Getenv("POD_IP")
	if podIp == "" {
		log.Fatalf("env [POD_IP] is empty")
	}
	return podIp
}
