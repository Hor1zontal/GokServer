/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/6/28
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package cache

import (
	"time"
	"aliens/common/character"
)

const (

	registerTotal = "reg:"
	dayLoginPrefix = "day:login:"
	dayRegisterPrefix = "day:register:"
	dayChargePrefix = "day:charge:"
	dayChargeFeePrefix = "day:chargeFee:"
)


func getDayPrefix(time time.Time) string {
	return time.Format("20060102")
}


func (this *UserCacheManager) IncrDayLogin(uid int32, time time.Time) {
	this.redisClient.HIncrby(dayLoginPrefix + getDayPrefix(time), character.Int32ToString(uid), 1)
}

func (this *UserCacheManager) IncrDayRegister(time time.Time) {
	this.redisClient.Incr(dayRegisterPrefix + getDayPrefix(time))
	this.redisClient.Incr(registerTotal)
}

func (this *UserCacheManager) IncrDayCharge(uid int32, time time.Time) {
	this.redisClient.HIncrby(dayChargePrefix + getDayPrefix(time), character.Int32ToString(uid), 1)
}

func (this *UserCacheManager) IncrByDayChargeFee(time time.Time, addValue float64) {
	this.redisClient.IncrBy(dayChargeFeePrefix + getDayPrefix(time), addValue)
}

func (this *UserCacheManager) GetDayLogin(time time.Time) int {
	result, _ := this.redisClient.HLen(dayLoginPrefix + getDayPrefix(time))
	return result
}

func (this *UserCacheManager) GetDayRegister(time time.Time) int {
	return this.redisClient.GetDataInt32(dayRegisterPrefix + getDayPrefix(time))
}

func (this *UserCacheManager) GetDayCharge(time time.Time) int {
	result, _ := this.redisClient.HLen(dayChargePrefix + getDayPrefix(time))
	return result
}

func (this *UserCacheManager) GetDayChargeFee(time time.Time) int {
	return this.redisClient.GetDataInt32(dayChargeFeePrefix + getDayPrefix(time))
}

func (this *UserCacheManager) GetRegisterTotal() int {
	return this.redisClient.GetDataInt32(registerTotal)
}

func (this *UserCacheManager) SetRegisterTotal(total int) {
	this.redisClient.SetData(registerTotal, total)
}