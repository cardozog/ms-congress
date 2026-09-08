package models

type PessoaFisicaResponseDTO struct {
	ID             uint   `json:"id"`
	UsuarioID      uint   `json:"usuarioID"`
	Nome           string `json:"nome"`
	DataNascimento string `json:"dataNascimento"`
}

type PessoaJuridicaResponseDTO struct {
	ID           uint   `json:"id"`
	UsuarioID    uint   `json:"usuarioID"`
	RazaoSocial  string `json:"razaoSocial"`
	NomeFantasia string `json:"nomeFantasia"`
}

type ImprensaResponseDTO struct {
	ID          uint   `json:"id"`
	UsuarioID   uint   `json:"usuarioID"`
	NomeVeiculo string `json:"nomeVeiculo"`
	Site        string `json:"site"`
}

type TipoPessoaResponseDTO struct {
	ID   uint   `json:"id"`
	Tipo string `json:"tipo"`
}

type RoleResponseDTO struct {
	ID   uint   `json:"id"`
	Nome string `json:"nome"`
}

type TelefoneResponseDTO struct {
	ID     uint   `json:"id"`
	Numero string `json:"numero"`
}

type EnderecoResponseDTO struct {
	ID               uint   `json:"id"`
	Logradouro       string `json:"logradouro"`
	Cidade           string `json:"cidade"`
	Estado           string `json:"estado"`
	CEP              string `json:"cep"`
	EnderecoCobranca bool   `json:"enderecoCobranca"`
}

type UsuarioDetalhadoResponseDTO struct {
	ID             uint                       `json:"id"`
	Email          string                     `json:"email"`
	TipoPessoa     TipoPessoaResponseDTO      `json:"tipoPessoa"`
	Papeis         []RoleResponseDTO          `json:"papeis"`
	Telefones      []TelefoneResponseDTO      `json:"telefones"`
	Enderecos      []EnderecoResponseDTO      `json:"enderecos"`
	PessoaFisica   *PessoaFisicaResponseDTO   `json:"pessoaFisica,omitempty"`
	PessoaJuridica *PessoaJuridicaResponseDTO `json:"pessoaJuridica,omitempty"`
	Imprensa       *ImprensaResponseDTO       `json:"imprensa,omitempty"`
}
