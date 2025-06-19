// Package tu is the test utility
package tu

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math"
	"math/big"
	"os"

	"github.com/acoshift/pgsql/pgctx"

	"user-service/schema"
)

// Context is the test context contain test's dependencies
type Context struct {
	database       string
	databaseSource string
	pDB            *sql.DB
	cleanupHooks   []func()

	DB *sql.DB
}

func (ctx *Context) setup() {
	var err error

	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	// prepare the database
	ctx.pDB, err = sql.Open("postgres", fmt.Sprintf(ctx.databaseSource, "postgres"))
	if err != nil {
		panic(err)
	}

	_, err = ctx.pDB.Exec(`create database ` + ctx.database)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err != nil {
			ctx.Teardown()
		}
	}()

	ctx.DB, err = sql.Open("postgres", fmt.Sprintf(ctx.databaseSource, ctx.database))
	if err != nil {
		return
	}

	// prepare schema
	err = schema.Migrate(context.Background(), ctx.DB)
	if err != nil {
		return
	}

	_, err = pgctx.Exec(ctx.Ctx(), `
		select 1;
	`)
	if err != nil {
		panic(err)
	}
}

func (ctx *Context) Teardown() {
	for _, f := range ctx.cleanupHooks {
		f()
	}

	if ctx.DB != nil {
		ctx.DB.Close()
	}

	if ctx.pDB != nil {
		ctx.pDB.Exec(`drop database if exists ` + ctx.database)
		ctx.pDB.Close()
	}
}

func (ctx *Context) Ctx() context.Context {
	c := context.Background()
	c = pgctx.NewContext(c, ctx.DB)
	return c
}

// Setup setups test dependencies
func Setup() *Context {
	ctx := &Context{
		database:       fmt.Sprintf("test_%d", randInt()),
		databaseSource: os.Getenv("TEST_DB_URL"),
	}
	if ctx.databaseSource == "" {
		panic("TEST_DB_URL env required")
	}
	ctx.setup()

	return ctx
}

func randInt() int {
	n, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		panic(err)
	}
	return int(n.Int64())
}
