package appeal

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate ifacemaker -f repo.go -o repo_if.go -i Repository -s repo -p appeal
type repo struct {
	psql *pgxpool.Pool
}

func NewRepo(psql *pgxpool.Pool) Repository {
	return &repo{
		psql: psql,
	}
}

func (r *repo) SaveAppeal(ctx context.Context, p *Appeal) (*NewAppeal, error) {
	const insertQuery = `INSERT INTO appeal (
                    division, 
                    subject, 
                    text, 
                    chat_id,
                    username,
                    admin_id
                    ) SELECT $1, $2, $3, $4, $5, admin.id
						FROM admin WHERE admin.division = $1
							LIMIT 1
					    RETURNING appeal.id`

	rows, err := r.psql.Query(ctx, insertQuery, p.Division, p.Subject, p.Text, p.ChatID, p.Username)
	if err != nil {
		return nil, err
	}

	appealID, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[int64])
	if err != nil {
		return nil, err
	}

	const selectQuery = `SELECT ad.chat_id FROM admin ad
							JOIN appeal ap on ad.id = ap.admin_id
							WHERE ap.id = $1`

	rows, err = r.psql.Query(ctx, selectQuery, appealID)
	if err != nil {
		return nil, err
	}

	adminChatID, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[int64])
	if err != nil {
		return nil, err
	}

	return &NewAppeal{
		AppealID:    appealID,
		AdminChatID: adminChatID,
	}, nil
}

func (r *repo) AnswerAppeal(ctx context.Context, id int64) error {
	const query = `UPDATE appeal SET
                    answered_at = CURRENT_TIMESTAMP
                    WHERE id = $1`

	_, err := r.psql.Exec(ctx, query, id)
	return err
}

func (r *repo) GetAppeal(ctx context.Context, id int64) (*Appeal, error) {
	const query = `SELECT
						division,
						subject, 
						text, 
						chat_id,
						username
				   FROM appeal
				   	WHERE id = $1`

	rows, err := r.psql.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Appeal])
}

func (r *repo) AdminExists(ctx context.Context, division string) (bool, error) {
	const query = `SELECT EXISTS(SELECT FROM admin WHERE division = $1)`

	rows, err := r.psql.Query(ctx, query, division)
	if err != nil {
		return false, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowTo[bool])
}
