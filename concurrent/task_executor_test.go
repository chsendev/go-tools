package concurrent

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	// 消费函数
	//var consumerFunc = func(ctx context.Context, task string) (string, error) {
	//	log.Println("start")
	//
	//	time.Sleep(time.Second * 5)
	//	log.Println("end")
	//	return "ok", nil
	//}
	//
	//taskList := []string{"1", "2", "3"}
	//
	//results, err := NewTaskExecutor(context.Background(), taskList, consumerFunc).
	//	Run().
	//	Wait()
	//
	//fmt.Println(results)
	//fmt.Println(err)
}

func TestDistinctFunc(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	persons := []*Person{
		{Name: "n1", Age: 1},
		//{Name: "n1", Age: 1},
		//{Name: "n1", Age: 1},
		{Name: "n2", Age: 2},
		{Name: "n3", Age: 3},
		{Name: "n4", Age: 4},
		{Name: "n5", Age: 5},
	}
	for _, v := range persons {
		fmt.Println(v)
	}

	fmt.Println("-------")
	persons = sliceTool.DistinctFunc(persons, func(a, b *Person) bool {
		return a.Name == b.Name && a.Age == b.Age
	})
	for _, v := range persons {
		fmt.Println(v)
	}

	var bb []int = sliceTool.Map(persons, func(a *Person) int {
		return a.Age
	})
	fmt.Println(bb)

}
