package goini

import (
	"testing"
	"strings"
	"fmt"
	"sort"
)

var config *IniConfig

func init() {
	config,_ = NewIniConfig("./confini/test.ini")
	config.loadFile()
}

func TestIniConfig_String(t *testing.T) {
	expectConfigStr := `[db]
ip=127.0.0.1
port=3306
charset=utf8mb4
dbname=test
poolsize=20
idleTimeOut=10
idleConn=10
uname=user
upasswd=password

[family]
province=guangdong
city=shenzhen
[.brother]
name=liming
age=87
[..son]
name=lixiaolong
age=64
[...son]
name=limingze
age=37
[...daughter]
name=lixiaoxiao
age=35
[..daughter]
name=lirui
age=62
[.sister]
name=liqiu
age=81
[..son]
name=liumingze
age=59
[...daughter]
name=liumeimei
age=55
[.young_sister]
name=lixiaomei
age=80
[..daughter]
name=zhenxiaolong
age=57
`
	if config.String() != expectConfigStr {
		t.Error("config.String() get wrong result")
	}
}

func TestIniConfig_GetKeyValue(t *testing.T) {
	if v, ok := config.GetKeyValue("db.ip"); ok != true || v != "127.0.0.1" {
		t.Error("config.getKeyValue(db.ip) get wrong result")
	}

	if v, ok := config.GetKeyValue("db.idleTimeOut"); ok != true || v != "10" {
		t.Error("config.getKeyValue(db.idleTimeOut) get wrong result")
	}

	if v, ok := config.GetKeyValue("family.brother.name"); ok != true || v != "liming" {
		t.Error("config.getKeyValue(family.brother.name) get wrong result")
	}

	if v, ok := config.GetKeyValue("family.brother.son.name"); ok != true || v != "lixiaolong" {
		t.Error("config.getKeyValue(family.brother.son.name) get wrong result")
	}

	if v, ok := config.GetKeyValue("family.brother.son.son.age"); ok != true || v != "37" {
		t.Error("config.getKeyValue(family.brother.son.son.age) get wrong result")
	}

	if v, ok := config.GetKeyValue("family.brother.son.daughter.name"); ok != true || v != "lixiaoxiao" {
		t.Error("config.getKeyValue(family.brother.son.daughter.name) get wrong result")
	}

	if v, ok := config.GetKeyValue("family.young_sister.daughter.name"); ok != true || v != "zhenxiaolong" {
		t.Error("config.getKeyValue(family.young_sister.daughter.name) get wrong result")
	}
}



func TestIniConfig_GetBlockKeyValues(t *testing.T) {
	var expectedStr = "charset=utf8mb4&dbname=test&idleConn=10&idleTimeOut=10&ip=127.0.0.1&poolsize=20&port=3306&uname=user&upasswd=password"
	v, ok := config.GetBlockKeyValues("db")
	if !ok || formatMapStr(v) != expectedStr {
		t.Error("config.GetBlockKeyValues(db) get wrong result")
	}
	expectedStr = "brother.age=87&brother.daughter.age=62&brother.daughter.name=lirui&brother.name=liming&brother.son.age=64&brother.son.daughter.age=35&brother.son.daughter.name=lixiaoxiao&brother.son.name=lixiaolong&brother.son.son.age=37&brother.son.son.name=limingze&city=shenzhen&province=guangdong&sister.age=81&sister.name=liqiu&sister.son.age=59&sister.son.daughter.age=55&sister.son.daughter.name=liumeimei&sister.son.name=liumingze&young_sister.age=80&young_sister.daughter.age=57&young_sister.daughter.name=zhenxiaolong&young_sister.name=lixiaomei"
	v, ok = config.GetBlockKeyValues("family")
	if !ok || formatMapStr(v) != expectedStr {
		t.Error("config.GetBlockKeyValues(family) get wrong result")
	}
	expectedStr = "age=64&daughter.age=35&daughter.name=lixiaoxiao&name=lixiaolong&son.age=37&son.name=limingze"
	v, ok = config.GetBlockKeyValues("family.brother.son")
	if !ok || formatMapStr(v) != expectedStr {
		t.Error("config.GetBlockKeyValues(family.brother.son) get wrong result")
	}
	expectedStr = "age=80&daughter.age=57&daughter.name=zhenxiaolong&name=lixiaomei"
	v, ok = config.GetBlockKeyValues("family.young_sister")
	if !ok || formatMapStr(v) != expectedStr {
		t.Error("config.GetBlockKeyValues(family.young_sister) get wrong result")
	}
	expectedStr = "age=57&name=zhenxiaolong"
	v, ok = config.GetBlockKeyValues("family.young_sister.daughter")
	if !ok || formatMapStr(v) != expectedStr {
		t.Error("config.GetBlockKeyValues(family.young_sister.daughter) get wrong result")
	}
	expectedStr = "age=57&name=zhenxiaolong"
	v, ok = config.GetBlockKeyValues("family.young_sister.daughters")
	if ok || formatMapStr(v) == expectedStr {
		t.Error("config.GetBlockKeyValues(family.young_sister.daughters) get wrong result")
	}
}

func formatMapStr(m map[string]string) string {
	var stringArr = make([]string, 0, len(m))
	for k, v := range m {
		stringArr = append(stringArr, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(stringArr)
	return strings.Join(stringArr, "&")
}













