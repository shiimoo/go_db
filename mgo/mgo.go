package mgo

// TODO: 同意通过Mgr调度, 该文件仅提供默认的实际操作接口
// TODO：mgr的默认数据库在创建时进行指定？亦或者回调操作? 该方案待定

// HasCollection 判定数据库database中是否存在集合collection
func HasCollection(database, collection string) bool {
	return getMgr().HasCollection(database, collection)
}

// 创建索引
func CreateIndex(database, collection string, index Indexs) {
	getMgr().CreateIndex(database, collection, index)
}

// 增删改查
