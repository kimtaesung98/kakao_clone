// internal/service/errors.go
package service

import (
	"errors"
	"fmt"
	"net/http"
)

// ─────────────────────────────────────────
// AppError — 커스텀 에러 타입
// ─────────────────────────────────────────
type AppError struct {
	Code    int    // HTTP 상태코드
	Message string // 클라이언트 응답 메시지
	Err     error  // 원본 에러 (서버 로그용)
}

// Error() — error interface 구현
// 이 메서드가 있으면 error 타입으로 쓸 수 있음
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap() — errors.Is/As 로 원본 에러 접근 가능하게
func (e *AppError) Unwrap() error {
	return e.Err
}

// ─────────────────────────────────────────
// 자주 쓰는 에러 생성 함수
// ─────────────────────────────────────────
func ErrNotFound(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg}
}

func ErrUnauthorized(msg string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: msg}
}

func ErrBadRequest(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func ErrInternal(msg string, err error) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg, Err: err}
}

// ─────────────────────────────────────────
// 에러 종류 확인
// errors.As — 특정 타입인지 확인
// Node.js: err instanceof AppError 와 동일
// ─────────────────────────────────────────
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
