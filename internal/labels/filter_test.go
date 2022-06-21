package labels

import "testing"

func Test_FilterRoute(t *testing.T) {
	examples := []struct {
		input  string
		output string
	}{
		{"/api?id=23", "/api?id=*"},
		{"/2323", "/*"},
		{"/", "/"},
		{"/api?id=23&ip=23", "/api?id=*&ip=*"},
		{"", ""},
	}

	for _, example := range examples {
		if example.output != Filter.FilterRoute(example.input) {
			t.Errorf("filter number failed, input: %s,expected ouput: %s, real output: %s",
				example.input, example.output, Filter.FilterRoute(example.input))
		}
	}
}

func Test_FilterSQL(t *testing.T) {
	examples := []struct {
		rawSQL string
		ok     bool
		cmd    string
		sql    string
	}{
		{
			rawSQL: "SELECT * from t_user_info where id (8) in (?,?, ?, ?) in (?,?, ?, ?) and id = 123",
			cmd:    "select",
			ok:     true,
			sql:    "select * from t_user_info where id (?) in (?) in (?) and id = ?",
		},
		{
			rawSQL: " select * from t_user_info123 where id in (3) and id in (2 , 3 , ? , ?) in (?,?, ?) and time = '2020-0'8-20 12:22:02'",
			ok:     true,
			cmd:    "select",
			sql:    "select * from t_user_info? where id in (?) and id in (?) in (?) and time = '?'",
		},
		{
			rawSQL: "  select *   from t_user_info?   where   id   in (?) and id in (*) in (*) and time = '?-?-? ?:?:?'   ",
			ok:     true,
			cmd:    "select",
			sql:    "select * from t_user_info? where id in (?) and id in (?) in (?) and time = '?'",
		},
		{
			rawSQL: "INSERT INTO  user_rights (user_id,right,group_id,value) VALUES ( '42',  '160',  '1',  '1' );",
			ok:     true,
			cmd:    "insert",
			sql:    "insert into user_rights (user_id,right,group_id,value)",
		},
		{
			rawSQL: " bad data ",
			ok:     false,
		},
		{
			rawSQL: "into insert",
			ok:     false,
		},
		{
			rawSQL: "insert-into",
			ok:     false,
		},
		{
			rawSQL: "UPDATE receipt_invoices dest, (  SELECT      `receipt_id`,     CAST((net * 100) " +
				"/ 112 AS DECIMAL (11, 2)) witoutvat     FROM     receipt   WHERE CAST((net * 100) / 112 " +
				"AS DECIMAL (11, 2)) != total       AND vat_percentage = 12  ) src SET  dest.price = src.witoutvat,  " +
				"dest.amount = src.witoutvat WHERE col_tobefixed = 1   AND dest.`receipt_id` = src.receipt_id ;",
			ok:  true,
			cmd: "update",
			sql: "update receipt_invoices dest, ( select `receipt_id`, cast((net * ?) / ? as decimal (?)) " +
				"witoutvat from receipt where cast((net * ?) / ? as decimal (?)) != total and vat_percentage = ? ) " +
				"src set dest.price = src.witoutvat, dest.amount = src.witoutvat " +
				"where col_tobefixed = ? and dest.`receipt_id` = src.receipt_id ;",
		},
		{
			rawSQL: "update `ui_user` set `id`=?,`ename`='san.zhang',`cname`='张'三',`email`='san.'zhang@dtl.cn',`department_id`=?," +
				"`productline_id`=?,`path`='技术中心>核心基础设施>应用运维组',`dtl_id`='dt?',`status`='active',`utime`='?-?-? ?:?:?' " +
				"where `id` = ?",
			ok:  true,
			cmd: "update",
			sql: "update `ui_user` set `id`=?,`ename`='?',`cname`='?',`email`='?',`department_id`=?,`productline_id`=?," +
				"`path`='?',`dtl_id`='?',`status`='?',`utime`='?' where `id` = ?",
		},
		{
			rawSQL: "UPDATE `users` SET `age`=age - 1 WHERE id IN (1,2,3,4) AND name IN (\"san.zhang\",\"si.li\") AND id = 1 AND age > 1",
			ok:     true,
			cmd:    "update",
			sql:    "update `users` set `age`=age - ? where id in (?) and name in (?) and id = ? and age > ?",
		},
		{
			rawSQL: "SELECT * FROM `users` WHERE id IN (1,2,3,4) AND name IN (\"san.zhang\",\"si.li\")",
			ok:     true,
			cmd:    "select",
			sql:    "select * from `users` where id in (?) and name in (?)",
		},
		{
			rawSQL: "select `vals` from `white_list` where type in (?) and vals like '   hyman.cheng  %  '",
			ok:     true,
			cmd:    "select",
			sql:    "select `vals` from `white_list` where type in (?) and vals like '?'",
		},
	}
	for _, example := range examples {
		cmd, sql, ok := Filter.FilterSQL(example.rawSQL)
		if ok != example.ok {
			t.Errorf("filter sql failed, expect: %t, real: %t", example.ok, ok)
		} else {
			if !example.ok {
				continue
			}
			if cmd != example.cmd {
				t.Errorf("filter sql failed, raw sql: %s, expect cmd: %s, real cmd: %s",
					example.rawSQL, example.cmd, cmd)
			}
			if sql != example.sql {
				t.Errorf("filter sql failed, raw sql: %s, expect sql: %s, real sql: %s",
					example.rawSQL, example.sql, sql)
			}
		}
	}
}
