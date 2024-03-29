package interfaces

import (
	"net/http"

	"github.com/FelipeAz/desafio-serasa/internal/pkg/app/usecases"
	"github.com/gin-gonic/gin"
)

// MainframeController representa uma instancia do controller.
type MainframeController struct {
	MainframeService usecases.MainframeService
}

// NewMainframeController retorna uma instancia do controller do Mainframe.
func NewMainframeController(sqlHandler SQLHandler, redis Redis, cryptoHandler CryptoHandler) *MainframeController {
	return &MainframeController{
		MainframeService: usecases.MainframeService{
			NegativacaoRepository: &NegativacaoRepository{
				SQLHandler: sqlHandler,
				Redis:      redis,
			},
			CryptoHandler: &cryptoHandler,
		},
	}
}

// Get retorna todas as Negativacoes do servidor JSON
func (mc *MainframeController) Get(c *gin.Context) {
	data, err := mc.MainframeService.Get()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// Integrate persiste todas as negativacoes do servidor JSON no banco de dados
func (mc *MainframeController) Integrate(c *gin.Context) {
	err := mc.MainframeService.Integrate()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Integration Finished")
}
