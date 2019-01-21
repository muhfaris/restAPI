package router

import (
	"github.com/globalsign/mgo"
	"gitlab.com/muhfaris/restAPI/internal/pkg/logging"
)

var (
	logger *logging.Logger
	dbPool *mgo.Database
)

func Init(lg *logging.Logger, db *mgo.Database) {
	logger = lg
	dbPool = db
}
