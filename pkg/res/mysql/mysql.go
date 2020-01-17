package mysql

import "github.com/saas/hostgolang/pkg/res"

func NewMysql(mem, cpu, ss float64, envs map[string]string) res.Res {
	return res.NewGenericResource(mem, cpu, ss, envs, 3306, "mysql", "mysql:8.0.19")
}
