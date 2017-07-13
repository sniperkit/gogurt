package packages

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/alexandrecarlton/gogurt"
)

// TODO: Link to our own zlib

type GCC struct{}

func (gcc GCC) Name() string {
	return "gcc"
}

func (gcc GCC) URL(version string) string {
	return fmt.Sprintf("https://ftp.gnu.org/pub/gnu/gcc/gcc-%s/gcc-%s.tar.gz", version, version)
}

func (gcc GCC) Build(config gogurt.Config) error {
	buildDir := filepath.Join(config.BuildDir(gcc), "build")
	os.MkdirAll(buildDir, 0755)

	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(gcc),
		Args: []string{
			"--with-gmp=" + config.InstallDir(GMP{}),
			"--with-mpfr=" + config.InstallDir(MPFR{}),
			"--with-mpc=" + config.InstallDir(MPC{}),
			// "--with-system-zlib", // can't seem to link in our own. :/
			"--disable-libada",
			"--disable-multilib", // no need for 32 bit support
			"--enable-languages=c,c++", // This isn't documented in ./configure --help. Huh.
		},
		Dir: buildDir,
		CFlags: []string{
			// "-I" + config.IncludeDir(Zlib{}),
		},
		CppFlags: []string{
			// "-I" + config.IncludeDir(Zlib{}),
		},
		CxxFlags: []string{
			// "-I" + config.IncludeDir(Zlib{}),
		},
		LdFlags: []string{
			// "-L" + config.LibDir(Zlib{}),
		},
		Libs: []string{
			// "-lz",
		},
	}.Cmd()
	// We're calling configure in the root of the source, but from inside the build directory.
	configure.Path = "../configure"
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
		Dir: buildDir,
	}.Cmd()
	return make.Run()
}

func (gcc GCC) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (gcc GCC) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		GMP{},
		MPC{},
		MPFR{},
		Zlib{},
	}
}
