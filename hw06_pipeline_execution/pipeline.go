package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = stage(wrapStage(in, done))
	}

	return wrapStage(in, done)
}

func wrapStage(stageDataCh Out, done In) Out {
	out := make(Bi)

	clearAllData := func() {
		for v := range stageDataCh {
			_ = v
		}
	}
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				go clearAllData()
				return
			case val, ok := <-stageDataCh:
				if !ok {
					return
				}
				select {
				case out <- val:
				case <-done:
					go clearAllData()
					return
				}
			}
		}
	}()

	return out
}
