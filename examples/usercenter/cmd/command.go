package cmd

//UserCenter module call cmd
const (
	// module http
	// MCHTTPWXLogin weixin login cmd
	MCHTTPWXLogin = "wxlogin"

	// module grpc

	// module logic
	// MCLogicWXLogin
	MCLogicWXLogin = MCHTTPWXLogin

	// module wx
	MCWXLogin = "wxlogin"

	// module redis
	// MCRedisGetUserInfo
	MCRedisGetUserInfo = "getUserInfo"
	// MCRedisUpdateUserInfo
	MCRedisUpdateUserInfo = "updateUserInfo"

	// module db
	// MCDBGetUIDByOpenID get uid by openid
	MCDBGetUIDByOpenID = "getUserByOpenID"
	// MCDBInsertOID2UID  make openid to uid mapping
	MCDBInsertOID2UID = "insertOID2UID"
	// MCDBCreateUser
	MCDBCreateUser = "createUser"
	// MCDBGetUserByUID
	MCDBGetUserByUID = "getUserByUID"
)
