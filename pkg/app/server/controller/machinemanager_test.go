package controller

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
)

type v1Suite struct {
	db         *gorm.DB
	mock       sqlmock.Sqlmock
	departnode model.DepartNode
}

func TestGORMV1(t *testing.T) {
	s := &v1Suite{}
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	if db == nil {
		t.Error("mock db is null")
	}

	if s.mock == nil {
		t.Error("sqlmock is null")
	}

	s.db, err = gorm.Open("postgres", db)
	if err != nil {
		t.Errorf("Failed to open gorm db, got error: %v", err)
	}

	if s.db == nil {
		t.Error("gorm db is null")
	}

	s.departnode = model.DepartNode{
		ID:           1,
		PID:          0,
		ParentDepart: "xx",
		Depart:       "xx",
		NodeLocate:   0}

	defer db.Close()

	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "depart_nodes" ("id","p_id","parent_depart","depart","node_locate")
		 VALUES ($1,$2,$3,$4,$5) RETURNING "depart_nodes"."id"`)).
		WithArgs(s.departnode.ID, s.departnode.PID, s.departnode.ParentDepart, s.departnode.Depart, s.departnode.NodeLocate).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(s.departnode.ID))

	s.mock.ExpectCommit()
	//s.db.Create(&s.departnode).Error
	if err = dao.AddDepart(s.db, &(s.departnode)); err != nil {
		t.Errorf("Failed to insert to gorm db, got error: %v", err)
	}

	err = s.mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}
