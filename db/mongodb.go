package db

import (
	"time"

	mgoUtil "github.com/mongodb/mongo-tools/common/util"
	mgo "gopkg.in/mgo.v2"
)

// NewMongoDBClient 建立一顆新芒果
func NewMongoDBClient(mgoEndpoints string, db string, workerSize int) (*mgo.Session, error) {
	mgoServerEndpoint, setName := mgoUtil.ParseConnectionString(mgoEndpoints)
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:     mgoServerEndpoint,
		Timeout:   60 * time.Second,
		Database:  db,
		PoolLimit: workerSize,
	}
	if setName != "" {
		mongoDBDialInfo.ReplicaSetName = setName
	}
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		return nil, err
	}
	session.SetSocketTimeout(1 * time.Hour)
	session.SetMode(mgo.Eventual, false)
	return session, err
}
