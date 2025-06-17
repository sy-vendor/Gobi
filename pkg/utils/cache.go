package utils

import (
	"time"

	"database/sql"
	"fmt"
	"gobi/internal/models"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	cache "github.com/patrickmn/go-cache"
)

var QueryCache *cache.Cache

func InitQueryCache(defaultExpiration, cleanupInterval time.Duration) {
	QueryCache = cache.New(defaultExpiration, cleanupInterval)
}

func GetQueryCache(key string) (interface{}, bool) {
	return QueryCache.Get(key)
}

func SetQueryCache(key string, value interface{}, ttl time.Duration) {
	QueryCache.Set(key, value, ttl)
}

func DeleteQueryCache(key string) {
	QueryCache.Delete(key)
}

// ExecuteSQL connects to the given data source and executes the SQL, returning the result as []map[string]interface{} or error
func ExecuteSQL(ds models.DataSource, sqlStr string) ([]map[string]interface{}, error) {
	var dsn, driver string
	switch ds.Type {
	case "mysql":
		driver = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", ds.Username, ds.Password, ds.Host, ds.Port, ds.Database)
	case "postgres":
		driver = "postgres"
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", ds.Host, ds.Port, ds.Username, ds.Password, ds.Database)
	case "sqlite":
		driver = "sqlite3"
		dsn = ds.Database
	default:
		return nil, fmt.Errorf("unsupported data source type: %s", ds.Type)
	}
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := []map[string]interface{}{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		scanArgs := make([]interface{}, len(cols))
		for i := range columns {
			scanArgs[i] = &columns[i]
		}
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}
		rowMap := make(map[string]interface{})
		for i, col := range cols {
			val := columns[i]
			b, ok := val.([]byte)
			if ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}
		results = append(results, rowMap)
	}
	return results, nil
}

// EncryptAES 加密明文，返回 base64 字符串
func EncryptAES(plaintext string) (string, error) {
	key := []byte(os.Getenv("DATA_SOURCE_SECRET"))
	if len(key) != 32 {
		return "", fmt.Errorf("DATA_SOURCE_SECRET must be 32 bytes (256 bit)")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES 解密 base64 字符串，返回明文
func DecryptAES(ciphertext string) (string, error) {
	key := []byte(os.Getenv("DATA_SOURCE_SECRET"))
	if len(key) != 32 {
		return "", fmt.Errorf("DATA_SOURCE_SECRET must be 32 bytes (256 bit)")
	}
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertextBytes := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
