package service

import (
	"github.com/diarikom/running-app/running-app-api/pkg/nsql"
	"github.com/jmoiron/sqlx"
)

type MileStatement struct {
	findById                          *sqlx.Stmt
	findChallengeById                 *sqlx.Stmt
	current                           *sqlx.Stmt
	findUserChallenge                 *sqlx.Stmt
	insertUserChallenge               *sqlx.NamedStmt
	findCurrentChallenges             *sqlx.Stmt
	findUnaccomplishedChallengeByUser *sqlx.Stmt
	getChallengesByStatus             *sqlx.Stmt
	getChallengesIn                   *sqlx.Stmt
	getUserChallengeByMilestoneId     *sqlx.Stmt
}

func initMilestoneStatement(db *nsql.SqlDatabase) MileStatement {
	return MileStatement{
		findById:                          db.Prepare(`SELECT id, name, period_start, period_end, period_tz, status, created_at, updated_at, version FROM milestone WHERE id = $1`),
		current:                           db.Prepare(`SELECT id, period_start, period_end, period_tz, status FROM milestone WHERE status = $2 AND period_start <= $1 AND period_end >= $1`),
		findChallengeById:                 db.Prepare(`SELECT id, milestone_id, title, description, level, status, rules, sort, created_at, updated_at, version FROM challenge WHERE id = $1`),
		findUserChallenge:                 db.Prepare(`SELECT id, user_id, milestone_id, milestone_snapshot, milestone_version, challenge_id, challenge_snapshot, challenge_version, challenge_result_snapshot, reward_snapshot, reward_type_id, reward_ref_id, reward_value, status, updated_at FROM user_challenge WHERE user_id = $1 AND challenge_id = $2`),
		insertUserChallenge:               db.PrepareNamed(`INSERT INTO user_challenge(id, user_id, milestone_id, milestone_snapshot, milestone_version, challenge_id, challenge_snapshot, challenge_version, challenge_result_snapshot, reward_snapshot, reward_type_id, reward_ref_id, reward_value, status, updated_at) VALUES (:id, :user_id, :milestone_id, :milestone_snapshot, :milestone_version, :challenge_id, :challenge_snapshot, :challenge_version, :challenge_result_snapshot, :reward_snapshot, :reward_type_id, :reward_ref_id, :reward_value, :status, :updated_at)`),
		findCurrentChallenges:             db.Prepare(`SELECT c.id, c.milestone_id, c.title, c.description, c.level, c.status, c.rules, c.sort, c.created_at, c.updated_at, c.version FROM challenge c INNER JOIN milestone m ON c.milestone_id = m.id WHERE c.status = $1 AND m.status = $2 ORDER BY c.level, c.sort`),
		findUnaccomplishedChallengeByUser: db.Prepare(`select c.id, c.milestone_id, c.title, c.description, c.level, c.status, c.rules, c.sort, c.created_at, c.updated_at, c.version FROM challenge as c left join user_challenge uc on c.id = uc.challenge_id and uc.user_id = $1 where uc.id is null and c.milestone_id = $2 and c.status = $3 ORDER BY level, sort;`),
		getChallengesByStatus:             db.Prepare(`SELECT challenge.id, challenge.title, COALESCE((SELECT status FROM user_challenge WHERE challenge_id = challenge.id AND user_id = $1), challenge.status) as status, COALESCE((SELECT updated_at FROM user_challenge WHERE challenge_id = challenge.id AND user_id = $1), challenge.updated_at) as updated_at, challenge.rules FROM challenge WHERE challenge.milestone_id = $2 AND status >= $3 ORDER BY challenge.level, challenge.sort`),
		getChallengesIn:                   db.Prepare(`SELECT challenge.id, challenge.title, challenge.rules, COALESCE(uc.status, challenge.status) as status, COALESCE(uc.updated_at, challenge.updated_at) as updated_at FROM challenge LEFT JOIN user_challenge uc on challenge.id = uc.challenge_id WHERE challenge.id IN ($1) ORDER BY challenge.level, challenge.sort`),
		getUserChallengeByMilestoneId:     db.Prepare(`SELECT id, user_id, milestone_id, milestone_snapshot, milestone_version, challenge_id, challenge_snapshot, challenge_version, challenge_result_snapshot, reward_snapshot, reward_type_id, reward_ref_id, reward_value, status, updated_at FROM user_challenge WHERE user_id = $1 AND milestone_id = $2 AND status = $3`),
	}
}
