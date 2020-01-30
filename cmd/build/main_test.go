package main_test

import (
	"os"
	"testing"

	"github.com/bstick12/git-metadata-buildpack/metadata"
	"github.com/bstick12/git-metadata-buildpack/utils"

	cmdBuild "github.com/bstick12/git-metadata-buildpack/cmd/build"
	"github.com/buildpack/libbuildpack/buildpackplan"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitBuild(t *testing.T) {
	spec.Run(t, "Build", testBuild, spec.Report(report.Terminal{}))
}

func testBuild(t *testing.T, when spec.G, it spec.S) {

	var factory *test.BuildFactory

	it.Before(func() {
		RegisterTestingT(t)
		factory = test.NewBuildFactory(t)
	})

	when("building source", func() {
		it("should pass if successful", func() {

			defer utils.ResetEnv(os.Environ())
			os.Clearenv()

			factory.AddPlan(buildpackplan.Plan{
				Name:    metadata.Dependency,
				Version: "",
				Metadata: buildpackplan.Metadata{
					"build":  false,
					"launch": true,
					"sha":    "7aa636e253c4115df34b1f2fab526739cbf27570",
					"branch": "fork/master",
					"remote": "git@github.com/example/example.git",
				},
			})
			code, err := cmdBuild.RunBuild(factory.Build)
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(build.SuccessStatusCode))
			metadataLayer := factory.Build.Layers.Layer(metadata.Dependency)
			Expect(metadataLayer).To(test.HaveLayerMetadata(false, false, true))
			md := metadata.GitMetadata{}
			metadataLayer.ReadMetadata(&md)
			Expect(md).To(Equal(metadata.GitMetadata{
				Sha:    "7aa636e253c4115df34b1f2fab526739cbf27570",
				Branch: "fork/master",
				Remote: "git@github.com/example/example.git",
			}))

		})

		it("should fail if it doesn't contribute", func() {
			defer utils.ResetEnv(os.Environ())
			os.Clearenv()
			code, err := cmdBuild.RunBuild(factory.Build)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to find build plan"))
			Expect(code).To(Equal(cmdBuild.FailureStatusCode))
		})
	})

}
