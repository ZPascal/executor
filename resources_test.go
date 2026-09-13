package executor_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/executor"
)

var _ = Describe("Container", func() {
	Describe("HasTags", func() {
		var container executor.Container

		Context("when tags are nil", func() {
			BeforeEach(func() {
				container = executor.Container{
					Tags: nil,
				}
			})

			It("returns true if requested tags are nil", func() {
				Expect(container.HasTags(nil)).To(BeTrue())
			})

			It("returns false if requested tags are not nil", func() {
				Expect(container.HasTags(executor.Tags{"a": "b"})).To(BeFalse())
			})
		})

		Context("when tags are not nil", func() {
			BeforeEach(func() {
				container = executor.Container{
					Tags: executor.Tags{"a": "b"},
				}
			})

			It("returns true when found", func() {
				Expect(container.HasTags(executor.Tags{"a": "b"})).To(BeTrue())
			})

			It("returns false when nil", func() {
				Expect(container.HasTags(nil)).To(BeFalse())
			})

			It("returns false when not found", func() {
				Expect(container.HasTags(executor.Tags{"a": "c"})).To(BeFalse())
			})
		})
	})

	Describe("Subtract", func() {
		const (
			defaultDiskMB     = 20
			defaultMemoryMB   = 30
			defaultContainers = 3
		)

		It("returns false when the number of containers is less than 1", func() {
			resources := executor.NewExecutorResources(defaultMemoryMB, defaultDiskMB, 0)
			resourceToSubtract := executor.NewResource(defaultMemoryMB-1, defaultDiskMB-1, -1)
			Expect(resources.Subtract(&resourceToSubtract)).To(BeFalse())
		})

		It("returns false when disk size exceeds total available disk size", func() {
			resources := executor.NewExecutorResources(defaultMemoryMB, 10, defaultContainers)
			resourceToSubtract := executor.NewResource(defaultMemoryMB-1, 20, -1)
			Expect(resources.Subtract(&resourceToSubtract)).To(BeFalse())
		})

		It("returns false when memory exceeds total available memory", func() {
			resources := executor.NewExecutorResources(10, defaultDiskMB, defaultContainers)
			resourceToSubtract := executor.NewResource(20, defaultDiskMB-1, -1)
			Expect(resources.Subtract(&resourceToSubtract)).To(BeFalse())
		})
	})

	Describe("TransitionToComplete", func() {
		var (
			container     *executor.Container
			failureReason string
		)

		BeforeEach(func() {
			container = &executor.Container{}
			failureReason = "foo"
		})

		JustBeforeEach(func() {
			container.TransitionToComplete(true, failureReason, false)
		})

		It("does not truncate short failure reasons", func() {
			Expect(container.RunResult.FailureReason).To(Equal("foo"))
		})

		When("failure reason is longer than 10000 characters", func() {
			BeforeEach(func() {
				failureReason = strings.Repeat("a", 6000) + strings.Repeat("b", 6000)
			})

			It("is truncated", func() {
				Expect(container.RunResult.FailureReason).To(HaveLen(10000))
				Expect(container.RunResult.FailureReason[4991:5008]).To(Equal("\n... (truncated)\n"))
			})
		})
	})
})

var _ = Describe("Resource and ExecutorResources GPU fields", func() {
	It("Resource carries a GPU request without disturbing existing fields", func() {
		r := executor.NewResource(128, 512, 100)
		r.GPULimit = 1
		r.GPUType = "nvidia"

		Expect(r.MemoryMB).To(Equal(128))
		Expect(r.DiskMB).To(Equal(512))
		Expect(r.MaxPids).To(Equal(100))
		Expect(r.GPULimit).To(Equal(1))
		Expect(r.GPUType).To(Equal("nvidia"))
	})

	It("ExecutorResources.Subtract still only accounts for memory/disk/containers, ignoring GPU fields", func() {
		total := executor.NewExecutorResources(1024, 1024, 4)
		total.GPUTotal = 2
		total.GPUFree = 2
		total.GPUType = "nvidia"

		req := executor.Resource{MemoryMB: 100, DiskMB: 100, GPULimit: 1, GPUType: "nvidia"}
		ok := total.Subtract(&req)

		Expect(ok).To(BeTrue())
		Expect(total.MemoryMB).To(Equal(924))
		Expect(total.Containers).To(Equal(3))
		// GPU accounting is NOT done via Subtract - GPUManager tracks it
		// separately (Task 6). GPUTotal/GPUFree are untouched here.
		Expect(total.GPUTotal).To(Equal(2))
		Expect(total.GPUFree).To(Equal(2))
	})
})
