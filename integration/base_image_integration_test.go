package integration_test

import (
	. "github.com/onsi/ginkgo/extensions/table"
	types2 "github.com/onsi/gomega/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab", func() {
	DescribeTable("generates base property", func(inputImage string, matchers ...types2.GomegaMatcher) {
		metadataLabel := runDeplabAgainstImage(inputImage)

		Expect(metadataLabel.Base).To(SatisfyAll(matchers...))
	},
		Entry("ubuntu:bionic image", "ubuntu:bionic-20190718",
			HaveKeyWithValue("name", "Ubuntu"),
			HaveKeyWithValue("version", "18.04.2 LTS (Bionic Beaver)"),
			HaveKeyWithValue("version_id", "18.04"),
			HaveKeyWithValue("id_like", "debian"),
			HaveKeyWithValue("version_codename", "bionic"),
			HaveKeyWithValue("pretty_name", "Ubuntu 18.04.2 LTS"),
		),
		Entry("a non-ubuntu:bionic image with /etc/os-release", "alpine:3.10.1",
			HaveKeyWithValue("name", "Alpine Linux"),
			HaveKeyWithValue("version_id", "3.10.1"),
		),
		Entry("an image that doesn't have an os-release", "pivotalnavcon/test-asset-no-os-release",
			HaveKeyWithValue("name", "unknown"),
			HaveKeyWithValue("version_codename", "unknown"),
			HaveKeyWithValue("version_id", "unknown"),
		),
		Entry("an image that doesn't have cat but has an os-release", "pivotalnavcon/test-asset-no-grep-no-cat",
			HaveKeyWithValue("name", "Ubuntu"),
			HaveKeyWithValue("version_codename", "bionic"),
			HaveKeyWithValue("version_id", "18.04"),
		),
	)
})