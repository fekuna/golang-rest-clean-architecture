package repository

import (
	"context"
	"fmt"

	"github.com/fekuna/go-rest-clean-architecture/internal/auth"
	"github.com/fekuna/go-rest-clean-architecture/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type authRepo struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) auth.Repository {
	return &authRepo{
		db: db,
	}
}

func (r *authRepo) Register(ctx context.Context, user *models.User) (*models.User, error) {
	// TODO: Tracing
	u := &models.User{}
	if err := r.db.QueryRowxContext(
		ctx, createUserQuery, &user.FirstName, &user.LastName, &user.Email,
		&user.Password, &user.Role, &user.About, &user.AvatarID, &user.PhoneNumber, &user.Address, &user.City,
		&user.Gender, &user.Postcode, &user.Birthday,
	).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "authRepo.Register.StructScan")
	}

	return u, nil
}

func (r *authRepo) FindByEmail(ctx context.Context, user *models.User) (*models.User, error) {
	// TODO: findByEmail

	foundUser := &models.User{}
	if err := r.db.QueryRowxContext(ctx, findUserByEmail, user.Email).StructScan(foundUser); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByEmail.QueryRowxContext")
	}

	return foundUser, nil
}

func (r *authRepo) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	foundUser := &models.User{}
	if err := r.db.QueryRowxContext(ctx, getUserById, userID).StructScan(foundUser); err != nil {
		return nil, errors.Wrap(err, "authRepo.GetUserByID.QueryRowxContext")
	}

	return foundUser, nil
}

func (r *authRepo) Update(ctx context.Context, user *models.User) (*models.User, error) {
	// TODO: Tracing

	u := &models.User{}
	if err := r.db.GetContext(ctx, u, updateUserQuery, &user.FirstName, &user.LastName, &user.Email,
		&user.Role, &user.About, &user.AvatarID, &user.PhoneNumber, &user.Address, &user.City, &user.Gender,
		&user.Postcode, &user.Birthday, &user.UserID,
	); err != nil {
		return nil, errors.Wrap(err, "authRepo.Update.GetContext")
	}

	return u, nil
}

func (r *authRepo) FindAvatarByFilePath(ctx context.Context, filePath string) (*models.Avatar, error) {
	foundAvatar := &models.Avatar{}
	if err := r.db.QueryRowxContext(ctx, findAvatarByFilePath, filePath).StructScan(foundAvatar); err != nil {
		return nil, errors.Wrap(err, "authRepo.findAvatarByFilePath.QueryRowxContext")
	}

	return foundAvatar, nil
}

func (r *authRepo) FindAvatarByID(ctx context.Context, avatarID uuid.UUID) (*models.Avatar, error) {
	foundAvatar := &models.Avatar{}
	if err := r.db.QueryRowxContext(ctx, findAvatarByID, avatarID).StructScan(foundAvatar); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindAvatarByID.QueryRowxContext")
	}

	return foundAvatar, nil
}

func (r *authRepo) CreateAvatar(ctx context.Context, avatar *models.Avatar) (*models.Avatar, error) {
	// TODO: Tracing
	a := &models.Avatar{}
	if err := r.db.QueryRowxContext(
		ctx, createAvatar, &avatar.Bucket, &avatar.FilePath).StructScan(a); err != nil {
		fmt.Println("mashok nich")
		return nil, errors.Wrap(err, "authRepo.CreateAvatar.StructScan")
	}

	return a, nil
}
