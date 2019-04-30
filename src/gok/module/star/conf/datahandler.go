package conf

import (
	"aliens/common/character"
	"aliens/common/set"
	"aliens/common/util"
	"gok/service/msg/protocol"
	"math/rand"
	"strconv"
)

//初始化基础信息
func InitBuildingLimit() {
	buildingLimit := make(map[int32]int32)
	for _, building := range Base.Building {
		LimitId := GenStarBuildingID(building.StarID, building.BuildID)
		if buildingLimit[LimitId] == 0 || building.Level > buildingLimit[LimitId] {
			buildingLimit[LimitId] = building.Level
		}
	}
	Base.BLimit = buildingLimit
}

//RandomResult map[int32]util2.WeightData `json:"-"`

func (this *BelieverUpgradeBase)RandomDataAddProb(weightData map[int32]util.WeightData, addProb float32) map[int32]util.WeightData {
	if addProb <= 0 {
		return weightData
	}
	result := make(map[int32]util.WeightData)
	for weight, data := range weightData {
		//不是最大的权重需要加系数
		if weight != this.MaxWeight {
			weight = int32(float32(weight) * (1 + addProb))
		}
		result[weight] = data
	}
	return result
}

func (this *BelieverUpgradeBase) AddProb(addProb float32) map[int32]util.WeightData {
	if addProb <= 0 {
		return this.RandomResult
	}
	result := make(map[int32]util.WeightData)
	for weight, data := range this.RandomResult {
		//不是最大的权重需要加系数
		if weight != this.MaxWeight {
			weight = int32(float32(weight) * (1 + addProb))
		}
		result[weight] = data
	}
	return result
}

func (this *BelieverUpgradeBase) GetRandomResult() map[int32]util.WeightData {
	return this.RandomResult
}


func GetBuffBase(buffID int32) (int32, float32) {
	buffBase := Base.Buff[buffID]
	if (buffBase == nil) {
		return 0,0
	}
	return buffBase.Type, buffBase.Ratio
}

func GetItemBuff(itemID int32) int32 {
	itemBase := Base.Item[itemID]
	if itemBase == nil {
		return 0
	}
	return itemBase.BuffID
}


//获取信徒
func GetUpgradeBelieverRatio(civilizationLevel int32, upgradeLevel int32) float64 {
	//信徒等级是1 - 6 2代表一级信徒  ，索引0代表1级信徒
	index := int(upgradeLevel - 2)
	if index < 0 {
		return 0
	}

	entry := Base.CivilizationRatio[civilizationLevel]
	if entry == nil {
		return 0
	}
	if index < len(entry.Modulus)  {
		return entry.Modulus[index]
	}
	return 0
}

func GetCivilConfigByLevel(seq int32, level int32) *Civilization {
	if seq > Base.CivilMaxSeq {
		seq = Base.CivilMaxSeq
	}

	levelMapping := Base.CivilData[seq]
	if levelMapping == nil {
		return  nil
	}
	return levelMapping[level]
}

func GetCivilRewardConfigById(id int32) *CivilizationReward{

	rewardMapping := Base.CivilizationReward[id]
	if rewardMapping == nil {
		return nil
	}
	return rewardMapping
}

func GetCivilMaxLevel(seq int32) int32 {
	if seq > Base.CivilMaxSeq {
		seq = Base.CivilMaxSeq
	}

	return Base.CivilMaxLevel[seq]
}


//获取星球建筑配置信息
func GetBuildingConf(starType int32, buildingType int32, buildingLevel int32) *BuildingBase {
	for _, building := range Base.Building {
		if building.StarID == starType && building.BuildID == buildingType &&
			building.Level == buildingLevel {
			return building
		}
	}
	return nil
}


func Contains(arrays []int32, data int32) bool {
	if (arrays == nil || len(arrays) == 0) {
		return false
	}
	for _, array := range arrays {
		if (data == array) {
			return true
		}
	}
	return false
}

//获取建筑加成的法力值上限
func GetBuildingPowerLimit(starType int32, buildingType int32, buildingLevel int32)  int32 {
	var powerLimit int32 = 0
	for _, building := range Base.Building {
		if building.StarID == starType && building.BuildID == buildingType {
			if building.Level <= buildingLevel {
				powerLimit += building.PowerLimit
			}
		}
	}
	return powerLimit
}

func GetBuildingPowerReward(starType int32, buildingType int32, buildingLevel int32) int32 {
	var powerReward int32 = 0
	for _, building := range Base.Building {
		if building.StarID == starType && building.BuildID == buildingType {
			if building.Level <= buildingLevel {
				powerReward += building.PowerReward
			}
		}
	}
	return powerReward
}

