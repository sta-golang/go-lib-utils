package example

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/algorithm"
)

/**
场景
cat 这个单词可以转化为 kat
kat 这个单词可以转化为 hat
那么 cat 就可以转化成 hat
现在有n个单词都是可以相互转换的。现在给你 m条转化的规则 让你来判断两个单词是否可以互相转换 时间复杂度小于 O(mlg2m)
 */
func unionFind()  {
	arr := []interface{}{"cat","kat","hat"}
	uf := algorithm.NewUnionFind(arr)
	uf.Union("cat","kat")
	uf.Union("kat","hat")
	fmt.Println(uf.Connected("cat","hat"))
}