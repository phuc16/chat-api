package errors

import "app/pkg/apperror"

const (
	CodeUserError = 20000 + iota
	CodeUserNotFound
	CodeUserExists
	CodeUserNameExists
	CodePhoneNumberExists
	CodeUserEmailExists
	CodePasswordIncorrect
	CodeUserNotRegister
	CodeUserInactive
)

func UserNotFound() *apperror.Error {
	return apperror.NewError(CodeUserNotFound, "Người dùng không tồn tại")
}

func UserExists() *apperror.Error {
	return apperror.NewError(CodeUserExists, "Người dùng đã tồn tại")
}

func UserNameExists() *apperror.Error {
	return apperror.NewError(CodeUserNameExists, "Tên đăng nhập đã tồn tại")
}

func PhoneNumberExists() *apperror.Error {
	return apperror.NewError(CodePhoneNumberExists, "Số điện thoại đã tồn tại")
}

func UserEmailExists() *apperror.Error {
	return apperror.NewError(CodeUserEmailExists, "Email đã tồn tại")
}

func PasswordIncorrect() *apperror.Error {
	return apperror.NewError(CodePasswordIncorrect, "Mật khẩu không chính xác")
}

func UserNotRegister() *apperror.Error {
	return apperror.NewError(CodeUserNotRegister, "Không tìm thấy tài khoản, vui lòng đăng ký trước khi tiếp tục")
}
