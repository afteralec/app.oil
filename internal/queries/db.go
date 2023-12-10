// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package queries

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.addCommentToRequestStmt, err = db.PrepareContext(ctx, addCommentToRequest); err != nil {
		return nil, fmt.Errorf("error preparing query AddCommentToRequest: %w", err)
	}
	if q.addCommentToRequestFieldStmt, err = db.PrepareContext(ctx, addCommentToRequestField); err != nil {
		return nil, fmt.Errorf("error preparing query AddCommentToRequestField: %w", err)
	}
	if q.addReplyToCommentStmt, err = db.PrepareContext(ctx, addReplyToComment); err != nil {
		return nil, fmt.Errorf("error preparing query AddReplyToComment: %w", err)
	}
	if q.addReplyToFieldCommentStmt, err = db.PrepareContext(ctx, addReplyToFieldComment); err != nil {
		return nil, fmt.Errorf("error preparing query AddReplyToFieldComment: %w", err)
	}
	if q.countEmailsStmt, err = db.PrepareContext(ctx, countEmails); err != nil {
		return nil, fmt.Errorf("error preparing query CountEmails: %w", err)
	}
	if q.countOpenRequestsStmt, err = db.PrepareContext(ctx, countOpenRequests); err != nil {
		return nil, fmt.Errorf("error preparing query CountOpenRequests: %w", err)
	}
	if q.createCharacterApplicationContentStmt, err = db.PrepareContext(ctx, createCharacterApplicationContent); err != nil {
		return nil, fmt.Errorf("error preparing query CreateCharacterApplicationContent: %w", err)
	}
	if q.createCharacterApplicationContentHistoryStmt, err = db.PrepareContext(ctx, createCharacterApplicationContentHistory); err != nil {
		return nil, fmt.Errorf("error preparing query CreateCharacterApplicationContentHistory: %w", err)
	}
	if q.createEmailStmt, err = db.PrepareContext(ctx, createEmail); err != nil {
		return nil, fmt.Errorf("error preparing query CreateEmail: %w", err)
	}
	if q.createPlayerStmt, err = db.PrepareContext(ctx, createPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query CreatePlayer: %w", err)
	}
	if q.createRequestStmt, err = db.PrepareContext(ctx, createRequest); err != nil {
		return nil, fmt.Errorf("error preparing query CreateRequest: %w", err)
	}
	if q.deleteEmailStmt, err = db.PrepareContext(ctx, deleteEmail); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteEmail: %w", err)
	}
	if q.getCharacterApplicationContentStmt, err = db.PrepareContext(ctx, getCharacterApplicationContent); err != nil {
		return nil, fmt.Errorf("error preparing query GetCharacterApplicationContent: %w", err)
	}
	if q.getCharacterApplicationContentForRequestStmt, err = db.PrepareContext(ctx, getCharacterApplicationContentForRequest); err != nil {
		return nil, fmt.Errorf("error preparing query GetCharacterApplicationContentForRequest: %w", err)
	}
	if q.getEmailStmt, err = db.PrepareContext(ctx, getEmail); err != nil {
		return nil, fmt.Errorf("error preparing query GetEmail: %w", err)
	}
	if q.getPlayerStmt, err = db.PrepareContext(ctx, getPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayer: %w", err)
	}
	if q.getPlayerByUsernameStmt, err = db.PrepareContext(ctx, getPlayerByUsername); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerByUsername: %w", err)
	}
	if q.getPlayerPWHashStmt, err = db.PrepareContext(ctx, getPlayerPWHash); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerPWHash: %w", err)
	}
	if q.getPlayerUsernameStmt, err = db.PrepareContext(ctx, getPlayerUsername); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerUsername: %w", err)
	}
	if q.getPlayerUsernameByIdStmt, err = db.PrepareContext(ctx, getPlayerUsernameById); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerUsernameById: %w", err)
	}
	if q.getRequestStmt, err = db.PrepareContext(ctx, getRequest); err != nil {
		return nil, fmt.Errorf("error preparing query GetRequest: %w", err)
	}
	if q.getRequestCommentStmt, err = db.PrepareContext(ctx, getRequestComment); err != nil {
		return nil, fmt.Errorf("error preparing query GetRequestComment: %w", err)
	}
	if q.getRoleStmt, err = db.PrepareContext(ctx, getRole); err != nil {
		return nil, fmt.Errorf("error preparing query GetRole: %w", err)
	}
	if q.getVerifiedEmailByAddressStmt, err = db.PrepareContext(ctx, getVerifiedEmailByAddress); err != nil {
		return nil, fmt.Errorf("error preparing query GetVerifiedEmailByAddress: %w", err)
	}
	if q.listCharacterApplicationContentForPlayerStmt, err = db.PrepareContext(ctx, listCharacterApplicationContentForPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query ListCharacterApplicationContentForPlayer: %w", err)
	}
	if q.listCharacterApplicationsForPlayerStmt, err = db.PrepareContext(ctx, listCharacterApplicationsForPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query ListCharacterApplicationsForPlayer: %w", err)
	}
	if q.listCommentsForRequestStmt, err = db.PrepareContext(ctx, listCommentsForRequest); err != nil {
		return nil, fmt.Errorf("error preparing query ListCommentsForRequest: %w", err)
	}
	if q.listEmailsStmt, err = db.PrepareContext(ctx, listEmails); err != nil {
		return nil, fmt.Errorf("error preparing query ListEmails: %w", err)
	}
	if q.listRepliesToCommentStmt, err = db.PrepareContext(ctx, listRepliesToComment); err != nil {
		return nil, fmt.Errorf("error preparing query ListRepliesToComment: %w", err)
	}
	if q.listRequestsForPlayerStmt, err = db.PrepareContext(ctx, listRequestsForPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query ListRequestsForPlayer: %w", err)
	}
	if q.listVerifiedEmailsStmt, err = db.PrepareContext(ctx, listVerifiedEmails); err != nil {
		return nil, fmt.Errorf("error preparing query ListVerifiedEmails: %w", err)
	}
	if q.markEmailVerifiedStmt, err = db.PrepareContext(ctx, markEmailVerified); err != nil {
		return nil, fmt.Errorf("error preparing query MarkEmailVerified: %w", err)
	}
	if q.updateCharacterApplicationContentStmt, err = db.PrepareContext(ctx, updateCharacterApplicationContent); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateCharacterApplicationContent: %w", err)
	}
	if q.updateCharacterApplicationContentBackstoryStmt, err = db.PrepareContext(ctx, updateCharacterApplicationContentBackstory); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateCharacterApplicationContentBackstory: %w", err)
	}
	if q.updateCharacterApplicationContentDescriptionStmt, err = db.PrepareContext(ctx, updateCharacterApplicationContentDescription); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateCharacterApplicationContentDescription: %w", err)
	}
	if q.updateCharacterApplicationContentGenderStmt, err = db.PrepareContext(ctx, updateCharacterApplicationContentGender); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateCharacterApplicationContentGender: %w", err)
	}
	if q.updateCharacterApplicationContentNameStmt, err = db.PrepareContext(ctx, updateCharacterApplicationContentName); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateCharacterApplicationContentName: %w", err)
	}
	if q.updateCharacterApplicationContentSdescStmt, err = db.PrepareContext(ctx, updateCharacterApplicationContentSdesc); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateCharacterApplicationContentSdesc: %w", err)
	}
	if q.updateCharacterApplicationContentVersionStmt, err = db.PrepareContext(ctx, updateCharacterApplicationContentVersion); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateCharacterApplicationContentVersion: %w", err)
	}
	if q.updatePlayerPasswordStmt, err = db.PrepareContext(ctx, updatePlayerPassword); err != nil {
		return nil, fmt.Errorf("error preparing query UpdatePlayerPassword: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addCommentToRequestStmt != nil {
		if cerr := q.addCommentToRequestStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addCommentToRequestStmt: %w", cerr)
		}
	}
	if q.addCommentToRequestFieldStmt != nil {
		if cerr := q.addCommentToRequestFieldStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addCommentToRequestFieldStmt: %w", cerr)
		}
	}
	if q.addReplyToCommentStmt != nil {
		if cerr := q.addReplyToCommentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addReplyToCommentStmt: %w", cerr)
		}
	}
	if q.addReplyToFieldCommentStmt != nil {
		if cerr := q.addReplyToFieldCommentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addReplyToFieldCommentStmt: %w", cerr)
		}
	}
	if q.countEmailsStmt != nil {
		if cerr := q.countEmailsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countEmailsStmt: %w", cerr)
		}
	}
	if q.countOpenRequestsStmt != nil {
		if cerr := q.countOpenRequestsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countOpenRequestsStmt: %w", cerr)
		}
	}
	if q.createCharacterApplicationContentStmt != nil {
		if cerr := q.createCharacterApplicationContentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createCharacterApplicationContentStmt: %w", cerr)
		}
	}
	if q.createCharacterApplicationContentHistoryStmt != nil {
		if cerr := q.createCharacterApplicationContentHistoryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createCharacterApplicationContentHistoryStmt: %w", cerr)
		}
	}
	if q.createEmailStmt != nil {
		if cerr := q.createEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createEmailStmt: %w", cerr)
		}
	}
	if q.createPlayerStmt != nil {
		if cerr := q.createPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createPlayerStmt: %w", cerr)
		}
	}
	if q.createRequestStmt != nil {
		if cerr := q.createRequestStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createRequestStmt: %w", cerr)
		}
	}
	if q.deleteEmailStmt != nil {
		if cerr := q.deleteEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteEmailStmt: %w", cerr)
		}
	}
	if q.getCharacterApplicationContentStmt != nil {
		if cerr := q.getCharacterApplicationContentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCharacterApplicationContentStmt: %w", cerr)
		}
	}
	if q.getCharacterApplicationContentForRequestStmt != nil {
		if cerr := q.getCharacterApplicationContentForRequestStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCharacterApplicationContentForRequestStmt: %w", cerr)
		}
	}
	if q.getEmailStmt != nil {
		if cerr := q.getEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEmailStmt: %w", cerr)
		}
	}
	if q.getPlayerStmt != nil {
		if cerr := q.getPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerStmt: %w", cerr)
		}
	}
	if q.getPlayerByUsernameStmt != nil {
		if cerr := q.getPlayerByUsernameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerByUsernameStmt: %w", cerr)
		}
	}
	if q.getPlayerPWHashStmt != nil {
		if cerr := q.getPlayerPWHashStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerPWHashStmt: %w", cerr)
		}
	}
	if q.getPlayerUsernameStmt != nil {
		if cerr := q.getPlayerUsernameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerUsernameStmt: %w", cerr)
		}
	}
	if q.getPlayerUsernameByIdStmt != nil {
		if cerr := q.getPlayerUsernameByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerUsernameByIdStmt: %w", cerr)
		}
	}
	if q.getRequestStmt != nil {
		if cerr := q.getRequestStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRequestStmt: %w", cerr)
		}
	}
	if q.getRequestCommentStmt != nil {
		if cerr := q.getRequestCommentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRequestCommentStmt: %w", cerr)
		}
	}
	if q.getRoleStmt != nil {
		if cerr := q.getRoleStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRoleStmt: %w", cerr)
		}
	}
	if q.getVerifiedEmailByAddressStmt != nil {
		if cerr := q.getVerifiedEmailByAddressStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getVerifiedEmailByAddressStmt: %w", cerr)
		}
	}
	if q.listCharacterApplicationContentForPlayerStmt != nil {
		if cerr := q.listCharacterApplicationContentForPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listCharacterApplicationContentForPlayerStmt: %w", cerr)
		}
	}
	if q.listCharacterApplicationsForPlayerStmt != nil {
		if cerr := q.listCharacterApplicationsForPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listCharacterApplicationsForPlayerStmt: %w", cerr)
		}
	}
	if q.listCommentsForRequestStmt != nil {
		if cerr := q.listCommentsForRequestStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listCommentsForRequestStmt: %w", cerr)
		}
	}
	if q.listEmailsStmt != nil {
		if cerr := q.listEmailsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listEmailsStmt: %w", cerr)
		}
	}
	if q.listRepliesToCommentStmt != nil {
		if cerr := q.listRepliesToCommentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listRepliesToCommentStmt: %w", cerr)
		}
	}
	if q.listRequestsForPlayerStmt != nil {
		if cerr := q.listRequestsForPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listRequestsForPlayerStmt: %w", cerr)
		}
	}
	if q.listVerifiedEmailsStmt != nil {
		if cerr := q.listVerifiedEmailsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listVerifiedEmailsStmt: %w", cerr)
		}
	}
	if q.markEmailVerifiedStmt != nil {
		if cerr := q.markEmailVerifiedStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing markEmailVerifiedStmt: %w", cerr)
		}
	}
	if q.updateCharacterApplicationContentStmt != nil {
		if cerr := q.updateCharacterApplicationContentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateCharacterApplicationContentStmt: %w", cerr)
		}
	}
	if q.updateCharacterApplicationContentBackstoryStmt != nil {
		if cerr := q.updateCharacterApplicationContentBackstoryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateCharacterApplicationContentBackstoryStmt: %w", cerr)
		}
	}
	if q.updateCharacterApplicationContentDescriptionStmt != nil {
		if cerr := q.updateCharacterApplicationContentDescriptionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateCharacterApplicationContentDescriptionStmt: %w", cerr)
		}
	}
	if q.updateCharacterApplicationContentGenderStmt != nil {
		if cerr := q.updateCharacterApplicationContentGenderStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateCharacterApplicationContentGenderStmt: %w", cerr)
		}
	}
	if q.updateCharacterApplicationContentNameStmt != nil {
		if cerr := q.updateCharacterApplicationContentNameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateCharacterApplicationContentNameStmt: %w", cerr)
		}
	}
	if q.updateCharacterApplicationContentSdescStmt != nil {
		if cerr := q.updateCharacterApplicationContentSdescStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateCharacterApplicationContentSdescStmt: %w", cerr)
		}
	}
	if q.updateCharacterApplicationContentVersionStmt != nil {
		if cerr := q.updateCharacterApplicationContentVersionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateCharacterApplicationContentVersionStmt: %w", cerr)
		}
	}
	if q.updatePlayerPasswordStmt != nil {
		if cerr := q.updatePlayerPasswordStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updatePlayerPasswordStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                                               DBTX
	tx                                               *sql.Tx
	addCommentToRequestStmt                          *sql.Stmt
	addCommentToRequestFieldStmt                     *sql.Stmt
	addReplyToCommentStmt                            *sql.Stmt
	addReplyToFieldCommentStmt                       *sql.Stmt
	countEmailsStmt                                  *sql.Stmt
	countOpenRequestsStmt                            *sql.Stmt
	createCharacterApplicationContentStmt            *sql.Stmt
	createCharacterApplicationContentHistoryStmt     *sql.Stmt
	createEmailStmt                                  *sql.Stmt
	createPlayerStmt                                 *sql.Stmt
	createRequestStmt                                *sql.Stmt
	deleteEmailStmt                                  *sql.Stmt
	getCharacterApplicationContentStmt               *sql.Stmt
	getCharacterApplicationContentForRequestStmt     *sql.Stmt
	getEmailStmt                                     *sql.Stmt
	getPlayerStmt                                    *sql.Stmt
	getPlayerByUsernameStmt                          *sql.Stmt
	getPlayerPWHashStmt                              *sql.Stmt
	getPlayerUsernameStmt                            *sql.Stmt
	getPlayerUsernameByIdStmt                        *sql.Stmt
	getRequestStmt                                   *sql.Stmt
	getRequestCommentStmt                            *sql.Stmt
	getRoleStmt                                      *sql.Stmt
	getVerifiedEmailByAddressStmt                    *sql.Stmt
	listCharacterApplicationContentForPlayerStmt     *sql.Stmt
	listCharacterApplicationsForPlayerStmt           *sql.Stmt
	listCommentsForRequestStmt                       *sql.Stmt
	listEmailsStmt                                   *sql.Stmt
	listRepliesToCommentStmt                         *sql.Stmt
	listRequestsForPlayerStmt                        *sql.Stmt
	listVerifiedEmailsStmt                           *sql.Stmt
	markEmailVerifiedStmt                            *sql.Stmt
	updateCharacterApplicationContentStmt            *sql.Stmt
	updateCharacterApplicationContentBackstoryStmt   *sql.Stmt
	updateCharacterApplicationContentDescriptionStmt *sql.Stmt
	updateCharacterApplicationContentGenderStmt      *sql.Stmt
	updateCharacterApplicationContentNameStmt        *sql.Stmt
	updateCharacterApplicationContentSdescStmt       *sql.Stmt
	updateCharacterApplicationContentVersionStmt     *sql.Stmt
	updatePlayerPasswordStmt                         *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                                    tx,
		tx:                                    tx,
		addCommentToRequestStmt:               q.addCommentToRequestStmt,
		addCommentToRequestFieldStmt:          q.addCommentToRequestFieldStmt,
		addReplyToCommentStmt:                 q.addReplyToCommentStmt,
		addReplyToFieldCommentStmt:            q.addReplyToFieldCommentStmt,
		countEmailsStmt:                       q.countEmailsStmt,
		countOpenRequestsStmt:                 q.countOpenRequestsStmt,
		createCharacterApplicationContentStmt: q.createCharacterApplicationContentStmt,
		createCharacterApplicationContentHistoryStmt: q.createCharacterApplicationContentHistoryStmt,
		createEmailStmt:                                  q.createEmailStmt,
		createPlayerStmt:                                 q.createPlayerStmt,
		createRequestStmt:                                q.createRequestStmt,
		deleteEmailStmt:                                  q.deleteEmailStmt,
		getCharacterApplicationContentStmt:               q.getCharacterApplicationContentStmt,
		getCharacterApplicationContentForRequestStmt:     q.getCharacterApplicationContentForRequestStmt,
		getEmailStmt:                                     q.getEmailStmt,
		getPlayerStmt:                                    q.getPlayerStmt,
		getPlayerByUsernameStmt:                          q.getPlayerByUsernameStmt,
		getPlayerPWHashStmt:                              q.getPlayerPWHashStmt,
		getPlayerUsernameStmt:                            q.getPlayerUsernameStmt,
		getPlayerUsernameByIdStmt:                        q.getPlayerUsernameByIdStmt,
		getRequestStmt:                                   q.getRequestStmt,
		getRequestCommentStmt:                            q.getRequestCommentStmt,
		getRoleStmt:                                      q.getRoleStmt,
		getVerifiedEmailByAddressStmt:                    q.getVerifiedEmailByAddressStmt,
		listCharacterApplicationContentForPlayerStmt:     q.listCharacterApplicationContentForPlayerStmt,
		listCharacterApplicationsForPlayerStmt:           q.listCharacterApplicationsForPlayerStmt,
		listCommentsForRequestStmt:                       q.listCommentsForRequestStmt,
		listEmailsStmt:                                   q.listEmailsStmt,
		listRepliesToCommentStmt:                         q.listRepliesToCommentStmt,
		listRequestsForPlayerStmt:                        q.listRequestsForPlayerStmt,
		listVerifiedEmailsStmt:                           q.listVerifiedEmailsStmt,
		markEmailVerifiedStmt:                            q.markEmailVerifiedStmt,
		updateCharacterApplicationContentStmt:            q.updateCharacterApplicationContentStmt,
		updateCharacterApplicationContentBackstoryStmt:   q.updateCharacterApplicationContentBackstoryStmt,
		updateCharacterApplicationContentDescriptionStmt: q.updateCharacterApplicationContentDescriptionStmt,
		updateCharacterApplicationContentGenderStmt:      q.updateCharacterApplicationContentGenderStmt,
		updateCharacterApplicationContentNameStmt:        q.updateCharacterApplicationContentNameStmt,
		updateCharacterApplicationContentSdescStmt:       q.updateCharacterApplicationContentSdescStmt,
		updateCharacterApplicationContentVersionStmt:     q.updateCharacterApplicationContentVersionStmt,
		updatePlayerPasswordStmt:                         q.updatePlayerPasswordStmt,
	}
}
