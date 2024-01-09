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
	if q.countEmailsStmt, err = db.PrepareContext(ctx, countEmails); err != nil {
		return nil, fmt.Errorf("error preparing query CountEmails: %w", err)
	}
	if q.countOpenCharacterApplicationsForPlayerStmt, err = db.PrepareContext(ctx, countOpenCharacterApplicationsForPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query CountOpenCharacterApplicationsForPlayer: %w", err)
	}
	if q.countOpenRequestsStmt, err = db.PrepareContext(ctx, countOpenRequests); err != nil {
		return nil, fmt.Errorf("error preparing query CountOpenRequests: %w", err)
	}
	if q.countUnresolvedCommentsStmt, err = db.PrepareContext(ctx, countUnresolvedComments); err != nil {
		return nil, fmt.Errorf("error preparing query CountUnresolvedComments: %w", err)
	}
	if q.createCharacterApplicationContentStmt, err = db.PrepareContext(ctx, createCharacterApplicationContent); err != nil {
		return nil, fmt.Errorf("error preparing query CreateCharacterApplicationContent: %w", err)
	}
	if q.createEmailStmt, err = db.PrepareContext(ctx, createEmail); err != nil {
		return nil, fmt.Errorf("error preparing query CreateEmail: %w", err)
	}
	if q.createHistoryForCharacterApplicationStmt, err = db.PrepareContext(ctx, createHistoryForCharacterApplication); err != nil {
		return nil, fmt.Errorf("error preparing query CreateHistoryForCharacterApplication: %w", err)
	}
	if q.createHistoryForRequestStatusChangeStmt, err = db.PrepareContext(ctx, createHistoryForRequestStatusChange); err != nil {
		return nil, fmt.Errorf("error preparing query CreateHistoryForRequestStatusChange: %w", err)
	}
	if q.createPlayerStmt, err = db.PrepareContext(ctx, createPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query CreatePlayer: %w", err)
	}
	if q.createPlayerPermissionStmt, err = db.PrepareContext(ctx, createPlayerPermission); err != nil {
		return nil, fmt.Errorf("error preparing query CreatePlayerPermission: %w", err)
	}
	if q.createPlayerPermissionIssuedChangeHistoryStmt, err = db.PrepareContext(ctx, createPlayerPermissionIssuedChangeHistory); err != nil {
		return nil, fmt.Errorf("error preparing query CreatePlayerPermissionIssuedChangeHistory: %w", err)
	}
	if q.createPlayerPermissionRevokedChangeHistoryStmt, err = db.PrepareContext(ctx, createPlayerPermissionRevokedChangeHistory); err != nil {
		return nil, fmt.Errorf("error preparing query CreatePlayerPermissionRevokedChangeHistory: %w", err)
	}
	if q.createPlayerSettingsStmt, err = db.PrepareContext(ctx, createPlayerSettings); err != nil {
		return nil, fmt.Errorf("error preparing query CreatePlayerSettings: %w", err)
	}
	if q.createRequestStmt, err = db.PrepareContext(ctx, createRequest); err != nil {
		return nil, fmt.Errorf("error preparing query CreateRequest: %w", err)
	}
	if q.createRequestCommentStmt, err = db.PrepareContext(ctx, createRequestComment); err != nil {
		return nil, fmt.Errorf("error preparing query CreateRequestComment: %w", err)
	}
	if q.deleteEmailStmt, err = db.PrepareContext(ctx, deleteEmail); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteEmail: %w", err)
	}
	if q.deletePlayerPermissionStmt, err = db.PrepareContext(ctx, deletePlayerPermission); err != nil {
		return nil, fmt.Errorf("error preparing query DeletePlayerPermission: %w", err)
	}
	if q.getCharacterApplicationStmt, err = db.PrepareContext(ctx, getCharacterApplication); err != nil {
		return nil, fmt.Errorf("error preparing query GetCharacterApplication: %w", err)
	}
	if q.getCharacterApplicationContentStmt, err = db.PrepareContext(ctx, getCharacterApplicationContent); err != nil {
		return nil, fmt.Errorf("error preparing query GetCharacterApplicationContent: %w", err)
	}
	if q.getCharacterApplicationContentForRequestStmt, err = db.PrepareContext(ctx, getCharacterApplicationContentForRequest); err != nil {
		return nil, fmt.Errorf("error preparing query GetCharacterApplicationContentForRequest: %w", err)
	}
	if q.getCommentWithAuthorStmt, err = db.PrepareContext(ctx, getCommentWithAuthor); err != nil {
		return nil, fmt.Errorf("error preparing query GetCommentWithAuthor: %w", err)
	}
	if q.getEmailStmt, err = db.PrepareContext(ctx, getEmail); err != nil {
		return nil, fmt.Errorf("error preparing query GetEmail: %w", err)
	}
	if q.getEmailByAddressForPlayerStmt, err = db.PrepareContext(ctx, getEmailByAddressForPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query GetEmailByAddressForPlayer: %w", err)
	}
	if q.getHelpStmt, err = db.PrepareContext(ctx, getHelp); err != nil {
		return nil, fmt.Errorf("error preparing query GetHelp: %w", err)
	}
	if q.getHelpRelatedStmt, err = db.PrepareContext(ctx, getHelpRelated); err != nil {
		return nil, fmt.Errorf("error preparing query GetHelpRelated: %w", err)
	}
	if q.getPlayerStmt, err = db.PrepareContext(ctx, getPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayer: %w", err)
	}
	if q.getPlayerByUsernameStmt, err = db.PrepareContext(ctx, getPlayerByUsername); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerByUsername: %w", err)
	}
	if q.getPlayerSettingsStmt, err = db.PrepareContext(ctx, getPlayerSettings); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerSettings: %w", err)
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
	if q.getVerifiedEmailByAddressStmt, err = db.PrepareContext(ctx, getVerifiedEmailByAddress); err != nil {
		return nil, fmt.Errorf("error preparing query GetVerifiedEmailByAddress: %w", err)
	}
	if q.incrementRequestVersionStmt, err = db.PrepareContext(ctx, incrementRequestVersion); err != nil {
		return nil, fmt.Errorf("error preparing query IncrementRequestVersion: %w", err)
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
	if q.listCommentsForRequestWithAuthorStmt, err = db.PrepareContext(ctx, listCommentsForRequestWithAuthor); err != nil {
		return nil, fmt.Errorf("error preparing query ListCommentsForRequestWithAuthor: %w", err)
	}
	if q.listEmailsStmt, err = db.PrepareContext(ctx, listEmails); err != nil {
		return nil, fmt.Errorf("error preparing query ListEmails: %w", err)
	}
	if q.listHelpHeadersStmt, err = db.PrepareContext(ctx, listHelpHeaders); err != nil {
		return nil, fmt.Errorf("error preparing query ListHelpHeaders: %w", err)
	}
	if q.listHelpSlugsStmt, err = db.PrepareContext(ctx, listHelpSlugs); err != nil {
		return nil, fmt.Errorf("error preparing query ListHelpSlugs: %w", err)
	}
	if q.listOpenCharacterApplicationsStmt, err = db.PrepareContext(ctx, listOpenCharacterApplications); err != nil {
		return nil, fmt.Errorf("error preparing query ListOpenCharacterApplications: %w", err)
	}
	if q.listPlayerPermissionsStmt, err = db.PrepareContext(ctx, listPlayerPermissions); err != nil {
		return nil, fmt.Errorf("error preparing query ListPlayerPermissions: %w", err)
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
	if q.markRequestCanceledStmt, err = db.PrepareContext(ctx, markRequestCanceled); err != nil {
		return nil, fmt.Errorf("error preparing query MarkRequestCanceled: %w", err)
	}
	if q.markRequestInReviewStmt, err = db.PrepareContext(ctx, markRequestInReview); err != nil {
		return nil, fmt.Errorf("error preparing query MarkRequestInReview: %w", err)
	}
	if q.markRequestReadyStmt, err = db.PrepareContext(ctx, markRequestReady); err != nil {
		return nil, fmt.Errorf("error preparing query MarkRequestReady: %w", err)
	}
	if q.markRequestSubmittedStmt, err = db.PrepareContext(ctx, markRequestSubmitted); err != nil {
		return nil, fmt.Errorf("error preparing query MarkRequestSubmitted: %w", err)
	}
	if q.searchHelpByContentStmt, err = db.PrepareContext(ctx, searchHelpByContent); err != nil {
		return nil, fmt.Errorf("error preparing query SearchHelpByContent: %w", err)
	}
	if q.searchHelpByTitleStmt, err = db.PrepareContext(ctx, searchHelpByTitle); err != nil {
		return nil, fmt.Errorf("error preparing query SearchHelpByTitle: %w", err)
	}
	if q.searchPlayersByUsernameStmt, err = db.PrepareContext(ctx, searchPlayersByUsername); err != nil {
		return nil, fmt.Errorf("error preparing query SearchPlayersByUsername: %w", err)
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
	if q.updateCharacterApplicationContentShortDescriptionStmt, err = db.PrepareContext(ctx, updateCharacterApplicationContentShortDescription); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateCharacterApplicationContentShortDescription: %w", err)
	}
	if q.updatePlayerPasswordStmt, err = db.PrepareContext(ctx, updatePlayerPassword); err != nil {
		return nil, fmt.Errorf("error preparing query UpdatePlayerPassword: %w", err)
	}
	if q.updatePlayerSettingsThemeStmt, err = db.PrepareContext(ctx, updatePlayerSettingsTheme); err != nil {
		return nil, fmt.Errorf("error preparing query UpdatePlayerSettingsTheme: %w", err)
	}
	if q.updateRequestStatusStmt, err = db.PrepareContext(ctx, updateRequestStatus); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateRequestStatus: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.countEmailsStmt != nil {
		if cerr := q.countEmailsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countEmailsStmt: %w", cerr)
		}
	}
	if q.countOpenCharacterApplicationsForPlayerStmt != nil {
		if cerr := q.countOpenCharacterApplicationsForPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countOpenCharacterApplicationsForPlayerStmt: %w", cerr)
		}
	}
	if q.countOpenRequestsStmt != nil {
		if cerr := q.countOpenRequestsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countOpenRequestsStmt: %w", cerr)
		}
	}
	if q.countUnresolvedCommentsStmt != nil {
		if cerr := q.countUnresolvedCommentsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countUnresolvedCommentsStmt: %w", cerr)
		}
	}
	if q.createCharacterApplicationContentStmt != nil {
		if cerr := q.createCharacterApplicationContentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createCharacterApplicationContentStmt: %w", cerr)
		}
	}
	if q.createEmailStmt != nil {
		if cerr := q.createEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createEmailStmt: %w", cerr)
		}
	}
	if q.createHistoryForCharacterApplicationStmt != nil {
		if cerr := q.createHistoryForCharacterApplicationStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createHistoryForCharacterApplicationStmt: %w", cerr)
		}
	}
	if q.createHistoryForRequestStatusChangeStmt != nil {
		if cerr := q.createHistoryForRequestStatusChangeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createHistoryForRequestStatusChangeStmt: %w", cerr)
		}
	}
	if q.createPlayerStmt != nil {
		if cerr := q.createPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createPlayerStmt: %w", cerr)
		}
	}
	if q.createPlayerPermissionStmt != nil {
		if cerr := q.createPlayerPermissionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createPlayerPermissionStmt: %w", cerr)
		}
	}
	if q.createPlayerPermissionIssuedChangeHistoryStmt != nil {
		if cerr := q.createPlayerPermissionIssuedChangeHistoryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createPlayerPermissionIssuedChangeHistoryStmt: %w", cerr)
		}
	}
	if q.createPlayerPermissionRevokedChangeHistoryStmt != nil {
		if cerr := q.createPlayerPermissionRevokedChangeHistoryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createPlayerPermissionRevokedChangeHistoryStmt: %w", cerr)
		}
	}
	if q.createPlayerSettingsStmt != nil {
		if cerr := q.createPlayerSettingsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createPlayerSettingsStmt: %w", cerr)
		}
	}
	if q.createRequestStmt != nil {
		if cerr := q.createRequestStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createRequestStmt: %w", cerr)
		}
	}
	if q.createRequestCommentStmt != nil {
		if cerr := q.createRequestCommentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createRequestCommentStmt: %w", cerr)
		}
	}
	if q.deleteEmailStmt != nil {
		if cerr := q.deleteEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteEmailStmt: %w", cerr)
		}
	}
	if q.deletePlayerPermissionStmt != nil {
		if cerr := q.deletePlayerPermissionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deletePlayerPermissionStmt: %w", cerr)
		}
	}
	if q.getCharacterApplicationStmt != nil {
		if cerr := q.getCharacterApplicationStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCharacterApplicationStmt: %w", cerr)
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
	if q.getCommentWithAuthorStmt != nil {
		if cerr := q.getCommentWithAuthorStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCommentWithAuthorStmt: %w", cerr)
		}
	}
	if q.getEmailStmt != nil {
		if cerr := q.getEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEmailStmt: %w", cerr)
		}
	}
	if q.getEmailByAddressForPlayerStmt != nil {
		if cerr := q.getEmailByAddressForPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEmailByAddressForPlayerStmt: %w", cerr)
		}
	}
	if q.getHelpStmt != nil {
		if cerr := q.getHelpStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getHelpStmt: %w", cerr)
		}
	}
	if q.getHelpRelatedStmt != nil {
		if cerr := q.getHelpRelatedStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getHelpRelatedStmt: %w", cerr)
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
	if q.getPlayerSettingsStmt != nil {
		if cerr := q.getPlayerSettingsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerSettingsStmt: %w", cerr)
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
	if q.getVerifiedEmailByAddressStmt != nil {
		if cerr := q.getVerifiedEmailByAddressStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getVerifiedEmailByAddressStmt: %w", cerr)
		}
	}
	if q.incrementRequestVersionStmt != nil {
		if cerr := q.incrementRequestVersionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing incrementRequestVersionStmt: %w", cerr)
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
	if q.listCommentsForRequestWithAuthorStmt != nil {
		if cerr := q.listCommentsForRequestWithAuthorStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listCommentsForRequestWithAuthorStmt: %w", cerr)
		}
	}
	if q.listEmailsStmt != nil {
		if cerr := q.listEmailsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listEmailsStmt: %w", cerr)
		}
	}
	if q.listHelpHeadersStmt != nil {
		if cerr := q.listHelpHeadersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listHelpHeadersStmt: %w", cerr)
		}
	}
	if q.listHelpSlugsStmt != nil {
		if cerr := q.listHelpSlugsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listHelpSlugsStmt: %w", cerr)
		}
	}
	if q.listOpenCharacterApplicationsStmt != nil {
		if cerr := q.listOpenCharacterApplicationsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listOpenCharacterApplicationsStmt: %w", cerr)
		}
	}
	if q.listPlayerPermissionsStmt != nil {
		if cerr := q.listPlayerPermissionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listPlayerPermissionsStmt: %w", cerr)
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
	if q.markRequestCanceledStmt != nil {
		if cerr := q.markRequestCanceledStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing markRequestCanceledStmt: %w", cerr)
		}
	}
	if q.markRequestInReviewStmt != nil {
		if cerr := q.markRequestInReviewStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing markRequestInReviewStmt: %w", cerr)
		}
	}
	if q.markRequestReadyStmt != nil {
		if cerr := q.markRequestReadyStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing markRequestReadyStmt: %w", cerr)
		}
	}
	if q.markRequestSubmittedStmt != nil {
		if cerr := q.markRequestSubmittedStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing markRequestSubmittedStmt: %w", cerr)
		}
	}
	if q.searchHelpByContentStmt != nil {
		if cerr := q.searchHelpByContentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing searchHelpByContentStmt: %w", cerr)
		}
	}
	if q.searchHelpByTitleStmt != nil {
		if cerr := q.searchHelpByTitleStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing searchHelpByTitleStmt: %w", cerr)
		}
	}
	if q.searchPlayersByUsernameStmt != nil {
		if cerr := q.searchPlayersByUsernameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing searchPlayersByUsernameStmt: %w", cerr)
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
	if q.updateCharacterApplicationContentShortDescriptionStmt != nil {
		if cerr := q.updateCharacterApplicationContentShortDescriptionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateCharacterApplicationContentShortDescriptionStmt: %w", cerr)
		}
	}
	if q.updatePlayerPasswordStmt != nil {
		if cerr := q.updatePlayerPasswordStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updatePlayerPasswordStmt: %w", cerr)
		}
	}
	if q.updatePlayerSettingsThemeStmt != nil {
		if cerr := q.updatePlayerSettingsThemeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updatePlayerSettingsThemeStmt: %w", cerr)
		}
	}
	if q.updateRequestStatusStmt != nil {
		if cerr := q.updateRequestStatusStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateRequestStatusStmt: %w", cerr)
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
	db                                                    DBTX
	tx                                                    *sql.Tx
	countEmailsStmt                                       *sql.Stmt
	countOpenCharacterApplicationsForPlayerStmt           *sql.Stmt
	countOpenRequestsStmt                                 *sql.Stmt
	countUnresolvedCommentsStmt                           *sql.Stmt
	createCharacterApplicationContentStmt                 *sql.Stmt
	createEmailStmt                                       *sql.Stmt
	createHistoryForCharacterApplicationStmt              *sql.Stmt
	createHistoryForRequestStatusChangeStmt               *sql.Stmt
	createPlayerStmt                                      *sql.Stmt
	createPlayerPermissionStmt                            *sql.Stmt
	createPlayerPermissionIssuedChangeHistoryStmt         *sql.Stmt
	createPlayerPermissionRevokedChangeHistoryStmt        *sql.Stmt
	createPlayerSettingsStmt                              *sql.Stmt
	createRequestStmt                                     *sql.Stmt
	createRequestCommentStmt                              *sql.Stmt
	deleteEmailStmt                                       *sql.Stmt
	deletePlayerPermissionStmt                            *sql.Stmt
	getCharacterApplicationStmt                           *sql.Stmt
	getCharacterApplicationContentStmt                    *sql.Stmt
	getCharacterApplicationContentForRequestStmt          *sql.Stmt
	getCommentWithAuthorStmt                              *sql.Stmt
	getEmailStmt                                          *sql.Stmt
	getEmailByAddressForPlayerStmt                        *sql.Stmt
	getHelpStmt                                           *sql.Stmt
	getHelpRelatedStmt                                    *sql.Stmt
	getPlayerStmt                                         *sql.Stmt
	getPlayerByUsernameStmt                               *sql.Stmt
	getPlayerSettingsStmt                                 *sql.Stmt
	getPlayerUsernameStmt                                 *sql.Stmt
	getPlayerUsernameByIdStmt                             *sql.Stmt
	getRequestStmt                                        *sql.Stmt
	getVerifiedEmailByAddressStmt                         *sql.Stmt
	incrementRequestVersionStmt                           *sql.Stmt
	listCharacterApplicationContentForPlayerStmt          *sql.Stmt
	listCharacterApplicationsForPlayerStmt                *sql.Stmt
	listCommentsForRequestStmt                            *sql.Stmt
	listCommentsForRequestWithAuthorStmt                  *sql.Stmt
	listEmailsStmt                                        *sql.Stmt
	listHelpHeadersStmt                                   *sql.Stmt
	listHelpSlugsStmt                                     *sql.Stmt
	listOpenCharacterApplicationsStmt                     *sql.Stmt
	listPlayerPermissionsStmt                             *sql.Stmt
	listRequestsForPlayerStmt                             *sql.Stmt
	listVerifiedEmailsStmt                                *sql.Stmt
	markEmailVerifiedStmt                                 *sql.Stmt
	markRequestCanceledStmt                               *sql.Stmt
	markRequestInReviewStmt                               *sql.Stmt
	markRequestReadyStmt                                  *sql.Stmt
	markRequestSubmittedStmt                              *sql.Stmt
	searchHelpByContentStmt                               *sql.Stmt
	searchHelpByTitleStmt                                 *sql.Stmt
	searchPlayersByUsernameStmt                           *sql.Stmt
	updateCharacterApplicationContentBackstoryStmt        *sql.Stmt
	updateCharacterApplicationContentDescriptionStmt      *sql.Stmt
	updateCharacterApplicationContentGenderStmt           *sql.Stmt
	updateCharacterApplicationContentNameStmt             *sql.Stmt
	updateCharacterApplicationContentShortDescriptionStmt *sql.Stmt
	updatePlayerPasswordStmt                              *sql.Stmt
	updatePlayerSettingsThemeStmt                         *sql.Stmt
	updateRequestStatusStmt                               *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:              tx,
		tx:              tx,
		countEmailsStmt: q.countEmailsStmt,
		countOpenCharacterApplicationsForPlayerStmt:           q.countOpenCharacterApplicationsForPlayerStmt,
		countOpenRequestsStmt:                                 q.countOpenRequestsStmt,
		countUnresolvedCommentsStmt:                           q.countUnresolvedCommentsStmt,
		createCharacterApplicationContentStmt:                 q.createCharacterApplicationContentStmt,
		createEmailStmt:                                       q.createEmailStmt,
		createHistoryForCharacterApplicationStmt:              q.createHistoryForCharacterApplicationStmt,
		createHistoryForRequestStatusChangeStmt:               q.createHistoryForRequestStatusChangeStmt,
		createPlayerStmt:                                      q.createPlayerStmt,
		createPlayerPermissionStmt:                            q.createPlayerPermissionStmt,
		createPlayerPermissionIssuedChangeHistoryStmt:         q.createPlayerPermissionIssuedChangeHistoryStmt,
		createPlayerPermissionRevokedChangeHistoryStmt:        q.createPlayerPermissionRevokedChangeHistoryStmt,
		createPlayerSettingsStmt:                              q.createPlayerSettingsStmt,
		createRequestStmt:                                     q.createRequestStmt,
		createRequestCommentStmt:                              q.createRequestCommentStmt,
		deleteEmailStmt:                                       q.deleteEmailStmt,
		deletePlayerPermissionStmt:                            q.deletePlayerPermissionStmt,
		getCharacterApplicationStmt:                           q.getCharacterApplicationStmt,
		getCharacterApplicationContentStmt:                    q.getCharacterApplicationContentStmt,
		getCharacterApplicationContentForRequestStmt:          q.getCharacterApplicationContentForRequestStmt,
		getCommentWithAuthorStmt:                              q.getCommentWithAuthorStmt,
		getEmailStmt:                                          q.getEmailStmt,
		getEmailByAddressForPlayerStmt:                        q.getEmailByAddressForPlayerStmt,
		getHelpStmt:                                           q.getHelpStmt,
		getHelpRelatedStmt:                                    q.getHelpRelatedStmt,
		getPlayerStmt:                                         q.getPlayerStmt,
		getPlayerByUsernameStmt:                               q.getPlayerByUsernameStmt,
		getPlayerSettingsStmt:                                 q.getPlayerSettingsStmt,
		getPlayerUsernameStmt:                                 q.getPlayerUsernameStmt,
		getPlayerUsernameByIdStmt:                             q.getPlayerUsernameByIdStmt,
		getRequestStmt:                                        q.getRequestStmt,
		getVerifiedEmailByAddressStmt:                         q.getVerifiedEmailByAddressStmt,
		incrementRequestVersionStmt:                           q.incrementRequestVersionStmt,
		listCharacterApplicationContentForPlayerStmt:          q.listCharacterApplicationContentForPlayerStmt,
		listCharacterApplicationsForPlayerStmt:                q.listCharacterApplicationsForPlayerStmt,
		listCommentsForRequestStmt:                            q.listCommentsForRequestStmt,
		listCommentsForRequestWithAuthorStmt:                  q.listCommentsForRequestWithAuthorStmt,
		listEmailsStmt:                                        q.listEmailsStmt,
		listHelpHeadersStmt:                                   q.listHelpHeadersStmt,
		listHelpSlugsStmt:                                     q.listHelpSlugsStmt,
		listOpenCharacterApplicationsStmt:                     q.listOpenCharacterApplicationsStmt,
		listPlayerPermissionsStmt:                             q.listPlayerPermissionsStmt,
		listRequestsForPlayerStmt:                             q.listRequestsForPlayerStmt,
		listVerifiedEmailsStmt:                                q.listVerifiedEmailsStmt,
		markEmailVerifiedStmt:                                 q.markEmailVerifiedStmt,
		markRequestCanceledStmt:                               q.markRequestCanceledStmt,
		markRequestInReviewStmt:                               q.markRequestInReviewStmt,
		markRequestReadyStmt:                                  q.markRequestReadyStmt,
		markRequestSubmittedStmt:                              q.markRequestSubmittedStmt,
		searchHelpByContentStmt:                               q.searchHelpByContentStmt,
		searchHelpByTitleStmt:                                 q.searchHelpByTitleStmt,
		searchPlayersByUsernameStmt:                           q.searchPlayersByUsernameStmt,
		updateCharacterApplicationContentBackstoryStmt:        q.updateCharacterApplicationContentBackstoryStmt,
		updateCharacterApplicationContentDescriptionStmt:      q.updateCharacterApplicationContentDescriptionStmt,
		updateCharacterApplicationContentGenderStmt:           q.updateCharacterApplicationContentGenderStmt,
		updateCharacterApplicationContentNameStmt:             q.updateCharacterApplicationContentNameStmt,
		updateCharacterApplicationContentShortDescriptionStmt: q.updateCharacterApplicationContentShortDescriptionStmt,
		updatePlayerPasswordStmt:                              q.updatePlayerPasswordStmt,
		updatePlayerSettingsThemeStmt:                         q.updatePlayerSettingsThemeStmt,
		updateRequestStatusStmt:                               q.updateRequestStatusStmt,
	}
}
