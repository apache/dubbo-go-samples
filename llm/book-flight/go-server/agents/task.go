/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package agents

import (
	"context"
	"strings"

	"github.com/apache/dubbo-go-samples/llm/book-flight/go-server/tools"
)

type TaskState int

const (
	TaskUndefined TaskState = 0
	TaskSubmitted TaskState = 1 << iota
	TaskWorking
	TaskInputRequired
	TaskCompleted
	TaskFailed
	TaskCanceled
	TaskUnrelated
)

// InitTaskState
func InitTaskState(value string) TaskState {
	rst := TaskUndefined
	switch strings.ToUpper(value) {
	case "TASKUNDEFINED":
		rst = TaskUndefined
	case "TASKSUBMITTED":
		rst = TaskSubmitted
	case "TASKWORKING":
		rst = TaskWorking
	case "TASKINPUTREQUIRED":
		rst = TaskInputRequired
	case "TASKCOMPLETED":
		rst = TaskCompleted
	case "TASKFAILED":
		rst = TaskFailed
	case "TASKCANCELED":
		rst = TaskCanceled
	case "TASKUNRELATED":
		rst = TaskUnrelated
	default:
		rst = TaskUndefined
	}
	return rst
}

// GetTaskStateTools
func GetTaskStateTools(taskFlag TaskState) []tools.Tool {
	var toolsTask = []tools.Tool{NewTaskUndefinedTool("TaskUnrelated", "不相关问题占位符工具")}
	if taskFlag&TaskSubmitted != 0 {
		toolsTask = append(toolsTask, NewTaskSubmittedTool("TaskSubmitted", "任务已提交"))
	}
	if taskFlag&TaskWorking != 0 {
		toolsTask = append(toolsTask, NewTaskWorkingTool("TaskWorking", "任务正在处理中"))
	}
	if taskFlag&TaskInputRequired != 0 {
		toolsTask = append(toolsTask, NewTaskInputRequiredTool("TaskInputRequired", "任务需要更多信息输入"))
	}
	if taskFlag&TaskCompleted != 0 {
		toolsTask = append(toolsTask, NewTaskCompletedTool("TaskCompleted", "任务已完成"))
	}
	if taskFlag&TaskFailed != 0 {
		toolsTask = append(toolsTask, NewTaskFailedTool("TaskFailed", "任务执行失败"))
	}
	if taskFlag&TaskCanceled != 0 {
		toolsTask = append(toolsTask, NewTaskCanceledTool("TaskCanceled", "任务已取消"))
	}
	if taskFlag&TaskUnrelated != 0 {
		toolsTask = append(toolsTask, NewTaskUnrelatedTool("TaskUnrelated", "不相关任务"))
	}

	return toolsTask
}

func InterruptTask(taskFlag TaskState) bool {
	return taskFlag&(TaskInputRequired|TaskCompleted|TaskUnrelated) != 0
}

/*
TaskUndefinedTool
*/
type TaskUndefinedTool struct {
	tools.BaseTool
}

func NewTaskUndefinedTool(name string, description string) TaskUndefinedTool {
	return TaskUndefinedTool{
		tools.NewBaseTool(
			name, description, tools.GetStructKeys(nil), "", "", ""),
	}
}

func (ptt TaskUndefinedTool) Call(ctx context.Context, input string) (string, error) {
	return "TaskUndefined", nil
}

/*
TaskSubmittedTool
*/
type TaskSubmittedTool struct {
	tools.BaseTool
}

func NewTaskSubmittedTool(name string, description string) TaskSubmittedTool {
	return TaskSubmittedTool{
		tools.NewBaseTool(
			name, description, tools.GetStructKeys(nil), "", "", ""),
	}
}

func (ptt TaskSubmittedTool) Call(ctx context.Context, input string) (string, error) {
	return "TaskSubmitted", nil
}

/*
TaskWorkingTool
*/
type TaskWorkingTool struct {
	tools.BaseTool
}

func NewTaskWorkingTool(name string, description string) TaskWorkingTool {
	return TaskWorkingTool{
		tools.NewBaseTool(
			name, description, tools.GetStructKeys(nil), "", "", ""),
	}
}

func (ptt TaskWorkingTool) Call(ctx context.Context, input string) (string, error) {
	return "TaskWorking", nil
}

/*
TaskInputRequiredTool
*/
type TaskInputRequiredTool struct {
	tools.BaseTool
}

type taskInputRequiredData struct {
	MissingInfo string `json:"missing_info" validate:"required"`
}

func NewTaskInputRequiredTool(name string, description string) TaskInputRequiredTool {
	return TaskInputRequiredTool{
		tools.NewBaseTool(
			name, description, tools.GetStructKeys(taskInputRequiredData{}), "", "", ""),
	}
}

func (mi TaskInputRequiredTool) Call(ctx context.Context, input string) (string, error) {
	return input, nil
}

/*
TaskCompletedTool
*/
type TaskCompletedTool struct {
	tools.BaseTool
}

func NewTaskCompletedTool(name string, description string) TaskCompletedTool {
	return TaskCompletedTool{
		tools.NewBaseTool(
			name, description, tools.GetStructKeys(nil), "", "", ""),
	}
}

func (ptt TaskCompletedTool) Call(ctx context.Context, input string) (string, error) {
	return "TaskCompleted", nil
}

/*
TaskFailedTool
*/
type TaskFailedTool struct {
	tools.BaseTool
}

func NewTaskFailedTool(name string, description string) TaskFailedTool {
	return TaskFailedTool{
		tools.NewBaseTool(
			name, description, tools.GetStructKeys(nil), "", "", ""),
	}
}

func (ptt TaskFailedTool) Call(ctx context.Context, input string) (string, error) {
	return "TaskFailed", nil
}

/*
TaskCanceledTool
*/
type TaskCanceledTool struct {
	tools.BaseTool
}

func NewTaskCanceledTool(name string, description string) TaskCanceledTool {
	return TaskCanceledTool{
		tools.NewBaseTool(
			name, description, tools.GetStructKeys(nil), "", "", ""),
	}
}

func (ptt TaskCanceledTool) Call(ctx context.Context, input string) (string, error) {
	return "TaskCanceled", nil
}

/*
TaskUnrelatedTool
*/
type TaskUnrelatedTool struct {
	tools.BaseTool
}

func NewTaskUnrelatedTool(name string, description string) TaskUnrelatedTool {
	return TaskUnrelatedTool{
		tools.NewBaseTool(
			name, description, tools.GetStructKeys(nil), "", "", ""),
	}
}

func (ptt TaskUnrelatedTool) Call(ctx context.Context, input string) (string, error) {
	return "TaskUnrelated", nil
}
