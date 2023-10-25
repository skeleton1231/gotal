// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package flag

import (
	"fmt"

	"github.com/spf13/pflag"
)

func AddGlobalFlags(fs *pflag.FlagSet, name string) {
	fs.BoolP("help", "h", false, fmt.Sprintf("help for %s", name))
}
