package controllers

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func Foo() error { return nil }
func Bar() error { return nil }
func Baz() error { return nil }

func printDataRace() {
	fmt.Println(`==================
WARNING: DATA RACE
Write at 0x00c001df1dd0 by goroutine 466:
runtime.mapassign_faststr()
/usr/lib/golang/src/internal/runtime/maps/runtime_faststr.go:263 +0x0
github.com/kuadrant/authorino/controllers.(*AuthConfigReconciler).Reconcile.func1()
/usr/src/authorino/controllers/auth_config_controller.go:129 +0x4a

Previous write at 0x00c001df1dd0 by goroutine 176:
runtime.mapassign_faststr()
/usr/lib/golang/src/internal/runtime/maps/runtime_faststr.go:263 +0x0
github.com/kuadrant/authorino/controllers.(*AuthConfigReconciler).Reconcile()
/usr/src/authorino/controllers/auth_config_controller.go:132 +0x272
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Reconcile()
/opt/app-root/src/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/internal/controller/controller.go:119 +0x194
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).reconcileHandler()
/opt/app-root/src/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/internal/controller/controller.go:316 +0x575
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).processNextWorkItem()
/opt/app-root/src/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/internal/controller/controller.go:266 +0x32f
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Start.func2.2()
/opt/app-root/src/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/internal/controller/controller.go:227 +0xae

Goroutine 466 (running) created at:
github.com/kuadrant/authorino/controllers.(*AuthConfigReconciler).Reconcile()
/usr/src/authorino/controllers/auth_config_controller.go:128 +0x252
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Reconcile()
/opt/app-root/src/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/internal/controller/controller.go:119 +0x194
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).reconcileHandler()
/opt/app-root/src/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/internal/controller/controller.go:316 +0x575
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).processNextWorkItem()
/opt/app-root/src/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/internal/controller/controller.go:266 +0x32f
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Start.func2.2()
/opt/app-root/src/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/internal/controller/controller.go:227 +0xae

Goroutine 176 (running) created at:
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Start.func2()
/opt/app-root/src/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/internal/controller/controller.go:223 +0x7b7
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).Start()
/opt/app-root/src/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/internal/controller/controller.go:234 +0x2a7
sigs.k8s.io/controller-runtime/pkg/manager.(*runnableGroup).reconcile.func1()
/opt/app-root/src/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/manager/runnable_group.go:223 +0x1b5
sigs.k8s.io/controller-runtime/pkg/manager.(*runnableGroup).reconcile.gowrap1()
/opt/app-root/src/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.16.3/pkg/manager/runnable_group.go:226 +0x38
==================`)
}

func Run(ctx context.Context) error {
	err := Foo()
	if err != nil {
		return err
	}

	printDataRace()

	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		err = Baz()
		if err != nil {
			return err
		}

		return nil
	})
	wg.Go(func() error {
		err = Bar()
		if err != nil {
			return err
		}

		return nil
	})

	return wg.Wait()
}
