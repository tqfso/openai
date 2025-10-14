package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type BandWidth int64

const (
	MaxBandwidth = 1250000000
)

// 根据带宽字符串创建 Bandwidth
func NewBandwidthFromString(s string) (BandWidth, error) {
	re := regexp.MustCompile(`^([+,-]?\d+)([kKmMgG]?[bB][pP][sS]?)$`)
	matches := re.FindStringSubmatch(s)
	if len(matches) != 3 {
		return 0, errors.New("invalid bandwidth format")
	}

	value, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, err
	}

	unit := strings.ToLower(matches[2])
	switch unit {
	case "kbps":
		value *= 1000
	case "mbps":
		value *= 1000 * 1000
	case "gbps":
		value *= 1000 * 1000 * 1000
	}

	return BandWidth(value / 8), nil
}

// 加
func (b *BandWidth) Add(data BandWidth) {
	*b += data
}

// 减
func (b *BandWidth) Sub(data BandWidth) {
	*b -= data
}

// String 将 Bandwidth 转换为字符串表示
func (b BandWidth) String() string {
	var prefix string
	value := int64(b * 8)

	unit := "bps"
	gbpsNum := int64(1000 * 1000 * 1000)
	mbpsNum := int64(1000 * 1000)
	absValue := int64(math.Abs(float64(value)))
	if absValue >= gbpsNum && absValue%gbpsNum == 0 {
		value /= gbpsNum
		unit = "Gbps"
	} else if absValue >= mbpsNum && absValue%mbpsNum == 0 {
		value /= mbpsNum
		unit = "Mbps"
	} else if absValue >= 1000 {
		value /= 1000
		unit = "kbps"
	}

	return fmt.Sprintf("%s%d%s", prefix, value, unit)
}

// MarshalJSON 实现 JSON 序列化
func (b BandWidth) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

// UnmarshalJSON 实现 JSON 反序列化
func (b *BandWidth) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	bandwidth, err := NewBandwidthFromString(s)
	if err != nil {
		return err
	}

	*b = bandwidth
	return nil
}
