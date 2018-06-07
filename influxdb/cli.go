package influxdb

import (
	"log"
	"fmt"
	"time"
	"github.com/influxdata/influxdb/client/v2"
)
type Cli struct{
	MyDB,Username,Password,Addr,Precision string
	session client.Client
}
func (c * Cli) InitHttp() (err error) {
	// Create a new HTTPClient
	c.session, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     c.Addr,	// "http://localhost:8086",
		Username: c.Username,
		Password: c.Password,
	})
	if c.Precision==""{	
		c.Precision="s"		// 默认设置为秒 
	}
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
// 存数据   这个数据库自动扩展字段
func (c * Cli) WriteDB(table string,tags map[string]string,fields map[string]interface{})error {
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  c.MyDB,
		Precision: c.Precision,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	// Create a point and add to batch
	// tags := map[string]string{"mem": "mem-total","mem2":"mem2-total"}
	// fields := map[string]interface{}{
	// 	"all": 4096,
	// 	"used": 3308,
	// }
	pt, err := client.NewPoint(table, tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
		return err
	}
	bp.AddPoint(pt)
	// Write the batch
	if err := c.session.Write(bp); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
// QueryDB convenience function to query the database
func (c * Cli) QueryDB(cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: c.MyDB,
	}
	if response, err := c.session.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
// 创建数据库
func (c * Cli) CreateDB(db string) error {
	_, err := c.QueryDB(fmt.Sprintf("CREATE DATABASE %s", db))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
