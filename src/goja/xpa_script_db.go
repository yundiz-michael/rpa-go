package goja

import (
	"go.uber.org/zap"
	"merkaba/common"
	"runtime"
	"time"
)

type ScriptDb struct {
	Client    *common.MysqlClient
	Instances []*ScriptInstance
}

type MerkabaNode struct {
	IP       string `db:"ip"`
	Port     string `db:"port"`
	MaxCount int    `db:"maxCount"`
}

type ScriptInfo struct {
	Content string `db:"content"`
	Version string `db:"version"`
}

func (db *ScriptDb) Init() {
	db.Instances = make([]*ScriptInstance, 0)
}

func (db *ScriptDb) ReadScriptByUri(uri string) (content string, version string) {
	var id []string
	db.Client.Select(&id, "select id from merkaba where uri=?", uri)
	if len(id) == 0 {
		return "", ""
	}
	return db.ReadScript(id[0])
}

func (db *ScriptDb) ReadScript(id string) (content string, version string) {
	var sql string
	if common.Env.Environment.Production {
		sql = `select content,version from merkaba_script_prd where id=?`
	} else {
		sql = `select content,version from merkaba_script_dev where id=?`
	}
	var result ScriptInfo
	var err = db.Client.Get(&result, sql, id)
	if err != nil {
		common.LoggerStd.Error(sql, zap.Error(err))
		return "", ""
	}
	return result.Content, result.Version
}

func (db *ScriptDb) RegisterMerkabaNode() {
	sql := `
            insert into merkaba_node(ip,name,port,platform,createdTime,state) values(?,?,?,?,?,?)
            on duplicate key update 
                name = values(name),port = values(port),platform = values(platform),
                createdTime=values(createdTime),state=values(state)                             
	`
	var _, err = db.Client.Update(sql, common.LocalIP, common.LocalName, common.LocalPort, runtime.GOOS, time.Now().UnixMilli(), 4)
	if err != nil {
		common.LoggerStd.Error(sql, zap.Error(err))
		return
	}
}

func (db *ScriptDb) UnRegisterMerkabaNode() {
	sql := `
            update merkaba_node set state=5 where ip=?                             
	`
	var _, err = db.Client.Update(sql, common.LocalIP)
	if err != nil {
		common.LoggerStd.Error(sql, zap.Error(err))
		return
	}
}

func (db *ScriptDb) ReadLocalMerkabaNode() *MerkabaNode {
	sql := `
            select ip,maxCount from merkaba_node where ip=?                            
	`
	var result MerkabaNode
	var err = db.Client.Get(&result, sql, common.LocalIP)
	if err != nil {
		common.LoggerStd.Error(sql, zap.Error(err))
	}
	return &result
}

func (db *ScriptDb) FindMemInstance(taskName string) *ScriptInstance {
	for i := 0; i < len(db.Instances); i++ {
		instance := db.Instances[i]
		if instance.Context.TaskName == taskName {
			return instance
		}
	}
	return nil
}

func (db *ScriptDb) InstanceCount() int {
	return len(db.Instances)
}

// ExistMemInstance 0=NotExist, Idle, Running */
func (db *ScriptDb) ExistMemInstance(scriptUri string, taskName string) int {
	for i := 0; i < len(db.Instances); i++ {
		instance := db.Instances[i]
		if instance.Context.ScriptUri == scriptUri && instance.Context.TaskName == taskName {
			if instance.Status == "Running" {
				return 2
			} else {
				return 1
			}
		}
	}
	return 0
}

func (db *ScriptDb) AddMemInstance(instance *ScriptInstance) {
	db.Instances = append(db.Instances, instance)
}

func (db *ScriptDb) RemoveMemInstance(instance *ScriptInstance) {
	for i := 0; i < len(db.Instances); i++ {
		if db.Instances[i] == instance {
			db.Instances = append(db.Instances[:i], db.Instances[i+1:]...)
			i--
		}
	}
}

func (db *ScriptDb) ReadInstanceCount() (runCount int, idleCount int) {
	runCount = 0
	idleCount = 0
	for _, i := range db.Instances {
		if i.Status == "Running" {
			runCount += 1
		} else {
			idleCount += 1
		}
	}
	return runCount, idleCount
}

func (db *ScriptDb) updateNodeStatus() (runCount int, idleCount int) {
	runCount, idleCount = db.ReadInstanceCount()
	sql := "update merkaba_node set runCount=?,idleCount=? where ip=?"
	db.Client.Update(sql, runCount, idleCount, common.LocalIP)
	return runCount, idleCount
}

func (db *ScriptDb) UseMemInstance(instance *ScriptInstance) (runCount int, idleCount int) {
	instance.Status = "Running"
	runCount, idleCount = db.updateNodeStatus()
	return runCount, idleCount
}

func (db *ScriptDb) FreeMemInstance(instance *ScriptInstance) (runCount int, idleCount int) {
	instance.Status = "Idle"
	runCount, idleCount = db.updateNodeStatus()
	return runCount, idleCount
}
