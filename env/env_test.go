package env_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/DamnWidget/VenGO/cache"
	"github.com/DamnWidget/VenGO/env"
)

var _ = Describe("Env", func() {

	// disable log output
	log.SetOutput(ioutil.Discard)

	AfterSuite(func() {
		path := filepath.Join(cache.VenGO_PATH, "goTest")
		os.RemoveAll(path)
	})

	Describe("NewEnvironment", func() {
		It("Return a configured Environment structure", func() {
			name, prompt := "goTest", "(goTest)"
			path := filepath.Join(cache.VenGO_PATH, name)
			e := env.NewEnvironment(name, prompt)
			Expect(e.Goroot).To(Equal(filepath.Join(path, "lib")))
			Expect(e.Gotooldir).To(Equal(filepath.Join(
				path, "lib", "pkg", "tool", fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH))))
			Expect(e.Gopath).To(Equal(path))
			Expect(e.PS1).To(Equal(prompt))
			Expect(e.VenGO_PATH).To(Equal(path))
		})
	})

	Describe("Generate", func() {
		It("Will generate a valid template file", func() {
			name := "goTest"
			prompt := "[{(goTest)}] "
			e := env.NewEnvironment(name, prompt)
			err := e.Generate()

			Expect(err).ToNot(HaveOccurred())
			activate, err := ioutil.ReadFile(filepath.Join(e.VenGO_PATH, "bin", "activate"))

			Expect(err).ToNot(HaveOccurred())
			byteLines := bytes.Split(activate, []byte("\n"))
			vengoPath := fmt.Sprintf(`VENGO_ENV="%s/goTest"`, cache.VenGO_PATH)
			sysPath := fmt.Sprint(`PATH="$GOROOT/bin:$PATH"`)
			goRoot := fmt.Sprintf(`GOROOT="%s"`, e.Goroot)
			goTooldir := fmt.Sprintf(`GOTOOLDIR="%s"`, e.Gotooldir)
			goPath := fmt.Sprintf(`GOPATH="%s"`, e.VenGO_PATH)
			ps1 := fmt.Sprintf(`PS1="%s ${_VENGO_PREV_PS1}"`, e.PS1)

			Expect(byteLines[53]).To(Equal([]byte(vengoPath)))
			Expect(byteLines[73]).To(Equal([]byte(goRoot)))
			Expect(byteLines[77]).To(Equal([]byte(sysPath)))
			Expect(byteLines[81]).To(Equal([]byte(goTooldir)))
			Expect(byteLines[84]).To(Equal([]byte(goPath)))
			Expect(byteLines[87]).To(Equal([]byte(ps1)))

		})
	})

	Describe("Install", func() {
		It("Will create a symboolic link into VenGO_PATH", func() {
			Expect(cache.CacheDonwloadMercurial("1.3.2")).To(Succeed())

			name := "goTest"
			prompt := "(gotest)"
			e := env.NewEnvironment(name, prompt)
			err := e.Generate()

			Expect(err).ToNot(HaveOccurred())
			Expect(e.Install("1.3.2")).To(Succeed())
		})
	})
})
