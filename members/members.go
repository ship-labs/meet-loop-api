// Package members
package members

import (
	"cmp"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Oudwins/zog"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ship-labs/meet-loop-api/internal"
	"github.com/ship-labs/meet-loop-api/internal/pkg/sqlc"
	"github.com/ship-labs/meet-loop-api/middleware"
)

const (
	defaultLimit = 20
)

func GetUserProfile(store *sqlc.Store) middleware.Handler {
	return func(w http.ResponseWriter, r *http.Request) middleware.Handler {
		limitParam := r.URL.Query().Get("limit")
		limit, _ := strconv.Atoi(limitParam)
		limit = cmp.Or(limit, int(defaultLimit))

		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			return middleware.Error(fmt.Errorf("getting user ID: %w", err))
		}

		claims, err := middleware.GetClaims(r.Context())
		if err != nil {
			return middleware.Error(fmt.Errorf("getting user ID: %w", err))
		}

		groups, err := store.GetUserGrops(r.Context(), sqlc.GetUserGropsParams{
			UserID: userID,
			Limit:  int32(limit),
		})
		if err != nil {
			return middleware.Error(fmt.Errorf("getting user groups: %w", err))
		}

		data := map[string]any{
			"user":   claims,
			"groups": groups,
		}

		return middleware.JSON(middleware.Response{
			Data:    data,
			Message: http.StatusText(http.StatusOK),
		})
	}
}

func CreateGroup(store *sqlc.Store) middleware.Handler {
	return func(w http.ResponseWriter, r *http.Request) middleware.Handler {
		type Body struct {
			GroupName        string `json:"group_name" zog:"group_name"`
			GroupDescription string `json:"group_description" zog:"group_description"`
		}

		v := zog.Struct(zog.Shape{
			"GroupName":        zog.String().Required(zog.Message("Group name is required")),
			"GroupDescription": zog.String().Optional(),
		})

		body, err := internal.Validate[Body](v, r.Body)
		if err != nil {
			var v internal.ValidationError
			if errors.As(err, &v) {
				return middleware.Error(v)
			}
			return middleware.Error(fmt.Errorf("validating onboarding data: %w", err))
		}

		var group sqlc.Group
		var member sqlc.Member

		user, err := middleware.GetUserMetadata(r.Context())
		if err != nil {
			return middleware.Error(fmt.Errorf("getting user metadata: %w", err))
		}

		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			return middleware.Error(fmt.Errorf("getting user ID: %w", err))
		}

		transactionError := store.ExecuteTransaction(r.Context(), func() error {
			group, err = store.CreateGroup(r.Context(), sqlc.CreateGroupParams{
				Name: body.GroupName,
				Description: pgtype.Text{
					String: body.GroupDescription,
					Valid:  body.GroupDescription != "",
				},
				UserID: userID,
			})
			if err != nil {
				return fmt.Errorf("creating group: %w", err)
			}

			member, err = store.CreateGroupMember(r.Context(), sqlc.CreateGroupMemberParams{
				GroupID: group.ID,
				Email: pgtype.Text{
					String: user.Email,
					Valid:  true,
				},
				Phone:  user.Phone,
				Name:   user.Name,
				UserID: userID,
			})
			if err != nil {
				return fmt.Errorf("creating group member: %w", err)
			}

			_, err = store.CreateGroupAdmin(r.Context(), sqlc.CreateGroupAdminParams{
				GroupID:  group.ID,
				MemberID: member.ID,
			})
			if err != nil {
				return fmt.Errorf("creating group admin: %w", err)
			}

			return nil
		})

		if transactionError != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == internal.UniqueViolationCode {
					return middleware.Error(fmt.Errorf("group %w", internal.ErrExists))
				}
			}
			return middleware.Error(fmt.Errorf("onboarding user: %w", err))
		}

		return middleware.JSON(middleware.Response{
			Message: http.StatusText(http.StatusCreated),
			Data: map[string]any{
				"member": member,
				"group":  group,
			},
		})
	}
}
