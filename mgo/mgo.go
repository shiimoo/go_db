package mgo

func SetDatabase(database string) {
	getMgr().SetDatabase(database)
}

// 判定数据库database中是否存在集合collection
func HasCollection(collection string) bool {
	m := getMgr()
	return m.GetConn().hasCollection(m.database, collection)
}

// 增删改查
