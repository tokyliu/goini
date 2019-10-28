# goini
> golang parse ini config file with multistage block configuration 


#### Features:
* support mutil level block configuration, for example: 
```
[level_one]
x1 = y1
x2 = y2
  [.level_two]
  m1 = n1
  m2 = n2
    [..level_three]
       s1 = t1
       s2 = t2
```
* configuration comment support, and continuously blank line more than 50 parse break out support
* support get all the <key,value> of a block with its child blocks content.
* load the config file contents to memory once


#### How to apply:

* install the code to your application:

`go get github.com/tokyliu/goini`

* init the config object:

`config, err := NewIniConfig(filePath)`

* get the config structed mirror to judge the config strcuted object is match your config file:

`configStr := config.String()`

* get the config item value with the key:

`value, exist := config.GetKeyValue(keyName)`

* get all the <key,value> pair of a block:

`kvMap,exist := config.GetBlockKeyValues(blockName)`

#### Examples:

* sample ini config file content as below:
````
[family]
province = guangdong
city = shenzhen
[.brother]
name = liming
age = 87
[..son]
name = lixiaolong
age = 64
[...son]
name = limingze
age = 37
[...daughter]
name = lixiaoxiao
age = 35
[..daughter]
name = lirui
age = 62
[.sister]
name = liqiu
age = 81
[..son]
name = liumingze
age = 59
[...daughter]
name = liumeimei
age = 55
[.young_sister]
name = lixiaomei
age = 80
[..daughter]
name = zhenxiaolong
age = 57
````

* use the `GetKeyValue()` to get the value of keyName:

```aidl
value, exist := config.GetKeyValue("family.brother.son.name") // value=lixiaolong
value, exist := config.GetKeyValue("family.brother.son.son.age") // value = 37
value, exist := config.GetKeyValue("family.young_sister.daughter.name") // value = zhenxiaolong
```


* use the `GetBlockKeyValues` to get all the <key,value> pair:

```aidl
kvMap, exist := config.GetBlockKeyValues(family.young_sister) 
```

















  
