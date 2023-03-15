package main

import (
	"github.com/diarikom/running-app/running-app-api/internal/api"
	"github.com/diarikom/running-app/running-app-api/pkg/nhttp"
)

const (
	// Middlewares
	AuthClientMiddleware          = "auth.client"
	AuthClientDashboardMiddleware = "auth.client.dash"
	AuthUserMiddleware            = "auth.user"
	ResetPasswordMiddleware       = "auth.one_time.reset_password"
	VerifyEmailMiddleware         = "auth.one_time.verify_email"
)

func initMiddlewares(router *nhttp.Router, services *api.Services) {
	router.RegisterMiddleware(AuthClientMiddleware, nhttp.NewClientAuthMiddleware(services.Auth.ValidateClient,
		nhttp.KeyAuthorization, log))
	router.RegisterMiddleware(AuthClientDashboardMiddleware, nhttp.NewClientAuthMiddleware(services.Auth.ValidateClientDashboard,
		nhttp.KeyAuthorization, log))
	router.RegisterMiddleware(AuthUserMiddleware, nhttp.NewUserAuthMiddleware(services.Auth.ValidateUserAccess,
		services.User.ValidateSession, nhttp.KeyAuthorization, log))
	router.RegisterMiddleware(ResetPasswordMiddleware, api.NewResetPasswordSessionMiddleware(
		services.Auth.ValidateResetPasswordToken, services.User.ValidateResetPasswordSignature, nhttp.KeyAuthorization, log))
	router.RegisterMiddleware(VerifyEmailMiddleware, api.NewVerifyEmailSessionMiddleware(
		services.Auth.ValidateVerifyEmailToken, services.User.ValidateVerifyEmailSignature, nhttp.KeyAuthorization, log))
}

func initRoutes(router *nhttp.Router, handlers Handlers) {
	// API Common
	router.Handle("", handlers.ApiStatus).Methods("GET")

	// Assets
	router.HandleWithMiddleware("/assets", AuthUserMiddleware, handlers.Asset.PostUploadFile).Methods("POST")

	// Discover Contents
	router.HandleWithMiddleware("/discover-contents", AuthUserMiddleware, handlers.DiscoverContent.GetContents).Methods("GET")
	router.HandleWithMiddleware("/discover/ad-tags", AuthUserMiddleware, handlers.AdTag.GetTags).Methods("GET")

	// Runs
	router.HandleWithMiddleware("/run-sessions", AuthUserMiddleware, handlers.Run.GetRunSessionHistory).Methods("GET")
	router.HandleWithMiddleware("/run-sessions/summary", AuthUserMiddleware, handlers.Run.PostRunSummary).Methods("POST")

	// Users
	router.Handle("/users/register", handlers.User.PostRegister).Methods("POST")
	router.Handle("/users/log-in", handlers.User.PostLogin).Methods("POST")
	router.Handle("/users/email", handlers.User.GetCheckEmail).Methods("GET")
	router.Handle("/users/reset-password", handlers.User.PostResetPassword).Methods("POST")
	router.Handle("/users/reset-password", handlers.User.GetResetPassword).Methods("GET")
	router.Handle("/users/verify-email", handlers.User.GetVerifyEmail).Methods("GET")
	router.HandleWithMiddleware("/users/log-out", AuthUserMiddleware, handlers.User.DeleteLogout).Methods("DELETE")
	router.HandleWithMiddleware("/users/profile", AuthUserMiddleware, handlers.User.GetProfile).Methods("GET")
	router.HandleWithMiddleware("/users/change-password", AuthUserMiddleware, handlers.User.PutChangePassword).Methods("PUT")
	router.HandleWithMiddleware("/users/reset-password", ResetPasswordMiddleware, handlers.User.PutResetPassword).Methods("PUT")
	router.HandleWithMiddleware("/users/verify-email", VerifyEmailMiddleware, handlers.User.PutVerifyEmail).Methods("PUT")
	router.HandleWithMiddleware("/users/refresh-session", AuthUserMiddleware, handlers.User.PostRefreshToken).Methods("PUT")
	router.HandleWithMiddleware("/users/credits", AuthUserMiddleware, handlers.User.GetCreditBalance).Methods("GET")
	router.HandleWithMiddleware("/users/donations", AuthUserMiddleware, handlers.Initiative.ListUserDonation).Methods("GET")
	router.HandleWithMiddleware("/users/providers/{providerId}/ref-id", AuthUserMiddleware, handlers.User.GetUserProviderRefId).Methods("GET")
	router.HandleWithMiddleware("/users/subscriptions", AuthUserMiddleware, handlers.User.PostUserSubscribe).Methods("POST")
	router.HandleWithMiddleware("/users/subscriptions", AuthUserMiddleware, handlers.User.DeleteUserCancelSubscription).Methods("DELETE")
	router.HandleWithMiddleware("/users/subscriptions", AuthUserMiddleware, handlers.User.GetUserSubscription).Methods("GET")
	router.HandleWithMiddleware("/users/subscriptions/vouchers", AuthUserMiddleware, handlers.User.GetValidateVoucher).Methods("GET")
	router.HandleWithMiddleware("/users", AuthUserMiddleware, handlers.User.PutUpdateProfile).Methods("PUT")
	router.HandleWithMiddleware("/milestones/current", AuthUserMiddleware, handlers.MilestoneHandler.Current).Methods("GET")

	router.HandleWithMiddleware("/admin/users/milestones/check-achievements", AuthClientDashboardMiddleware, handlers.MilestoneHandler.CheckChallengeAchieve).Methods("PUT")
	router.HandleWithMiddleware("/admin/users/milestones/reload", AuthClientDashboardMiddleware, handlers.MilestoneHandler.ReloadMilestone).Methods("PUT")
	router.HandleWithMiddleware("/admin/advertisers/resend-activation", AuthClientDashboardMiddleware, handlers.User.PostSendAdvertiserActivation).Methods("POST")
	router.HandleWithMiddleware("/challenges/{id}/claim", AuthUserMiddleware, handlers.User.GetClaimCredit).Methods("POST")

	// Initiatives
	router.HandleWithMiddleware("/initiatives/{id}/donate", AuthUserMiddleware, handlers.Initiative.Donate).Methods("POST")
	router.HandleWithMiddleware("/initiatives", AuthUserMiddleware, handlers.Initiative.List).Methods("GET")

	// Subscription plans
	router.HandleWithMiddleware("/subscriptions/plans", AuthUserMiddleware, handlers.SubscriptionPlan.List).Methods("GET")

	// SiteSetting
	router.Handle("/static/contents/{type}", handlers.SiteSetting.StaticContent).Methods("GET")
}
