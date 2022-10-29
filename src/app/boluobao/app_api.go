package boluobao

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/request"
	"github.com/VeronicaAlexia/pineapple-backups/struct/sfacg_structs/task"
)

// "https://api.sfacg.com/"
func GET_TASK_LIST() string {
	s := new(task.Task)
	request.Get(new(request.Context).Init("user/tasks/").Query(
		"taskCategory", "1").Query(
		"package", "com.sfacg").Query(
		"deviceToken", "").Query(
		"page", "0").Query(
		"size", "20").
		QueryToString(), s)
	if s != nil && s.Status.HttpCode == 200 {
		for _, task_info := range s.Data {
			fmt.Println("任务号:", task_info.TaskId, "\t\t\t任务名:", task_info.Name)
		}
	} else {
		fmt.Println("task error:", s.Status.Msg)
	}
	return ""
}
