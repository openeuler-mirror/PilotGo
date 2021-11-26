package utils

/**
 * @Author: zhang han
 * @Date: 2021/11/4 11:00
 * @Description:
 */

import (
	"math/rand"
	"time"
)

func RandomString(n int) string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for k := range result {
		result[k] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
