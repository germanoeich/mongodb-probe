package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func getMongoURI(node string) string {
	uri := "mongodb://"
	if MongoUser != "" {
		uri += MongoUser + ":" + MongoPass + "@"
	}

	uri += node + "/?directConnect=true"
	return uri
}

type ReplSetGetStatusMemberOptime struct {
	Timestamp *primitive.Timestamp `bson:"ts,omitempty"`
}

type ReplSetGetStatusMember struct {
	Name string `bson:"name,omitempty"`
	Uptime int32 `bson:"uptime,omitempty"`
	Optime ReplSetGetStatusMemberOptime `bson:"optime,omitempty"`
	Self bool `bson:"self,omitempty"`
	State int32 `bson:"state,omitempty"`
}

type ReplSetGetStatusResult struct {
	MyState int32 `bson:"myState,omitempty"`
	Members []ReplSetGetStatusMember `bson:"members,omitempty"`
}

func check() {
	for _, node := range MongoNodes {
		uri := getMongoURI(node)

		opts := options.Client().ApplyURI(uri).SetDirect(true)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, opts)
		if err != nil {
			logrus.Error(err)
			ErrorCounter.Inc()
			StatusGauge.With(map[string]string{
				"node": node,
			}).Set(-1)
			continue
		}

		logrus.Trace("Connected to " + node)

		var replStatus ReplSetGetStatusResult
		err = client.Database("admin").RunCommand(ctx, bson.D{{"replSetGetStatus", "1" }}).Decode(&replStatus)

		if err != nil {
			logrus.Error(err)
			ErrorCounter.Inc()
			StatusGauge.With(map[string]string{
				"node": node,
			}).Set(-1)
			continue
		}

		var myOptime int64
		var primaryOptime int64

		for _, e := range replStatus.Members {
			if e.Self && e.Optime.Timestamp != nil {
				myOptime = int64(e.Optime.Timestamp.T)
			}

			if e.State == 1 && e.Optime.Timestamp != nil{
				primaryOptime = int64(e.Optime.Timestamp.T)
			}
		}

		var lag float64 = -10
		if myOptime != 0 && primaryOptime != 0 {
			lag = float64(myOptime - primaryOptime)
		}

		ReplicationLagGauge.With(map[string]string{
			"node": node,
		}).Set(lag)

		StatusGauge.With(map[string]string{
			"node": node,
		}).Set(float64(replStatus.MyState))
	}
}

func tick() {
	t := time.NewTicker(30 * time.Second)

	for range t.C {
		check()
	}
}