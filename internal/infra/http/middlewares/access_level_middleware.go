package middlewares

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/controllers"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type FindableMember interface {
	FindMember(uint64, uint64) (domain.GroupMember, error)
}

func CheckRoleMiddleware(accessLevels []domain.AccessLevel, groupService app.GroupService, service FindableMember, groupPathKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			var (
				err           error
				accessGranted bool
			)

			ctx := r.Context()
			groupId, err := strconv.ParseUint(chi.URLParam(r, groupPathKey), 10, 64)
			if err != nil {
				err = fmt.Errorf("invalid %s parameter(only non-negative integers)", groupPathKey)
				log.Print(err)
				controllers.BadRequest(w, err)
				return
			}

			grp, err := groupService.Find(groupId)
			if err != nil {
				log.Print(err)
				controllers.BadRequest(w, err)
				return
			}

			user := ctx.Value(controllers.UserKey).(domain.User)
			grpDomain := grp.(domain.Group)
			if grpDomain.UserId == user.Id {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			member, err := service.FindMember(user.Id, grpDomain.Id)
			if err != nil {
				err = errors.New("access denied. You are not the member of the group")
				log.Print(err)
				controllers.BadRequest(w, err)
				return
			}

			for _, accessLevel := range accessLevels {
				if member.AccessLevel == accessLevel.GetRole() {
					accessGranted = true
					break
				}
			}

			if !accessGranted {
				err = errors.New("access denied. You have no access level to this route")
				controllers.Forbidden(w, err)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
