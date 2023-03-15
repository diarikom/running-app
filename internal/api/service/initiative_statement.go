package service

import (
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"github.com/jmoiron/sqlx"
)

type InitiativeStatement struct {
	findActive         *sqlx.Stmt
	findById           *sqlx.Stmt
	findDonationByUser *sqlx.Stmt
	insertDonation     *sqlx.NamedStmt
	insertDonationLog  *sqlx.NamedStmt
}

func initInitiativeStatement(db *nsql.SqlDatabase) InitiativeStatement {
	return InitiativeStatement{
		findActive:         db.Prepare(`select id, organization_id, name, description, image_files, external_urls, price, currency_id, donation_conversion, status_id, tags, stat_donation_count, created_at, updated_at, "version", headline from initiative where status_id = 2 order by updated_at desc limit $1 offset $2`),
		findById:           db.Prepare(`select id, organization_id, name, description, image_files, external_urls, price, currency_id, status_id, tags, stat_donation_count, created_at, updated_at, "version", headline from initiative where id = $1`),
		findDonationByUser: db.Prepare(`select id, initiative_id, initiative_snapshot, user_id, user_snapshot, payment_method_id, payment_snapshot, payment_trx_ref, qty, total_price, currency_id, status_id, notes, created_at, updated_at, modified_by, version from donation where user_id = $1 order by updated_at desc limit $2 offset $3`),
		insertDonation:     db.PrepareNamed(`INSERT INTO donation(id, initiative_id, initiative_snapshot, user_id, user_snapshot, payment_method_id, payment_snapshot, payment_trx_ref, qty, total_price, currency_id, status_id, notes, created_at, updated_at, modified_by, version) VALUES (:id, :initiative_id, :initiative_snapshot, :user_id, :user_snapshot, :payment_method_id, :payment_snapshot, :payment_trx_ref, :qty, :total_price, :currency_id, :status_id, :notes, :created_at, :updated_at, :modified_by, :version);`),
		insertDonationLog:  db.PrepareNamed(`INSERT INTO donation_log(log_id, changelog, id, payment_method_id, payment_snapshot, payment_trx_ref, status_id, updated_at, modified_by, version, notes) VALUES (:log_id, :changelog, :id, :payment_method_id, :payment_snapshot, :payment_trx_ref, :status_id, :updated_at, :modified_by, :version, :notes);`),
	}
}
