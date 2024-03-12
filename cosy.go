package cosy

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type Ctx[T any] struct {
	*gin.Context
	ID                       int
	rules                    gin.H
	Payload                  map[string]interface{}
	Model                    T
	OriginModel              T
	table                    string
	tableArgs                []interface{}
	abort                    bool
	nextHandler              *gin.HandlerFunc
	skipAssociationsOnCreate bool
	beforeDecodeHookFunc     []func(ctx *Ctx[T])
	beforeExecuteHookFunc    []func(ctx *Ctx[T])
	executedHookFunc         []func(ctx *Ctx[T])
	gormScopes               []func(tx *gorm.DB) *gorm.DB
	preloads                 []string
	scan                     func(tx *gorm.DB) any
	transformer              func(*T) any
	permanentlyDelete        bool
	SelectedFields           []string
	itemKey                  string
	columnWhiteList          map[string]bool
	in                       []string
	eq                       []string
	fussy                    []string
	orIn                     []string
	orEq                     []string
	orFussy                  []string
	preload                  []string
	search                   []string
}

func Core[T any](c *gin.Context) *Ctx[T] {
	return &Ctx[T]{
		Context:                  c,
		gormScopes:               make([]func(tx *gorm.DB) *gorm.DB, 0),
		beforeExecuteHookFunc:    make([]func(ctx *Ctx[T]), 0),
		beforeDecodeHookFunc:     make([]func(ctx *Ctx[T]), 0),
		itemKey:                  "id",
		skipAssociationsOnCreate: true,
		columnWhiteList:          make(map[string]bool),
	}
}

func (c *Ctx[T]) SetTable(table string, args ...interface{}) *Ctx[T] {
	c.table = table
	c.tableArgs = args
	return c
}

func (c *Ctx[T]) SetItemKey(key string) *Ctx[T] {
	c.itemKey = key
	return c
}

func (c *Ctx[T]) SetValidRules(rules gin.H) *Ctx[T] {
	c.rules = rules

	return c
}

func (c *Ctx[T]) SetPreloads(args ...string) *Ctx[T] {
	c.preloads = append(c.preloads, args...)
	return c
}

func (c *Ctx[T]) SetScan(scan func(tx *gorm.DB) any) *Ctx[T] {
	c.scan = scan
	return c
}

func (c *Ctx[T]) SetTransformer(t func(m *T) any) *Ctx[T] {
	c.transformer = t
	return c
}

func (c *Ctx[T]) GetParamID() int {
	return cast.ToInt(c.Param("id"))
}

func (c *Ctx[T]) AddColWhiteList(cols ...string) *Ctx[T] {
	for _, col := range cols {
		c.columnWhiteList[col] = true
	}
	return c
}
