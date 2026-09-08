package models

type Usuario struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"unique;not null"`
	Senha        string `gorm:"not null"`
	TipoPessoaID uint   `gorm:"not null"`

	TipoPessoa TipoPessoa        `gorm:"foreignKey:TipoPessoaID"`
	Telefones  []TelefoneUsuario `gorm:"foreignKey:UsuarioID"`
	Enderecos  []EnderecoUsuario `gorm:"foreignKey:UsuarioID"`
	Papeis     []UsuarioRole     `gorm:"foreignKey:UsuarioID"`
}

func (*Usuario) TableName() string {
	return "tb_usuario"
}

type TipoPessoa struct {
	ID   uint   `gorm:"primaryKey"`
	Tipo string `gorm:"unique;not null"`
}

func (*TipoPessoa) TableName() string {
	return "tb_tipo_pessoa"
}

type PessoaFisica struct {
	ID        uint `gorm:"primaryKey"`
	UsuarioID uint `gorm:"unique;not null"`

	Nome           string `gorm:"not null"`
	CPF            string `gorm:"unique;not null"`
	DataNascimento string

	Usuario Usuario `gorm:"foreignKey:UsuarioID"`
}

func (*PessoaFisica) TableName() string {
	return "tb_pessoa_fisica"
}

type PessoaJuridica struct {
	ID        uint `gorm:"primaryKey"`
	UsuarioID uint `gorm:"unique;not null"`

	RazaoSocial  string `gorm:"not null"`
	NomeFantasia string
	CNPJ         string `gorm:"unique;not null"`

	Usuario Usuario `gorm:"foreignKey:UsuarioID"`
}

func (*PessoaJuridica) TableName() string {
	return "tb_pessoa_juridica"
}

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Nome string `gorm:"unique;not null"`
}

func (*Role) TableName() string {
	return "tb_role"
}

type UsuarioRole struct {
	UsuarioID uint    `gorm:"primaryKey"`
	RoleID    uint    `gorm:"primaryKey"`
	Usuario   Usuario `gorm:"foreignKey:UsuarioID"`
	Role      Role    `gorm:"foreignKey:RoleID"`
}

func (*UsuarioRole) TableName() string {
	return "tb_usuario_role"
}

type TelefoneUsuario struct {
	ID        uint   `gorm:"primaryKey"`
	UsuarioID uint   `gorm:"not null"`
	Numero    string `gorm:"not null"`

	Usuario Usuario `gorm:"foreignKey:UsuarioID"`
}

func (*TelefoneUsuario) TableName() string {
	return "tb_telefone_usuario"
}

type EnderecoUsuario struct {
	ID uint `gorm:"primaryKey"`

	UsuarioID uint `gorm:"not null"`

	Logradouro string `gorm:"not null"`
	Cidade     string `gorm:"not null"`
	Estado     string `gorm:"not null"`
	CEP        string `gorm:"not null"`

	EnderecoCobranca bool `gorm:"not null"`

	Usuario Usuario `gorm:"foreignKey:UsuarioID"`
}

func (*EnderecoUsuario) TableName() string {
	return "tb_endereco_usuario"
}

type Imprensa struct {
	ID        uint `gorm:"primaryKey"`
	UsuarioID uint `gorm:"unique;not null"`

	NomeVeiculo string
	Site        string
	CNPJ        string

	Usuario Usuario `gorm:"foreignKey:UsuarioID"`
}

func (*Imprensa) TableName() string {
	return "tb_imprensa"
}
