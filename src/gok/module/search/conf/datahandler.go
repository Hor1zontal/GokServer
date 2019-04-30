/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/7/11
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package conf


//获取事件的随机目标规则
func GetEventRandomTargets(eventType int32) map[int32]*RandomTargetBase {
	return Base.RandomTarget[eventType]
}
