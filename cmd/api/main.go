package main

import (
	accountServer "TransactionManager/internal/account_service"
	accountrepo "TransactionManager/internal/account_service/repo"
	accountservice "TransactionManager/internal/account_service/service"
	transactionServer "TransactionManager/internal/transaction_service"
	transactionrepo "TransactionManager/internal/transaction_service/repo"
	transactionservice "TransactionManager/internal/transaction_service/service"
	"TransactionManager/packages/configloader"
	"TransactionManager/packages/database"
	"TransactionManager/packages/logger"
	"TransactionManager/packages/server"
	"fmt"
	"log"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	cfg, err := configloader.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	isProduction := cfg.App.Env == "prod"
	appLogger, err := logger.New(isProduction)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()

	db, err := database.NewGormDB(cfg.Database)
	if err != nil {
		appLogger.Fatal("failed to connect to database", zap.Error(err))
	}

	accountRepo := accountrepo.NewRepository(db)
	accountSvc := accountservice.NewAccountService(accountRepo, appLogger)

	transactionRepo := transactionrepo.NewRepository(db)
	transactionSvc := transactionservice.NewTransactionService(transactionRepo, accountRepo, appLogger)

	srv := server.New()
	root := chi.NewRouter()
	srv.Router().Mount("/", root)
	accountServer.NewAccountHandlerServer(root, accountSvc)
	transactionServer.NewTransactionHandlerServer(root, transactionSvc)

	serverPort := fmt.Sprintf(":%s", cfg.App.Port)
	srv.Start(serverPort)
}
