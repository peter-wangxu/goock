package linux

import (
	"github.com/peter-wangxu/goock/pkg/exec"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger = logrus.New()

func SetLogger(l *logrus.Logger) {
	log = l
}

var executor = exec.New()

func SetExecutor(e exec.Interface) {
	executor = e
}
