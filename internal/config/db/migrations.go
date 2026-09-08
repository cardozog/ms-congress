package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	eventos "ms-congress/internal/models/evento"
	models "ms-congress/internal/models/usuario"

	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Migration struct {
	Version   uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	AppliedAt time.Time `gorm:"not null;autoCreateTime"`
}

func (Migration) TableName() string {
	return "schema_migrations"
}

func RunMigrations(ctx context.Context, database *gorm.DB) error {
	database = database.WithContext(ctx)
	if err := database.AutoMigrate(&Migration{}); err != nil {
		return err
	}

	if err := database.AutoMigrate(
		&models.Usuario{},
		&models.TipoPessoa{},
		&models.PessoaFisica{},
		&models.PessoaJuridica{},
		&models.Role{},
		&models.UsuarioRole{},
		&models.TelefoneUsuario{},
		&models.EnderecoUsuario{},
		&models.Imprensa{},
		&eventos.Evento{},
		&eventos.TipoIngresso{},
		&eventos.EventoIngresso{},
		&eventos.Ingresso{},
	); err != nil {
		return err
	}
	if err := seedData(database); err != nil {
		return err
	}
	log.Default().Println("Migrations applied successfully")
	return database.Where(Migration{Version: 1}).FirstOrCreate(&Migration{
		Version: 1,
		Name:    "initialize_schema_migrations",
	}).Error
}

func seedData(database *gorm.DB) error {
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	if email == "" || password == "" {
		return fmt.Errorf("ADMIN_EMAIL and ADMIN_PASSWORD must be set")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return database.Transaction(func(tx *gorm.DB) error {
		tiposPessoa := map[string]models.TipoPessoa{}
		for _, nome := range []string{"FISICA", "JURIDICA"} {
			tipo := models.TipoPessoa{Tipo: nome}
			if err := tx.Where("tipo = ?", nome).FirstOrCreate(&tipo).Error; err != nil {
				return err
			}
			tiposPessoa[nome] = tipo
		}

		roles := map[string]models.Role{}
		for _, nome := range []string{"ADMIN", "ORGANIZADOR", "PARTICIPANTE", "IMPRENSA"} {
			role := models.Role{Nome: nome}
			if err := tx.Where("nome = ?", nome).FirstOrCreate(&role).Error; err != nil {
				return err
			}
			roles[nome] = role
		}

		admin, err := seedUsuario(tx, email, string(hashedPassword), tiposPessoa["FISICA"].ID)
		if err != nil {
			return err
		}
		if err := seedUsuarioRole(tx, admin.ID, roles["ADMIN"].ID); err != nil {
			return err
		}
		if err := seedPessoaFisica(tx, admin.ID, "Administrador", "00000000001", "1980-01-01"); err != nil {
			return err
		}

		organizador, err := seedUsuario(tx, "organizador@local.com", string(hashedPassword), tiposPessoa["JURIDICA"].ID)
		if err != nil {
			return err
		}
		if err := seedUsuarioRole(tx, organizador.ID, roles["ORGANIZADOR"].ID); err != nil {
			return err
		}
		organizacao, err := seedPessoaJuridica(tx, organizador.ID, "Congresso Local LTDA", "Congresso Local", "00000000000191")
		if err != nil {
			return err
		}

		participante, err := seedUsuario(tx, "participante@local.com", string(hashedPassword), tiposPessoa["FISICA"].ID)
		if err != nil {
			return err
		}
		if err := seedUsuarioRole(tx, participante.ID, roles["PARTICIPANTE"].ID); err != nil {
			return err
		}
		if err := seedPessoaFisica(tx, participante.ID, "Participante Exemplo", "00000000002", "1990-02-02"); err != nil {
			return err
		}

		imprensa, err := seedUsuario(tx, "imprensa@local.com", string(hashedPassword), tiposPessoa["JURIDICA"].ID)
		if err != nil {
			return err
		}
		if err := seedUsuarioRole(tx, imprensa.ID, roles["IMPRENSA"].ID); err != nil {
			return err
		}
		if err := seedImprensa(tx, imprensa.ID, "Portal Congresso Local", "https://congresso.local", "00000000000272"); err != nil {
			return err
		}

		for _, usuarioID := range []uint{admin.ID, organizador.ID, participante.ID, imprensa.ID} {
			if err := seedTelefone(tx, usuarioID, "11999990000"); err != nil {
				return err
			}
			if err := seedEndereco(tx, usuarioID); err != nil {
				return err
			}
		}

		evento, err := seedEvento(tx, uint64(organizacao.ID))
		if err != nil {
			return err
		}
		inteira, err := seedTipoIngresso(tx, "Inteira")
		if err != nil {
			return err
		}
		meia, err := seedTipoIngresso(tx, "Meia-entrada")
		if err != nil {
			return err
		}
		loteInteira, err := seedEventoIngresso(tx, evento.ID, inteira.ID, 150, 100)
		if err != nil {
			return err
		}
		loteMeia, err := seedEventoIngresso(tx, evento.ID, meia.ID, 75, 50)
		if err != nil {
			return err
		}
		if err := seedIngresso(tx, loteInteira.ID, "MS-CONGRESS-001"); err != nil {
			return err
		}
		return seedIngresso(tx, loteMeia.ID, "MS-CONGRESS-002")
	})
}

func seedUsuario(database *gorm.DB, email, senha string, tipoPessoaID uint) (models.Usuario, error) {
	usuario := models.Usuario{}
	err := database.Where("email = ?", email).First(&usuario).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		usuario = models.Usuario{Email: email, Senha: senha, TipoPessoaID: tipoPessoaID}
		err = database.Create(&usuario).Error
	}
	return usuario, err
}

