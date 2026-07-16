// internal/service/user_service.go
package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"kakao-clone/internal/model"
	"kakao-clone/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// ─────────────────────────────────────────
// 회원가입
// ─────────────────────────────────────────
func (s *UserService) Register(email, password, name string) (*model.User, error) {
	// 이메일 중복 확인
	existing, _ := s.repo.FindByEmail(email)
	if existing != nil {
		return nil, ErrBadRequest("이미 사용 중인 이메일입니다")
	}

	// 비밀번호 해싱
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, ErrInternal("비밀번호 처리 실패", err)
	}

	hashedStr := string(hashed)
	user := &model.User{
		Email:    email,
		Password: &hashedStr,
		Name:     name,
		Provider: "local",
	}

	if err := s.repo.Create(user); err != nil {
		return nil, ErrInternal("회원가입 실패", err)
	}

	return user, nil
}

// ─────────────────────────────────────────
// 로그인
// ─────────────────────────────────────────
func (s *UserService) Login(email, password string) (*model.User, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", ErrUnauthorized("이메일 또는 비밀번호가 틀렸습니다")
		}
		return nil, "", ErrInternal("유저 조회 실패", err)
	}

	if user.Password == nil {
		return nil, "", ErrUnauthorized("소셜 로그인 계정입니다")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(*user.Password), []byte(password),
	); err != nil {
		return nil, "", ErrUnauthorized("이메일 또는 비밀번호가 틀렸습니다")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// ─────────────────────────────────────────
// Google 로그인 (소셜)
// ─────────────────────────────────────────
func (s *UserService) GoogleLogin(email, name, googleID string) (*model.User, string, error) {
	user, err := s.repo.FindByEmail(email)

	if err != nil {
		// 없으면 자동 생성
		user = &model.User{
			Email:    email,
			Name:     name,
			Provider: "google",
			GoogleID: googleID,
			Password: nil,
		}
		if err := s.repo.Create(user); err != nil {
			return nil, "", ErrInternal("유저 생성 실패", err)
		}
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// ─────────────────────────────────────────
// 프로필 수정
// ─────────────────────────────────────────
func (s *UserService) UpdateProfile(userID uint, name, statusMsg string) (*model.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, ErrNotFound("유저를 찾을 수 없습니다")
	}

	user.Name = name
	user.StatusMsg = statusMsg

	if err := s.repo.Update(user); err != nil {
		return nil, ErrInternal("프로필 수정 실패", err)
	}

	return user, nil
}

// ─────────────────────────────────────────
// 유저 검색
// ─────────────────────────────────────────
func (s *UserService) Search(keyword string, myID uint) ([]model.User, error) {
	users, err := s.repo.Search(keyword)
	if err != nil {
		return nil, ErrInternal("검색 실패", err)
	}

	// 나 자신 제외
	result := make([]model.User, 0)
	for _, u := range users {
		if u.ID != myID {
			result = append(result, u)
		}
	}

	return result, nil
}

// ─────────────────────────────────────────
// JWT 발급 (내부 함수)
// ─────────────────────────────────────────
func (s *UserService) generateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", ErrInternal("토큰 발급 실패", err)
	}
	return tokenStr, nil
}
