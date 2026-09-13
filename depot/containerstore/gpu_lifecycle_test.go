package containerstore_test

import (
	"encoding/json"
	"time"

	"code.cloudfoundry.org/clock/fakeclock"
	mfakes "code.cloudfoundry.org/diego-logging-client/testhelpers"
	"code.cloudfoundry.org/executor"
	"code.cloudfoundry.org/executor/depot/containerstore"
	"code.cloudfoundry.org/executor/depot/containerstore/containerstorefakes"
	eventfakes "code.cloudfoundry.org/executor/depot/event/fakes"
	"code.cloudfoundry.org/executor/depot/transformer/faketransformer"
	"code.cloudfoundry.org/executor/initializer/configuration/configurationfakes"
	"code.cloudfoundry.org/garden/gardenfakes"
	"code.cloudfoundry.org/lager/v3"
	"code.cloudfoundry.org/lager/v3/lagertest"
	"code.cloudfoundry.org/volman/volmanfakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GPU container lifecycle", func() {
	var (
		logger              lager.Logger
		gpuManager          *executor.GPUManager
		containerStore      containerstore.ContainerStore
		gardenClient        *gardenfakes.FakeClient
		gardenClientFactory *containerstorefakes.FakeGardenClientFactory
		gardenContainer     *gardenfakes.FakeContainer
		totalCapacity       executor.ExecutorResources
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("gpu-lifecycle-test")

		gpuManager = &executor.GPUManager{
			GPUs: []executor.GPUInfo{{Index: 0, UUID: "GPU-0", Name: "Tesla V100", MemoryMiB: 16160}},
		}

		gardenContainer = new(gardenfakes.FakeContainer)
		gardenContainer.HandleReturns("gpu-container")
		gardenClient = new(gardenfakes.FakeClient)
		gardenClient.CreateReturns(gardenContainer, nil)

		gardenClientFactory = new(containerstorefakes.FakeGardenClientFactory)
		gardenClientFactory.NewGardenClientReturns(gardenClient)

		totalCapacity = executor.NewExecutorResources(1024, 1024, 4)
		totalCapacity.GPUTotal = 1
		totalCapacity.GPUType = "nvidia"

		containerStore = containerstore.New(
			containerstore.ContainerConfig{OwnerName: "test-executor"},
			&totalCapacity,
			gardenClientFactory,
			new(containerstorefakes.FakeDependencyManager),
			new(volmanfakes.FakeManager),
			new(containerstorefakes.FakeCredManager),
			new(containerstorefakes.FakeLogManager),
			fakeclock.NewFakeClock(time.Now()),
			new(eventfakes.FakeHub),
			new(faketransformer.FakeTransformer),
			"",
			new(mfakes.FakeIngressClient),
			new(configurationfakes.FakeRootFSSizer),
			"",
			containerstore.NewNoopProxyConfigHandler(),
			"cell-id",
			false,
			false,
			new(containerstorefakes.FakeVolumeMountedFilesImplementor),
			json.Marshal,
			gpuManager,
		)
	})

	It("allocates a GPU and passes it to Garden as a CDI device on create", func() {
		req := executor.NewAllocationRequest("gpu-guid", &executor.Resource{
			MemoryMB: 128, DiskMB: 512, GPULimit: 1, GPUType: "nvidia",
		}, false, executor.Tags{})

		_, err := containerStore.Reserve(logger, "trace-id", &req)
		Expect(err).NotTo(HaveOccurred())

		runReq := executor.NewRunRequest("gpu-guid", &executor.RunInfo{}, executor.Tags{})
		Expect(containerStore.Initialize(logger, &runReq)).To(Succeed())

		_, err = containerStore.Create(logger, "trace-id", "gpu-guid")
		Expect(err).NotTo(HaveOccurred())

		Expect(gardenClient.CreateCallCount()).To(Equal(1))
		spec := gardenClient.CreateArgsForCall(0)
		Expect(spec.CDIDevices).To(Equal([]string{"nvidia.com/gpu=0"}))
		Expect(spec.Env).To(ContainElement("CUDA_VISIBLE_DEVICES=0"))
		Expect(gpuManager.Available()).To(Equal(0))
	})

	It("releases the GPU on destroy", func() {
		req := executor.NewAllocationRequest("gpu-guid", &executor.Resource{
			MemoryMB: 128, DiskMB: 512, GPULimit: 1, GPUType: "nvidia",
		}, false, executor.Tags{})
		_, err := containerStore.Reserve(logger, "trace-id", &req)
		Expect(err).NotTo(HaveOccurred())
		runReq := executor.NewRunRequest("gpu-guid", &executor.RunInfo{}, executor.Tags{})
		Expect(containerStore.Initialize(logger, &runReq)).To(Succeed())
		_, err = containerStore.Create(logger, "trace-id", "gpu-guid")
		Expect(err).NotTo(HaveOccurred())
		Expect(gpuManager.Available()).To(Equal(0))

		Expect(containerStore.Destroy(logger, "trace-id", "gpu-guid")).To(Succeed())

		Expect(gpuManager.Available()).To(Equal(1))
	})
})
