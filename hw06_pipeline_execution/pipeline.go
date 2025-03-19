package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = wrapStage(stage(orDone(in, done)), done)
	}

	return in
}

func wrapStage(stageDataCh Out, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				for range stageDataCh {
					<-stageDataCh
				}
				return
			case val, ok := <-stageDataCh:
				if !ok {
					return
				}
				select {
				case out <- val:
				case <-done:
					for range stageDataCh {
						<-stageDataCh
					}
					return
				}
			}
		}
	}()

	return out
}

func orDone(in Out, done In) Out {
	outCh := make(Bi)

	go func() {
		defer close(outCh)
		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case outCh <- val:
				}
			}
		}
	}()
	return outCh
}
