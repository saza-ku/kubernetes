package benchmark

import "k8s.io/kubernetes/test/utils/ktesting"

type Option func(ktesting.TContext)

type PrepareFn func(ktesting.TContext) error

func WithPrepareFn(prepareFn PrepareFn) Option {
	return func(tCtx ktesting.TContext) {
		if err := prepareFn(tCtx); err != nil {
			tCtx.Errorf("failed to run prepareFn: %v", err)
		}
	}
}
