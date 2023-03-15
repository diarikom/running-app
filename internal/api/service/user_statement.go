package service

import (
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"github.com/jmoiron/sqlx"
)

type userStatements struct {
	deleteAllSession                      *sqlx.Stmt
	deleteSessionById                     *sqlx.Stmt
	findAuthByEmail                       *sqlx.Stmt
	findAuthById                          *sqlx.Stmt
	findAuthByThirdParty                  *sqlx.Stmt
	findEmailById                         *sqlx.Stmt
	findProfileByEmail                    *sqlx.Stmt
	findProfileById                       *sqlx.Stmt
	findProviderRefId                     *sqlx.Stmt
	findSessionById                       *sqlx.Stmt
	isExistByEmail                        *sqlx.Stmt
	isExistBy3rdPartyAcc                  *sqlx.Stmt
	insertUserAuth                        *sqlx.NamedStmt
	insertUserAuthThirdParty              *sqlx.NamedStmt
	insertUserProfile                     *sqlx.NamedStmt
	insertUserSession                     *sqlx.NamedStmt
	updateEmailVerified                   *sqlx.Stmt
	updatePassword                        *sqlx.Stmt
	findProviderSubscriptionPlanTypeRefId *sqlx.Stmt
	insertProviderUserRefId               *sqlx.NamedStmt
	insertSubscription                    *sqlx.NamedStmt
	findLatestSubscriptionByUser          *sqlx.Stmt
	updateSubscriptionStatus              *sqlx.NamedStmt
	findActiveSubscription                *sqlx.Stmt
	insertAdminInvitation                 *sqlx.NamedStmt
	isUserHasSubscribed                   *sqlx.Stmt
}

