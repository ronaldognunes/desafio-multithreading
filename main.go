package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilCep struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func main() {
	fmt.Println("digite o CEP:")
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Erro ao ler o CEP")
		panic("Erro ao ler o CEP")
	}

	input = strings.TrimSpace(input)
	fmt.Println("CEP digitado2: ", input)

	chanViaCEP := make(chan *ViaCEP)
	chanBrasilCep := make(chan *BrasilCep)
	go func() {
		resultcep, _ := ConsultaCepViaSep(input)
		//time.Sleep(2 * time.Second) //Simulando um tempo de resposta
		chanViaCEP <- resultcep
	}()

	go func() {
		resultcep, _ := ConsultaCepBrasilApi(input)
		//time.Sleep(2 * time.Second) //Simulando um tempo de resposta
		chanBrasilCep <- resultcep
	}()

	select {
	case retornoViaCep := <-chanViaCEP:
		fmt.Printf("API -> ViaCEP: %+v", retornoViaCep)
	case retornoBasilApi := <-chanBrasilCep:
		fmt.Printf("API ->BrasilAPI: %+v", retornoBasilApi)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout")
	}

}

func ConsultaCepViaSep(cep string) (*ViaCEP, error) {
	req, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	retorno, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	var dados ViaCEP
	err = json.Unmarshal(retorno, &dados)
	if err != nil {
		return nil, err
	}
	return &dados, nil
}

func ConsultaCepBrasilApi(cep string) (*BrasilCep, error) {
	req, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	retorno, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	var dados BrasilCep
	err = json.Unmarshal(retorno, &dados)
	if err != nil {
		return nil, err
	}
	return &dados, nil
}
