package db

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/globalsign/mgo"
)

//defining the collection name as article
const (
	User     = `user`
	Article  = `article`
	Comments = `comments`
)

var (
	// Obj defines the mongodb session, which connects mongodb instance.
	Obj  *mgo.Session
	once sync.Once
)

// ConnectDB connect to specified MongoDB instance.
// Obj is a Singleton DB object
func ConnectDB() {
	once.Do(func() {
		Obj = connectLocalDB()
	})
}

// connectLocalDB connects to local MongoDB for development and testing
func connectLocalDB() *mgo.Session {

	dialInfo := &mgo.DialInfo{
		Addrs:         []string{"127.0.0.1:27017"},
		Timeout:       30 * time.Second,
		PoolLimit:     10, // per node
		MinPoolSize:   50,
		PoolTimeout:   time.Minute * 10,
		MaxIdleTimeMS: 30000,
	}
	log.Printf("INFO: Mgo Dialinfo: %v\n", dialInfo)

	//Dialing MongoDB with DialInfo
	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		log.Printf("Error in dialing mongo server: %s", err.Error())
		os.Exit(2)
	}

	log.Printf("INFO: Mgo connected to : 127.0.0.1:27017")

	session.SetMode(mgo.Primary, true)
	session.SetSafe(&mgo.Safe{})

	return session
}
