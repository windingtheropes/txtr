package main

import (
	"fmt"
	"iter"
	"os"
	"strings"
)

type keyValue struct {
	key string
	value string
}

func get_kvs(lines iter.Seq[string]) []keyValue {
	var kvs []keyValue;
	for line := range lines {
		parts := strings.Split(line, "=")
		if len(parts) < 2 { continue }
		kvs = append(kvs, keyValue{
			key: strings.Trim(parts[0], `\n`),
			value: strings.TrimSpace(strings.Join(parts[1:], " ")),
		})
	}
	return kvs
}
func sub_kvs(inbytes []byte, kvs []keyValue) ([]byte) {
	b := inbytes;
	for i := range kvs {
		kv := kvs[i]
		b = []byte(
			strings.Replace(
				string(b), 
				fmt.Sprintf("${%v}", kv.key), 
				kv.value, 
				-1,
			),
		)
	}
	return b
}

func Kv_Run(opt Option, c *Command) (error) {
	if len(opt.args) == 0 {
		return fmt.Errorf("kvinput has no value")
	}
	
	kv_bytes, err := os.ReadFile(opt.args[0])
	if err != nil {
		return err
	}
	
	c.bytes = sub_kvs(c.bytes, get_kvs(strings.Lines(string(kv_bytes))))
	return nil
}