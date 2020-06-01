package common

const (
	// CSLoginCmd login cmd
	CSLoginCmd = "cs_login"
	// SSLoginCmd login cmd
	SSLoginCmd = "ss_login"
	// SSGetUserInfoCmd get userinfo
	SSGetUserInfoCmd = "ss_get_userinfo"
)

const (
	// ErrorCodeRedisKeyNotExist redis key not exist
	ErrorCodeRedisKeyNotExist = 1000
)

const (
	// CallIDLoginCmd login cmd
	CallIDLoginCmd = "login"
	// CallIDGetUserInfoCmd get userinfo
	CallIDGetUserInfoCmd = "getUser"
	//CallIDInsertUserInfoCmd insert user info
	CallIDInsertUserInfoCmd = "insertUser"
	//CallIDUpdateUserInfoCmd update user info
	CallIDUpdateUserInfoCmd = "updateUser"
)

//SnowFlake cmd
const (
	// CallIDGenUIDCmd create uniq uid cmd
	CallIDGenUIDCmd = "gen_uid"

	// CallIDGenUIDsCmd create uniq uid cmd
	CallIDGenUIDsCmd = "gen_uids"
)
