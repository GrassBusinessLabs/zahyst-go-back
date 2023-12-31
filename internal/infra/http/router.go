package http

import (
	"boilerplate/config"
	"boilerplate/config/container"
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/controllers"
	"boilerplate/internal/infra/http/middlewares"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Router(cont container.Container) http.Handler {

	router := chi.NewRouter()

	router.Use(middleware.RedirectSlashes, middleware.Logger, cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Route("/api", func(apiRouter chi.Router) {
		// Health
		apiRouter.Route("/ping", func(healthRouter chi.Router) {
			healthRouter.Get("/", PingHandler())
			healthRouter.Handle("/*", NotFoundJSON())
		})

		apiRouter.Route("/v1", func(apiRouter chi.Router) {
			// Public routes
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Route("/auth", func(apiRouter chi.Router) {
					AuthRouter(apiRouter, cont.AuthController, cont.AuthMw)
				})
			})

			// Protected routes
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Use(cont.AuthMw)

				LocationRouter(apiRouter, cont.LocationController, cont.LocationService)
				UserRouter(apiRouter, cont.UserController)
				GroupRouter(apiRouter, cont.GroupController, cont.GroupService)
				GroupMemberRouter(apiRouter, cont.GroupMemberController, cont.GroupMemberService, cont.GroupService)

				apiRouter.Handle("/*", NotFoundJSON())
			})
		})
	})

	router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, config.GetConfiguration().FileStorageLocation))
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
		fs.ServeHTTP(w, r)
	})

	return router
}

func AuthRouter(r chi.Router, ac controllers.AuthController, amw func(http.Handler) http.Handler) {
	r.Route("/", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/register",
			ac.Register(),
		)
		apiRouter.Post(
			"/login",
			ac.Login(),
		)
		apiRouter.With(amw).Post(
			"/change-pwd",
			ac.ChangePassword(),
		)
		apiRouter.With(amw).Post(
			"/logout",
			ac.Logout(),
		)
	})
}

func UserRouter(r chi.Router, uc controllers.UserController) {
	r.Route("/users", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/",
			uc.FindMe(),
		)
		apiRouter.Put(
			"/",
			uc.Update(),
		)
		apiRouter.Delete(
			"/",
			uc.Delete(),
		)
		apiRouter.Put(
			"/coordinates",
			uc.SetCoordinates(),
		)
		apiRouter.Get(
			"/coordinates",
			uc.GetCoordinates(),
		)
	})
}

func LocationRouter(r chi.Router, lc controllers.LocationController, ls app.LocationService) {
	r.Route("/locations", func(apiRouter chi.Router) {
		lpom := middlewares.PathObject("locationId", controllers.LocationKey, ls)
		omw := middlewares.IsOwnerMiddleware[domain.Location](controllers.LocationKey)
		apiRouter.Post(
			"/",
			lc.Save(),
		)
		apiRouter.Get(
			"/my",
			lc.FindByUserId(),
		)
		apiRouter.Post(
			"/in-area",
			lc.FindByArea(),
		)
		apiRouter.With(lpom).Get(
			"/{locationId}",
			lc.Detail(),
		)
		apiRouter.With(lpom, omw).Put(
			"/{locationId}",
			lc.Update(),
		)
		apiRouter.With(lpom, omw).Delete(
			"/{locationId}",
			lc.Delete(),
		)
	})
}

func GroupRouter(r chi.Router, gc controllers.GroupController, gs app.GroupService) {
	r.Route("/groups", func(apiRouter chi.Router) {
		gpom := middlewares.PathObject("groupId", controllers.GroupKey, gs)
		omw := middlewares.IsOwnerMiddleware[domain.Group](controllers.GroupKey)
		apiRouter.Post(
			"/",
			gc.Save(),
		)
		apiRouter.Get(
			"/list",
			gc.GetList(),
		)
		apiRouter.With(gpom, omw).Get(
			"/access_code/{groupId}",
			gc.GetAccessCode(),
		)
		apiRouter.With(gpom).Get(
			"/{groupId}",
			gc.Detail(),
		)
		apiRouter.With(gpom, omw).Put(
			"/{groupId}",
			gc.Update(),
		)
		apiRouter.With(gpom, omw).Delete(
			"/{groupId}",
			gc.Delete(),
		)
	})
}

func GroupMemberRouter(r chi.Router, gmc controllers.GroupMemberController, gms app.GroupMemberService, gs app.GroupService) {
	r.Route("/members", func(apiRouter chi.Router) {
		gmpom := middlewares.PathObject("groupMemberId", controllers.GroupMemberKey, gms)
		ismoderator := middlewares.CheckRoleMiddleware([]domain.AccessLevel{domain.ModeratorAccessLevel{}, domain.AdminAccessLevel{}}, gs, gms, "groupId")
		isadmin := middlewares.CheckRoleMiddleware([]domain.AccessLevel{domain.AdminAccessLevel{}}, gs, gms, "groupId")
		apiRouter.Post(
			"/",
			gmc.AddGroupMember(),
		)
		apiRouter.With(gmpom, isadmin).Put(
			"/{groupId}/{groupMemberId}",
			gmc.ChangeAccessLevel(),
		)
		apiRouter.With(gmpom, isadmin).Delete(
			"/{groupId}/{groupMemberId}",
			gmc.DeleteGroupMember(),
		)
		apiRouter.With(ismoderator).Get(
			"/{groupId}",
			gmc.GetMembersList(),
		)
		apiRouter.With(ismoderator).Post(
			"/{groupId}",
			gmc.FindMembersByArea(),
		)
	})
}

func NotFoundJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode("Resource Not Found")
		if err != nil {
			fmt.Printf("writing response: %s", err)
		}
	}
}

func PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode("Ok")
		if err != nil {
			fmt.Printf("writing response: %s", err)
		}
	}
}
