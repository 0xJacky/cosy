package cosy

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
)

type Ctx[T any] struct {
	*gin.Context
	ID                       uint64
	rules                    gin.H
	Payload                  map[string]interface{}
	Model                    T
	OriginModel              T
	BatchUpdateEffectedIDs   []uint64
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
	joins                    []string
	scan                     func(tx *gorm.DB) any
	transformer              func(*T) any
	permanentlyDelete        bool
	selectedFields           map[string]bool
	itemKey                  string
	columnWhiteList          map[string]bool
	in                       []string
	eq                       []string
	fussy                    []string
	orIn                     []string
	orEq                     []string
	orFussy                  []string
	search                   []string
	between                  []string
	unique                   []string
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
		selectedFields:           make(map[string]bool),
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

	for k, rule := range rules {
		rule := cast.ToString(rule)
		if strings.Contains(rule, "db_unique") {
			c.unique = append(c.unique, k)
			rules[k] = strings.ReplaceAll(rule, "db_unique", "")
			rules[k] = strings.TrimRight(rule, ",")
		}
	}

	return c
}

func (c *Ctx[T]) SetUnique(keys ...string) {
	c.unique = append(c.unique, keys...)
}

func (c *Ctx[T]) SetPreloads(args ...string) *Ctx[T] {
	c.preloads = append(c.preloads, args...)
	return c
}

func (c *Ctx[T]) SetJoins(args ...string) *Ctx[T] {
	c.joins = append(c.joins, args...)
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

func (c *Ctx[T]) GetParamID() uint64 {
	return cast.ToUint64(c.Param("id"))
}

func (c *Ctx[T]) AddColWhiteList(cols ...string) *Ctx[T] {
	for _, col := range cols {
		c.columnWhiteList[col] = true
	}
	return c
}

func (c *Ctx[T]) AddSelectedFields(fields ...string) *Ctx[T] {
	for _, field := range fields {
		c.selectedFields[field] = true
	}
	return c
}

func (c *Ctx[T]) GetSelectedFields() []string {
	var fields []string
	for field := range c.selectedFields {
		fields = append(fields, field)
	}
	return fields
}
