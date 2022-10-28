package daos

import (
	"database/sql"

	"github.com/nhuphuocnguyen/vcs_sms/models"
)

type ServerDAO struct {
	Db *sql.DB
}

func (sd *ServerDAO) CreateServer(server models.Server)  error {
	sqlStatement := `
	INSERT INTO vcs_server (server_id,server_name, status, created_time,last_updated,ipv4)
	VALUES ($1, $2, $3, $4,$5,$6)`
	err := sd.Db.QueryRow(sqlStatement, server.Server_id, server.Server_name, server.Status, server.Created_time, server.Last_updated, server.Ipv4).Err()
	if err != nil {
		return err
	}
	return  nil
}

func (sd *ServerDAO) UpdateServer(server models.Server, ids string) (string, error) {
	sqlStatement := `
	UPDATE vcs_server
	SET server_name=$2, status=$3,last_updated=$4,ipv4=$5
	WHERE server_id = $1;`
	_, err := sd.Db.Exec(sqlStatement, ids, server.Server_name, server.Status, server.Last_updated, server.Ipv4)
	if err != nil {
		return "", err
	}
	return ids, nil
}

func (sd *ServerDAO) DeleteServer(server models.Server, ids string) (string, error) {
	sqlStatement := `
	DELETE from vcs_server
	WHERE server_id = $1;`
	_, err := sd.Db.Exec(sqlStatement, ids)
	if err != nil {
		return "", err
	}
	return ids, nil
}
func (sd *ServerDAO) Get(id string) (*models.Server, error) {
	sqlStatement := `
	Select server_id,server_name, status, created_time,last_updated,ipv4 from vcs_server
	Where server_id=$1;`
	row := sd.Db.QueryRow(sqlStatement, id)
	var server models.Server
	err := row.Scan(&server.Server_id, &server.Server_name, &server.Status, &server.Created_time, &server.Last_updated, &server.Ipv4)
	if err != nil {
		return nil, err
	}
	return &server, nil

}

func (sd *ServerDAO) Count() (int, error) {
	var count int
	row := sd.Db.QueryRow("SELECT COUNT(*) FROM vcs_server")
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil

}
func (sd *ServerDAO) Listserver(sort string, option string, from int, size int) ([]models.Server, error) {
	var rows *sql.Rows
	var err error
	if option == "DESC" {
		rows, err = sd.Db.Query(`
	SELECT server_id,server_name, status, created_time,last_updated,ipv4 
	FROM vcs_server order by $1 DESC
	Offset $2
	limit $3;`, sort, from, size)
	}
	if option == "ASC" {
		rows, err = sd.Db.Query(`
	SELECT server_id,server_name, status, created_time,last_updated,ipv4 
	FROM vcs_server order by $1 ASC
	Offset $2
	limit $3;`, sort, from, size)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var array []models.Server
	for rows.Next() {
		sever := models.Server{}

		err = rows.Scan(&sever.Server_id, &sever.Server_name, &sever.Status, &sever.Created_time, &sever.Last_updated, &sever.Ipv4)
		if err != nil {
			return nil, err
		}

		array = append(array, sever)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return array, nil
}
