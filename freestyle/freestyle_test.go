package freestyle

import (
	"mongo-crasher/common"
	"mongo-crasher/db"
	"sync"
	"testing"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func fetch(mongoDBClient *mgo.Session, wait *sync.WaitGroup) {
	mgoSession := mongoDBClient.Copy()
	var result []interface{}
	mgoSession.DB("O3O").C("CoCoChan").Find(bson.M{}).All(&result)
	mgoSession.Close()
	wait.Done()
}

func TestFreeStyle(t *testing.T) {
	var wait sync.WaitGroup
	var done chan bool
	mongoDBClient, _ := db.NewMongoDBClient("localhost:27017", "O3O", 128)
	common.CoreRunLoop(mongoDBClient, &wait, done, fetch)
	<-done
	mongoDBClient.Close()
}
