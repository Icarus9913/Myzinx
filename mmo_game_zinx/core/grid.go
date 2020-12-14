package core

import "sync"

/*
	一个AOI地图中的格子类型
*/
type Grid struct {
	//格子ID
	GID int
	//格子的左边边界坐标
	MinX int
	//格子的右边边界坐标
	MaxX int
	//格子的上边边界坐标
	MinY int
	//格子的下边边界坐标
	MaxY int
	//当前各自内玩家或者物体成员的ID集合
	playerIDs map[int]bool
	//保护当前集合的锁
	pIDLock sync.RWMutex
}


//初始化当前的格子的方法
func NewGrid(gID,minX,MaxX,minY,maxY int) *Grid {
	return &Grid{
		GID: gID,
		MinX: minX,
		MaxX: MaxX,
		MinY: minY,
		MaxY: maxY,
		playerIDs: make(map[int]bool),
	}
}


//
//
//
//
//