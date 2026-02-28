package auction

import (
	"context"
	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"github.com/joho/godotenv"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCreateAuction_AutoCloseAfterDuration(t *testing.T) {
	_ = godotenv.Load("cmd/auction/.env")

	os.Setenv("AUCTION_DURATION", "100ms")

	mongoURL := os.Getenv("MONGODB_URL")
	if mongoURL == "" {
		mongoURL = "mongodb://localhost:27017"
	}
	if strings.Contains(mongoURL, "@mongodb:") {
		mongoURL = strings.Replace(mongoURL, "@mongodb:", "@localhost:", 1)
	}
	os.Setenv("MONGODB_URL", mongoURL)

	mongoDB := os.Getenv("MONGODB_DB")
	if mongoDB == "" {
		mongoDB = "auctions"
	}
	os.Setenv("MONGODB_DB", mongoDB)

	ctx := context.Background()
	database, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		t.Skip("MongoDB não disponível, pulando teste de integração:", err)
	}

	repo := NewAuctionRepository(database)

	auction, errEntity := auction_entity.CreateAuction(
		"Produto Teste",
		"Categoria",
		"Descrição com mais de dez caracteres",
		auction_entity.New,
	)
	if errEntity != nil {
		t.Fatalf("Erro ao criar entidade de leilão: %v", errEntity)
	}

	if err := repo.CreateAuction(ctx, auction); err != nil {
		t.Fatalf("Erro ao criar leilão: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	found, errFind := repo.FindAuctionById(ctx, auction.Id)
	if errFind != nil {
		t.Fatalf("Erro ao buscar leilão: %v", errFind)
	}

	if found.Status != auction_entity.Completed {
		t.Errorf("Status esperado: %v, obtido: %v. O leilão deveria ter sido fechado automaticamente.", auction_entity.Completed, found.Status)
	}
}
