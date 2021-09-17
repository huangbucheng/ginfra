package config

import (
	"encoding/json"
	"errors"
	"time"
)

//Duration 配置文件中时间字段类型
type Duration struct {
	time.Duration
}

//Marshal 时间字段序列化
func (d Duration) Marshal() ([]byte, error) {
	return json.Marshal(d.String())
}

//Unmarshal 时间字段反序列化
func (d *Duration) Unmarshal(v interface{}) error {
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
