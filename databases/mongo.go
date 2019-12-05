package databases

import "gopkg.in/mgo.v2"

var MgoSession *mgo.Session

func init() {
	var err error
	MgoSession, err = mgo.Dial("mongo")
	if err != nil {
		panic(err)
	}
}