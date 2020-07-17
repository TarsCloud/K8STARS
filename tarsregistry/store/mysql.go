package store

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlDriver struct {
	db *sql.DB
}

func NewMysqlDB(dsn string) (Store, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &mysqlDriver{db: db}, nil
}

func (m *mysqlDriver) RegisterNode(ctx context.Context, nodeName string) error {
	nodeObj := fmt.Sprintf("tars.tarsnode.NodeObj@tcp -h %s -p 0", nodeName)
	sql := `insert into t_node_info(node_name, node_obj, endpoint_ip, endpoint_port, 
				last_reg_time, last_heartbeat, setting_state, present_state) 
				values(?, ?, ?, ?, CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP(), ?, ?)
				ON DUPLICATE KEY UPDATE last_heartbeat = CURRENT_TIMESTAMP()`
	_, err := m.db.ExecContext(ctx, sql, nodeName, nodeObj, nodeName, 0, "active", "active")
	return err
}

func (m *mysqlDriver) RegisterServer(ctx context.Context, conf *ServerConf) error {
	sql := `insert into t_server_conf(
		application, server_name, node_name, patch_version, present_state,
		enable_set, set_name, set_area, set_group, grid_flag,
		server_type, setting_state, registry_timestamp, patch_time, posttime
	 )
	 values(
		?, ?, ?, ?, ?,
		?, ?, ?, ?, ?,
		"tars_cpp", "active", CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP()
	 )
	 ON DUPLICATE KEY UPDATE patch_version=?, present_state=?, 
	 enable_set=?, set_name=?, set_area=?, set_group=?, grid_flag=?,
	 server_type="tars_cpp", present_state="active", registry_timestamp=CURRENT_TIMESTAMP(), patch_time=CURRENT_TIMESTAMP(), posttime=CURRENT_TIMESTAMP()
	 `
	_, err := m.db.ExecContext(ctx, sql, conf.Application, conf.Server, conf.NodeName, conf.Version, conf.State,
		conf.EnableSet, conf.SetName, conf.SetGroup, conf.SetArea, conf.GridFlag,
		conf.Version, conf.State,
		conf.EnableSet, conf.SetName, conf.SetGroup, conf.SetArea, conf.GridFlag,
	)
	return err
}

func (m *mysqlDriver) RegistryAdapter(ctx context.Context, confs []*AdapterConf) error {
	for _, conf := range confs {
		sql := `insert into t_adapter_conf(
			application, server_name, node_name,
			adapter_name, servant, thread_num, endpoint, 
			protocol, max_connections, queuecap, queuetimeout
		 )
		 values(
			?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?, ?
		 )
		 ON DUPLICATE KEY UPDATE
		 adapter_name=?, servant=?, thread_num=?, endpoint=?,
		 protocol=?, max_connections=?, queuecap=?, queuetimeout=?`
		_, err := m.db.ExecContext(ctx, sql, conf.Application, conf.Server, conf.NodeName,
			conf.AdapterName, conf.Servant, conf.ThreadNum, conf.Endpoint,
			conf.Protocol, conf.MaxConns, conf.QueueCap, conf.QueueTimeout,
			conf.AdapterName, conf.Servant, conf.ThreadNum, conf.Endpoint,
			conf.Protocol, conf.MaxConns, conf.QueueCap, conf.QueueTimeout,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *mysqlDriver) DeleteNodeConf(ctx context.Context, nodeName string) error {
	sql := "delete from t_adapter_conf where node_name=?"
	if _, err := m.db.ExecContext(ctx, sql, nodeName); err != nil {
		return err
	}
	sql = "delete from t_server_conf where node_name=?"
	if _, err := m.db.ExecContext(ctx, sql, nodeName); err != nil {
		return err
	}
	sql = "delete from t_node_info where node_name=?"
	if _, err := m.db.ExecContext(ctx, sql, nodeName); err != nil {
		return err
	}
	return nil
}

func (m *mysqlDriver) KeepAliveNode(ctx context.Context, nodeName string) error {
	sql := "update t_node_info set last_heartbeat = CURRENT_TIMESTAMP() where node_name = ?"
	_, err := m.db.ExecContext(ctx, sql, nodeName)
	return err
}

func (m *mysqlDriver) SetServerState(ctx context.Context, nodeName, state string) error {
	sql := "update t_server_conf set present_state = ? where node_name = ?"
	_, err := m.db.ExecContext(ctx, sql, state, nodeName)
	sql = "update t_node_info set present_state = ? where node_name = ?"
	_, err = m.db.ExecContext(ctx, sql, state, nodeName)
	return err
}
