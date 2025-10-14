package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type VpcID uint64

func NewVpcIdWithTopAndVni(topId, netId uint32) VpcID {
	return VpcID(uint64(topId)<<32 + uint64(netId))
}

func NewVpcId(id uint64) VpcID {
	return VpcID(id)
}

func (v VpcID) MarshalJSON() (buf []byte, err error) {
	return json.Marshal(v.String())
}

func (v *VpcID) UnmarshalJSON(buf []byte) (err error) {
	var vpcId uint64
	vpcString := strings.Trim(string(buf), "\"")
	if err := json.Unmarshal([]byte(vpcString), &vpcId); err != nil {
		return err
	}

	*v = VpcID(vpcId)
	return nil
}

func (v VpcID) IsValid() bool {
	return v.Number() > 0
}

func (v VpcID) String() string {
	return strconv.FormatUint(v.Number(), 10)
}

func (v VpcID) Number() uint64 {
	return uint64(v)
}

func (v VpcID) BridgeName() string {
	return fmt.Sprintf("bridge_%d", v.GetVni())
}

func (v VpcID) VxLanName() string {
	return fmt.Sprintf("vxlan_%d", v.GetVni())
}

func (v VpcID) GetTop() uint32 {
	return uint32(uint64(v) >> 32)
}

func (v VpcID) GetVni() uint32 {
	return uint32(v)
}
