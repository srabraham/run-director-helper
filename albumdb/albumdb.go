package albumdb

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"regexp"
	"strconv"
	"time"
)

var (
	yearRe = regexp.MustCompile("[0-9]{4}")
)

func TestDbConnection(dbProjectId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := firestore.NewClient(ctx, dbProjectId)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
}

func AddAlbumToDb(parkrunDate, albumName, albumUrl, dbProjectId, dbCollectionName string, runNum int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := firestore.NewClient(ctx, dbProjectId)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	yearKey := yearRe.FindString(parkrunDate)
	doc, err := client.Collection(dbCollectionName).Doc(yearKey).Get(ctx)
	var yearAlbumMap map[string]interface{}
	if status.Code(err) == codes.NotFound {
		log.Printf("no album for %v. Creating one", yearKey)
		yearAlbumMap = make(map[string]interface{})
	} else if err != nil {
		log.Fatal(err)
	} else {
		yearAlbumMap = doc.Data()
	}
	res := map[string]interface{}{
		"name": albumName,
		"date": parkrunDate,
		"num":  runNum,
		"url":  albumUrl,
	}
	yearAlbumMap[parkrunDate+"#"+strconv.FormatInt(runNum, 10)] = res

	result, err := client.Collection(dbCollectionName).Doc(yearKey).Set(ctx, yearAlbumMap)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("got result %v", result)
}
