package usecases

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/FelipeAz/desafio-serasa/internal/pkg/app/entity"
)

// MainframeService representa uma instancia do servico relacionado ao mainframe.
type MainframeService struct {
	NegativacaoRepository NegativacaoRepository
	CryptoHandler         CryptoHandler
}

// ConnectJSONServer conecta com o JSONServer
func (ms *MainframeService) ConnectJSONServer() (*http.Response, error) {
	resp, err := http.Get(os.Getenv("MAINFRAME_URL"))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Get retorna todas as negativacoes do mainframe
func (ms *MainframeService) Get() ([]entity.Negativacao, error) {
	var data []entity.Negativacao
	resp, err := ms.ConnectJSONServer()
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil || err == io.EOF {
		log.Println(err)
		return nil, err
	}

	if (len(data) == 1 && data[0] == entity.Negativacao{}) {
		return nil, fmt.Errorf("record not found")
	}

	return data, nil
}

// Integrate realiza a integracao com o mainframe persistindo as negativacoes no BD.
func (ms *MainframeService) Integrate() error {
	negativacoes, err := ms.Get()
	if err != nil {
		return err
	}

	// Persist Negativacoes
	for i := 0; i < len(negativacoes); i++ {
		_, err := ms.persistNegativacao(negativacoes[i])
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (ms *MainframeService) persistNegativacao(negativacao entity.Negativacao) (id uint, err error) {
	negativacao.CustomerDocument, err = ms.CryptoHandler.EncryptString(negativacao.CustomerDocument)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	id, err = ms.NegativacaoRepository.Create(negativacao)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return
}
