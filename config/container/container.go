package container

import (
	"boilerplate/config"
	"boilerplate/internal/app"
	"boilerplate/internal/infra/database"
	"boilerplate/internal/infra/http/controllers"
	"boilerplate/internal/infra/http/middlewares"

	"github.com/go-chi/jwtauth/v5"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"

	//"github.com/upper/db/v4/adapter/sqlite"
	"log"
	"net/http"
)

type Container struct {
	Middlewares
	Services
	Controllers
}

type Middlewares struct {
	AuthMw func(http.Handler) http.Handler
}

type Services struct {
	app.AuthService
	app.UserService
	app.LocationService
	app.GroupService
	app.GroupMemberService
}

type Controllers struct {
	controllers.AuthController
	controllers.UserController
	controllers.LocationController
	controllers.GroupController
	controllers.GroupMemberController
}

func New(conf config.Configuration) Container {
	tknAuth := jwtauth.New("HS256", []byte(conf.JwtSecret), nil)
	sess := getDbSess(conf)

	userRepository := database.NewUserRepository(sess)
	sessionRepository := database.NewSessRepository(sess)
	locationRepository := database.NewLocationRepository(sess)
	groupRepository := database.NewGroupRepository(sess)
	groupMemberRepository := database.NewGroupMemberRepository(sess)

	userService := app.NewUserService(userRepository)
	authService := app.NewAuthService(sessionRepository, userService, conf, tknAuth)
	locationService := app.NewLocationService(locationRepository)
	groupService := app.NewGroupService(groupRepository)
	groupMemberService := app.NewGroupMemberService(groupMemberRepository, groupRepository, userRepository)

	authController := controllers.NewAuthController(authService, userService)
	userController := controllers.NewUserController(userService)
	locationController := controllers.NewLocationController(locationService)
	groupController := controllers.NewGroupController(groupService)
	groupMemberController := controllers.NewGroupMemberController(groupMemberService)

	authMiddleware := middlewares.AuthMiddleware(tknAuth, authService, userService)

	return Container{
		Middlewares: Middlewares{
			AuthMw: authMiddleware,
		},
		Services: Services{
			authService,
			userService,
			locationService,
			groupService,
			groupMemberService,
		},
		Controllers: Controllers{
			authController,
			userController,
			locationController,
			groupController,
			groupMemberController,
		},
	}
}

func getDbSess(conf config.Configuration) db.Session {
	sess, err := postgresql.Open(
		postgresql.ConnectionURL{
			User:     conf.DatabaseUser,
			Host:     conf.DatabaseHost,
			Password: conf.DatabasePassword,
			Database: conf.DatabaseName,
		})
	//sess, err := sqlite.Open(
	//	sqlite.ConnectionURL{
	//		Database: conf.DatabasePath,
	//	})
	if err != nil {
		log.Fatalf("Unable to create new DB session: %q\n", err)
	}
	return sess
}
