package integration_test

import (
	"sort"

	"github.com/pivotal/deplab/metadata"
	"github.com/pivotal/deplab/providers"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deplab dpkg", func() {
	var (
		inputImage    string
		metadataLabel metadata.Metadata
	)

	Context("with an ubuntu:bionic image", func() {
		BeforeEach(func() {
			inputImage = "pivotalnavcon/test-asset-additional-sources"
			metadataLabel = runDeplabAgainstImage(inputImage)
		})

		It("applies a metadata label", func() {
			dependencyMetadata := metadataLabel.Dependencies[0].Source.Metadata
			dpkgMetadata := dependencyMetadata.(map[string]interface{})

			By("listing the dpkg sources, sorted alphabetically")
			sources, ok := dpkgMetadata["apt_sources"].([]interface{})
			Expect(ok).To(BeTrue())
			Expect(len(sources)).To(BeNumerically(">", 0))
			Expect(sources).To(ConsistOf(
				"deb http://archive.ubuntu.com/ubuntu/ bionic main restricted",
				"deb http://archive.ubuntu.com/ubuntu/ bionic universe",
				"deb http://archive.ubuntu.com/ubuntu/ bionic-updates main restricted",
				"deb http://archive.ubuntu.com/ubuntu/ bionic-updates universe",
				"deb http://archive.ubuntu.com/ubuntu/ bionic multiverse",
				"deb http://archive.ubuntu.com/ubuntu/ bionic-updates multiverse",
				"deb http://archive.ubuntu.com/ubuntu/ bionic-backports main restricted universe multiverse",
				"deb http://security.ubuntu.com/ubuntu/ bionic-security main restricted",
				"deb http://security.ubuntu.com/ubuntu/ bionic-security universe",
				"deb http://security.ubuntu.com/ubuntu/ bionic-security multiverse",
				"deb http://example.com/ubuntu getdeb example",
			))
			Expect(AreSourcesSorted(sources)).To(BeTrue())

			By("listing debian package dependencies in the image, sorted by name")
			Expect(metadataLabel.Dependencies[0].Type).To(Equal(providers.DebianPackageListSourceType))

			pkgs := dpkgMetadata["packages"].([]interface{})
			Expect(pkgs).To(HaveLen(89))
			Expect(ArePackagesSorted(pkgs)).To(BeTrue())

			By("generating a sha256 digest of the metadata content as version")
			Expect(metadataLabel.Dependencies[0].Source.Version["sha256"]).To(MatchRegexp(`^[0-9a-f]{64}$`))
		})
	})

	Context("with an image without dpkg", func() {
		BeforeEach(func() {
			inputImage = "alpine:latest"
			metadataLabel = runDeplabAgainstImage(inputImage)
		})

		It("does not return a dpkg list", func() {
			_, ok := filterDpkgDependency(metadataLabel.Dependencies)
			Expect(ok).To(BeFalse())
		})
	})

	Context("with an image with dpkg, but no apt sources", func() {
		BeforeEach(func() {
			inputImage = "pivotalnavcon/test-asset-no-sources"
			metadataLabel = runDeplabAgainstImage(inputImage)
		})

		It("does not return a dpkg list", func() {
			dependencyMetadata := metadataLabel.Dependencies[0].Source.Metadata
			dpkgMetadata := dependencyMetadata.(map[string]interface{})

			sources, ok := dpkgMetadata["apt_sources"].([]interface{})

			Expect(ok).To(BeTrue())
			Expect(sources).To(BeEmpty())
		})
	})

	Context("with an ubuntu:bionic based image with a non-shell entrypoint", func() {
		BeforeEach(func() {
			inputImage = "pivotalnavcon/test-asset-entrypoint-return-stdout"
			metadataLabel = runDeplabAgainstImage(inputImage)
		})

		It("should return the apt source list", func() {
			dependencyMetadata := metadataLabel.Dependencies[0].Source.Metadata
			dpkgMetadata := dependencyMetadata.(map[string]interface{})

			sources, ok := dpkgMetadata["apt_sources"].([]interface{})

			Expect(ok).To(BeTrue())
			Expect(sources).ToNot(BeEmpty())
			Expect(sources[0].(string)).To(ContainSubstring("deb http://archive.ubuntu.com/ubuntu/ bionic main restricted"))
		})
	})

	Context("with Pivotal Tiny", func() {
		BeforeEach(func() {
			inputImage = "cloudfoundry/run:tiny"
			metadataLabel = runDeplabAgainstImage(inputImage)
		})

		It("returns a dpkg list", func() {
			By("listing debian package dependencies in the image alphabetically")
			Expect(metadataLabel.Dependencies[0].Type).To(Equal(providers.DebianPackageListSourceType))

			dependencyMetadata := metadataLabel.Dependencies[0].Source.Metadata
			dpkgMetadata := dependencyMetadata.(map[string]interface{})

			pkgs := dpkgMetadata["packages"].([]interface{})
			Expect(pkgs).To(HaveLen(7))
			Expect(ArePackagesSorted(pkgs)).To(BeTrue())
		})
	})
})

func filterDpkgDependency(dependencies []metadata.Dependency) (metadata.Dependency, bool) {
	for _, dependency := range dependencies {
		if dependency.Source.Type == providers.DebianPackageListSourceType {
			return dependency, true
		}
	}
	return metadata.Dependency{}, false //should never be reached
}

func ArePackagesSorted(pkgs []interface{}) bool {
	collator := collate.New(language.BritishEnglish)
	return sort.SliceIsSorted(pkgs, func(p, q int) bool {
		lhs := pkgs[p].(map[string]interface{})
		rhs := pkgs[q].(map[string]interface{})
		return collator.CompareString(lhs["package"].(string), rhs["package"].(string)) <= 0
	})
}

func AreSourcesSorted(sources []interface{}) bool {
	collator := collate.New(language.BritishEnglish)
	return sort.SliceIsSorted(sources, func(p, q int) bool {
		return collator.CompareString(sources[p].(string), sources[q].(string)) <= 0
	})
}
