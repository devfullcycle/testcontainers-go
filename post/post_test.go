package post

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

type databaseContainer struct {
	*mysql.MySQLContainer
	connectionString string
}

func setupContainer(ctx context.Context) (*databaseContainer, error) {
	mysqlContainer, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8.0.36"),
		mysql.WithDatabase("posts"),
		mysql.WithUsername("root"),
		mysql.WithPassword("root"),
		mysql.WithScripts("../schema.sql"),
	)
	if err != nil {
		log.Fatalf("Could not start mysql container: %s", err)
	}

	connStr, err := mysqlContainer.ConnectionString(ctx, "parseTime=true")
	if err != nil {
		log.Fatalf("Could not get connection string: %s", err)
	}

	return &databaseContainer{
		MySQLContainer:   mysqlContainer,
		connectionString: connStr,
	}, nil
}

func TestAddPost(t *testing.T) {
	ctx := context.Background()
	container, err := setupContainer(ctx)
	if err != nil {
		t.Fatalf("Could not setup container: %s", err)
	}

	post := &Post{
		ID:      1,
		Content: "Hello World!",
		Author:  "Sau Sheong",
	}

	post.DB, err = sql.Open("mysql", container.connectionString)
	if err != nil {
		t.Fatalf("Could not open a connection to the database: %s", err)
	}
	defer post.DB.Close()
	defer container.Terminate(ctx)

	if err := post.Create(); err != nil {
		t.Fatalf("Could not create post: %s", err)
	}

	var postResult Post
	postResult.DB = post.DB
	postResult, err = postResult.GetPost(1)
	if err != nil {
		t.Fatalf("Could not retrieve post: %s", err)
	}

	assert.Equal(t, postResult.ID, post.ID)
	assert.Equal(t, postResult.Content, post.Content)
	assert.Equal(t, postResult.Author, post.Author)
}
