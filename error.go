package lightinggo

import (
	"github.com/pkg/errors"
)

// CustomError return custom error.
type CustomError struct {
	Code    ErrCode `json:"code"`    // 业务码
	Message string  `json:"message"` // 描述信息
}

// Error method override.
func (e *CustomError) Error() string {
	return e.Code.String()
}

type ErrCode int64 //错误码

// define errorCode
//go:generate stringer -type ErrCode -linecomment
const (
	// 业务
	NotReadyError ErrCode = 10001 // 服务未就绪
	NotFoundError ErrCode = 10002 // 未找到错误
	ExpiredError  ErrCode = 10003 // 已失效错误

	// 服务级错误码
	ServerError        ErrCode = 10101 // Internal Server Error
	TooManyRequests    ErrCode = 10102 // Too Many Requests
	ParamBindError     ErrCode = 10103 // 参数信息有误
	AuthorizationError ErrCode = 10104 // 签名信息有误
	CallHTTPError      ErrCode = 10105 // 调用第三方HTTP接口失败
	ResubmitError      ErrCode = 10106 // ResubmitError
	ResubmitMsg        ErrCode = 10107 // 请勿重复提交
	HashIdsDecodeError ErrCode = 10108 // ID参数有误
	SignatureError     ErrCode = 10109 // SignatureError

	// 业务模块级错误码
	// 用户模块
	IllegalUserName ErrCode = 20101 // 非法用户名
	UserCreateError ErrCode = 20102 // 创建用户失败
	UserUpdateError ErrCode = 20103 // 更新用户失败
	UserSearchError ErrCode = 20104 // 查询用户失败

	// 配置
	ConfigParseError        ErrCode = 20201 // 配置文件解析失败
	ConfigSaveError         ErrCode = 20202 // 配置文件写入失败
	ConfigRedisConnectError ErrCode = 20203 // Redis连接失败
	ConfigMySQLConnectError ErrCode = 20204 // MySQL连接失败
	ConfigMySQLInitError    ErrCode = 20205 // MySQL初始化数据失败
	ConfigGoVersionError    ErrCode = 20206 // GoVersion不满足要求

	// 实用工具箱
	SearchRedisError ErrCode = 20501 // 查询RedisKey失败
	ClearRedisError  ErrCode = 20502 // 清空RedisKey失败
	SearchRedisEmpty ErrCode = 20503 // 查询的RedisKey不存在
	SearchMySQLError ErrCode = 20504 // 查询mysql失败

	// 菜单栏
	MenuCreateError ErrCode = 20601 // 创建菜单失败
	MenuUpdateError ErrCode = 20602 // 更新菜单失败
	MenuListError   ErrCode = 20603 // 删除菜单失败
	MenuDeleteError ErrCode = 20604 // 获取菜单列表页失败
	MenuDetailError ErrCode = 20605 // 获取菜单详情失败
)

func Text(code ErrCode) string {
	return code.String()
}

// InitError custom error initialization.
func InitError(code ErrCode) error {
	// 初次调用得用Wrap方法，进行实例化
	return errors.Wrap(&CustomError{
		Code:    code,
		Message: code.String(),
	}, "")
}
