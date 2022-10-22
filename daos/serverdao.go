package daos

import (
	"database/sql"

	"github.com/nhuphuocnguyen/vcs_sms/models"
)

type ServerDAO struct {
	Db *sql.DB
}

func (sd *ServerDAO) CreateServer(server models.Server) (int, error) {
	sqlStatement := `
	INSERT INTO vcs_server (server_name, status, created_time,last_updated,ipv4)
	VALUES ($1, $2, $3, $4,$5)
	RETURNING server_id`
	id := 0
	err := sd.Db.QueryRow(sqlStatement, server.Server_name, server.Status, server.Created_time, server.Last_updated, server.Ipv4).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
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
func (sd *ServerDAO) Listserver() ([]models.Server,error){
    rows, err := sd.Db.Query("SELECT server_id,server_name, status, created_time,last_updated,ipv4 FROM vcs_server order by server_id ")
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
	return array,nil
}

