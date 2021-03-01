package interfaces

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/FelipeAz/desafio-serasa/app/entity"
	"github.com/FelipeAz/desafio-serasa/app/usecases"
	"github.com/dgrijalva/jwt-go"
)

// JWTAuth contem variaveis do servico.
type JWTAuth struct {
	SQLHandler SQLHandler
	issure     string
}

// NewJWTAuth instancia do servico.
func NewJWTAuth(sqlHandler SQLHandler) *JWTAuth {
	return &JWTAuth{
		SQLHandler: sqlHandler,
		issure:     "Felipe",
	}
}

// CreateToken cria um token JWT.
func (jwtAuth *JWTAuth) CreateToken(auth entity.Access) (td *usecases.TokenDetails, err error) {
	td = &usecases.TokenDetails{
		AtExpires: time.Now().Add(time.Minute * 15).Unix(),
		RtExpires: time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["token"] = auth.AccessToken
	claims["user_id"] = auth.UserID
	claims["issuer"] = jwtAuth.issure
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	td.AccessToken, err = token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err

	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_token"] = auth.RefreshToken
	rtClaims["user_id"] = auth.UserID
	rtClaims["sub"] = 1
	rtClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))

	return
}

// TokenValid retorna se o token eh valido.
func (jwtAuth *JWTAuth) TokenValid(r *http.Request, tokenString string) error {
	token, err := jwtAuth.VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	if jwtAuth.FetchToken(tokenString) == false {
		return fmt.Errorf("refresh token invalid")
	}

	return nil
}

// VerifyToken extrai o token e verifica.
func (jwtAuth *JWTAuth) VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := jwtAuth.ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// ExtractToken extrai o token da requisicao.
func (jwtAuth *JWTAuth) ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearToken := r.Header.Get("Authorization")

	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// FetchToken verifica se existe um usuario logado no BD atraves do token da requisicao.
// Se o usuario estiver logado mas o token for invalido, revalidamos sua sessao SE o refresh token estiver valido.
// Caso o usuario deslogue, sera necessario logar novamente e utilizar um novo token, pois o mesmo eh apagado quando
// damos logout.
func (jwtAuth *JWTAuth) FetchToken(token string) bool {
	var access entity.Access
	db := jwtAuth.SQLHandler.GetGorm()
	if err := db.Where("access_token=?", token).First(&access).Error; err != nil {
		return false
	}

	return jwtAuth.refreshToken(access.RefreshToken)
}

// refreshToken atualiza o token do usuario
func (jwtAuth *JWTAuth) refreshToken(refreshToken string) bool {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if int(claims["sub"].(float64)) == 1 {
			return true
		}

		return false
	}

	if err != nil {
		return false
	}

	return false
}
