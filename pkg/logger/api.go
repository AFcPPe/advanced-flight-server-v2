package logger

import "go.uber.org/zap"

// Debug 输出 debug 级别日志
func Debug(msg string, fields ...zap.Field) {
	L().Debug(msg, fields...)
}

// Info 输出 info 级别日志
func Info(msg string, fields ...zap.Field) {
	L().Info(msg, fields...)
}

// Warn 输出 warn 级别日志
func Warn(msg string, fields ...zap.Field) {
	L().Warn(msg, fields...)
}

// Error 输出 error 级别日志
func Error(msg string, fields ...zap.Field) {
	L().Error(msg, fields...)
}

// Fatal 输出 fatal 级别日志并退出程序
func Fatal(msg string, fields ...zap.Field) {
	L().Fatal(msg, fields...)
}

// Debugf 格式化输出 debug 级别日志
func Debugf(template string, args ...interface{}) {
	S().Debugf(template, args...)
}

// Infof 格式化输出 info 级别日志
func Infof(template string, args ...interface{}) {
	S().Infof(template, args...)
}

// Warnf 格式化输出 warn 级别日志
func Warnf(template string, args ...interface{}) {
	S().Warnf(template, args...)
}

// Errorf 格式化输出 error 级别日志
func Errorf(template string, args ...interface{}) {
	S().Errorf(template, args...)
}

// Fatalf 格式化输出 fatal 级别日志并退出程序
func Fatalf(template string, args ...interface{}) {
	S().Fatalf(template, args...)
}

// With 创建带有额外字段的 logger
func With(fields ...zap.Field) *zap.Logger {
	return L().With(fields...)
}

// WithOptions 创建带有额外选项的 logger
func WithOptions(opts ...zap.Option) *zap.Logger {
	return L().WithOptions(opts...)
}
