SDK para integração com [Melhor Envio](https://melhorenvio.com.br).

Código em qualidade `alpha`. Não está completo e pode mudar a qualquer momento.

## Exemplo de uso

### Primeira autenticação

```go
	client = melhorenvio.NewClient(ctx, melhorenvio.Config{
		Credentials: melhorenvio.Credentials{
			ClientId:     1234,
			ClientSecret: "{secret}",
			Code:         "{code do oauth2}",
		},
		RedirectUri:     "{url}",
		ApplicationName: "{nome do app}",
		Email:           "{email de contato técnico}",
		CredentialsChangedCallback: func(credentials melhorenvio.Credentials) error {
			// função executada de forma síncrona
			// executada em toda atualização de token
			// o salvamento das credenciais pode ser feito por aqui
			// o erro retornado aqui é repassado para a chamada da função
			// que executou a atualização de token
			return nil
		},
	})

	// executa a requisição de obtenção de token a partir do code
	err := client.AutenticateByCode()
	if err != nil {
		panic(err)
	}
```

### Cotação de Frete

```go
	client = melhorenvio.NewClient(ctx, melhorenvio.Config{
		Credentials: melhorenvio.Credentials{
			ClientId:     1234,
			ClientSecret: "{secret}",
			AccessToken:  "{access token}",
			RefreshToken: "{refresh token}",
			ExpiresAt:    {data de expiração do token (time.Time)},
		},
		RedirectUri:     "{url}",
		ApplicationName: "{nome do app}",
		Email:           "{email de contato técnico}",
		CredentialsChangedCallback: func(credentials melhorenvio.Credentials) error {
			// caso ocorra atualização do token (refresh), essa função será chamada
			return nil
		},
	})

	resp, err := client.CotarFrete(&melhorenvio.CotacaoRequest{
		From: melhorenvio.ToFrom{
			PostalCode: "{cep origem}",
		},
		To: melhorenvio.ToFrom{
			PostalCode: "{cep destino}",
		},
		Products: []melhorenvio.Product{
			{
				ID:             "1",
				Width:          15,
				Height:         4,
				Length:         12,
				Weight:         0.3,
				InsuranceValue: 15.75,
				Quantity:       1,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	for _, r := range resp {
		fmt.Printf("%+v\n", r)
	}
```
