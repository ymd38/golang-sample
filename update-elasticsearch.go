/*
	mysqlから指定された時間以降で更新のあったデータをelastic searchへ更新をかける
	ex)10分前以降に更新のあったドキュメントを更新する by mac
	 go lodge_update_elasticsearch -t `date -v-10M +'%Y%m%d%H%M%S'`

	 MEMO:
	  httpコネクションはデフォルトで同一ホストで最大2で制限　http.DefaultMaxIdleConnsPerHostを変更し最大コネクション数を変更可能
		dbコネクションはいい感じにpoolしてくれるぽい
*/
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	ES_URL = "https://xxxxxxxxxx/sample/articles/"
	DB_DSN = "<username>:<password>@tcp(<host>)/<dbname>"
)

/* ドキュメント情報 */
type Article struct {
	Id         int      `json:"id"`
	User       string   `json:"user"`
	Title      string   `json:"title"`
	Body       string   `json:"body"`
	Url        string   `json:"url"`
	Updated_at string   `json:"updated_at"`
	Tags       []string `json:"tags,omitempty"`
}

//コネクション系
var db *sql.DB
var http_client *http.Client

/* 更新対象となるIDを取得 */
func getUpdateIds(criteria_time string) (id_list []string) {
	rows, err := db.Query("select id from articles where updated_at >=" + criteria_time)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		id_list = append(id_list, id)
		if err != nil {
			panic(err.Error())
		}
	}
	return id_list
}

/* lodgeデータの詳細取得 */
func getArticle(id string) (article Article) {
	rows, err := db.Query("select a.id,b.name,a.title,a.body,a.updated_at from articles a INNER JOIN users b ON a.user_id = b.id where a.id =" + id)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&article.Id, &article.User, &article.Title, &article.Body, &article.Updated_at)
		if err != nil {
			panic(err.Error())
		}
		article.Url = "http://lodge.mediba.jp/articles/" + id
	}

	return article
}

/* ドキュメントに紐づくタグを取得 */
func getTags(id string) []string {
	var tags string

	rows, err := db.Query("select t.name from taggings ti INNER JOIN tags t ON ti.tag_id = t.id and ti.taggable_type='Article' INNER JOIN articles a ON ti.taggable_id=a.id AND a.id=" + id)
	if err != nil {
		panic(err.Error())
	}

	for i := 0; rows.Next(); i++ {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			panic(err.Error())
		}
		if i == 0 {
			tags += tag
		} else {
			tags += "," + tag
		}
	}

	return strings.Split(tags, ",")
}

/* lodgeデータをelastic searchに更新 */
func upsertES(article Article) {
	/* jsonエンコード */
	article_json, err := json.Marshal(article)
	if err != nil {
		panic(err.Error())
	}

	/* httpリクエスト作成 */
	req_body := strings.NewReader(string(article_json))
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%d", ES_URL, article.Id), req_body)
	if err != nil {
		panic(err.Error())
	}

	/* elastic searchへ更新処理 */
	resp, _ := http_client.Do(req)
	defer resp.Body.Close()

	/* responseの表示 */
	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
}

/* main関数 */
func main() {
	/* 更新対象となる基準日時をパラメータで受け取る */
	criteria_time := flag.String("t", "YYYYMMDDHHMMSS", "The reference date and time.")
	flag.Parse()
	fmt.Println(*criteria_time)

	/* パラメータのフォーマットチェック */
	_, err := time.Parse("20060102150405 -0700", *criteria_time+" +0900")
	if err != nil {
		fmt.Println("Please check --help option.")
		return
	}

	/* DB初期処理 */
	db, err = sql.Open("mysql", DB_DSN)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	/* httpコネクション */
	http_client = new(http.Client)

	/* 更新対象となるIDを取得 */
	target_id_list := getUpdateIds(*criteria_time)

	/* IDから詳細を取得し、elastic search へ更新処理を並列処理で行なう */
	var wg sync.WaitGroup
	for i := 0; i < len(target_id_list); i++ {
		wg.Add(1)
		go func(id string) {
			fmt.Println("goroutine->", id)
			defer wg.Done()
			/* ドキュメント詳細取得 */
			article := getArticle(id)
			article.Tags = getTags(id)
			//elastic search へ更新処理を行なう
			upsertES(article)
		}(target_id_list[i])
	}
	wg.Wait()
}
