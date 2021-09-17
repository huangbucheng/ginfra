package utils

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
	"github.com/sony/sonyflake"
)

//NewSnowFlakeNode snowflake node
func NewSnowFlakeNode(nodeidx int64) (*snowflake.Node, error) {
	// Create a new Node with a Node number of ipaddr
	node, err := snowflake.NewNode(nodeidx)
	if err != nil {
		return nil, err
	}
	return node, nil
}

//NewSonyFlake new sony flake instance
func NewSonyFlake(machineId uint16) *sonyflake.Sonyflake {
	var st sonyflake.Settings
	if machineId > 0 {
		st.MachineID = func() (uint16, error) { return machineId, nil }
	}
	return sonyflake.NewSonyflake(st)
}

func main() {

	// Create a new Node with a Node number of 1
	node, err := NewSnowFlakeNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate a snowflake ID.
	id := node.Generate()

	// Print out the ID in a few different ways.
	fmt.Printf("Int64  ID: %d\n", id)
	fmt.Printf("String ID: %s\n", id)
	fmt.Printf("Base2  ID: %s\n", id.Base2())
	fmt.Printf("Base64 ID: %s\n", id.Base64())

	// Print out the ID's timestamp
	fmt.Printf("ID Time  : %d\n", id.Time())

	// Print out the ID's node number
	fmt.Printf("ID Node  : %d\n", id.Node())

	// Print out the ID's sequence number
	fmt.Printf("ID Step  : %d\n", id.Step())

	// Generate and print, all in one.
	fmt.Printf("ID       : %d\n", node.Generate().Int64())

	sf := NewSonyFlake(1)
	sfid, err := sf.NextID()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Int64  ID: %d\n", sfid)
}
