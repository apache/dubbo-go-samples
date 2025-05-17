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
)

import (
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/tools"
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
	var rst TaskState
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

func CreateToolkitByVariadic(description string, taskFlag TaskState, ts ...tools.Tool) tools.Tools {
	return CreateToolkit(description, taskFlag, ts)
}

// CreateToolkit
func CreateToolkit(description string, taskFlag TaskState, ts []tools.Tool) tools.Tools {
	var toolsTask []tools.Tool = ts
	if taskFlag&TaskSubmitted != 0 {
		if t, err := tools.CreateTool[TaskSubmittedTool]("TaskSubmitted", "任务已提交", ""); err == nil {
			toolsTask = append(toolsTask, t)
		}
	}
	if taskFlag&TaskWorking != 0 {
		if t, err := tools.CreateTool[TaskWorkingTool]("TaskWorking", "任务正在处理中", ""); err == nil {
			toolsTask = append(toolsTask, t)
		}
	}
	if taskFlag&TaskInputRequired != 0 {
		if t, err := tools.CreateTool[TaskInputRequiredTool]("TaskInputRequired", "任务需要更多信息输入", ""); err == nil {
			toolsTask = append(toolsTask, t)
		}
	}
	if taskFlag&TaskCompleted != 0 {
		if t, err := tools.CreateTool[TaskCompletedTool]("TaskCompleted", "任务已完成", ""); err == nil {
			toolsTask = append(toolsTask, t)
		}
	}
	if taskFlag&TaskFailed != 0 {
		if t, err := tools.CreateTool[TaskFailedTool]("TaskFailed", "任务执行失败", ""); err == nil {
			toolsTask = append(toolsTask, t)
		}
	}
	if taskFlag&TaskCanceled != 0 {
		if t, err := tools.CreateTool[TaskCanceledTool]("TaskCanceled", "任务已取消", ""); err == nil {
			toolsTask = append(toolsTask, t)
		}
	}
	if taskFlag&TaskUnrelated != 0 {
		if t, err := tools.CreateTool[TaskUnrelatedTool]("TaskUnrelated", "不相关任务", ""); err == nil {
			toolsTask = append(toolsTask, t)
		}
	}

	return tools.NewToolkit(toolsTask, description)
}

func InterruptTask(taskFlag TaskState) bool {
	return taskFlag&(TaskInputRequired|TaskCompleted|TaskFailed|TaskCanceled|TaskUnrelated) != 0
}

/*
TaskUndefinedTool
*/
type TaskUndefinedTool struct {
	tools.BaseTool
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

func (ptt TaskSubmittedTool) Call(ctx context.Context, input string) (string, error) {
	return "TaskSubmitted", nil
}

/*
TaskWorkingTool
*/
type TaskWorkingTool struct {
	tools.BaseTool
}

func (ptt TaskWorkingTool) Call(ctx context.Context, input string) (string, error) {
	return "TaskWorking", nil
}

/*
TaskInputRequiredTool
*/
type TaskInputRequiredTool struct {
	tools.BaseTool
	MissingInfo string `json:"missing_info" validate:"required"`
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

func (ptt TaskCompletedTool) Call(ctx context.Context, input string) (string, error) {
	return "TaskCompleted", nil
}

/*
TaskFailedTool
*/
type TaskFailedTool struct {
	tools.BaseTool
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

func (ptt TaskCanceledTool) Call(ctx context.Context, input string) (string, error) {
	return "TaskCanceled", nil
}

/*
TaskUnrelatedTool
*/
type TaskUnrelatedTool struct {
	tools.BaseTool
}

func (ptt TaskUnrelatedTool) Call(ctx context.Context, input string) (string, error) {
	return "TaskUnrelated", nil
}
