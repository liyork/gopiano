package clientgo

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	v2 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

func TestIndex() {
	index := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{"testmodes": testIndexFunc})

	pod1 := &v1.Pod{ObjectMeta: v2.ObjectMeta{Name: "one", Labels: map[string]string{"foo": "bar"}}}
	pod2 := &v1.Pod{ObjectMeta: v2.ObjectMeta{Name: "two", Labels: map[string]string{"foo": "bar"}}}
	pod3 := &v1.Pod{ObjectMeta: v2.ObjectMeta{Name: "tre", Labels: map[string]string{"foo": "biz"}}}

	index.Add(pod1)
	index.Add(pod2)
	index.Add(pod3)

	keys := index.ListIndexFuncValues("testmodes")
	if len(keys) != 2 {
		fmt.Println("Expected 2 keys but got ", len(keys))
	}

	for _, key := range keys {
		if key != "bar" && key != "biz" {
			fmt.Println("Expected only 'bar' or 'biz' but got ", key)
		}
	}
}

func testIndexFunc(obj interface{}) ([]string, error) {
	pod := obj.(*v1.Pod)
	return []string{pod.Labels["foo"]}, nil
}
