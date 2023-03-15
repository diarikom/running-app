package api

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
)

type SnowflakeGen struct {
	node *snowflake.Node
}

func (g *SnowflakeGen) New() string {
	return g.node.Generate().String()
}

func (g *SnowflakeGen) NewInt64() int64 {
	return g.node.Generate().Int64()
}

func NewSnowflakeGenerator(instanceNo int64) *SnowflakeGen {
	node, err := snowflake.NewNode(instanceNo)
	if err != nil {
		panic(fmt.Errorf("running-app-api: unable to init snowflake node (%s)", err))
	}

	gen := SnowflakeGen{node: node}
	return &gen
}
