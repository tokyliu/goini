/**
  * load config from file
  * @copyright   tokyliu
 */
package goini

import (
	"os"
	"bufio"
	"io"
	"fmt"
	"strings"
)

type kvitem struct {
	key string
	value string
}

type block struct {
	name string
	items []kvitem
	childBlocks []*block
}

func (b block)formatString(level int) string {
	var strBuilder strings.Builder
	strBuilder.WriteByte('[')
	strBuilder.WriteString(strings.Repeat(".", level))
	strBuilder.WriteString(b.name)
	strBuilder.WriteByte(']')
	strBuilder.WriteByte('\n')
	for _, item := range b.items {
		strBuilder.WriteString(item.key)
		strBuilder.WriteByte('=')
		strBuilder.WriteString(item.value)
		strBuilder.WriteByte('\n')
	}
	for _, cb := range b.childBlocks {
		strBuilder.WriteString(cb.formatString(level+1))
	}
	return strBuilder.String()
}

type IniConfig struct {
	filePath string
	rootBlocks []*block
}


func NewIniConfig(filePath string) (config *IniConfig, err error) {
	if _, err = os.Stat(filePath); err != nil {
		return
	}
	config = &IniConfig{
		filePath: filePath,
	}
	return
}

func (c IniConfig)String() string{
	if c.rootBlocks == nil || len(c.rootBlocks) == 0 {
		return ""
	}
	var strBuilder strings.Builder
	for _, b := range c.rootBlocks {
		strBuilder.WriteString(b.formatString(0))
		strBuilder.WriteByte('\n')
	}
	return strBuilder.String()[:strBuilder.Len()-1]
}


func (c *IniConfig)loadFile() error {
	if c.rootBlocks != nil {
		return nil
	}

	c.rootBlocks = make([]*block, 0, 4)
	fi, err := os.Open(c.filePath);
	if err != nil {
		return err
	}
	defer fi.Close()

	br := bufio.NewReader(fi)

	var upLevelBlocks = make([]*block, 0, 4)
	var blankLineCnt int
	for {
		bytes, _, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read config line error:", err)
			continue
		}

		line := strings.TrimSpace(string(bytes))
		switch {
		case len(line) == 0 :
			blankLineCnt ++
			//since cointinuous blank line more than 50, stop to parse the config file
			if blankLineCnt > 50 {
				break
			}
		case line[0] == '#':
			blankLineCnt = 0
		case line[0] == '[' && line[len(line)-1] == ']':
			blankLineCnt = 0
			blockItem := new(block)
			var level int
			for level=0; level<len(line)-1; level++ {
				if line[level+1] != '.' {
					break
				}
			}
			blockItem.name = line[level+1:len(line)-1]
			if level > len(upLevelBlocks) {
				 return fmt.Errorf("wrong level block name set: %s", line)
			}

			if level == 0 {
				c.rootBlocks = append(c.rootBlocks, blockItem)
				upLevelBlocks = []*block{blockItem}
			}else{
				upLevelBlocks = upLevelBlocks[:level]
				fatherBlock := upLevelBlocks[len(upLevelBlocks)-1]
				if fatherBlock.childBlocks == nil {
					fatherBlock.childBlocks = make([]*block, 0, 4)
				}
				fatherBlock.childBlocks = append(fatherBlock.childBlocks, blockItem)
				upLevelBlocks = append(upLevelBlocks, blockItem)
			}
		default:
			blankLineCnt = 0
			if  len(upLevelBlocks) == 0 {
				fmt.Println("[warning]config item content {", line, "} no block belongs to, promise no blank line with the block belongs to")
				continue
			}
			lineParseArr := strings.Split(line, "=")
			if len(lineParseArr) != 2 {
				fmt.Println("invalid item content {", line, "}, promise the format {key=value}")
				continue
			}
			belongBlock := upLevelBlocks[len(upLevelBlocks)-1]
			if belongBlock.items == nil {
				belongBlock.items = make([]kvitem, 0, 4)
			}
			item := kvitem{
				key: strings.TrimSpace(lineParseArr[0]),
				value: strings.TrimSpace(lineParseArr[1]),
			}
			belongBlock.items = append(belongBlock.items, item)
		}
	}

	return nil
}


func (c *IniConfig)GetKeyValue(keyName string) (string, bool) {
	c.loadFile()
	keyName = strings.TrimSpace(keyName)
	if keyName == "" {
		return "",false
	}

	keyParseArr := strings.Split(keyName, ".")
	if len(keyParseArr) < 2 {
		return "",false
	}

	var blockItem *block
	searchBlocks := c.rootBlocks
	for i:=0; i<len(keyParseArr)-1; i++ {
		if searchBlocks == nil || len(searchBlocks) == 0 {
			return "",false
		}
		blockItem = nil
		for _, b := range searchBlocks {
			if b.name == keyParseArr[i] {
				blockItem = b
				searchBlocks = b.childBlocks
				break
			}
		}
		if blockItem == nil {
			return "", false
		}
	}
	for _, item := range blockItem.items {
		if item.key == keyParseArr[len(keyParseArr)-1] {
			return item.value, true
		}
	}

	return "",false
}


func (c *IniConfig)GetBlockKeyValues(keyName string) (map[string]string, bool) {
	c.loadFile()
	if keyName == "" {
		return nil, false
	}

	keyParseArr := strings.Split(keyName, ".")

	var blockItem *block
	searchBlocks := c.rootBlocks

	for _, k := range keyParseArr {
		blockItem = nil
		for _, b := range searchBlocks {
			if b.name == k {
				blockItem = b
				searchBlocks = b.childBlocks
				break
			}
		}
		if blockItem == nil {
			return nil, false
		}
	}

	blockAllKvItems := loadBlockAllItems(blockItem, "")
	res := make(map[string]string, len(blockAllKvItems)+1)
	for _, kv := range blockAllKvItems {
		res[kv.key] = kv.value
	}

	return res, true
}


func loadBlockAllItems(b *block, prefix string) []kvitem {
	result := make([]kvitem, 0, (len(b.items)+1) * (len(b.childBlocks)+1))
	if b.items != nil && len(b.items) > 0 {
		for _, kv := range b.items {
			if prefix != "" {
				result = append(result, kvitem{
					key: strings.Join([]string{prefix, kv.key}, "."),
					value: kv.value,
				})
			}else{
				result = append(result, kv)
			}
		}
	}

	if b.childBlocks != nil && len(b.childBlocks) > 0 {
		for _, cb := range b.childBlocks {
			if prefix == "" {
				result = append(result, loadBlockAllItems(cb, cb.name)...)
			}else{
				result = append(result, loadBlockAllItems(cb, strings.Join([]string{prefix, cb.name}, "."))...)
			}
		}
	}
	return result
}
























