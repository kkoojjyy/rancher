package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type PipelineExecutionLifecycle interface {
	Create(obj *PipelineExecution) (runtime.Object, error)
	Remove(obj *PipelineExecution) (runtime.Object, error)
	Updated(obj *PipelineExecution) (runtime.Object, error)
}

type pipelineExecutionLifecycleAdapter struct {
	lifecycle PipelineExecutionLifecycle
}

func (w *pipelineExecutionLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*PipelineExecution))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *pipelineExecutionLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*PipelineExecution))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *pipelineExecutionLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*PipelineExecution))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewPipelineExecutionLifecycleAdapter(name string, clusterScoped bool, client PipelineExecutionInterface, l PipelineExecutionLifecycle) PipelineExecutionHandlerFunc {
	adapter := &pipelineExecutionLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *PipelineExecution) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
