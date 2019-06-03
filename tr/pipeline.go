package tr

// Pipeline ...
type Pipeline struct {
	ID   string
	Name string
	Node TaskNode
}

// Pipelines ...
var Pipelines = make(map[string]Pipeline)

// RegisterPipeline ...
var RegisterPipeline = make(map[string]string)

func startPipeline(id string) {
	if pipelineName, ok := RegisterPipeline[id]; ok {
		if pipeline, ok := Pipelines[pipelineName]; ok {
			currentNode := pipeline.Node
			chTasks := make(chan TaskFeedBack)
		LoopNode:
			for {
				println("STARTING TASK ID ", currentNode.TaskID)
				task := GetTask(currentNode.TaskID)
				if task == nil {
					println("ERROR TASK NON TROUVER")
					break LoopNode
				}
				go task.Execute(id, currentNode.Params, chTasks)
			LoopTask:
				for {
					select {
					case feedback := <-chTasks:
						switch feedback.Event {
						case "DONE":
							println("Task done passing to next")
							if len(currentNode.NextNode) == 0 {
								break LoopNode
							} else {
								currentNode = currentNode.NextNode[0]
							}
							break LoopTask
						case "ERROR":
							println("ERROR RUNNING TASK ",currentNode.TaskID, " error : ",feedback.Message.(error).Error())
							break LoopNode
						}
					}
				}
			}
			println("FIN DE LA PIPEPINE ", id)
		} else {
			println("PIPELINE NOT FOUND")
		}
	} else {
		println("REGISTER PIPELINE NOT FOUND ", id)
	}

}
