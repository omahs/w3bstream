package redistestutil_test

import (
	"os"
	"testing"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	. "github.com/machinefi/w3bstream/pkg/depends/testutil/redistestutil"
	. "github.com/onsi/gomega"
)

func TestInit(t *testing.T) {
	os.Setenv(consts.GoRuntimeEnv, consts.DevelopEnv)

	NewWithT(t).Expect(Endpoint.Endpoint.String()).To(Equal("tcp://127.0.0.1:6379"))
	NewWithT(t).Expect(Endpoint.Prefix).To(Equal("dev:test:"))

	NewWithT(t).Expect(Redis.Host).To(Equal("127.0.0.1"))
	NewWithT(t).Expect(Redis.Port).To(Equal(6379))
}

func DISABLE_TestLivenessCheck(t *testing.T) {
	kvs := Redis.LivenessCheck()
	NewWithT(t).Expect(kvs[Redis.Host]).To(Equal("ok"))

	kvs = Endpoint.LivenessCheck()
	NewWithT(t).Expect(kvs[Endpoint.Endpoint.Host()]).To(Equal("ok"))

	op := Redis.Get()
	repl, err := op.Do("set", Redis.Key("test"), 100)
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(repl).To(Equal("OK"))

	repl, err = op.Do("get", Redis.Key("test"))
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(repl).To(Equal([]byte("100")))
}
