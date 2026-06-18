package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alizarazot/2026-i-dbm/internal/auth"
	"github.com/alizarazot/2026-i-dbm/internal/database"
	"github.com/alizarazot/2026-i-dbm/internal/mail"
	"github.com/alizarazot/2026-i-dbm/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func handlerManagerListUsers(userStore *database.UserStore) echo.HandlerFunc {
	return func(c *echo.Context) error {
		const pageParam = "page"
		page, err := echo.QueryParam[uint](c, pageParam)
		if err != nil {
			return echo.ErrBadRequest.Wrap(fmt.Errorf("error reading param %q: %w", pageParam, err))
		}

		const limitParam = "limit"
		limit, err := echo.QueryParam[uint](c, limitParam)
		if err != nil {
			return echo.ErrBadRequest.Wrap(fmt.Errorf("error reading param %q: %w", limitParam, err))
		}
		if limit < 1 {
			return echo.ErrBadRequest.Wrap(fmt.Errorf("the provided limit (%d) must be >= 1", limit))
		}

		const roleParam = "role"
		roleRaw := c.QueryParam(roleParam)

		role := model.NewUserRole(roleRaw)
		if roleRaw != "" && role == model.UserRoleInvalid {
			return echo.ErrBadRequest.Wrap(fmt.Errorf("unknown role %q", roleRaw))
		}

		if role == model.UserRoleInvalid {
			users, totalUsers, err := userStore.GetUsers(c.Request().Context(), page, limit)
			if err != nil {
				return echo.ErrInternalServerError.Wrap(err)
			}

			return c.JSON(http.StatusOK, map[string]any{"users": users, "totalUsers": totalUsers})
		}

		users, totalUsers, err := userStore.GetUsersByRole(c.Request().Context(), role, page, limit)
		if err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}

		return c.JSON(http.StatusOK, map[string]any{"users": users, "totalUsers": totalUsers})
	}
}

func handlerManagerAddUser(userStore *database.UserStore) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var data struct {
			User struct {
				Email string         `json:"email"`
				Role  model.UserRole `json:"role"`
				Info  model.UserInfo `json:"info"`
			} `json:"user"`
			InitialPassword string `json:"initialPassword"`
		}
		if err := echo.BindBody(c, &data); err != nil {
			return echo.ErrBadRequest.Wrap(err)
		}

		user := model.User{
			Email: data.User.Email,
			Role:  data.User.Role,
			Info:  data.User.Info,
		}
		if err := userStore.AddUser(c.Request().Context(), &user, data.InitialPassword); err != nil {
			if errors.Is(err, database.ErrUserAlreadyExists) {
				return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("user %q already exists", user.Email)).Wrap(err)
			}

			return echo.ErrInternalServerError.Wrap(err)
		}

		return c.NoContent(http.StatusOK)
	}
}

func handlerManagerEditUser(userStore *database.UserStore) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var data struct {
			User struct {
				Email string         `json:"email"`
				Role  model.UserRole `json:"role"`
				Info  model.UserInfo `json:"info"`
			} `json:"user"`
		}
		if err := echo.BindBody(c, &data); err != nil {
			return echo.ErrBadRequest.Wrap(err)
		}

		user := model.User{
			Email: data.User.Email,
			Role:  data.User.Role,
			Info:  data.User.Info,
		}
		if err := userStore.EditUser(c.Request().Context(), &user); err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}

		return c.NoContent(http.StatusOK)
	}
}

func handlerManagerDeleteUser(userStore *database.UserStore) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var data struct {
			Email string `query:"email"`
		}
		if err := c.Bind(&data); err != nil {
			return echo.ErrBadRequest.Wrap(err)
		}

		if err := userStore.DeleteUser(c.Request().Context(), data.Email); err != nil {
			if errors.Is(err, database.ErrUserNotExists) {
				return echo.ErrNotFound.Wrap(err)
			}

			return echo.ErrInternalServerError.Wrap(err)
		}

		return nil
	}
}

func handlerManagerListCFCs(userStore *database.UserStore, cfcStore *database.CFCStore) echo.HandlerFunc {
	return func(c *echo.Context) error {
		incompleteCFCs, err := cfcStore.ListCFCs(c.Request().Context())
		if err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}

		cfcs := make([]*model.CFC, len(incompleteCFCs))
		for idx, ic := range incompleteCFCs {
			email, err := userStore.GetUserEmail(c.Request().Context(), ic.UserID)
			if err != nil {
				return echo.ErrInternalServerError.Wrap(err)
			}
			ic.CFC.UserEmail = email
			cfcs[idx] = ic.CFC
		}

		return c.JSON(http.StatusOK, map[string]any{"cfcs": cfcs})
	}
}

func handlerManagerAddCFCAnswer(mailService *mail.Service, userStore *database.UserStore, cfcStore *database.CFCStore) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var data struct {
			CFCID  string `json:"cfcId"`
			Answer string
		}
		if err := echo.BindBody(c, &data); err != nil {
			return echo.ErrBadRequest.Wrap(err)
		}

		token, err := echo.ContextGet[*jwt.Token](c, "user")
		if err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}
		user := auth.ExtractUser(token)

		id, err := userStore.GetUserID(c.Request().Context(), user.Email)
		if err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}

		if err := cfcStore.AddCFCAnswer(c.Request().Context(), data.Answer, data.CFCID, id); err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}

		inc, err := cfcStore.GetCFC(c.Request().Context(), data.CFCID)
		if err != nil {
			c.Logger().Error("couldn't get cfc", "err", err)
			return c.NoContent(http.StatusOK)
		}

		email, err := userStore.GetUserEmail(c.Request().Context(), inc.UserID)
		if err != nil {
			c.Logger().Error("couldn't user email", "err", err)
			return c.NoContent(http.StatusOK)
		}

		if err := mailService.SendPDF(email, fmt.Sprintf("Your %s has been answered!", inc.CFC.Category.CanonicalString()), fmt.Sprintf("Answer: %s", data.Answer), fmt.Sprintf("Subject: %s\n\nCategory: %s\n\nDetails: %s\n\nAnswer: %s", inc.CFC.Subject, inc.CFC.Category.CanonicalString(), inc.CFC.Details, data.Answer)); err != nil {
			c.Logger().Error("error sending mail", "err", err)
		}

		return c.NoContent(http.StatusOK)
	}
}

func handlerManagerGetCFCAnswer(userStore *database.UserStore, cfcStore *database.CFCStore) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var data struct {
			ID string `query:"id"`
		}
		if err := c.Bind(&data); err != nil {
			return echo.ErrBadRequest.Wrap(err)
		}

		incomplete, err := cfcStore.GetCFCAnswer(c.Request().Context(), data.ID)
		if err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}
		email, err := userStore.GetUserEmail(c.Request().Context(), incomplete.UserID)
		if err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}
		incomplete.CFCAnswer.UserEmail = email

		return c.JSON(http.StatusOK, map[string]any{"answer": incomplete.CFCAnswer})
	}
}
