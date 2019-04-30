/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/11/6
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package core

import (
	"time"
	"math/rand"
	"aliens/common/util"
)

type StarCategory struct {
	level  int32 //等级
	min  int32
	max  int32
	mapping  map[int32]time.Time //
	//filters  map[int32]*StarState //过滤列表
}

//随机星球
func RandomStar(mapping map[int32]time.Time, filter []int32) (int32, *time.Time) {
	if mapping == nil {
		return 0, nil
	}
	randLen := len(mapping)
	if randLen == 0 {
		return 0, nil
	}
	randomIndex := rand.Intn(randLen)
	index := -1
	for starID, activeTime  := range mapping {
		index ++
		if randomIndex != index {
			continue
		}
		if filter != nil && util.ContainsInt32(starID, filter) {
			randomIndex ++
			index ++
			continue
		}

		if randomIndex == index {
			return starID, &activeTime
		}
	}
	return 0, nil
}




func (this *StarCategory) InScope(level int32 ) bool {
	return level >= this.min && level <= this.max
}
