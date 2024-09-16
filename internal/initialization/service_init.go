package initialization

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"os"
	"tender-service/internal/config"
	"tender-service/internal/database/postgres"
	"tender-service/internal/database/repositories"
	"tender-service/internal/handlers"
	"tender-service/internal/service"
)

type Repos struct {
	TenderRepo       repositories.TenderRepository
	UserRepo         repositories.UserRepository
	OrganizationRepo repositories.OrgRepository
	BidResp          repositories.BidRepository
}

type Services struct {
	TenderService *service.TenderService
	UserService   *service.UserService
	OrgService    *service.OrganizationService
	BidService    *service.BidService
}

type Handlers struct {
	PingHandler   *handlers.PingHandler
	TenderHandler *handlers.TenderHandler
	BidHandler    *handlers.BidHandler
}

func RepositoriesInit(config *config.Config) *Repos {
	db := mustCreateDatabase(config)

	return &Repos{
		postgres.NewTenderRepository(db),
		postgres.NewUserRepository(db),
		postgres.NewOrganizationRepository(db),
		postgres.NewBidRepository(db),
	}
}

func mustCreateDatabase(config *config.Config) *sql.DB {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.PostgresHost,
		config.PostgresPort,
		config.PostgresUsername,
		config.PostgresPassword,
		config.PostgresDatabase)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(fmt.Errorf("%s: %w", "database.storage.MustCreateDataBase", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("%s: %w", "database.storage.MustCreateDataBase", err))
	}

	op1 := `drop type if exists bid_status; 
			create type BID_STATUS as ENUM ('Created', 'Published', 'Canceled');
	`
	op2 := `drop type if exists service_type; 
			create type SERVICE_TYPE as ENUM('Delivery', 'Manufacture', 'Construction');`

	op22 := `drop type if exists tender_status; 
			create type TENDER_STATUS as ENUM('Published' ,'Open', 'Canceled', 'Created', 'Closed');`

	op3 := `
	create table IF NOT EXISTS Tenders (
    id UUID DEFAULT uuid_generate_v4() NOT NULL,
    version int NOT NULL,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    creator_username VARCHAR(50) references employee(username),
    organization_id UUID references organization(id) NOT NULL,
    status TENDER_STATUS NOT NULL,
    type SERVICE_TYPE NOT NULL,
    last_version BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	op4 := `create table IF NOT EXISTS Bids (
    id UUID DEFAULT uuid_generate_v4() NOT NULL,
    tender_id UUID NOT NULL,


    
    name VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    
    creator_username VARCHAR(50) references employee(username) NOT NULL,
    organization_id UUID references organization(id) NOT NULL,

    status BID_STATUS DEFAULT 'Created',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(op1)
	if err != nil {
		panic(fmt.Errorf("%s: %w", "Op1", err))
	}

	_, err = db.Exec(op2)
	if err != nil {
		panic(fmt.Errorf("%s: %w", "Op2", err))
	}

	_, err = db.Exec(op22)
	if err != nil {
		panic(fmt.Errorf("%s: %w", "Op22", err))
	}

	_, err = db.Exec(op3)
	if err != nil {
		panic(fmt.Errorf("%s: %w", "Op3", err))
	}

	_, err = db.Exec(op4)
	if err != nil {
		panic(fmt.Errorf("%s: %w", "Op3", err))
	}

	return db
}

func ServiceInit(repos *Repos) *Services {
	userService := service.NewUserService(repos.UserRepo)
	organService := service.InitOrganizationService(repos.OrganizationRepo)
	tenderService := service.NewTenderService(repos.TenderRepo, userService, organService)
	bidService := service.InitBidService(repos.BidResp, userService, organService, tenderService)
	return &Services{
		tenderService,
		userService,
		organService,
		bidService,
	}
}

func HandlersInit(services *Services, logger *slog.Logger) *Handlers {
	pingHandler := handlers.InitPingHandler()
	tenderHandler := handlers.InitTenderHandler(services.TenderService, services.UserService, services.OrgService, logger)
	bidHandler := handlers.InitBidHandler(services.TenderService, services.UserService, services.OrgService, services.BidService, logger)
	return &Handlers{
		pingHandler,
		tenderHandler,
		bidHandler,
	}
}

func InitChiServer(handlers *Handlers) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Get("/api/ping", handlers.PingHandler.Ping())
	router.Post("/api/tenders/new", handlers.TenderHandler.CreateTender())
	router.Patch("/api/tenders/{id}/edit", handlers.TenderHandler.Edit())
	router.Get("/api/tenders", handlers.TenderHandler.Tenders())
	router.Get("/api/tenders/my", handlers.TenderHandler.MyTenders())
	router.Put("/api/tenders/{id}/rollback/{version}", handlers.TenderHandler.TenderRollback())
	router.Post("/api/bids/new", handlers.BidHandler.CreateBid())
	router.Get("/api/bids/{id}/list", handlers.BidHandler.TenderBids())
	router.Get("/api/bids/my", handlers.BidHandler.MyBid())
	router.Get("/api/bids/{id}/status", handlers.BidHandler.BidStatus())
	router.Post("/api/bids/{id}/status", handlers.BidHandler.ChangeBidStatus())
	return router
}

func InitLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}