func seedUsuarioRole(database *gorm.DB, usuarioID, roleID uint) error {
	return database.Where(models.UsuarioRole{UsuarioID: usuarioID, RoleID: roleID}).FirstOrCreate(&models.UsuarioRole{}).Error
}

func seedPessoaFisica(database *gorm.DB, usuarioID uint, nome, cpf, dataNascimento string) error {
	pessoa := models.PessoaFisica{}
	err := database.Where("usuario_id = ?", usuarioID).First(&pessoa).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return database.Create(&models.PessoaFisica{UsuarioID: usuarioID, Nome: nome, CPF: cpf, DataNascimento: dataNascimento}).Error
	}
	return err
}

func seedPessoaJuridica(database *gorm.DB, usuarioID uint, razaoSocial, nomeFantasia, cnpj string) (models.PessoaJuridica, error) {
	pessoa := models.PessoaJuridica{}
	err := database.Where("usuario_id = ?", usuarioID).First(&pessoa).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		pessoa = models.PessoaJuridica{UsuarioID: usuarioID, RazaoSocial: razaoSocial, NomeFantasia: nomeFantasia, CNPJ: cnpj}
		err = database.Create(&pessoa).Error
	}
	return pessoa, err
}

func seedImprensa(database *gorm.DB, usuarioID uint, nomeVeiculo, site, cnpj string) error {
	imprensa := models.Imprensa{}
	err := database.Where("usuario_id = ?", usuarioID).First(&imprensa).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return database.Create(&models.Imprensa{UsuarioID: usuarioID, NomeVeiculo: nomeVeiculo, Site: site, CNPJ: cnpj}).Error
	}
	return err
}

func seedTelefone(database *gorm.DB, usuarioID uint, numero string) error {
	telefone := models.TelefoneUsuario{}
	err := database.Where("usuario_id = ? AND numero = ?", usuarioID, numero).First(&telefone).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return database.Create(&models.TelefoneUsuario{UsuarioID: usuarioID, Numero: numero}).Error
	}
	return err
}

func seedEndereco(database *gorm.DB, usuarioID uint) error {
	endereco := models.EnderecoUsuario{}
	err := database.Where("usuario_id = ?", usuarioID).First(&endereco).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return database.Create(&models.EnderecoUsuario{UsuarioID: usuarioID, Logradouro: "Avenida Central, 100", Cidade: "Sao Paulo", Estado: "SP", CEP: "01000000", EnderecoCobranca: true}).Error
	}
	return err
}

func seedEvento(database *gorm.DB, organizadorID uint64) (eventos.Evento, error) {
	evento := eventos.Evento{}
	err := database.Where("nome = ?", "Congresso de Tecnologia 2026").First(&evento).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		inicio := time.Date(2026, time.November, 10, 9, 0, 0, 0, time.UTC)
		evento = eventos.Evento{Nome: "Congresso de Tecnologia 2026", Descricao: "Evento demonstrativo do MS Congress", DataInicio: inicio, DataFim: inicio.Add(72 * time.Hour), OrganizadorID: organizadorID, Logradouro: "Avenida Paulista, 1000", Cidade: "Sao Paulo", Estado: "SP", CEP: "01310000"}
		err = database.Create(&evento).Error
	}
	return evento, err
}

func seedTipoIngresso(database *gorm.DB, nome string) (eventos.TipoIngresso, error) {
	tipo := eventos.TipoIngresso{}
	err := database.Where("nome = ?", nome).First(&tipo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		tipo = eventos.TipoIngresso{Nome: nome}
		err = database.Create(&tipo).Error
	}
	return tipo, err
}

func seedEventoIngresso(database *gorm.DB, eventoID, tipoIngressoID uint64, preco float64, quantidade int) (eventos.EventoIngresso, error) {
	lote := eventos.EventoIngresso{}
	err := database.Where("evento_id = ? AND tipo_ingresso_id = ?", eventoID, tipoIngressoID).First(&lote).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		inicio := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
		lote = eventos.EventoIngresso{EventoID: eventoID, TipoIngressoID: tipoIngressoID, Preco: preco, Quantidade: quantidade, Disponivel: quantidade, DataInicioVenda: inicio, DataFimVenda: time.Date(2026, time.November, 9, 23, 59, 59, 0, time.UTC)}
		err = database.Create(&lote).Error
	}
	return lote, err
}

func seedIngresso(database *gorm.DB, eventoIngressoID uint64, codigo string) error {
	ingresso := eventos.Ingresso{}
	err := database.Where("codigo = ?", codigo).First(&ingresso).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return database.Create(&eventos.Ingresso{EventoIngressoID: eventoIngressoID, Codigo: codigo, Status: "DISPONIVEL"}).Error
	}
	return err
}
