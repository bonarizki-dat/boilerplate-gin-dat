package mocks

import (
	"sync"

	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/domain/models"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/domain/repositories"
)

// MockUserRepository is an in-memory UserRepository for unit tests.
// Use AddUserByEmail / SetUserByRefreshToken / SetUserByPasswordResetToken to control responses.
type MockUserRepository struct {
	mu sync.RWMutex
	// key: email
	byEmail map[string]*models.User
	// key: refresh token
	byRefreshToken map[string]*models.User
	// key: password reset token
	byResetToken map[string]*models.User
}

// NewMockUserRepository returns a new MockUserRepository.
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		byEmail:        make(map[string]*models.User),
		byRefreshToken: make(map[string]*models.User),
		byResetToken:   make(map[string]*models.User),
	}
}

// Ensure MockUserRepository implements repositories.UserRepository.
var _ repositories.UserRepository = (*MockUserRepository)(nil)

// AddUserByEmail makes GetUserByEmail(email) return the given user (copy stored).
func (m *MockUserRepository) AddUserByEmail(email string, user *models.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.byEmail == nil {
		m.byEmail = make(map[string]*models.User)
	}
	u := *user
	u.Email = email
	m.byEmail[email] = &u
}

// SetUserByRefreshToken makes GetUserByRefreshToken(token) return the given user.
func (m *MockUserRepository) SetUserByRefreshToken(token string, user *models.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.byRefreshToken == nil {
		m.byRefreshToken = make(map[string]*models.User)
	}
	m.byRefreshToken[token] = user
}

// SetUserByPasswordResetToken makes GetUserByPasswordResetToken(token) return the given user.
func (m *MockUserRepository) SetUserByPasswordResetToken(token string, user *models.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.byResetToken == nil {
		m.byResetToken = make(map[string]*models.User)
	}
	m.byResetToken[token] = user
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if u, ok := m.byEmail[email]; ok {
		cp := *u
		return &cp, nil
	}
	return nil, nil
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.byEmail == nil {
		m.byEmail = make(map[string]*models.User)
	}
	u := *user
	m.byEmail[user.Email] = &u
	return nil
}

func (m *MockUserRepository) UpdateUser(user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.byEmail != nil {
		m.byEmail[user.Email] = user
	}
	if user.RefreshToken != "" {
		if m.byRefreshToken == nil {
			m.byRefreshToken = make(map[string]*models.User)
		}
		m.byRefreshToken[user.RefreshToken] = user
	}
	if user.PasswordResetToken != "" {
		if m.byResetToken == nil {
			m.byResetToken = make(map[string]*models.User)
		}
		m.byResetToken[user.PasswordResetToken] = user
	}
	return nil
}

func (m *MockUserRepository) GetUserByRefreshToken(token string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if u, ok := m.byRefreshToken[token]; ok {
		cp := *u
		return &cp, nil
	}
	return nil, nil
}

func (m *MockUserRepository) GetUserByPasswordResetToken(token string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if u, ok := m.byResetToken[token]; ok {
		cp := *u
		return &cp, nil
	}
	return nil, nil
}