func initUserStatement(db *nsql.SqlDatabase) userStatements {
	return userStatements{
		deleteAllSession:                      db.Prepare(`DELETE FROM user_session WHERE user_id = $1`),
		deleteSessionById:                     db.Prepare(`DELETE FROM user_session WHERE id = $1`),
		findAuthByEmail:                       db.Prepare(`SELECT id, username, password, status_id, created_at, updated_at FROM user_auth WHERE username = $1`),
		findAuthById:                          db.Prepare(`SELECT id, username, password, status_id, created_at, updated_at FROM user_auth WHERE id = $1`),
		findAuthByThirdParty:                  db.Prepare(`SELECT ua.id, ua.username, ua.password, ua.status_id, ua.created_at, ua.updated_at FROM user_auth_third_party uatp INNER JOIN user_auth ua ON uatp.user_id = ua.id WHERE uatp.access_key = $1 AND uatp.auth_provider_id = $2`),
		findEmailById:                         db.Prepare(`SELECT email FROM user_profile WHERE id = $1`),
		findProfileByEmail:                    db.Prepare(`SELECT id, full_name, avatar_file, gender_id, date_of_birth, email, created_at, updated_at, email_verified FROM user_profile WHERE email = $1`),
		findProfileById:                       db.Prepare(`SELECT id, full_name, avatar_file, gender_id, date_of_birth, email, created_at, updated_at, email_verified FROM user_profile WHERE id = $1`),
		findProviderRefId:                     db.Prepare(`SELECT provider_ref FROM provider_user_mapping WHERE provider_id = $1 AND user_id = $2`),
		findSessionById:                       db.Prepare(`SELECT id, user_id, auth_provider_id, device_platform_id, device_id, device_manufacturer, device_model, notification_channel_id, notification_token, signature, expired_at, created_at, updated_at FROM user_session WHERE id = $1`),
		isExistByEmail:                        db.Prepare(`SELECT COUNT(id) > 0 "is_exist" FROM user_profile WHERE email = $1`),
		isExistBy3rdPartyAcc:                  db.Prepare(`SELECT COUNT(id) > 0 "is_exist" FROM user_auth_third_party WHERE access_key = $1 AND auth_provider_id = $2`),
		insertUserAuth:                        db.PrepareNamed(`INSERT INTO user_auth(id, username, password, status_id, created_at, updated_at) VALUES (:id, :username, :password, :status_id, :created_at, :updated_at)`),
		insertUserAuthThirdParty:              db.PrepareNamed(`INSERT INTO user_auth_third_party(id, user_id, auth_provider_id, access_key, created_at, updated_at) VALUES (:id, :user_id, :auth_provider_id, :access_key, :created_at, :updated_at)`),
		insertUserProfile:                     db.PrepareNamed(`INSERT INTO user_profile(id, full_name, avatar_file, gender_id, date_of_birth, email, created_at, updated_at) VALUES (:id, :full_name, :avatar_file, :gender_id, :date_of_birth, :email, :created_at, :updated_at)`),
		insertUserSession:                     db.PrepareNamed(`INSERT INTO user_session(id, user_id, auth_provider_id, device_platform_id, device_id, device_manufacturer, device_model, notification_channel_id, notification_token, signature, expired_at, created_at, updated_at) VALUES (:id, :user_id, :auth_provider_id, :device_platform_id, :device_id, :device_manufacturer, :device_model, :notification_channel_id, :notification_token, :signature, :expired_at, :created_at, :updated_at)`),
		updateEmailVerified:                   db.Prepare(`UPDATE user_profile SET email_verified = $1, updated_at = $2 WHERE id = $3`),
		updatePassword:                        db.Prepare(`UPDATE user_auth SET password = $1, updated_at = $2 WHERE id = $3`),
		findProviderSubscriptionPlanTypeRefId: db.Prepare(`SELECT provider_trx_ref FROM provider_subscription_plan WHERE provider_id = $1 AND plan_type_id = $2`),
		insertProviderUserRefId:               db.PrepareNamed(`INSERT INTO provider_user_mapping(id, user_id, provider_id, provider_ref, created_at, updated_at) VALUES (:id, :user_id, :provider_id, :provider_ref, :created_at, :updated_at)`),
		insertSubscription:                    db.PrepareNamed(`INSERT INTO user_subscription(id, user_id, plan_type_id, provider_id, provider_subscription_ref, provider_options, period_start, period_end, status_id, metadata, created_at, updated_at, modified_by) VALUES (:id, :user_id, :plan_type_id, :provider_id, :provider_subscription_ref, :provider_options, :period_start, :period_end, :status_id, :metadata, :created_at, :updated_at, :modified_by)`),
		findLatestSubscriptionByUser:          db.Prepare(`SELECT id, user_id, plan_type_id, provider_id, provider_subscription_ref, provider_options, period_start, period_end, status_id, metadata, created_at, updated_at, modified_by FROM user_subscription WHERE user_id = $1 AND provider_id = $2 AND period_end > $3 ORDER BY status_id, created_at DESC LIMIT 1`),
		updateSubscriptionStatus:              db.PrepareNamed(`UPDATE user_subscription SET status_id = :status_id, updated_at = :updated_at, modified_by = :modified_by WHERE id = :id`),
		findActiveSubscription:                db.Prepare(`SELECT id, user_id, plan_type_id, provider_id, provider_subscription_ref, provider_options, period_start, period_end, status_id, metadata, created_at, updated_at, modified_by FROM user_subscription WHERE user_id = $1 AND period_end > $2 ORDER BY status_id, created_at DESC LIMIT 1`),
		insertAdminInvitation:                 db.PrepareNamed(`INSERT INTO adm_invitation(id, user_id, email, token, expired_at, created_at) VALUES (:id, :user_id, :email, :token, :expired_at, :created_at);`),
		isUserHasSubscribed:                   db.Prepare(`SELECT COUNT(id) > 0 as has_subscribed FROM user_subscription WHERE user_id = $1 AND provider_id = $2 AND status_id IN (1,2,3) LIMIT 1;`),
	}
}
