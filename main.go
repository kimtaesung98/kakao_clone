package main

import (
	"fmt"
	"kakao-clone/internal/service"
)

func main() {
	svc := service.NewUserService()

	// 존재하지 않는 유저 조회
	_, err := svc.GetByID(999)
	if err != nil {
		// 에러 메시지 출력
		fmt.Println(err) // [404] 유저를 찾을 수 없습니다

		// AppError인지 확인
		if appErr, ok := service.AsAppError(err); ok {
			fmt.Println("HTTP 코드:", appErr.Code) // 404
			fmt.Println("메시지:", appErr.Message)  // 유저를 찾을 수 없습니다
		}
	}
}
