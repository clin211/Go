// Copyright 2023 Forest-211 &lt;767425412lin@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Forest-211/go/tree/miniblog.

package main

import (
	"os"

	// 使程序自动设置 GOMAXPROCS 以匹配 Linux 容器 CPU 配额。通过正确设置容器的 CPU 配额，可以解决 GOMAXPROCS 可能设置过大，导致生成线程过多，从而导致严重的上下文切换，浪费 CPU，降低程序性能的潜在问题
	_ "go.uber.org/automaxprocs"

	"github.com/Forest-211/miniblog/internal/miniblog"
)

func main() {
	command := miniblog.NewMiniBlogCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
