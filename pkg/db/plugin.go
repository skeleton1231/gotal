// Copyright 2023 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package db

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Constants used to identify callback names and context values.
const (
	callBackBeforeName = "core:before" // Identifier for the 'before' callback.
	callBackAfterName  = "core:after"  // Identifier for the 'after' callback.
	startTime          = "_start_time" // Context key to store the start time of SQL execution.
)

// TracePlugin is a GORM plugin to trace the execution time of SQL statements.
type TracePlugin struct{}

// Name returns the name of the trace plugin.
func (op *TracePlugin) Name() string {
	return "tracePlugin"
}

// Initialize registers the callbacks to measure SQL execution time.
func (op *TracePlugin) Initialize(db *gorm.DB) (err error) {
	// Register callbacks that will be triggered before SQL execution.
	_ = db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	_ = db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	_ = db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	_ = db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before)
	_ = db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	_ = db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)

	// Register callbacks that will be triggered after SQL execution to measure the duration.
	_ = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	_ = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	_ = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	_ = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	_ = db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	_ = db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)

	return
}

// Ensure TracePlugin implements the gorm.Plugin interface.
var _ gorm.Plugin = &TracePlugin{}

// before is a callback function that sets the start time before SQL execution.
func before(db *gorm.DB) {
	// Record the current time as the start time of the SQL execution.
	db.InstanceSet(startTime, time.Now())
}

// after is a callback function that calculates and logs the time taken for SQL execution.
func after(db *gorm.DB) {
	// Retrieve the start time set in the 'before' callback.
	_ts, isExist := db.InstanceGet(startTime)
	if !isExist {
		return
	}

	// Ensure the retrieved value is of type time.Time.
	ts, ok := _ts.(time.Time)
	if !ok {
		return
	}

	// Calculate the duration of SQL execution and log it.
	logrus.Infof("sql cost time: %fs", time.Since(ts).Seconds())
}
