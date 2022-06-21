package runtime

func dtlSystemNamespace(s string) string {
	return "dtl_" + s
}

func memstatNamespace(s string) string {
	return dtlSystemNamespace("go_memstats_" + s)
}
