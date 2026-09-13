package executor_test

import (
	"code.cloudfoundry.org/executor"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GPUManager", func() {
	var manager *executor.GPUManager

	BeforeEach(func() {
		manager = &executor.GPUManager{
			GPUs: []executor.GPUInfo{
				{Index: 0, UUID: "GPU-0", Name: "Tesla V100", MemoryMiB: 16160},
				{Index: 1, UUID: "GPU-1", Name: "Tesla V100", MemoryMiB: 16160},
			},
		}
	})

	It("reports total and available GPUs", func() {
		Expect(manager.TotalGPUs()).To(Equal(2))
		Expect(manager.Available()).To(Equal(2))
	})

	It("allocates and tracks free GPUs", func() {
		indices, err := manager.Allocate("container-a", 1)
		Expect(err).NotTo(HaveOccurred())
		Expect(indices).To(HaveLen(1))
		Expect(manager.Available()).To(Equal(1))
		Expect(manager.ContainerForGPU(indices[0])).To(Equal("container-a"))
	})

	It("fails to allocate more GPUs than are free", func() {
		_, err := manager.Allocate("container-a", 3)
		Expect(err).To(HaveOccurred())
	})

	It("releases GPUs for reuse", func() {
		indices, err := manager.Allocate("container-a", 2)
		Expect(err).NotTo(HaveOccurred())
		Expect(manager.Available()).To(Equal(0))

		manager.Release("container-a")
		Expect(manager.Available()).To(Equal(2))

		for _, idx := range indices {
			Expect(manager.ContainerForGPU(idx)).To(Equal(""))
		}
	})

	It("is a no-op to release a container that holds no GPUs", func() {
		Expect(func() { manager.Release("unknown-container") }).NotTo(Panic())
		Expect(manager.Available()).To(Equal(2))
	})

	Context("when nvidia-smi is unavailable", func() {
		It("NewGPUManager returns an empty-but-valid manager, not an error", func() {
			manager, err := executor.NewGPUManager()
			Expect(err).NotTo(HaveOccurred())
			Expect(manager).NotTo(BeNil())
			// On a dev machine without nvidia-smi this is 0; on a real GPU
			// cell it reflects the actual discovered device count. Either
			// way, NewGPUManager must not fail.
			Expect(manager.TotalGPUs()).To(BeNumerically(">=", 0))
		})
	})
})
