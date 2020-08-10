package cmd

//UserCenter module call cmd
const (
	// module http
	// MCHTTPWXLogin weixin login cmd
	MCHTTPWXLogin = "WXLogin"

	// module grpc

	// module logic
	// MCLogicWXLogin
	MCLogicWXLogin = "WXLogin"
	// MCLogicGetUIDByOpenID
	MCLogicGetUIDByOpenID = "GetUIDByOpenID"

	// module wx
	MCWXLogin = "WXLogin"

	// module redis
	// MCRedisGetUserInfo
	MCRedisGetUserInfo = "GetUserInfo"
	// MCRedisUpdateUserInfo
	MCRedisUpdateUserInfo = "UpdateUserInfo"
	// MCRedisGetUIDByOpenID
	MCRedisGetUIDByOpenID = "GetUIDByOpenID"
	// MCRedisSetUIDByOpenID
	MCRedisSetUIDByOpenID = "SetUIDByOpenID"

	// module db
	// MCDBGetUIDByOpenID get uid by openid
	MCDBGetUIDByOpenID = "GetUserByOpenID"
	// MCDBInsertOID2UID  make openid to uid mapping
	MCDBInsertOID2UID = "InsertOID2UID"
	// MCDBCreateUser
	MCDBCreateUser = "CreateUser"
	// MCDBGetUserByUID
	MCDBGetUserByUID = "GetUserByUID"
)
