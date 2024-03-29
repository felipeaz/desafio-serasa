package infrastructure

import (
	"log"

	"github.com/FelipeAz/desafio-serasa/internal/pkg/app/interfaces"
	"github.com/FelipeAz/desafio-serasa/internal/pkg/app/middleware"
	"github.com/gin-gonic/gin"
)

// Router .
type Router struct {
	router *gin.Engine
}

// NewRouter retorna uma instancia do Router
func NewRouter() interfaces.Router {
	return Router{router: gin.Default()}
}

// Dispatch inicializa as rotas e redireciona aos controllers.
func (r Router) Dispatch(sqlHandler interfaces.SQLHandler, rds interfaces.Redis) {
	jwt := interfaces.NewJWTAuth(sqlHandler)
	cryptoHandler := *interfaces.NewCryptoHandler()

	userController := interfaces.NewUserController(sqlHandler, *jwt)
	r.router.POST("/signup", userController.SignUp)
	r.router.POST("/logout", userController.Logout)
	r.router.POST("/login", userController.Login)

	negativacaoController := interfaces.NewNegativacaoController(sqlHandler, rds, cryptoHandler)
	rg := r.router.Group("/negativacao")
	rg.GET("/", middleware.AuthorizeJWT(jwt), negativacaoController.Get)
	rg.GET("/:cpf", middleware.AuthorizeJWT(jwt), negativacaoController.GetByCPF)
	rg.POST("/", middleware.AuthorizeJWT(jwt), negativacaoController.Create)
	rg.PUT("/:id", middleware.AuthorizeJWT(jwt), negativacaoController.Update)
	rg.DELETE("/:id", middleware.AuthorizeJWT(jwt), negativacaoController.Delete)

	mainframeController := interfaces.NewMainframeController(sqlHandler, rds, cryptoHandler)
	rg = r.router.Group("/mainframe")
	rg.GET("/", middleware.AuthorizeJWT(jwt), mainframeController.Get)
	rg.GET("/integrate", middleware.AuthorizeJWT(jwt), mainframeController.Integrate)

	r.listen()
}

func (r Router) listen() {
	err := r.router.Run(":8080")
	if err != nil {
		log.Println(err)
	}
}
