package cmd

//UserCenter module call cmd
const (
	// module http
	// HTTPWXLogin weixin login cmd
	HTTPWXLogin = "WXLogin"

	// module grpc

	// module logic
	// LogicWXLogin
	LogicWXLogin = "WXLogin"
	// LogicGetUIDByOpenID
	LogicGetUIDByOpenID = "GetUIDByOpenID"
	//LogicGetUser get user
	LogicGetUser = "GetUser"
	//LogicCreateUser create user
	LogicCreateUser = "CreateUser"
	// module wx
	WXLogin = "WXLogin"

	// module redis
	// RedisGetUserInfo
	RedisGetUserInfo = "GetUserInfo"
	// RedisUpdateUserInfo
	RedisUpdateUserInfo = "UpdateUserInfo"
	// RedisGetUIDByOpenID
	RedisGetUIDByOpenID = "GetUIDByOpenID"
	// RedisSetUIDByOpenID
	RedisSetUIDByOpenID = "SetUIDByOpenID"

	// module db
	// DBGetUIDByOpenID get uid by openid
	DBGetUIDByOpenID = "GetUserByOpenID"
	// DBInsertOID2UID  make openid to uid mapping
	DBInsertOID2UID = "InsertOID2UID"
	// DBCreateUser
	DBCreateUser = "CreateUser"
	// DBGetUserByUID
	DBGetUserByUID = "GetUserByUID"

	//module snowflake
	//SFGetUID
	SFGetUID = "GetUID"
)
