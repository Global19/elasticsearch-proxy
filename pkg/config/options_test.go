package config_test

import (
	"net/url"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/openshift/elasticsearch-proxy/pkg/config"
)

func errorMessage(msgs ...string) string {
	result := make([]string, 0)
	result = append(result, "Invalid configuration:")
	result = append(result, msgs...)
	return strings.Join(result, "\n  ")
}

var _ = Describe("Initializing Config options", func() {

	Describe("when defining tls-client-ca without key or certs", func() {
		It("should fail", func() {
			args := []string{"--tls-client-ca=/foo/bar", "--metrics-tls-cert=/foo/bar", "--metrics-tls-key=/foo/bar"}
			options, err := config.Init(args)
			Expect(options).Should(BeNil())

			Expect(err.Error()).Should(
				Equal(errorMessage("tls-client-ca requires tls-key-file or tls-cert-file to be set to listen on tls")))
		})
	})

	Describe("when defining tls-client-ca and key without certs", func() {
		It("should fail", func() {
			args := []string{"--tls-client-ca=/foo/bar", "--tls-key=/foo/bar", "--metrics-tls-cert=/foo/bar", "--metrics-tls-key=/foo/bar"}
			options, err := config.Init(args)
			Expect(options).Should(BeNil())
			Expect(err.Error()).Should(
				Equal(errorMessage("tls-client-ca requires tls-key-file or tls-cert-file to be set to listen on tls")))
		})
	})

	Describe("when defining metrics-listening-address", func() {
		It("should fail without metrics-tls-cert", func() {
			args := []string{"--metrics-listening-address=:60001", "--metrics-tls-key=/foo/bar"}
			options, err := config.Init(args)
			Expect(options).Should(BeNil())
			Expect(err.Error()).Should(
				Equal(errorMessage("metrics-listening-address requires metrics-tls-cert and metrics-tls-key to be set")))
		})

		It("should fail without metrics-tls-key", func() {
			args := []string{"--metrics-listening-address=:60001", "--metrics-tls-cert=/foo/bar"}
			options, err := config.Init(args)
			Expect(options).Should(BeNil())
			Expect(err.Error()).Should(
				Equal(errorMessage("metrics-listening-address requires metrics-tls-cert and metrics-tls-key to be set")))
		})
	})

	Describe("when defining no options", func() {
		It("should not fail", func() {
			args := []string{}
			options, err := config.Init(args)
			Expect(err).Should(BeNil())
			Expect(options).Should(Not(BeNil()))
			Expect(&url.URL{Scheme: "https", Host: "localhost:9200", Path: "/"}).Should(Equal(options.ElasticsearchURL))
		})
	})

	Describe("when defining the admin role", func() {
		It("should succeed", func() {
			args := []string{"--auth-admin-role=foo"}
			options, err := config.Init(args)
			Expect(err).Should(BeNil())
			Expect(options).Should(Not(BeNil()))
			Expect(options.AuthAdminRole).Should(Equal("foo"))
		})
	})

	Describe("when defining the default role", func() {
		It("should succeed", func() {
			args := []string{"--auth-default-role=foo"}
			options, err := config.Init(args)
			Expect(err).Should(BeNil())
			Expect(options).Should(Not(BeNil()))
			Expect(options.AuthDefaultRole).Should(Equal("foo"))
		})

	})

	Describe("when defining whitelisted names", func() {
		It("should succeed", func() {
			args := []string{"--auth-whitelisted-name=foo", "--auth-whitelisted-name=bar"}
			options, err := config.Init(args)
			Expect(err).Should(BeNil())
			Expect(options).Should(Not(BeNil()))
			Expect(options.AuthWhiteListedNames).Should(Equal([]string{"foo", "bar"}))
		})

	})
	Describe("when defining auth backend role", func() {
		Describe("without a valid backendname", func() {

			It("should fail", func() {
				args := []string{"--auth-backend-role={'verb':'get'}"}
				options, err := config.Init(args)
				Expect(options).Should(BeNil())
				Expect(err.Error()).Should(
					Equal(errorMessage("auth-backend-role \"{'verb':'get'}\" should be name=SAR")))
			})
		})
		Describe("that is the same as one that exists", func() {

			It("should fail", func() {
				args := []string{"--auth-backend-role=foo={\"verb\":\"get\"}", "--auth-backend-role=foo={\"verb\":\"get\"}"}
				options, err := config.Init(args)
				Expect(options).Should(BeNil())
				Expect(err.Error()).Should(
					Equal(errorMessage("Backend role with that name \"foo={\\\"verb\\\":\\\"get\\\"}\" already exists")))
			})
		})
		Describe("with unique backend roles", func() {

			It("should succeed", func() {
				args := []string{"--auth-backend-role=foo={\"verb\":\"get\"}", "--auth-backend-role=bar={\"verb\":\"get\"}"}
				options, err := config.Init(args)
				Expect(err).Should(BeNil())
				Expect(options).Should(Not(BeNil()))
				exp := map[string]config.BackendRoleConfig{
					"foo": config.BackendRoleConfig{Verb: "get"},
					"bar": config.BackendRoleConfig{Verb: "get"},
				}
				Expect(options.AuthBackEndRoles).Should(Equal(exp))
			})
		})
	})
})
