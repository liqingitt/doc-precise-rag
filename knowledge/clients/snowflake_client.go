package clients

import (
	"github.com/bwmarrin/snowflake"
)

var SnowflakeClient *snowflake.Node

func init() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	SnowflakeClient = node
}
