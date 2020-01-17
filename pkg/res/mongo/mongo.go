package mongo

import "github.com/saas/hostgolang/pkg/res"

func NewMongo(mem, cpu, ss float64, envs map[string]string) res.Res {
	return res.NewGenericResource(mem, cpu, ss, envs, 27017, "mongo", "mongo:4.2.2")
}