//获取指定星球的所有建筑类型
func GetStarBuildingTypes(starId int32) []interface{} {
	resultSet := set.NewHashSet()
	for _, building := range Base.Building {
		if (building.StarID == starId) {
			resultSet.Add(building.BuildID)
		}
	}
	return resultSet.Elements()
}


func GetMaxBuildingLevel(starType int32, buildingType int32) int32 {
	id := GenStarBuildingID(starType, buildingType)
	return Base.BLimit[id]
}

//构建星球建筑id
func GenStarBuildingID(starType int32, buildingType int32) int32 {
	result,_ := strconv.Atoi(strconv.Itoa(int(starType)) + strconv.Itoa(int(buildingType)))
	return int32(result)
}


//获取信徒升级结果
func GetBelieverUpgradeResult(selectID string, matchID string) []*BelieverUpgradeBase {
	upgradeCondition := Base.BelieverUpgrade[selectID + matchID]
	if upgradeCondition == nil {
		return Base.BelieverUpgrade[matchID + selectID]
	}
	return upgradeCondition
}

func RandomStarType() int32 {
	randomLen := len(Base.StarWeCanArrive)
	var starType int32 = 1
	if randomLen != 0 {
		starType = Base.StarWeCanArrive[rand.Intn(randomLen)]
	}
	return starType
}

func RandomStarsType(confStars []int32, num int32, starsInfo []*protocol.StarInfo) []int32 {

	//if num == -1 {
	//	result := make([]int32, len(confStars))
	//	util.DeepCopy(confStars, result)
	//	return result
	//}

	result := make([]int32,num)
	//starsRandom := make([]int32, len(confStars))

	starsRandom := character.CopyArray(confStars)
	//if starsRandom = confStars

	//随机可选的星球去除已拥有的星球
	for _, starInfo := range starsInfo {
		for index, starType := range starsRandom {
			if starType == starInfo.StarType {
				starsRandom = append(starsRandom[:index], starsRandom[index+1:]...)
				continue
			}
		}
	}

	var starType int32 = 0
	if len(starsRandom) <= int(num) {
		return starsRandom
	}
	for i := 0 ; i < int(num); i++ {
		randomLen := len(starsRandom)
		randNum := rand.Intn(randomLen)
		starType = starsRandom[randNum]
		starsRandom = append(starsRandom[:randNum], starsRandom[randNum+1:]...)
		result[i] = starType
	}
	return result
}

func GetNextStarType(starType int32) int32 {
	matchIndex := -1
	for index, star := range Base.StarWeCanArrive {
		if star == starType {
			matchIndex = index
		}
	}

	if matchIndex == -1 {
		return RandomStarType()
	}

	if matchIndex + 1 < len(Base.StarWeCanArrive) {
		return Base.StarWeCanArrive[matchIndex + 1]
	}
	return Base.StarWeCanArrive[0]
}

func GetStarInitBeliever(starType int32) []string {
	star := Base.Star[starType]
	if (star == nil) {
		return nil
	}
	return star.Believers
}

func GetGameObjectLevel(objectID string) int32 {
	gameObject := Base.GameObject[objectID]
	if gameObject == nil {
		return -1
	}
	return gameObject.LV
}

//获取信徒的加速配置
func GetBuildingUpgradeAcc(believerLevel int32, buildingLevel int32) *BelieverBuildingExpedite{
	for _, result := range Base.BelieverBuildingExp {
		if (result.BelieverLevel == believerLevel && result.BuildingLevel == buildingLevel) {
			return result
		}
	}
	return nil
}

//func CompareLevelNumToConfig(levelNum []int32, require map[int32]*ItemGroupUnlockBase) int {
//	buildLen := len(levelNum)
//	requireLen := len(require)
//	var resultLen = 0
//	for i := 1 ; i <= requireLen ; i++ {
//		for j := 0 ; j < buildLen ; j++ {
//			if levelNum[j] < require[int32(i)].BuildingRequire[j] {
//				return resultLen
//			}
//		}
//		resultLen = i
//	}
//	return resultLen
//}

func GetUnlockGroupConfigID (level int32) int {
	var result int32 = 0
	for _, requireData := range Base.GroupUnlockRequire {
		if level < requireData.BuildingRequire {
			break
		}
		result = requireData.ID
	}
	return int(result)
}

func GetBelieverLevel(believerID string) int32 {
	return character.StringToInt32(believerID[3:4])
}

func GetIndexOfFirstGroupItem(level int32) int32 {
	for index, relicLevel := range Base.FirstGroupRelicDrop {
		if level == relicLevel {
			return int32(index)
		}
	}
	return -1
}