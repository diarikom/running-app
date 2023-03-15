package service

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/internal/api/model"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"time"
)

func NewUserRepository(db *nsql.SqlDatabase, logger nlog.Logger) api.UserRepository {
	r := userRepository{
		Db:      db,
		Differs: initUserDiffer(),
		Stmt:    initUserStatement(db),
		Logger:  logger,
	}

	return &r
}

type userRepository struct {
	Db      *nsql.SqlDatabase
	Differs userDiffer
	Stmt    userStatements
	Logger  nlog.Logger
}

func (u *userRepository) IsUserHasSubscribed(userId string, providerId int8) (bool, error) {
	var result bool
	err := u.Stmt.isUserHasSubscribed.Get(&result, userId, providerId)
	return result, err
}

func (u *userRepository) InsertAdminInvitation(invitation model.AdminInvitation) error {
	_, err := u.Stmt.insertAdminInvitation.Exec(invitation)
	return err
}

func (u *userRepository) FindActiveSubscription(userId string, now time.Time) (*model.UserSubscription, error) {
	var result model.UserSubscription
	err := u.Stmt.findActiveSubscription.Get(&result, userId, now)
	return &result, err
}

func (u *userRepository) UpdateSubscriptionStatus(subscription model.UserSubscription) error {
	_, err := u.Stmt.updateSubscriptionStatus.Exec(subscription)
	return err
}

func (u *userRepository) FindLatestSubscriptionByUser(userId string, providerId int8, now time.Time) (
	*model.UserSubscription, error) {
	var result model.UserSubscription
	err := u.Stmt.findLatestSubscriptionByUser.Get(&result, userId, providerId, now)
	return &result, err
}

func (u *userRepository) InsertSubscription(subscription model.UserSubscription) error {
	_, err := u.Stmt.insertSubscription.Exec(subscription)
	return err
}

func (u *userRepository) InsertProviderRefId(mapping model.ProviderUserMapping) error {
	_, err := u.Stmt.insertProviderUserRefId.Exec(mapping)
	return err
}

func (u *userRepository) FindProviderSubscriptionPlanTypeRefId(providerId, subscriptionPlanTypeId int8) (string, error) {
	var refId string
	err := u.Stmt.findProviderSubscriptionPlanTypeRefId.Get(&refId, providerId, subscriptionPlanTypeId)
	return refId, err
}

func (u *userRepository) FindEmailById(userId string) (string, error) {
	var email string
	err := u.Stmt.findEmailById.Get(&email, userId)
	return email, err
}

func (u *userRepository) FindProviderRefId(providerId int8, userId string) (string, error) {
	var refId string
	err := u.Stmt.findProviderRefId.Get(&refId, providerId, userId)
	return refId, err
}

func (u *userRepository) DeleteSessionById(id string) error {
	_, err := u.Stmt.deleteSessionById.Exec(id)
	return err
}

func (u *userRepository) UpdateProfile(oldUser, newUser model.UserProfile, trackChanges []string) error {
	// Get differ
	differ := u.Differs.profile

	// Compare instance
	diff, err := differ.Compare(oldUser, newUser, trackChanges)
	if err != nil {
		return err
	}

	// If no changes, return
	if diff.Count == 0 {
		u.Logger.Debug("no user profile changes detected")
		return nil
	}

	// Generate query
	q, args, err := differ.UpdateQuery(diff)
	if err != nil {
		u.Logger.Error("unable to generate update query", err)
		return err
	}

	// Rebind query
	q = u.Db.Conn.Rebind(q)

	// Execute
	_, err = u.Db.Conn.Exec(q, args...)
	return err
}

func (u *userRepository) UpdateVerifyEmail(userId string, isVerified bool, timestamp time.Time) error {
	_, err := u.Stmt.updateEmailVerified.Exec(isVerified, timestamp, userId)
	return err
}

func (u *userRepository) UpdatePassword(userId, password string, timestamp time.Time) error {
	_, err := u.Stmt.updatePassword.Exec(password, timestamp, userId)
	return err
}

func (u *userRepository) FindAuthById(userId string) (*model.UserAuth, error) {
	var auth model.UserAuth
	err := u.Stmt.findAuthById.Get(&auth, userId)
	return &auth, err
}

func (u *userRepository) FindProfileById(userId string) (*model.UserProfile, error) {
	var profile model.UserProfile
	err := u.Stmt.findProfileById.Get(&profile, userId)
	return &profile, err
}

func (u *userRepository) FindSessionById(sessionId string) (*model.UserSession, error) {
	var session model.UserSession
	err := u.Stmt.findSessionById.Get(&session, sessionId)
	return &session, err
}

func (u *userRepository) DeleteAllSession(userId string) error {
	_, err := u.Stmt.deleteAllSession.Exec(userId)
	return err
}

func (u *userRepository) InsertSession(userSession model.UserSession) error {
	_, err := u.Stmt.insertUserSession.Exec(&userSession)
	return err
}

func (u *userRepository) FindProfileByEmail(email string) (*model.UserProfile, error) {
	var profile model.UserProfile
	err := u.Stmt.findProfileByEmail.Get(&profile, email)
	return &profile, err
}

func (u *userRepository) FindAuthByEmail(email string) (*model.UserAuth, error) {
	var auth model.UserAuth
	err := u.Stmt.findAuthByEmail.Get(&auth, email)
	return &auth, err
}

func (u *userRepository) FindAuthByThirdParty(userId string, providerId int) (*model.UserAuth, error) {
	var auth model.UserAuth
	err := u.Stmt.findAuthByThirdParty.Get(&auth, userId, providerId)
	return &auth, err
}

func (u *userRepository) Insert(userProfile model.UserProfile, userAuth model.UserAuth) error {
	// Begin transaction
	tx, err := u.Db.Conn.Beginx()
	if err != nil {
		return err
	}
	defer nsql.ReleaseTx(tx, &err, u.Logger)

	// Insert user
	_, err = tx.NamedStmt(u.Stmt.insertUserProfile).Exec(&userProfile)
	if err != nil {
		return err
	}

	// Insert user auth
	_, err = tx.NamedStmt(u.Stmt.insertUserAuth).Exec(&userAuth)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) InsertWithThirdParty(userProfile model.UserProfile, userAuth model.UserAuth, userAuthThirdParty model.UserAuthThirdParty) error {
	// Begin transaction
	tx, err := u.Db.Conn.Beginx()
	if err != nil {
		return err
	}
	defer nsql.ReleaseTx(tx, &err, u.Logger)

	// Insert user
	_, err = tx.NamedStmt(u.Stmt.insertUserProfile).Exec(&userProfile)
	if err != nil {
		return err
	}

	// Insert user auth
	_, err = tx.NamedStmt(u.Stmt.insertUserAuth).Exec(&userAuth)
	if err != nil {
		return err
	}

	// Insert user auth third party
	_, err = tx.NamedStmt(u.Stmt.insertUserAuthThirdParty).Exec(&userAuthThirdParty)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) IsExistByEmail(email string) (bool, error) {
	var isExist bool
	err := u.Stmt.isExistByEmail.Get(&isExist, email)
	return isExist, err
}

func (u *userRepository) IsExistBy3rdPartyAcc(authProviderId int, accessKey string) (bool, error) {
	var isExist bool
	err := u.Stmt.isExistBy3rdPartyAcc.Get(&isExist, accessKey, authProviderId)
	return isExist, err
}
