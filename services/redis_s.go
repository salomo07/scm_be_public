package services

import (
	"context"
	"fmt"
	"log"
	"scm/config"
	"scm/utils"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type ChannelObj struct {
	NamaChannel string
	PubSub      *redis.PubSub
}

var (
	ctx    = context.Background()
	client *redis.Client
)

func init() {
	var err error
	client, err = createRedisClient()
	if err != nil {
		panic(err)
	}
}
func createRedisClient() (*redis.Client, error) {
	opt, err := redis.ParseURL(config.GetCredRedis())
	if err != nil {
		return nil, err
	}
	return redis.NewClient(opt), nil
}

var channelsMap = make(map[string]*redis.PubSub)

// key string, value string, expired time.Duration
func SaveValueRedis(data ...string) {
	if len(data) < 2 {
		fmt.Println("Insufficient arguments provided.")
		return
	}

	key := data[0]
	value := data[1]
	// fmt.Println("Saving value redis key (" + key + ")")

	// Cek apakah kunci sudah ada
	// exists, err := client.Exists(ctx, key).Result()
	// if err != nil {
	// 	fmt.Printf("Error checking if key exists: %v\n", err)
	// 	return
	// }

	// Tentukan waktu kedaluwarsa
	var expiration time.Duration
	if len(data) == 3 && data[2] != "" {
		// 🔹 Anggap data[2] = epoch second (hasil Unix timestamp)
		expTimeSec, err := strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			fmt.Printf("Error parsing expiration time: %v\n", err)
			return
		}
		expirationTime := time.Unix(expTimeSec, 0) // epoch second → time.Time
		expiration = time.Until(expirationTime)    // TTL dihitung dari sekarang
		if expiration <= 0 {
			fmt.Println("Expiration sudah lewat, pakai default 4 jam")
			expiration = 4 * time.Hour
		}
	} else {
		expiration = 4 * time.Hour
	}

	// Simpan atau perbarui nilai
	err := client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		fmt.Printf("Error setting value: %v\n", err)
		return
	}
	fmt.Printf("%s is saved in Redis\n", key)
}

func GetValueRedis(key string) (val string, errStr string) {
	print("\nGetting value redis key (" + key + ")\n")
	val, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ""
	} else if err != nil {
		return "", err.Error()
	}
	return val, ""
}
func RemoveValueRedis(key string) (val string, errStr string) {
	fmt.Printf("\nRemoving value redis key (%s)\n", key)

	// Cek apakah kunci ada
	exists, err := client.Exists(ctx, key).Result()
	if err != nil {
		fmt.Printf("Error checking if key exists: %v\n", err)
		return "", err.Error()
	}

	if exists == 0 {
		print("Key does not exist in Redis")
	}

	// Hapus key
	err = client.Del(ctx, key).Err()
	if err != nil {
		msg := fmt.Sprintf("Error deleting key: %v\n", err)
		return "", msg
	}

	msg := fmt.Sprintf("Key %s has been removed from Redis\n", key)
	return msg, ""
}

func SubscribeRedis(channelname string) {
	opt, _ := redis.ParseURL(config.GetCredRedis())
	client := redis.NewClient(opt)
	pubsub := client.Subscribe(context.Background(), channelname)

	channel := pubsub.Channel()
	go func() {
		for msg := range channel {
			log.Println("Terima Pesan: " + msg.Payload)
			println("\nMunkin bisa send to client pakai websocket")
		}
	}()
	println("Redis channel " + channelname + " is active")
	channelsMap[channelname] = pubsub
}
func Publish(channelName string, data any) {
	opt, _ := redis.ParseURL(config.GetCredRedis())
	client := redis.NewClient(opt)
	err := client.Publish(context.Background(), channelName, utils.StructToJson(data)).Err()
	if err != nil {
		log.Fatal(err)
	}
	Unsubscribe(channelName)
}
func Unsubscribe(channelName string) {
	pubsub := channelsMap[channelName]
	err := pubsub.Unsubscribe(context.Background(), channelName)
	if err != nil {
		log.Fatal(err)
	}
	errClose := pubsub.Close()
	if errClose != nil {
		log.Fatal(errClose)
	}
	println("\nUnsubscribe channel " + channelName + "\n")
}
