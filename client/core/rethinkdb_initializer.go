package core

import (
	log "github.com/sirupsen/logrus"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

var (
	DbName                             = "rma_subscribes"
	TableName_SubscribeTunnelRealFund  = "SubscribeTunnelRealFund"
	TableName_SubscribeCorpHoldMon     = "SubscribeCorpHoldMon"
	TableName_SubscribeQuoteMon        = "SubscribeQuoteMon"
	TableName_SubscribeCustRisk        = "SubscribeCustRisk"
	TableName_SubscribeCustHold        = "SubscribeCustHold"
	TableName_SubscribeCustGroupHold   = "SubscribeCustGroupHold"
	TableName_SubscribeProuctGroupRisk = "SubscribeProuctGroupRisk"
	TableName_SubscribeNearDediveHold  = "SubscribeNearDediveHold"
)

type RethinkDB interface {
	Init() (*r.Session, error)
}

type rethinkdb struct {
	address string
}

func NewRethinkDB(address string) RethinkDB {
	return &rethinkdb{
		address,
	}
}

func (i *rethinkdb) Init() (*r.Session, error) {
	session, err := r.Connect(r.ConnectOpts{
		Addresses: i.address, // endpoint without http
	})
	if err != nil {
		return nil, err
	}

	err = createDbAndTable(session)
	if err != nil {
		return nil, err
	}
	session.Use(DbName)

	return session, nil
}

func createDbAndTable(session *r.Session) error {
	err := createDb(session, DbName)
	if err != nil {
		log.Error(err)
		return err
	}
	session.Use(DbName)

	primaryKey := "ActionKey"
	err = createTable(session, TableName_SubscribeTunnelRealFund, primaryKey)
	if err != nil {
		log.Error(err)
		return err
	}
	err = createTable(session, TableName_SubscribeCorpHoldMon, primaryKey)
	if err != nil {
		log.Error(err)
		return err
	}
	err = createTable(session, TableName_SubscribeQuoteMon, primaryKey)
	if err != nil {
		log.Error(err)
		return err
	}
	err = createTable(session, TableName_SubscribeCustRisk, primaryKey)
	if err != nil {
		log.Error(err)
		return err
	}
	err = createTable(session, TableName_SubscribeCustHold, primaryKey)
	if err != nil {
		log.Error(err)
		return err
	}
	err = createTable(session, TableName_SubscribeCustGroupHold, primaryKey)
	if err != nil {
		log.Error(err)
		return err
	}
	err = createTable(session, TableName_SubscribeProuctGroupRisk, primaryKey)
	if err != nil {
		log.Error(err)
		return err
	}
	err = createTable(session, TableName_SubscribeNearDediveHold, primaryKey)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func createDb(session *r.Session, dbName string) error {
	res, err := r.DBList().Run(session)
	if err != nil {
		return err
	}
	defer res.Close()

	var rows []string
	err = res.All(&rows)
	if err != nil {
		return err
	}

	existed := false
	for _, d := range rows {
		if d == dbName {
			existed = true
			break
		}
	}
	if !existed {
		_, err = r.DBCreate(dbName).Run(session)
		if err != nil {
			return err
		}
		log.Infof("DB %s created.", dbName)

	} else {
		log.Infof("DB %s existed.", dbName)
	}

	return nil
}

func createTable(session *r.Session, tableName string, primaryKey string) error {
	res, err := r.TableList().Run(session)
	if err != nil {
		return err
	}
	defer res.Close()

	var rows []string
	err = res.All(&rows)
	if err != nil {
		return err
	}

	existed := false
	for _, t := range rows {
		if t == tableName {
			existed = true
			break
		}
	}
	if !existed {
		_, err = r.TableCreate(tableName, r.TableCreateOpts{PrimaryKey: primaryKey}).Run(session)
		if err != nil {
			return err
		}
		log.Infof("Table %s created.", tableName)

	} else {
		log.Infof("Table %s existed.", tableName)
	}

	return nil
}
