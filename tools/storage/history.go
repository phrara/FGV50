package storage

import (
	"fgv50/err"
	"fmt"
	"os"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var histDBPath string

func init() {
	histDBPath, _ = os.Getwd()
	histDBPath = filepath.Join(histDBPath, "/db/history.db")
}

type HisDB struct {
	db *leveldb.DB
}

func NewHistDB() (*HisDB, error) {
	db, err1 := leveldb.OpenFile(histDBPath, &opt.Options{
		BlockSize: 64 * opt.KiB,
		CompactionTotalSize: 12 * opt.MiB,
	})
	if err1 != nil {
		return nil, fmt.Errorf("%s: %s", err.ErrLevelDBInit, err1)
	}
	return &HisDB{
		db: db,
	}, nil
}

func (h *HisDB) PutResRecord(kTime, vRes []byte) error {
	kTime = append(kTime, []byte("RES")...)
	// println("==============",string(kTime), "==================================")
	return h.db.Put(kTime, vRes, nil)
}

func (h *HisDB) PutVulRecord(kTime, vVul []byte) error {
	kTime = append(kTime, []byte("VUL")...)
	// println("==============",string(kTime), "==================================")

	return h.db.Put(kTime, vVul, nil)
}

func (h *HisDB) PutHDRecord(kTime, vHD []byte) error {
	kTime = append(kTime, []byte("HD")...)
	// println("==============",string(kTime), "==================================")

	return h.db.Put(kTime, vHD, nil)
}

func (h *HisDB) GetRecord(kTime []byte) (vRes []byte, vVul []byte, vHD []byte) {
	kTimeRes := []byte(string(kTime)+"RES")
	kTimeVul := []byte(string(kTime)+"VUL")
	kTimeHD := []byte(string(kTime)+"HD")
	// println(string(kTimeRes), string(kTimeVul))
	vRes, _ = h.db.Get(kTimeRes, nil)
	vVul, _ = h.db.Get(kTimeVul, nil)
	vHD, _ = h.db.Get(kTimeHD, nil)
	return vRes, vVul, vHD
}

func (h *HisDB) Close() {
	h.db.Close()
}