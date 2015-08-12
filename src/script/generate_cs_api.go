package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func convert(s string, public bool) string {
	visibility := ""
	if public {
		visibility = "public "
	}
	s = strings.Replace(s, "const ", "", -1)
	if strings.Contains(s, "char_t *") {
		return visibility + strings.Replace(s, "char_t *", "string ", -1)
	}
	if strings.Contains(s, "uint64_t") {
		return visibility + strings.Replace(s, "uint64_t", "UInt64", -1)
	} 
	if strings.Contains(s, "int64_t") {
		return visibility + strings.Replace(s, "int64_t", "Int64", -1)
	}
	if strings.Contains(s, "void *") {
		return visibility + strings.Replace(s, "void *", "IntPtr ", -1)
	}
	if strings.Contains(s, "bool_t") {
		return visibility + strings.Replace(s, "bool_t", "bool", -1)
	}
	return s
}

func main() {
	// input
	b, err := ioutil.ReadFile("src/toggl_api.h")
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	l := strings.Split(string(b), "\n")
	var out []string
	write := func(s string) {
		out = append(out, s)
	}
	// start class
	write("// Warning: this file is autogenerated! do not modify, all changes will be lost!")
	write("// Copyright 2015 Toggl Desktop developers.")
	write("")
	write("using System;")
	write("using System.Runtime.InteropServices;")
	write("")
	write("namespace TogglDesktop")
	write("{")
	write("    public static partial class Toggl")
	write("    {")
	write("    		private const string dll = \"TogglDesktopDLL.dll\";");
	write("    		private const CharSet charset = CharSet.Unicode;");
	write("    		private const CallingConvention convention = CallingConvention.Cdecl;");
	write("    		private const int structPackingBytes = 8;");
	write("")
	csclass, csfunc := "", ""
	for _, s := range l {
		// line feeds
		if len(s) == 0 {
			write(s)
			continue
		}
		if strings.Contains(s, "//") {
			if strings.Contains(s, "in C#: ") {
				csclass = strings.Split(s, ":")[1]
			} else if !strings.Contains(s, "#") && !strings.Contains(s, "Copyright") {
				write(s)
			}
		} else if strings.Contains(s, "#define k") {
			w := strings.Split(s, " ")
			name, value := w[1], w[2]
			write("private const int " + name + " = " + value + ";")
		}  else if strings.Contains(s, "typedef struct {") {
			write("[StructLayout(LayoutKind.Sequential, Pack = structPackingBytes, CharSet = CharSet.Unicode)]")
			write("public struct" + csclass)
			write("{")
		} else if len(csclass) != 0 {
			if strings.Contains(s, "} Toggl") {
				csclass = ""
				write("}")
			} else if strings.Contains(s, "char_t *") {
				write("[MarshalAs(UnmanagedType.LPWStr)]")
				write(convert(s, true))
			} else if strings.Contains(s, "uint64_t") {
				write(convert(s, true))
			} else if strings.Contains(s, "int64_t") {
				write(convert(s, true))
			} else if strings.Contains(s, "void *") {
				write(convert(s, true))
			} else if strings.Contains(s, "bool_t") {
				write("[MarshalAs(UnmanagedType.I1)]")
				write(convert(s, true))
			}
		} else if strings.Contains(s, "TOGGL_EXPORT") && !strings.Contains(s, "#") {
			write("[DllImport(dll, CharSet = charset, CallingConvention = convention)]")
			s = strings.Trim(s, " ")
			s = convert(s, false)
			w := strings.Split(s, " ")
			for i, k := range w {
				println(i, k)
			}
			t := w[1]
			csfunc = w[len(w) - 1]
			write("private static extern " + t + " " + csfunc)
		} else if len(csfunc) != 0 {
			if strings.Contains(s, ");") {
				csfunc = ""
			}
			write(convert(s, false))
		}
	}
	// finish class
	write("    }")
	write("")
	write("}")

	// output
	b = []byte(strings.Join(out, "\n"))
	err = ioutil.WriteFile("src/ui/windows/TogglDesktop/TogglDesktop/TogglApi.cs", b, 0644)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}