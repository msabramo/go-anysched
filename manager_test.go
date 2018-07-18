package hyperion_test

import (
	"git.corp.adobe.com/abramowi/hyperion"
	"git.corp.adobe.com/abramowi/hyperion/managers/marathon"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("manager.go", func() {
	Context("a manager constructor function is registered under a type", func() {
		var (
			myNewManagerFuncCalled bool
			receivedManagerAddress string
			myNewManagerFunc       func(managerAddress string) (hyperion.Manager, error)
		)

		BeforeEach(func() {
			myNewManagerFuncCalled = false
			receivedManagerAddress = ""

			myNewManagerFunc = func(managerAddress string) (hyperion.Manager, error) {
				myNewManagerFuncCalled = true
				receivedManagerAddress = managerAddress
				return marathon.NewManager(managerAddress)
			}

			hyperion.ClearManagerTypeRegistry()
			hyperion.RegisterManagerType("foo", myNewManagerFunc)
		})

		Describe("RegisterManagerType", func() {
			It("panics if trying to register an already registered manager type", func() {
				Expect(func() { hyperion.RegisterManagerType("foo", myNewManagerFunc) }).To(Panic())
			})
		})

		Describe("NewManager", func() {
			It("calls that function and returns a non-nil manager if that type is passed in", func() {
				managerConfig := hyperion.ManagerConfig{Type: "foo", Address: "http://1.2.3.4:5678"}
				manager, err := hyperion.NewManager(managerConfig)
				Expect(myNewManagerFuncCalled).To(BeTrue())
				Expect(receivedManagerAddress).To(Equal("http://1.2.3.4:5678"))
				Expect(err).ToNot(HaveOccurred())
				Expect(manager).ToNot(BeNil())
			})

			It("returns an error if an unknown type is passed in", func() {
				managerConfig := hyperion.ManagerConfig{Type: "unknown_type", Address: "http://1.2.3.4:5678"}
				manager, err := hyperion.NewManager(managerConfig)
				Expect(myNewManagerFuncCalled).To(BeFalse())
				Expect(receivedManagerAddress).To(Equal(""))
				Expect(err).To(HaveOccurred())
				Expect(manager).To(BeNil())
			})
		})
	})
})
