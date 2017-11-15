package worker

import (
	"mongo-crasher/common"
	"mongo-crasher/db"
	"sync"
	"testing"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type JobResult struct {
	Result []interface{}
	Err    error
}

type Job struct {
	Input  bson.M
	Result chan JobResult
}

var jobs = make(chan Job, 100)

func testWorker(mongoDBClient *mgo.Session, jobs <-chan Job) {
	for j := range jobs {
		mgoSession := mongoDBClient.Copy()
		var result []interface{}
		err := mgoSession.DB("O3O").C("CoCoChan").Find(j.Input).All(&result)
		mgoSession.Close()
		j.Result <- JobResult{result, err}
	}
}

func fetch(mongoDBClient *mgo.Session, wait *sync.WaitGroup) {
	input := bson.M{}
	result := make(chan JobResult)
	jobs <- Job{input, result}
	<-result
	wait.Done()
}

func TestWorker(t *testing.T) {
	var wait sync.WaitGroup
	var done chan bool
	mongoDBClient, _ := db.NewMongoDBClient("localhost:27017", "O3O", 128)
	for i := 0; i < 5; i++ {
		go testWorker(mongoDBClient, jobs)
	}
	common.CoreRunLoop(mongoDBClient, &wait, done, fetch)
	<-done
	mongoDBClient.Close()
}
